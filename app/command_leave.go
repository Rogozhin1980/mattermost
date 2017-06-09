// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/platform/model"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type LeaveProvider struct {
}

const (
	CMD_LEAVE = "leave"
)

func init() {
	RegisterCommandProvider(&LeaveProvider{})
}

func (me *LeaveProvider) GetTrigger() string {
	return CMD_LEAVE
}

func (me *LeaveProvider) GetCommand(T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_LEAVE,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_leave.desc"),
		DisplayName:      T("api.command_leave.name"),
	}
}

func (me *LeaveProvider) DoCommand(args *model.CommandArgs, message string) *model.CommandResponse {
	err := LeaveChannel(args.ChannelId, args.UserId)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
	}

	team, err := GetTeam(args.TeamId)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
	}

	return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + model.DEFAULT_CHANNEL, Text: args.T("api.command_leave.success"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
}
