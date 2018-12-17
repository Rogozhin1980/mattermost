// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/store"
	"github.com/mattermost/mattermost-server/utils"
)

const (
	PENDING_POST_IDS_CACHE_SIZE = 25000
	PENDING_POST_IDS_CACHE_TTL  = 30 * time.Second
)

func (a *App) CreatePostAsUser(post *model.Post, clearPushNotifications bool) (*model.Post, *model.AppError) {
	// Check that channel has not been deleted
	result := <-a.Srv.Store.Channel().Get(post.ChannelId, true)
	if result.Err != nil {
		err := model.NewAppError("CreatePostAsUser", "api.context.invalid_param.app_error", map[string]interface{}{"Name": "post.channel_id"}, result.Err.Error(), http.StatusBadRequest)
		return nil, err
	}
	channel := result.Data.(*model.Channel)

	if strings.HasPrefix(post.Type, model.POST_SYSTEM_MESSAGE_PREFIX) {
		err := model.NewAppError("CreatePostAsUser", "api.context.invalid_param.app_error", map[string]interface{}{"Name": "post.type"}, "", http.StatusBadRequest)
		return nil, err
	}

	if channel.DeleteAt != 0 {
		err := model.NewAppError("createPost", "api.post.create_post.can_not_post_to_deleted.error", nil, "", http.StatusBadRequest)
		return nil, err
	}

	rp, err := a.CreatePost(post, channel, true)
	if err != nil {
		if err.Id == "api.post.create_post.root_id.app_error" ||
			err.Id == "api.post.create_post.channel_root_id.app_error" ||
			err.Id == "api.post.create_post.parent_id.app_error" {
			err.StatusCode = http.StatusBadRequest
		}

		if err.Id == "api.post.create_post.town_square_read_only" {
			result := <-a.Srv.Store.User().Get(post.UserId)
			if result.Err != nil {
				return nil, result.Err
			}
			user := result.Data.(*model.User)

			T := utils.GetUserTranslations(user.Locale)
			a.SendEphemeralPost(
				post.UserId,
				&model.Post{
					ChannelId: channel.Id,
					ParentId:  post.ParentId,
					RootId:    post.RootId,
					UserId:    post.UserId,
					Message:   T("api.post.create_post.town_square_read_only"),
					CreateAt:  model.GetMillis() + 1,
				},
			)
		}
		return nil, err
	}

	// Update the LastViewAt only if the post does not have from_webhook prop set (eg. Zapier app)
	if _, ok := post.Props["from_webhook"]; !ok {
		if _, err := a.MarkChannelsAsViewed([]string{post.ChannelId}, post.UserId, clearPushNotifications); err != nil {
			mlog.Error(fmt.Sprintf("Encountered error updating last viewed, channel_id=%s, user_id=%s, err=%v", post.ChannelId, post.UserId, err))
		}
	}

	return rp, nil
}

func (a *App) CreatePostMissingChannel(post *model.Post, triggerWebhooks bool) (*model.Post, *model.AppError) {
	result := <-a.Srv.Store.Channel().Get(post.ChannelId, true)
	if result.Err != nil {
		return nil, result.Err
	}
	channel := result.Data.(*model.Channel)

	return a.CreatePost(post, channel, triggerWebhooks)
}

