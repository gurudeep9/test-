// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/shared/request"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/require"
)

func TestGetAppLimits(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)

		// InitBasic creates 3 users by default
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
		require.Equal(t, int64(10000), appLimits.MaxUsersLimit)

		// 5 posts are created by default
		require.Equal(t, int64(5), appLimits.PostCount)
		require.Equal(t, int64(5_000_000), appLimits.MaxPostLimit)
	})

	t.Run("user count should increase on creating new user and decrease on permanently deleting", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we create a new user
		newUser := th.CreateUser()

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(4), appLimits.ActiveUserCount)

		// now we'll delete the user
		_ = th.App.PermanentDeleteUser(th.Context, newUser)
		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
	})

	t.Run("user count should increase on creating new guest user and decrease on permanently deleting", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we create a new user
		newGuestUser := th.CreateGuest()

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(4), appLimits.ActiveUserCount)

		// now we'll delete the user
		_ = th.App.PermanentDeleteUser(th.Context, newGuestUser)
		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
	})

	t.Run("user count should increase on creating new user and decrease on soft deleting", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we create a new user
		newUser := th.CreateUser()

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(4), appLimits.ActiveUserCount)

		// now we'll delete the user
		_, appErr = th.App.UpdateActive(th.Context, newUser, false)
		require.Nil(t, appErr)
		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
	})

	t.Run("user count should increase on creating new guest user and decrease on soft deleting", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we create a new user
		newGuestUser := th.CreateGuest()

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(4), appLimits.ActiveUserCount)

		// now we'll delete the user
		_, appErr = th.App.UpdateActive(th.Context, newGuestUser, false)
		require.Nil(t, appErr)
		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
	})

	t.Run("user count should not change on creating or deleting bots", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we create a new bot
		newBot := th.CreateBot()

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)

		// now we'll delete the bot
		_ = th.App.PermanentDeleteBot(th.Context, newBot.UserId)
		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(3), appLimits.ActiveUserCount)
	})

	t.Run("limits should be empty when there is a license", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		th.App.Srv().SetLicense(model.NewTestLicense())

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)

		require.Equal(t, int64(0), appLimits.ActiveUserCount)
		require.Equal(t, int64(0), appLimits.MaxUsersLimit)
	})

	t.Run("post count should increase on creating new post and NOT decrease on deleting post", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		appLimits, appErr := th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(5), appLimits.PostCount)

		// now we create a new post
		team := th.CreateTeam()
		channel := th.CreateChannel(request.TestContext(t), team)
		post := th.CreatePost(channel)

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(6), appLimits.PostCount)

		// now we'll delete the post
		_, appErr = th.App.DeletePost(request.TestContext(t), post.Id, "")
		require.Nil(t, appErr)

		appLimits, appErr = th.App.GetAppLimits()
		require.Nil(t, appErr)
		require.Equal(t, int64(6), appLimits.PostCount)
	})
}
