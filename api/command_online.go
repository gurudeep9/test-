// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
)

type OnlineProvider struct {
}

const (
	CMD_ONLINE = "online"
)

func init() {
	RegisterCommandProvider(&OnlineProvider{})
}

func (me *OnlineProvider) GetTrigger() string {
	return CMD_ONLINE
}

func (me *OnlineProvider) GetCommand(c *Context) *model.Command {
	return &model.Command{
		Trigger:          CMD_ONLINE,
		AutoComplete:     true,
		AutoCompleteDesc: c.T("api.command_online.desc"),
		AutoCompleteHint: c.T("api.command_online.hint"),
		DisplayName:      c.T("api.command_online.name"),
	}
}

func (me *OnlineProvider) DoCommand(c *Context, channelId string, message string) *model.CommandResponse {
	err := UpdateStatus(c.Session.UserId, model.USER_ONLINE)
	if err == nil {
		message = c.T("api.command_online.ok")
	} else {
		message = c.T("api.command_online.error")
		l4g.Error(err.ToJson())
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         message}
}