// deduplicateCreatePost attempts to make posting idempotent within a caching window.
func (a *App) deduplicateCreatePost(post *model.Post) (foundPost *model.Post, err *model.AppError) {
	// We rely on the client sending the pending post id across "duplicate" requests. If there
	// isn't one, we can't deduplicate, so allow creation normally.
	if post.PendingPostId == "" {
		return nil, nil
	}

	const unknownPostId = ""

	for attempt := 1; attempt <= 3; attempt++ {
		// Query the cache atomically for the given pending post id, saving a record if
		// it hasn't previously been seen.
		value, loaded := a.Srv.seenPendingPostIdsCache.GetOrAdd(post.PendingPostId, unknownPostId, PENDING_POST_IDS_CACHE_TTL)

		// If we were the first thread to save this pending post id into the cache,
		// proceed with create post normally.
		if !loaded {
			return nil, nil
		}

		postId := value.(string)

		// If another thread saved the cache record, but hasn't yet updated it with the
		// actual post id, wait a bit and try again. If the cache record expires, we'll
		// still end up with a duplicate post as expected. Note that we'll only ever be
		// doing any of this when racing on duplicate create post requests.
		if postId == unknownPostId {
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Return the created post back to the client, making the API call feel idempotent,
		// and avoiding having special-case logic in the client.
		post, err := a.GetSinglePost(postId)
		if err != nil {
			return nil, model.NewAppError("deduplicateCreatePost", "api.post.deduplicate_create_post.failed_to_get", nil, err.Error(), http.StatusInternalServerError)
		}

		return post, nil
	}

	// If we failed to create a cache entry above, or resolve to a created post, the attempt to
	// deduplicate failed altogether.
	return nil, model.NewAppError("deduplicateCreatePost", "api.post.deduplicate_create_post.failed", nil, "", http.StatusInternalServerError)
}

func (a *App) CreatePost(post *model.Post, channel *model.Channel, triggerWebhooks bool) (savedPost *model.Post, err *model.AppError) {
	if foundPost, err := a.deduplicateCreatePost(post); err != nil {
		return nil, err
	} else if foundPost != nil {
		return foundPost, nil
	}

	// If we get this far, we've recorded the client-provided pending post id to the cache.
	// Remove it if we fail below, allowing a proper retry by the client.
	defer func() {
		if post.PendingPostId == "" {
			return
		}

		if err != nil {
			a.Srv.seenPendingPostIdsCache.Remove(post.PendingPostId)
			return
		}

		a.Srv.seenPendingPostIdsCache.AddWithExpiresInSecs(post.PendingPostId, savedPost.Id, int64(PENDING_POST_IDS_CACHE_TTL.Seconds()))
	}()

	post.SanitizeProps()

	var pchan store.StoreChannel
	if len(post.RootId) > 0 {
		pchan = a.Srv.Store.Post().Get(post.RootId)
	}

	result := <-a.Srv.Store.User().Get(post.UserId)
	if result.Err != nil {
		return nil, result.Err
	}
	user := result.Data.(*model.User)

	if a.License() != nil && *a.Config().TeamSettings.ExperimentalTownSquareIsReadOnly &&
		!post.IsSystemMessage() &&
		channel.Name == model.DEFAULT_CHANNEL &&
		!a.RolesGrantPermission(user.GetRoles(), model.PERMISSION_MANAGE_SYSTEM.Id) {
		return nil, model.NewAppError("createPost", "api.post.create_post.town_square_read_only", nil, "", http.StatusForbidden)
	}

	// Verify the parent/child relationships are correct
	var parentPostList *model.PostList
	if pchan != nil {
		result = <-pchan
		if result.Err != nil {
			return nil, model.NewAppError("createPost", "api.post.create_post.root_id.app_error", nil, "", http.StatusBadRequest)
		}
		parentPostList = result.Data.(*model.PostList)
		if len(parentPostList.Posts) == 0 || !parentPostList.IsChannelId(post.ChannelId) {
			return nil, model.NewAppError("createPost", "api.post.create_post.channel_root_id.app_error", nil, "", http.StatusInternalServerError)
		}

		if post.ParentId == "" {
			post.ParentId = post.RootId
		}

		if post.RootId != post.ParentId {
			parent := parentPostList.Posts[post.ParentId]
			if parent == nil {
				return nil, model.NewAppError("createPost", "api.post.create_post.parent_id.app_error", nil, "", http.StatusInternalServerError)
			}
		}
	}

	post.Hashtags, _ = model.ParseHashtags(post.Message)

	if err := a.FillInPostProps(post, channel); err != nil {
		return nil, err
	}

	// Temporary fix so old plugins don't clobber new fields in SlackAttachment struct, see MM-13088
	if attachments, ok := post.Props["attachments"].([]*model.SlackAttachment); ok {
		jsonAttachments, err := json.Marshal(attachments)
		if err == nil {
			attachmentsInterface := []interface{}{}
			err = json.Unmarshal(jsonAttachments, &attachmentsInterface)
			post.Props["attachments"] = attachmentsInterface
		}
		if err != nil {
			mlog.Error("Could not convert post attachments to map interface, err=%s" + err.Error())
		}
	}

	if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
		var rejectionError *model.AppError
		pluginContext := a.PluginContext()
		pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
			replacementPost, rejectionReason := hooks.MessageWillBePosted(pluginContext, post)
			if rejectionReason != "" {
				rejectionError = model.NewAppError("createPost", "Post rejected by plugin. "+rejectionReason, nil, "", http.StatusBadRequest)
				return false
			}
			if replacementPost != nil {
				post = replacementPost
			}

			return true
		}, plugin.MessageWillBePostedId)
		if rejectionError != nil {
			return nil, rejectionError
		}
	}

	result = <-a.Srv.Store.Post().Save(post)
	if result.Err != nil {
		return nil, result.Err
	}
	rpost := result.Data.(*model.Post)

	// Update the mapping from pending post id to the actual post id, for any clients that
	// might be duplicating requests.
	a.Srv.seenPendingPostIdsCache.AddWithExpiresInSecs(post.PendingPostId, rpost.Id, int64(PENDING_POST_IDS_CACHE_TTL.Seconds()))

	if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
		a.Srv.Go(func() {
			pluginContext := a.PluginContext()
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				hooks.MessageHasBeenPosted(pluginContext, rpost)
				return true
			}, plugin.MessageHasBeenPostedId)
		})
	}

	esInterface := a.Elasticsearch
	if esInterface != nil && *a.Config().ElasticsearchSettings.EnableIndexing {
		a.Srv.Go(func() {
			esInterface.IndexPost(rpost, channel.TeamId)
		})
	}

	if a.Metrics != nil {
		a.Metrics.IncrementPostCreate()
	}

	if len(post.FileIds) > 0 {
		// There's a rare bug where the client sends up duplicate FileIds so protect against that
		post.FileIds = utils.RemoveDuplicatesFromStringArray(post.FileIds)

		for _, fileId := range post.FileIds {
			if result := <-a.Srv.Store.FileInfo().AttachToPost(fileId, post.Id); result.Err != nil {
				mlog.Error(fmt.Sprintf("Encountered error attaching files to post, post_id=%s, user_id=%s, file_ids=%v, err=%v", post.Id, post.FileIds, post.UserId, result.Err), mlog.String("post_id", post.Id))
			}
		}

		if a.Metrics != nil {
			a.Metrics.IncrementPostFileAttachment(len(post.FileIds))
		}
	}

	// Normally, we would let the API layer call PreparePostForClient, but we do it here since it also needs
	// to be done when we send the post over the websocket in handlePostEvents
	rpost = a.PreparePostForClient(rpost)

	if err := a.handlePostEvents(rpost, user, channel, triggerWebhooks, parentPostList); err != nil {
		mlog.Error("Failed to handle post events", mlog.Err(err))
	}

	return rpost, nil
}

