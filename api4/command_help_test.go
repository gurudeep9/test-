// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/model"
)

func TestHelpCommand(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	Client := th.Client
	channel := th.BasicChannel

	HelpLink := *th.App.Config().SupportSettings.HelpLink
	defer func() {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.SupportSettings.HelpLink = HelpLink })
	}()

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.SupportSettings.HelpLink = "" })
	rs1, _ := Client.ExecuteCommand(channel.Id, "/help ")
	require.Equal(t, rs1.GotoLocation, model.SUPPORT_SETTINGS_DEFAULT_HELP_LINK, "failed to default help link")

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.SupportSettings.HelpLink = "https://docs.mattermost.com/guides/user.html"
	})
	rs2, _ := Client.ExecuteCommand(channel.Id, "/help ")
	require.Equal(t, rs2.GotoLocation, "https://docs.mattermost.com/guides/user.html", "failed to help link")
}
