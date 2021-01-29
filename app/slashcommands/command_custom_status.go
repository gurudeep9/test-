// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"strings"

	goi18n "github.com/mattermost/go-i18n/i18n"

	"github.com/mattermost/mattermost-server/v5/app"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
)

type CustomStatusProvider struct {
}

const (
	CMD_CUSTOM_STATUS = "status"
	CMD_CUSTOM_STATUS_CLEAR = "clear"

	DEFAULT_CUSTOM_STATUS_EMOJI = "speech_balloon"
)

func init() {
	app.RegisterCommandProvider(&CustomStatusProvider{})
}

func (*CustomStatusProvider) GetTrigger() string {
	return CMD_CUSTOM_STATUS
}

func (*CustomStatusProvider) GetCommand(a *app.App, T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_CUSTOM_STATUS,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_custom_status.desc"),
		AutoCompleteHint: T("api.command_custom_status.hint"),
		DisplayName:      T("api.command_custom_status.name"),
	}
}

func (*CustomStatusProvider) DoCommand(a *app.App, args *model.CommandArgs, message string) *model.CommandResponse {
	if message == CMD_CUSTOM_STATUS_CLEAR {
		if err := a.RemoveCustomStatus(args.UserId); err != nil {
			mlog.Error(err.Error())
			return &model.CommandResponse{Text: args.T("api.command_custom_status.clear.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         args.T("api.command_custom_status.clear.success"),
		}
	}

	customStatus := &model.CustomStatus{
		Emoji: DEFAULT_CUSTOM_STATUS_EMOJI,
		Text:  message,
	}
	firstEmojiLocations := model.EMOJI_PATTERN.FindIndex([]byte(message))
	if len(firstEmojiLocations) > 0 && firstEmojiLocations[0] == 0 {
		// emoji found at starting index
		customStatus.Emoji = message[firstEmojiLocations[0]+1:firstEmojiLocations[1]-1]
		customStatus.Text = strings.TrimSpace(message[firstEmojiLocations[1]:])
	}

	if err := a.SetCustomStatus(args.UserId, customStatus); err != nil {
		mlog.Error(err.Error())
		return &model.CommandResponse{Text: args.T("api.command_custom_status.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         args.T("api.command_custom_status.success", map[string]interface{}{
			"EmojiName":     ":"+customStatus.Emoji+":",
			"StatusMessage": customStatus.Text,
		}),
	}
}