// FillInPostProps should be invoked before saving posts to fill in properties such as
// channel_mentions.
//
// If channel is nil, FillInPostProps will look up the channel corresponding to the post.
func (a *App) FillInPostProps(post *model.Post, channel *model.Channel) *model.AppError {
	channelMentions := post.ChannelMentions()
	channelMentionsProp := make(map[string]interface{})

	if len(channelMentions) > 0 {
		if channel == nil {
			result := <-a.Srv.Store.Channel().GetForPost(post.Id)
			if result.Err != nil {
				return model.NewAppError("FillInPostProps", "api.context.invalid_param.app_error", map[string]interface{}{"Name": "post.channel_id"}, result.Err.Error(), http.StatusBadRequest)
			}
			channel = result.Data.(*model.Channel)
		}

		mentionedChannels, err := a.GetChannelsByNames(channelMentions, channel.TeamId)
		if err != nil {
			return err
		}

		for _, mentioned := range mentionedChannels {
			if mentioned.Type == model.CHANNEL_OPEN {
				channelMentionsProp[mentioned.Name] = map[string]interface{}{
					"display_name": mentioned.DisplayName,
				}
			}
		}
	}

	if len(channelMentionsProp) > 0 {
		post.AddProp("channel_mentions", channelMentionsProp)
	} else if post.Props != nil {
		delete(post.Props, "channel_mentions")
	}

	return nil
}

