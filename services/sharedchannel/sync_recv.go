// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannel

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/services/remotecluster"
)

func (scs *Service) onReceiveSyncMessage(msg model.RemoteClusterMsg, rc *model.RemoteCluster, response remotecluster.Response) error {
	if msg.Topic != TopicSync {
		return fmt.Errorf("wrong topic, expected `%s`, got `%s`", TopicSync, msg.Topic)
	}

	if len(msg.Payload) == 0 {
		return errors.New("empty sync message")
	}

	var syncMessages []syncMsg

	if err := json.Unmarshal(msg.Payload, &syncMessages); err != nil {
		return fmt.Errorf("invalid sync message: %w", err)
	}

	scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceDebug, "Sync message received",
		mlog.String("remote", rc.DisplayName),
		mlog.Int("count", len(syncMessages)),
	)

	return scs.processSyncMessagesViaAppAddUsers(syncMessages, rc, response)
}

func (scs *Service) processSyncMessagesViaAppAddUsers(syncMessages []syncMsg, rc *model.RemoteCluster, response remotecluster.Response) error {
	var channel *model.Channel
	postErrors := make([]string, 0)
	usersSyncd := make([]string, 0)
	var lastSyncAt int64
	var err error

	for _, sm := range syncMessages {

		// TODO: modify perma-links (MM-31596)

		scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceDebug, "Sync post received",
			mlog.String("post_id", sm.Post.Id),
			mlog.String("channel_id", sm.Post.ChannelId),
			mlog.Int("reaction_count", len(sm.Reactions)),
			mlog.Int("user_count", len(sm.Users)))

		// add/update users first
		for _, user := range sm.Users {
			if userSaved, err := scs.upsertSyncUser(user, rc); err != nil {
				scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Error upserting sync user",
					mlog.String("post_id", sm.Post.Id),
					mlog.String("channel_id", sm.Post.ChannelId),
					mlog.String("user_id", user.Id))
			} else {
				usersSyncd = append(usersSyncd, userSaved.Id)
			}
		}

		if channel == nil {
			if channel, err = scs.server.GetStore().Channel().Get(sm.Post.ChannelId, true); err != nil {
				// if the channel doesn't exist then none of these sync messages are going to work.
				return fmt.Errorf("channel not found processing sync messages: %w", err)
			}
		}

		if sm.Post == nil {
			continue
		}

		rpost, err := scs.upsertSyncPost(sm.Post, channel, rc)
		if err != nil {
			postErrors = append(postErrors, sm.Post.Id)
			scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Error upserting sync post",
				mlog.String("post_id", sm.Post.Id),
				mlog.String("channel_id", sm.Post.ChannelId),
				mlog.Err(err))
			continue
		}

		for _, reaction := range sm.Reactions {
			if _, err := scs.upsertSyncReaction(reaction, rc); err != nil {
				scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Error creating/deleting sync reaction",
					mlog.String("remote", rc.DisplayName),
					mlog.String("ChannelId", sm.Post.ChannelId),
					mlog.String("PostId", sm.Post.Id),
					mlog.Err(err),
				)
			}
		}

		if lastSyncAt < rpost.UpdateAt {
			lastSyncAt = rpost.UpdateAt
		}
	}

	if lastSyncAt > 0 {
		response[ResponseLastUpdateAt] = lastSyncAt
	}

	response[ResponsePostErrors] = postErrors

	return nil
}

func (scs *Service) upsertSyncUser(user *model.User, rc *model.RemoteCluster) (*model.User, error) {
	var err error
	var userSaved *model.User

	user.RemoteId = rc.RemoteId

	// does the user already exist?
	euser, err := scs.server.GetStore().User().Get(user.Id)
	if err != nil {
		if _, ok := err.(errNotFound); !ok {
			return nil, fmt.Errorf("error checking sync user: %w", err)
		}
	}

	if euser == nil {
		if userSaved, err = scs.server.GetStore().User().Save(user); err != nil {
			if _, ok := err.(errConflict); !ok {
				return nil, fmt.Errorf("error inserting sync user: %w", err)
			}
			// probably a username or email collision
			// TODO: handle collision by modifying username/email (MM-32133)
			return nil, fmt.Errorf("username or email collision inserting sync user: %w", err)
		}
	} else {
		patch := &model.UserPatch{
			Nickname:  &user.Nickname,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
			Position:  &user.Position,
			Locale:    &user.Locale,
			Timezone:  user.Timezone,
		}
		euser.Patch(patch)
		userUpdated, err := scs.server.GetStore().User().Update(euser, false)
		if err != nil {
			return nil, fmt.Errorf("error updating sync user: %w", err)
		}
		userSaved = userUpdated.New
	}
	return userSaved, nil
}

func (scs *Service) upsertSyncPost(post *model.Post, channel *model.Channel, rc *model.RemoteCluster) (*model.Post, error) {
	var appErr *model.AppError

	post.RemoteId = rc.RemoteId

	rpost, err := scs.server.GetStore().Post().GetSingle(post.Id)
	if err != nil {
		if _, ok := err.(errNotFound); !ok {
			return nil, fmt.Errorf("error checking sync post: %w", err)
		}
	}

	if rpost == nil {
		// post doesn't exist; create new one
		rpost, appErr = scs.app.CreatePost(post, channel, true, true)
		scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Creating sync post",
			mlog.String("post_id", post.Id),
			mlog.String("channel_id", post.ChannelId))
	} else {
		// update post
		rpost, appErr = scs.app.UpdatePost(post, false)
		scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Updating sync post",
			mlog.String("post_id", post.Id),
			mlog.String("channel_id", post.ChannelId))
	}
	return rpost, appErr
}

func (scs *Service) upsertSyncReaction(reaction *model.Reaction, rc *model.RemoteCluster) (*model.Reaction, error) {
	savedReaction := reaction
	var err error

	if reaction.DeleteAt == 0 {
		savedReaction, err = scs.app.SaveReactionForPost(reaction)
		if err != nil {
			return nil, err
		}
	} else {
		err = scs.app.DeleteReactionForPost(reaction)
		if err != nil {
			return nil, err
		}
	}
	return savedReaction, nil
}
