// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/utils"
)

// GetPostsUsage returns the total posts count rounded down to the most
// significant digit
func (a *App) GetPostsUsage() (int64, *model.AppError) {
	count, err := a.Srv().Store.Post().AnalyticsPostCount(&model.PostCountOptions{ExcludeDeleted: true})
	if err != nil {
		return 0, model.NewAppError("GetPostsUsage", "app.post.analytics_posts_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return utils.RoundOffToZeroes(float64(count)), nil
}

// GetStorageUsage returns the sum of files stored
func (a *App) GetStorageUsage() (int64, *model.AppError) {
	usage, err := a.Srv().Store.FileInfo().GetStorageUsage(false)
	if err != nil {
		return 0, model.NewAppError("GetStorageUsage", "app.usage.get_storage_usage.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return usage, nil
}