func (a *App) handlePostEvents(post *model.Post, user *model.User, channel *model.Channel, triggerWebhooks bool, parentPostList *model.PostList) *model.AppError {
	var team *model.Team
	if len(channel.TeamId) > 0 {
		result := <-a.Srv.Store.Team().Get(channel.TeamId)
		if result.Err != nil {
			return result.Err
		}
		team = result.Data.(*model.Team)
	} else {
		// Blank team for DMs
		team = &model.Team{}
	}

	a.InvalidateCacheForChannel(channel)
	a.InvalidateCacheForChannelPosts(channel.Id)

	if _, err := a.SendNotifications(post, team, channel, user, parentPostList); err != nil {
		return err
	}

	if triggerWebhooks {
		a.Srv.Go(func() {
			if err := a.handleWebhookEvents(post, team, channel, user); err != nil {
				mlog.Error(err.Error())
			}
		})
	}

	return nil
}

func (a *App) SendEphemeralPost(userId string, post *model.Post) *model.Post {
	post.Type = model.POST_EPHEMERAL

	// fill in fields which haven't been specified which have sensible defaults
	if post.Id == "" {
		post.Id = model.NewId()
	}
	if post.CreateAt == 0 {
		post.CreateAt = model.GetMillis()
	}
	if post.Props == nil {
		post.Props = model.StringInterface{}
	}

	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_EPHEMERAL_MESSAGE, "", post.ChannelId, userId, nil)
	message.Add("post", a.PreparePostForClient(post).ToJson())
	a.Publish(message)

	return post
}

func (a *App) UpdatePost(post *model.Post, safeUpdate bool) (*model.Post, *model.AppError) {
	post.SanitizeProps()

	result := <-a.Srv.Store.Post().Get(post.Id)
	if result.Err != nil {
		return nil, result.Err
	}
	oldPost := result.Data.(*model.PostList).Posts[post.Id]

	if oldPost == nil {
		err := model.NewAppError("UpdatePost", "api.post.update_post.find.app_error", nil, "id="+post.Id, http.StatusBadRequest)
		return nil, err
	}

	if oldPost.DeleteAt != 0 {
		err := model.NewAppError("UpdatePost", "api.post.update_post.permissions_details.app_error", map[string]interface{}{"PostId": post.Id}, "", http.StatusBadRequest)
		return nil, err
	}

	if oldPost.IsSystemMessage() {
		err := model.NewAppError("UpdatePost", "api.post.update_post.system_message.app_error", nil, "id="+post.Id, http.StatusBadRequest)
		return nil, err
	}

	if a.License() != nil {
		if *a.Config().ServiceSettings.PostEditTimeLimit != -1 && model.GetMillis() > oldPost.CreateAt+int64(*a.Config().ServiceSettings.PostEditTimeLimit*1000) && post.Message != oldPost.Message {
			err := model.NewAppError("UpdatePost", "api.post.update_post.permissions_time_limit.app_error", map[string]interface{}{"timeLimit": *a.Config().ServiceSettings.PostEditTimeLimit}, "", http.StatusBadRequest)
			return nil, err
		}
	}

	newPost := &model.Post{}
	*newPost = *oldPost

	if newPost.Message != post.Message {
		newPost.Message = post.Message
		newPost.EditAt = model.GetMillis()
		newPost.Hashtags, _ = model.ParseHashtags(post.Message)
	}

	if !safeUpdate {
		newPost.IsPinned = post.IsPinned
		newPost.HasReactions = post.HasReactions
		newPost.FileIds = post.FileIds
		newPost.Props = post.Props
	}

	if err := a.FillInPostProps(post, nil); err != nil {
		return nil, err
	}

	if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
		var rejectionReason string
		pluginContext := a.PluginContext()
		pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
			newPost, rejectionReason = hooks.MessageWillBeUpdated(pluginContext, newPost, oldPost)
			return post != nil
		}, plugin.MessageWillBeUpdatedId)
		if newPost == nil {
			return nil, model.NewAppError("UpdatePost", "Post rejected by plugin. "+rejectionReason, nil, "", http.StatusBadRequest)
		}
	}

	result = <-a.Srv.Store.Post().Update(newPost, oldPost)
	if result.Err != nil {
		return nil, result.Err
	}
	rpost := result.Data.(*model.Post)

	if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
		a.Srv.Go(func() {
			pluginContext := a.PluginContext()
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				hooks.MessageHasBeenUpdated(pluginContext, newPost, oldPost)
				return true
			}, plugin.MessageHasBeenUpdatedId)
		})
	}

	esInterface := a.Elasticsearch
	if esInterface != nil && *a.Config().ElasticsearchSettings.EnableIndexing {
		a.Srv.Go(func() {
			rchannel := <-a.Srv.Store.Channel().GetForPost(rpost.Id)
			if rchannel.Err != nil {
				mlog.Error(fmt.Sprintf("Couldn't get channel %v for post %v for Elasticsearch indexing.", rpost.ChannelId, rpost.Id))
				return
			}
			esInterface.IndexPost(rpost, rchannel.Data.(*model.Channel).TeamId)
		})
	}

	rpost = a.PreparePostForClient(rpost)

	a.sendUpdatedPostEvent(rpost)

	a.InvalidateCacheForChannelPosts(rpost.ChannelId)

	return rpost, nil
}

func (a *App) PatchPost(postId string, patch *model.PostPatch) (*model.Post, *model.AppError) {
	post, err := a.GetSinglePost(postId)
	if err != nil {
		return nil, err
	}

	post.Patch(patch)

	updatedPost, err := a.UpdatePost(post, false)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (a *App) sendUpdatedPostEvent(post *model.Post) {
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_POST_EDITED, "", post.ChannelId, "", nil)
	message.Add("post", post.ToJson())
	a.Publish(message)
}

func (a *App) GetPostsPage(channelId string, page int, perPage int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetPosts(channelId, page*perPage, perPage, true)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetPosts(channelId string, offset int, limit int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetPosts(channelId, offset, limit, true)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetPostsEtag(channelId string) string {
	return (<-a.Srv.Store.Post().GetEtag(channelId, true)).Data.(string)
}

func (a *App) GetPostsSince(channelId string, time int64) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetPostsSince(channelId, time, true)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetSinglePost(postId string) (*model.Post, *model.AppError) {
	result := <-a.Srv.Store.Post().GetSingle(postId)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.Post), nil
}

func (a *App) GetPostThread(postId string) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().Get(postId)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetFlaggedPosts(userId string, offset int, limit int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetFlaggedPosts(userId, offset, limit)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetFlaggedPostsForTeam(userId, teamId string, offset int, limit int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetFlaggedPostsForTeam(userId, teamId, offset, limit)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetFlaggedPostsForChannel(userId, channelId string, offset int, limit int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetFlaggedPostsForChannel(userId, channelId, offset, limit)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetPermalinkPost(postId string, userId string) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().Get(postId)
	if result.Err != nil {
		return nil, result.Err
	}
	list := result.Data.(*model.PostList)

	if len(list.Order) != 1 {
		return nil, model.NewAppError("getPermalinkTmp", "api.post_get_post_by_id.get.app_error", nil, "", http.StatusNotFound)
	}
	post := list.Posts[list.Order[0]]

	channel, err := a.GetChannel(post.ChannelId)
	if err != nil {
		return nil, err
	}

	if err = a.JoinChannel(channel, userId); err != nil {
		return nil, err
	}

	return list, nil
}

func (a *App) GetPostsBeforePost(channelId, postId string, page, perPage int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetPostsBefore(channelId, postId, perPage, page*perPage)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetPostsAfterPost(channelId, postId string, page, perPage int) (*model.PostList, *model.AppError) {
	result := <-a.Srv.Store.Post().GetPostsAfter(channelId, postId, perPage, page*perPage)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) GetPostsAroundPost(postId, channelId string, offset, limit int, before bool) (*model.PostList, *model.AppError) {
	var pchan store.StoreChannel
	if before {
		pchan = a.Srv.Store.Post().GetPostsBefore(channelId, postId, limit, offset)
	} else {
		pchan = a.Srv.Store.Post().GetPostsAfter(channelId, postId, limit, offset)
	}

	result := <-pchan
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.(*model.PostList), nil
}

func (a *App) DeletePost(postId, deleteByID string) (*model.Post, *model.AppError) {
	result := <-a.Srv.Store.Post().GetSingle(postId)
	if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	}
	post := result.Data.(*model.Post)

	if result := <-a.Srv.Store.Post().Delete(postId, model.GetMillis(), deleteByID); result.Err != nil {
		return nil, result.Err
	}

	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_POST_DELETED, "", post.ChannelId, "", nil)
	message.Add("post", a.PreparePostForClient(post).ToJson())
	a.Publish(message)

	a.Srv.Go(func() {
		a.DeletePostFiles(post)
	})
	a.Srv.Go(func() {
		a.DeleteFlaggedPosts(post.Id)
	})

	esInterface := a.Elasticsearch
	if esInterface != nil && *a.Config().ElasticsearchSettings.EnableIndexing {
		a.Srv.Go(func() {
			esInterface.DeletePost(post)
		})
	}

	a.InvalidateCacheForChannelPosts(post.ChannelId)

	return post, nil
}

func (a *App) DeleteFlaggedPosts(postId string) {
	if result := <-a.Srv.Store.Preference().DeleteCategoryAndName(model.PREFERENCE_CATEGORY_FLAGGED_POST, postId); result.Err != nil {
		mlog.Warn(fmt.Sprintf("Unable to delete flagged post preference when deleting post, err=%v", result.Err))
		return
	}
}

func (a *App) DeletePostFiles(post *model.Post) {
	if len(post.FileIds) == 0 {
		return
	}

	if result := <-a.Srv.Store.FileInfo().DeleteForPost(post.Id); result.Err != nil {
		mlog.Warn(fmt.Sprintf("Encountered error when deleting files for post, post_id=%v, err=%v", post.Id, result.Err), mlog.String("post_id", post.Id))
	}
}

func (a *App) parseAndFetchChannelIdByNameFromInFilter(channelName, userId, teamId string, includeDeleted bool) (*model.Channel, error) {
	if strings.HasPrefix(channelName, "@") && strings.Contains(channelName, ",") {
		var userIds []string
		users, err := a.GetUsersByUsernames(strings.Split(channelName[1:], ","), false)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			userIds = append(userIds, user.Id)
		}

		channel, err := a.GetGroupChannel(userIds)
		if err != nil {
			return nil, err
		}
		return channel, nil
	}

	if strings.HasPrefix(channelName, "@") && !strings.Contains(channelName, ",") {
		user, err := a.GetUserByUsername(channelName[1:])
		if err != nil {
			return nil, err
		}
		channel, err := a.GetOrCreateDirectChannel(userId, user.Id)
		if err != nil {
			return nil, err
		}
		return channel, nil
	}

	channel, err := a.GetChannelByName(channelName, teamId, includeDeleted)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (a *App) SearchPostsInTeam(terms string, userId string, teamId string, isOrSearch bool, includeDeletedChannels bool, timeZoneOffset int, page, perPage int) (*model.PostSearchResults, *model.AppError) {
	paramsList := model.ParseSearchParams(terms, timeZoneOffset)
	includeDeleted := includeDeletedChannels && *a.Config().TeamSettings.ExperimentalViewArchivedChannels

	esInterface := a.Elasticsearch
	license := a.License()
	if esInterface != nil && *a.Config().ElasticsearchSettings.EnableSearching && license != nil && *license.Features.Elasticsearch {
		finalParamsList := []*model.SearchParams{}

		for _, params := range paramsList {
			params.OrTerms = isOrSearch
			// Don't allow users to search for "*"
			if params.Terms != "*" {
				// Convert channel names to channel IDs
				for idx, channelName := range params.InChannels {
					channel, err := a.parseAndFetchChannelIdByNameFromInFilter(channelName, userId, teamId, includeDeletedChannels)
					if err != nil {
						mlog.Error(fmt.Sprint(err))
						continue
					}
					params.InChannels[idx] = channel.Id
				}

				// Convert usernames to user IDs
				for idx, username := range params.FromUsers {
					if user, err := a.GetUserByUsername(username); err != nil {
						mlog.Error(fmt.Sprint(err))
					} else {
						params.FromUsers[idx] = user.Id
					}
				}

				finalParamsList = append(finalParamsList, params)
			}
		}

		// If the processed search params are empty, return empty search results.
		if len(finalParamsList) == 0 {
			return model.MakePostSearchResults(model.NewPostList(), nil), nil
		}

		// We only allow the user to search in channels they are a member of.
		userChannels, err := a.GetChannelsForUser(teamId, userId, includeDeleted)
		if err != nil {
			mlog.Error(fmt.Sprint(err))
			return nil, err
		}

		postIds, matches, err := a.Elasticsearch.SearchPosts(userChannels, finalParamsList, page, perPage)
		if err != nil {
			return nil, err
		}

		// Get the posts
		postList := model.NewPostList()
		if len(postIds) > 0 {
			presult := <-a.Srv.Store.Post().GetPostsByIds(postIds)
			if presult.Err != nil {
				return nil, presult.Err
			}
			for _, p := range presult.Data.([]*model.Post) {
				if p.DeleteAt == 0 {
					postList.AddPost(p)
					postList.AddOrder(p.Id)
				}
			}
		}

		return model.MakePostSearchResults(postList, matches), nil
	}

	if !*a.Config().ServiceSettings.EnablePostSearch {
		return nil, model.NewAppError("SearchPostsInTeam", "store.sql_post.search.disabled", nil, fmt.Sprintf("teamId=%v userId=%v", teamId, userId), http.StatusNotImplemented)
	}

	// Since we don't support paging we just return nothing for later pages
	if page > 0 {
		return model.MakePostSearchResults(model.NewPostList(), nil), nil
	}

	channels := []store.StoreChannel{}

	for _, params := range paramsList {
		params.IncludeDeletedChannels = includeDeleted
		params.OrTerms = isOrSearch
		// don't allow users to search for everything
		if params.Terms != "*" {
			for idx, channelName := range params.InChannels {
				if strings.HasPrefix(channelName, "@") {
					channel, err := a.parseAndFetchChannelIdByNameFromInFilter(channelName, userId, teamId, includeDeletedChannels)
					if err != nil {
						mlog.Error(fmt.Sprint(err))
						continue
					}
					params.InChannels[idx] = channel.Name
				}
			}
			channels = append(channels, a.Srv.Store.Post().Search(teamId, userId, params))
		}
	}

	posts := model.NewPostList()
	for _, channel := range channels {
		result := <-channel
		if result.Err != nil {
			return nil, result.Err
		}
		data := result.Data.(*model.PostList)
		posts.Extend(data)
	}

	posts.SortByCreateAt()

	return model.MakePostSearchResults(posts, nil), nil
}

func (a *App) GetFileInfosForPostWithMigration(postId string) ([]*model.FileInfo, *model.AppError) {
	pchan := a.Srv.Store.Post().GetSingle(postId)

	infos, err := a.GetFileInfosForPost(postId)
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		// No FileInfos were returned so check if they need to be created for this post
		result := <-pchan
		if result.Err != nil {
			return nil, result.Err
		}
		post := result.Data.(*model.Post)

		if len(post.Filenames) > 0 {
			a.Srv.Store.FileInfo().InvalidateFileInfosForPostCache(postId)
			// The post has Filenames that need to be replaced with FileInfos
			infos = a.MigrateFilenamesToFileInfos(post)
		}
	}

	return infos, nil
}

func (a *App) GetFileInfosForPost(postId string) ([]*model.FileInfo, *model.AppError) {
	result := <-a.Srv.Store.FileInfo().GetForPost(postId, false, true)
	if result.Err != nil {
		return nil, result.Err
	}

	return result.Data.([]*model.FileInfo), nil
}

func (a *App) PostWithProxyAddedToImageURLs(post *model.Post) *model.Post {
	if f := a.ImageProxyAdder(); f != nil {
		return post.WithRewrittenImageURLs(f)
	}
	return post
}

func (a *App) PostWithProxyRemovedFromImageURLs(post *model.Post) *model.Post {
	if f := a.ImageProxyRemover(); f != nil {
		return post.WithRewrittenImageURLs(f)
	}
	return post
}

func (a *App) PostPatchWithProxyRemovedFromImageURLs(patch *model.PostPatch) *model.PostPatch {
	if f := a.ImageProxyRemover(); f != nil {
		return patch.WithRewrittenImageURLs(f)
	}
	return patch
}

func (a *App) imageProxyConfig() (proxyType, proxyURL, options, siteURL string) {
	cfg := a.Config()

	if cfg.ServiceSettings.ImageProxyURL == nil || cfg.ServiceSettings.ImageProxyType == nil || cfg.ServiceSettings.SiteURL == nil {
		return
	}

	proxyURL = *cfg.ServiceSettings.ImageProxyURL
	proxyType = *cfg.ServiceSettings.ImageProxyType
	siteURL = *cfg.ServiceSettings.SiteURL

	if proxyURL == "" || proxyType == "" {
		return "", "", "", ""
	}

	if proxyURL[len(proxyURL)-1] != '/' {
		proxyURL += "/"
	}

	if siteURL == "" || siteURL[len(siteURL)-1] != '/' {
		siteURL += "/"
	}

	if cfg.ServiceSettings.ImageProxyOptions != nil {
		options = *cfg.ServiceSettings.ImageProxyOptions
	}

	return
}

func (a *App) ImageProxyAdder() func(string) string {
	proxyType, proxyURL, options, siteURL := a.imageProxyConfig()
	if proxyType == "" {
		return nil
	}

	return func(url string) string {
		if url == "" || url[0] == '/' || strings.HasPrefix(url, siteURL) || strings.HasPrefix(url, proxyURL) {
			return url
		}

		switch proxyType {
		case "atmos/camo":
			mac := hmac.New(sha1.New, []byte(options))
			mac.Write([]byte(url))
			digest := hex.EncodeToString(mac.Sum(nil))
			return proxyURL + digest + "/" + hex.EncodeToString([]byte(url))
		}

		return url
	}
}

func (a *App) ImageProxyRemover() (f func(string) string) {
	proxyType, proxyURL, _, _ := a.imageProxyConfig()
	if proxyType == "" {
		return nil
	}

	return func(url string) string {
		switch proxyType {
		case "atmos/camo":
			if strings.HasPrefix(url, proxyURL) {
				if slash := strings.IndexByte(url[len(proxyURL):], '/'); slash >= 0 {
					if decoded, err := hex.DecodeString(url[len(proxyURL)+slash+1:]); err == nil {
						return string(decoded)
					}
				}
			}
		}

		return url
	}
}

func (a *App) MaxPostSize() int {
	result := <-a.Srv.Store.Post().GetMaxPostSize()
	if result.Err != nil {
		mlog.Error(fmt.Sprint(result.Err))
		return model.POST_MESSAGE_MAX_RUNES_V1
	}
	return result.Data.(int)
}
