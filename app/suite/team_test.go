// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package suite

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/app/email"
	emailmocks "github.com/mattermost/mattermost-server/v6/app/email/mocks"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/store"
	"github.com/mattermost/mattermost-server/v6/store/sqlstore"
	"github.com/mattermost/mattermost-server/v6/store/storetest/mocks"
)

func TestCreateTeam(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	id := model.NewId()
	team := &model.Team{
		DisplayName: "dn_" + id,
		Name:        "name" + id,
		Email:       "success+" + id + "@simulator.amazonses.com",
		Type:        model.TeamOpen,
	}

	_, err := th.Suite.CreateTeam(th.Context, team)
	require.Nil(t, err, "Should create a new team")

	_, err = th.Suite.CreateTeam(th.Context, th.BasicTeam)
	require.NotNil(t, err, "Should not create a new team - team already exist")
}

func TestCreateTeamWithUser(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	id := model.NewId()
	team := &model.Team{
		DisplayName: "dn_" + id,
		Name:        "name" + id,
		Email:       "success+" + id + "@simulator.amazonses.com",
		Type:        model.TeamOpen,
	}

	_, err := th.Suite.CreateTeamWithUser(th.Context, team, th.BasicUser.Id)
	require.Nil(t, err, "Should create a new team with existing user")

	_, err = th.Suite.CreateTeamWithUser(th.Context, team, model.NewId())
	require.NotNil(t, err, "Should not create a new team - user does not exist")
}

func TestUpdateTeam(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	th.BasicTeam.DisplayName = "Testing 123"

	updatedTeam, err := th.Suite.UpdateTeam(th.BasicTeam)
	require.Nil(t, err, "Should update the team")
	require.Equal(t, "Testing 123", updatedTeam.DisplayName, "Wrong Team DisplayName")
}

func TestAddUserToTeam(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	t.Run("add user", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err, "Should add user to the team")
	})

	t.Run("allow user by domain", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "example.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err, "Should have allowed whitelisted user")
	})

	t.Run("block user by domain but allow bot", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "example.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user := model.User{Email: strings.ToLower(model.NewId()) + "test@invalid.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, err := th.Suite.CreateUser(th.Context, &user)
		require.Nil(t, err, "Error creating user: %s", err)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.NotNil(t, err, "Should not add restricted user")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")

		user = model.User{Email: strings.ToLower(model.NewId()) + "test@invalid.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), AuthService: "notnil", AuthData: model.NewString("notnil")}
		ruser, err = th.Suite.CreateUser(th.Context, &user)
		require.Nil(t, err, "Error creating authservice user: %s", err)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.NotNil(t, err, "Should not add authservice user")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")

		bot, err := th.Suite.CreateBot(th.Context, &model.Bot{
			Username:    "somebot",
			Description: "a bot",
			OwnerId:     th.BasicUser.Id,
		})
		require.Nil(t, err)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, bot.UserId, "")
		assert.Nil(t, err, "should be able to add bot to domain restricted team")
	})

	t.Run("block user with subdomain", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "example.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user := model.User{Email: strings.ToLower(model.NewId()) + "test@invalid.example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.NotNil(t, err, "Should not add restricted user")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")
	})

	t.Run("allow users by multiple domains", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "foo.com, bar.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user1 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@foo.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser1, _ := th.Suite.CreateUser(th.Context, &user1)

		user2 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@bar.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser2, _ := th.Suite.CreateUser(th.Context, &user2)

		user3 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@invalid.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser3, _ := th.Suite.CreateUser(th.Context, &user3)

		defer th.Suite.PermanentDeleteUser(th.Context, &user1)
		defer th.Suite.PermanentDeleteUser(th.Context, &user2)
		defer th.Suite.PermanentDeleteUser(th.Context, &user3)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser1.Id, "")
		require.Nil(t, err, "Should have allowed whitelisted user1")

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser2.Id, "")
		require.Nil(t, err, "Should have allowed whitelisted user2")

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser3.Id, "")
		require.NotNil(t, err, "Should not have allowed restricted user3")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")
	})

	t.Run("should set up initial sidebar categories when joining a team", func(t *testing.T) {
		user := th.CreateUser()
		team := th.CreateTeam()

		_, _, err := th.Suite.AddUserToTeam(th.Context, team.Id, user.Id, "")
		require.Nil(t, err)

		res, err := th.Suite.channels.GetSidebarCategoriesForTeamForUser(th.Context, user.Id, team.Id)
		require.Nil(t, err)
		assert.Len(t, res.Categories, 3)
		assert.Equal(t, model.SidebarCategoryFavorites, res.Categories[0].Type)
		assert.Equal(t, model.SidebarCategoryChannels, res.Categories[1].Type)
		assert.Equal(t, model.SidebarCategoryDirectMessages, res.Categories[2].Type)
	})
}

func TestAddUserToTeamByToken(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
	ruser, _ := th.Suite.CreateUser(th.Context, &user)
	rguest := th.CreateGuest()

	t.Run("invalid token", func(t *testing.T) {
		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, "123")
		require.NotNil(t, err, "Should fail on unexisting token")
	})

	t.Run("invalid token type", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeVerifyEmail,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)

		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		defer th.Suite.DeleteToken(token)

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.NotNil(t, err, "Should fail on bad token type")
	})

	t.Run("expired token", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)

		token.CreateAt = model.GetMillis() - InvitationExpiryTime - 1
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		defer th.Suite.DeleteToken(token)

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.NotNil(t, err, "Should fail on expired token")
	})

	t.Run("invalid team id", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": model.NewId()}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		defer th.Suite.DeleteToken(token)

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.NotNil(t, err, "Should fail on bad team id")
	})

	t.Run("invalid user id", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		defer th.Suite.DeleteToken(token)

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, model.NewId(), token.Token)
		require.NotNil(t, err, "Should fail on bad user id")
	})

	t.Run("valid request", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.Nil(t, err, "Should add user to the team")

		_, nErr := th.Suite.platform.Store.Token().GetByToken(token.Token)
		require.Error(t, nErr, "The token must be deleted after be used")

		members, err := th.Suite.channels.GetChannelMembersForUser(th.Context, th.BasicTeam.Id, ruser.Id)
		require.Nil(t, err)
		assert.Len(t, members, 2)
	})

	t.Run("invalid add a guest using a regular invite", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, rguest.Id, token.Token)
		assert.NotNil(t, err)
	})

	t.Run("invalid add a regular user using a guest invite", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeGuestInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id, "channels": th.BasicChannel.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		assert.NotNil(t, err)
	})

	t.Run("invalid add a guest user with a non-granted email domain", func(t *testing.T) {
		restrictedDomain := *th.Suite.platform.Config().GuestAccountsSettings.RestrictCreationToDomains
		defer func() {
			th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.GuestAccountsSettings.RestrictCreationToDomains = &restrictedDomain })
		}()
		th.Suite.platform.UpdateConfig(func(cfg *model.Config) { *cfg.GuestAccountsSettings.RestrictCreationToDomains = "restricted.com" })
		token := model.NewToken(
			TokenTypeGuestInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id, "channels": th.BasicChannel.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, rguest.Id, token.Token)
		require.NotNil(t, err)
		assert.Equal(t, "api.team.join_user_to_team.allowed_domains.app_error", err.Id)
	})

	t.Run("add a guest user with a granted email domain", func(t *testing.T) {
		restrictedDomain := *th.Suite.platform.Config().GuestAccountsSettings.RestrictCreationToDomains
		defer func() {
			th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.GuestAccountsSettings.RestrictCreationToDomains = &restrictedDomain })
		}()
		th.Suite.platform.UpdateConfig(func(cfg *model.Config) { *cfg.GuestAccountsSettings.RestrictCreationToDomains = "restricted.com" })
		token := model.NewToken(
			TokenTypeGuestInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id, "channels": th.BasicChannel.Id}),
		)
		guestEmail := rguest.Email
		rguest.Email = "test@restricted.com"
		_, err := th.Suite.platform.Store.User().Update(rguest, false)
		th.Suite.InvalidateCacheForUser(rguest.Id)
		require.NoError(t, err)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, appErr := th.Suite.AddUserToTeamByToken(th.Context, rguest.Id, token.Token)
		require.Nil(t, appErr)
		rguest.Email = guestEmail
		_, err = th.Suite.platform.Store.User().Update(rguest, false)
		require.NoError(t, err)
	})

	t.Run("add a guest user even though there are team and system domain restrictions", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "restricted-team.com"
		_, err := th.Suite.platform.Store.Team().Update(th.BasicTeam)
		require.NoError(t, err)
		restrictedDomain := *th.Suite.platform.Config().TeamSettings.RestrictCreationToDomains
		defer func() {
			th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.TeamSettings.RestrictCreationToDomains = &restrictedDomain })
		}()
		th.Suite.platform.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.RestrictCreationToDomains = "restricted.com" })
		token := model.NewToken(
			TokenTypeGuestInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id, "channels": th.BasicChannel.Id}),
		)
		_, err = th.Suite.platform.Store.User().Update(rguest, false)
		require.NoError(t, err)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))
		_, _, appErr := th.Suite.AddUserToTeamByToken(th.Context, rguest.Id, token.Token)
		require.Nil(t, appErr)
		th.BasicTeam.AllowedDomains = ""
		_, err = th.Suite.platform.Store.Team().Update(th.BasicTeam)
		require.NoError(t, err)
	})

	t.Run("valid request from guest invite", func(t *testing.T) {
		token := model.NewToken(
			TokenTypeGuestInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id, "channels": th.BasicChannel.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, rguest.Id, token.Token)
		require.Nil(t, err, "Should add user to the team")

		_, nErr := th.Suite.platform.Store.Token().GetByToken(token.Token)
		require.Error(t, nErr, "The token must be deleted after be used")

		members, err := th.Suite.channels.GetChannelMembersForUser(th.Context, th.BasicTeam.Id, rguest.Id)
		require.Nil(t, err)
		require.Len(t, members, 1)
		assert.Equal(t, members[0].ChannelId, th.BasicChannel.Id)
	})

	t.Run("group-constrained team", func(t *testing.T) {
		th.BasicTeam.GroupConstrained = model.NewBool(true)
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))

		_, _, err = th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.NotNil(t, err, "Should return an error when trying to join a group-constrained team.")
		require.Equal(t, "app.team.invite_token.group_constrained.error", err.Id)

		th.BasicTeam.GroupConstrained = model.NewBool(false)
		_, err = th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")
	})

	t.Run("block user", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "example.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user := model.User{Email: strings.ToLower(model.NewId()) + "test@invalid.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": th.BasicTeam.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))

		_, _, err = th.Suite.AddUserToTeamByToken(th.Context, ruser.Id, token.Token)
		require.NotNil(t, err, "Should not add restricted user")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")
	})

	t.Run("should set up initial sidebar categories when joining a team by token", func(t *testing.T) {
		user := th.CreateUser()
		team := th.CreateTeam()

		token := model.NewToken(
			TokenTypeTeamInvitation,
			model.MapToJSON(map[string]string{"teamId": team.Id}),
		)
		require.NoError(t, th.Suite.platform.Store.Token().Save(token))

		_, _, err := th.Suite.AddUserToTeamByToken(th.Context, user.Id, token.Token)
		require.Nil(t, err)

		res, err := th.Suite.channels.GetSidebarCategoriesForTeamForUser(th.Context, user.Id, team.Id)
		require.Nil(t, err)
		assert.Len(t, res.Categories, 3)
		assert.Equal(t, model.SidebarCategoryFavorites, res.Categories[0].Type)
		assert.Equal(t, model.SidebarCategoryChannels, res.Categories[1].Type)
		assert.Equal(t, model.SidebarCategoryDirectMessages, res.Categories[2].Type)
	})
}

func TestAddUserToTeamByTeamId(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	t.Run("add user", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)

		err := th.Suite.AddUserToTeamByTeamId(th.Context, th.BasicTeam.Id, ruser)
		require.Nil(t, err, "Should add user to the team")
	})

	t.Run("block user", func(t *testing.T) {
		th.BasicTeam.AllowedDomains = "example.com"
		_, err := th.Suite.UpdateTeam(th.BasicTeam)
		require.Nil(t, err, "Should update the team")

		user := model.User{Email: strings.ToLower(model.NewId()) + "test@invalid.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		err = th.Suite.AddUserToTeamByTeamId(th.Context, th.BasicTeam.Id, ruser)
		require.NotNil(t, err, "Should not add restricted user")
		require.Equal(t, "JoinUserToTeam", err.Where, "Error should be JoinUserToTeam")
	})

}

func TestSoftDeleteAllTeamsExcept(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	teams := []*model.Team{
		{
			DisplayName: "team-1",
			Name:        "team-1",
			Email:       "foo@foo.com",
			Type:        model.TeamOpen,
		},
	}
	teamId := ""
	for _, create := range teams {
		team, err := th.Suite.CreateTeam(th.Context, create)
		require.Nil(t, err)
		teamId = team.Id
	}

	err := th.Suite.SoftDeleteAllTeamsExcept(teamId)
	assert.Nil(t, err)
	allTeams, err := th.Suite.GetAllTeams()
	require.Nil(t, err)
	for _, team := range allTeams {
		if team.Id == teamId {
			require.Equal(t, int64(0), team.DeleteAt)
			require.Equal(t, false, team.CloudLimitsArchived)
		} else {
			require.NotEqual(t, int64(0), team.DeleteAt)
			require.Equal(t, true, team.CloudLimitsArchived)
		}
	}

}

func TestAdjustTeamsFromProductLimits(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	teams := []*model.Team{
		{
			DisplayName: "team-1",
			Name:        "team-1",
			Email:       "foo@foo.com",
			Type:        model.TeamOpen,
		},
		{
			DisplayName: "team-2",
			Name:        "team-2",
			Email:       "foo@foo.com",
			Type:        model.TeamOpen,
		},
		{
			DisplayName: "team-3",
			Name:        "team-3",
			Email:       "foo@foo.com",
			Type:        model.TeamOpen,
		},
	}
	teamIds := []string{}
	for _, create := range teams {
		team, err := th.Suite.CreateTeam(th.Context, create)
		require.Nil(t, err)
		teamIds = append(teamIds, team.Id)
	}
	t.Run("Should soft delete teams if there are more teams than the limit", func(t *testing.T) {
		activeLimit := 1
		teamLimits := &model.TeamsLimits{Active: &activeLimit}

		err := th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)

		teamsList, err := th.Suite.GetTeams(teamIds)

		require.Nil(t, err)

		// Sort the list of teams based on their creation date
		sort.Slice(teamsList, func(i, j int) bool {
			return teamsList[i].CreateAt < teamsList[j].CreateAt
		})

		for i := range teamsList {
			require.Equal(t, teamsList[i].DisplayName, teams[i].DisplayName)
			require.NotEqual(t, 0, teamsList[i].DeleteAt)
			require.Equal(t, true, teamsList[i].CloudLimitsArchived)
		}
	})

	t.Run("Should not do anything if the amount of teams is equal to the limit", func(t *testing.T) {

		expectedTeamsList, err := th.Suite.GetAllTeams()

		var expectedActiveTeams []*model.Team
		var expectedCloudArchivedTeams []*model.Team
		for _, team := range expectedTeamsList {
			if team.DeleteAt == 0 {
				expectedActiveTeams = append(expectedActiveTeams, team)
			}
			if team.DeleteAt > 0 && team.CloudLimitsArchived {
				expectedCloudArchivedTeams = append(expectedCloudArchivedTeams, team)
			}
		}

		require.Nil(t, err)

		activeLimit := len(expectedActiveTeams)
		teamLimits := &model.TeamsLimits{Active: &activeLimit}
		err = th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)

		actualTeamsList, err := th.Suite.GetAllTeams()

		require.Nil(t, err)
		var actualActiveTeams []*model.Team
		var actualCloudArchivedTeams []*model.Team
		for _, team := range actualTeamsList {
			if team.DeleteAt == 0 {
				actualActiveTeams = append(actualActiveTeams, team)
			}
			if team.DeleteAt > 0 && team.CloudLimitsArchived {
				actualCloudArchivedTeams = append(actualCloudArchivedTeams, team)
			}
		}

		require.Equal(t, len(expectedActiveTeams), len(actualActiveTeams))
		require.Equal(t, len(expectedCloudArchivedTeams), len(actualCloudArchivedTeams))
	})

	t.Run("Should restore archived teams if limit increases", func(t *testing.T) {
		activeLimit := 1
		teamLimits := &model.TeamsLimits{Active: &activeLimit}

		err := th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)
		activeLimit = 10000 // make the limit extremely high so all teams are enabled
		teamLimits = &model.TeamsLimits{Active: &activeLimit}

		err = th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)

		teamsList, err := th.Suite.GetTeams(teamIds)

		require.Nil(t, err)

		// Sort the list of teams based on their creation date
		sort.Slice(teamsList, func(i, j int) bool {
			return teamsList[i].CreateAt < teamsList[j].CreateAt
		})

		for i := range teamsList {
			require.Equal(t, teamsList[i].DisplayName, teams[i].DisplayName)
			require.Equal(t, int64(0), teamsList[i].DeleteAt)
			require.Equal(t, false, teamsList[i].CloudLimitsArchived)
		}
	})

	t.Run("Should only restore teams that were archived by cloud limits", func(t *testing.T) {

		activeLimit := 1
		teamLimits := &model.TeamsLimits{Active: &activeLimit}

		err := th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)

		cloudLimitsArchived := false
		patch := &model.TeamPatch{CloudLimitsArchived: &cloudLimitsArchived}
		team, err := th.Suite.PatchTeam(teamIds[0], patch)
		require.Nil(t, err)
		require.Equal(t, false, team.CloudLimitsArchived)

		activeLimit = 10000 // make the limit extremely high so all teams are enabled
		teamLimits = &model.TeamsLimits{Active: &activeLimit}

		err = th.Suite.AdjustTeamsFromProductLimits(teamLimits)
		require.Nil(t, err)

		teamsList, err := th.Suite.GetTeams(teamIds)

		require.Nil(t, err)

		// Sort the list of teams based on their creation date
		sort.Slice(teamsList, func(i, j int) bool {
			return teamsList[i].CreateAt < teamsList[j].CreateAt
		})

		require.NotEqual(t, int64(0), teamsList[0].DeleteAt)
		require.Equal(t, int64(0), teamsList[1].DeleteAt)
		require.Equal(t, int64(0), teamsList[2].DeleteAt)
	})

}

func TestPermanentDeleteTeam(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	team, err := th.Suite.CreateTeam(th.Context, &model.Team{
		DisplayName: "deletion-test",
		Name:        "deletion-test",
		Email:       "foo@foo.com",
		Type:        model.TeamOpen,
	})
	require.Nil(t, err, "Should create a team")

	defer func() {
		th.Suite.PermanentDeleteTeam(th.Context, team)
	}()

	// TODO: suite: move to commands.
	// command, err := th.Suite.CreateCommand(&model.Command{
	// 	CreatorId: th.BasicUser.Id,
	// 	TeamId:    team.Id,
	// 	Trigger:   "foo",
	// 	URL:       "http://foo",
	// 	Method:    model.CommandMethodPost,
	// })
	// require.Nil(t, err, "Should create a command")
	// defer th.Suite.DeleteCommand(command.Id)

	// command, err = th.Suite.GetCommand(command.Id)
	// require.NotNil(t, command, "command should not be nil")
	// require.Nil(t, err, "unable to get new command")

	// err = th.Suite.PermanentDeleteTeam(th.Context, team)
	// require.Nil(t, err)

	// command, err = th.Suite.GetCommand(command.Id)
	// require.Nil(t, command, "command wasn't deleted")
	// require.NotNil(t, err, "should not return an error")

	// Test deleting a team with no channels.
	team = th.CreateTeam()
	defer func() {
		th.Suite.PermanentDeleteTeam(th.Context, team)
	}()

	channels, err := th.Suite.channels.GetPublicChannelsForTeam(th.Context, team.Id, 0, 1000)
	require.Nil(t, err)

	for _, channel := range channels {
		err2 := th.Suite.channels.PermanentDeleteChannel(th.Context, channel)
		require.Nil(t, err2)
	}

	err = th.Suite.PermanentDeleteTeam(th.Context, team)
	require.Nil(t, err)
}

func TestSanitizeTeam(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	team := &model.Team{
		Id:             model.NewId(),
		Email:          th.MakeEmail(),
		InviteId:       model.NewId(),
		AllowedDomains: "example.com",
	}

	copyTeam := func() *model.Team {
		copy := &model.Team{}
		*copy = *team
		return copy
	}

	t.Run("not a user of the team", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: model.NewId(),
					Roles:  model.TeamUserRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.Empty(t, sanitized.Email, "should've sanitized team")
		require.Empty(t, sanitized.InviteId, "should've sanitized inviteid")
	})

	t.Run("user of the team", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: team.Id,
					Roles:  model.TeamUserRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.Empty(t, sanitized.Email, "should've sanitized team")
		require.NotEmpty(t, sanitized.InviteId, "should have not sanitized inviteid")
	})

	t.Run("team admin", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: team.Id,
					Roles:  model.TeamUserRoleId + " " + model.TeamAdminRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.NotEmpty(t, sanitized.Email, "shouldn't have sanitized team")
		require.NotEmpty(t, sanitized.InviteId, "shouldn't have sanitized inviteid")
	})

	t.Run("team admin of another team", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: model.NewId(),
					Roles:  model.TeamUserRoleId + " " + model.TeamAdminRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.Empty(t, sanitized.Email, "should've sanitized team")
		require.Empty(t, sanitized.InviteId, "should've sanitized inviteid")
	})

	t.Run("system admin, not a user of team", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId + " " + model.SystemAdminRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: model.NewId(),
					Roles:  model.TeamUserRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.NotEmpty(t, sanitized.Email, "shouldn't have sanitized team")
		require.NotEmpty(t, sanitized.InviteId, "shouldn't have sanitized inviteid")
	})

	t.Run("system admin, user of team", func(t *testing.T) {
		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId + " " + model.SystemAdminRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: team.Id,
					Roles:  model.TeamUserRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeam(session, copyTeam())
		require.NotEmpty(t, sanitized.Email, "shouldn't have sanitized team")
		require.NotEmpty(t, sanitized.InviteId, "shouldn't have sanitized inviteid")
	})
}

func TestSanitizeTeams(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	t.Run("not a system admin", func(t *testing.T) {
		teams := []*model.Team{
			{
				Id:             model.NewId(),
				Email:          th.MakeEmail(),
				AllowedDomains: "example.com",
			},
			{
				Id:             model.NewId(),
				Email:          th.MakeEmail(),
				AllowedDomains: "example.com",
			},
		}

		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: teams[0].Id,
					Roles:  model.TeamUserRoleId,
				},
				{
					UserId: userID,
					TeamId: teams[1].Id,
					Roles:  model.TeamUserRoleId + " " + model.TeamAdminRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeams(session, teams)

		require.Empty(t, sanitized[0].Email, "should've sanitized first team")
		require.NotEmpty(t, sanitized[1].Email, "shouldn't have sanitized second team")
	})

	t.Run("system admin", func(t *testing.T) {
		teams := []*model.Team{
			{
				Id:             model.NewId(),
				Email:          th.MakeEmail(),
				AllowedDomains: "example.com",
			},
			{
				Id:             model.NewId(),
				Email:          th.MakeEmail(),
				AllowedDomains: "example.com",
			},
		}

		userID := model.NewId()
		session := model.Session{
			Roles: model.SystemUserRoleId + " " + model.SystemAdminRoleId,
			TeamMembers: []*model.TeamMember{
				{
					UserId: userID,
					TeamId: teams[0].Id,
					Roles:  model.TeamUserRoleId,
				},
			},
		}

		sanitized := th.Suite.SanitizeTeams(session, teams)
		assert.NotEmpty(t, sanitized[0].Email, "shouldn't have sanitized first team")
		assert.NotEmpty(t, sanitized[1].Email, "shouldn't have sanitized second team")
	})
}

func TestJoinUserToTeam(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	id := model.NewId()
	team := &model.Team{
		DisplayName: "dn_" + id,
		Name:        "name" + id,
		Email:       "success+" + id + "@simulator.amazonses.com",
		Type:        model.TeamOpen,
	}

	_, err := th.Suite.CreateTeam(th.Context, team)
	require.Nil(t, err, "Should create a new team")

	maxUsersPerTeam := th.Suite.platform.Config().TeamSettings.MaxUsersPerTeam
	defer func() {
		th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.TeamSettings.MaxUsersPerTeam = maxUsersPerTeam })
		th.Suite.PermanentDeleteTeam(th.Context, team)
	}()
	one := 1
	th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.TeamSettings.MaxUsersPerTeam = &one })

	t.Run("new join", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, appErr := th.Suite.JoinUserToTeam(th.Context, team, ruser, "")
		require.Nil(t, appErr, "Should return no error")
	})

	t.Run("new join with limit problem", func(t *testing.T) {
		user1 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser1, _ := th.Suite.CreateUser(th.Context, &user1)
		user2 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser2, _ := th.Suite.CreateUser(th.Context, &user2)

		defer th.Suite.PermanentDeleteUser(th.Context, &user1)
		defer th.Suite.PermanentDeleteUser(th.Context, &user2)

		_, appErr := th.Suite.JoinUserToTeam(th.Context, team, ruser1, ruser2.Id)
		require.Nil(t, appErr, "Should return no error")

		_, appErr = th.Suite.JoinUserToTeam(th.Context, team, ruser2, ruser1.Id)
		require.NotNil(t, appErr, "Should fail")
	})

	t.Run("re-join after leaving with limit problem", func(t *testing.T) {
		user1 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser1, _ := th.Suite.CreateUser(th.Context, &user1)

		user2 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser2, _ := th.Suite.CreateUser(th.Context, &user2)

		defer th.Suite.PermanentDeleteUser(th.Context, &user1)
		defer th.Suite.PermanentDeleteUser(th.Context, &user2)

		_, appErr := th.Suite.JoinUserToTeam(th.Context, team, ruser1, ruser2.Id)
		require.Nil(t, appErr, "Should return no error")
		appErr = th.Suite.LeaveTeam(th.Context, team, ruser1, ruser1.Id)
		require.Nil(t, appErr, "Should return no error")
		_, appErr = th.Suite.JoinUserToTeam(th.Context, team, ruser2, ruser2.Id)
		require.Nil(t, appErr, "Should return no error")

		_, appErr = th.Suite.JoinUserToTeam(th.Context, team, ruser1, ruser2.Id)
		require.NotNil(t, appErr, "Should fail")
	})

	t.Run("new join with correct scheme_admin value from group syncable", func(t *testing.T) {
		user1 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser1, _ := th.Suite.CreateUser(th.Context, &user1)
		defer th.Suite.PermanentDeleteUser(th.Context, &user1)

		group := th.CreateGroup()

		_, err = th.Suite.UpsertGroupMember(group.Id, user1.Id)
		require.Nil(t, err)

		gs, err := th.Suite.UpsertGroupSyncable(&model.GroupSyncable{
			AutoAdd:     true,
			SyncableId:  team.Id,
			Type:        model.GroupSyncableTypeTeam,
			GroupId:     group.Id,
			SchemeAdmin: false,
		})
		require.Nil(t, err)

		th.Suite.platform.UpdateConfig(func(cfg *model.Config) { cfg.TeamSettings.MaxUsersPerTeam = model.NewInt(999) })

		tm1, appErr := th.Suite.JoinUserToTeam(th.Context, team, ruser1, "")
		require.Nil(t, appErr)
		require.False(t, tm1.SchemeAdmin)

		user2 := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser2, _ := th.Suite.CreateUser(th.Context, &user2)
		defer th.Suite.PermanentDeleteUser(th.Context, &user2)

		_, err = th.Suite.UpsertGroupMember(group.Id, user2.Id)
		require.Nil(t, err)

		gs.SchemeAdmin = true
		_, err = th.Suite.UpdateGroupSyncable(gs)
		require.Nil(t, err)

		tm2, appErr := th.Suite.JoinUserToTeam(th.Context, team, ruser2, "")
		require.Nil(t, appErr)
		require.True(t, tm2.SchemeAdmin)
	})
}

func TestLeaveTeamPanic(t *testing.T) {
	th := SetupWithStoreMock(t)
	defer th.TearDown()

	mockStore := th.Suite.platform.Store.(*mocks.Store)
	mockUserStore := mocks.UserStore{}
	mockUserStore.On("Get", context.Background(), "userID").Return(&model.User{Id: "userID"}, nil)
	mockUserStore.On("Count", mock.Anything).Return(int64(10), nil)

	mockChannelStore := mocks.ChannelStore{}
	mockChannelStore.On("Get", "channelID", true).Return(&model.Channel{Id: "channelID"}, nil)
	mockChannelStore.On("GetMember", context.Background(), "channelID", "userID").Return(&model.ChannelMember{
		NotifyProps: model.StringMap{
			model.PushNotifyProp: model.ChannelNotifyDefault,
		}}, nil)
	mockChannelStore.On("GetChannels", "myteam", "userID", mock.Anything).Return(model.ChannelList{}, nil)

	th.Suite.platform.Store = mockStore

	mockPreferenceStore := mocks.PreferenceStore{}
	mockPreferenceStore.On("Get", "userID", model.PreferenceCategoryDisplaySettings, model.PreferenceNameCollapsedThreadsEnabled).Return(&model.Preference{Value: "on"}, nil)

	mockPostStore := mocks.PostStore{}
	mockPostStore.On("GetMaxPostSize").Return(65535, nil)

	mockSystemStore := mocks.SystemStore{}
	mockSystemStore.On("GetByName", "UpgradedFromTE").Return(&model.System{Name: "UpgradedFromTE", Value: "false"}, nil)
	mockSystemStore.On("GetByName", "InstallationDate").Return(&model.System{Name: "InstallationDate", Value: "10"}, nil)
	mockSystemStore.On("GetByName", "FirstServerRunTimestamp").Return(&model.System{Name: "FirstServerRunTimestamp", Value: "10"}, nil)
	mockLicenseStore := mocks.LicenseStore{}
	mockLicenseStore.On("Get", "").Return(&model.LicenseRecord{}, nil)

	mockTeamStore := mocks.TeamStore{}
	mockTeamStore.On("GetMember", sqlstore.WithMaster(context.Background()), "myteam", "userID").Return(&model.TeamMember{TeamId: "myteam", UserId: "userID"}, nil)
	mockTeamStore.On("UpdateMember", mock.Anything).Return(nil, errors.New("repro error")) // This is the line that triggers the error

	mockStore.On("Channel").Return(&mockChannelStore)
	mockStore.On("Preference").Return(&mockPreferenceStore)
	mockStore.On("Post").Return(&mockPostStore)
	mockStore.On("User").Return(&mockUserStore)
	mockStore.On("System").Return(&mockSystemStore)
	mockStore.On("License").Return(&mockLicenseStore)
	mockStore.On("Team").Return(&mockTeamStore)
	mockStore.On("GetDBSchemaVersion").Return(1, nil)

	team := &model.Team{Id: "myteam"}
	user := &model.User{Id: "userID"}

	th.Suite.platform.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.ExperimentalEnableDefaultChannelLeaveJoinMessages = false
	})

	th.Suite.platform.Store = mockStore

	require.NotPanics(t, func() {
		th.Suite.LeaveTeam(th.Context, team, user, user.Id)
	}, "unexpected panic from LeaveTeam")
}

func TestAppUpdateTeamScheme(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	team := th.BasicTeam
	mockID := model.NewString("x")
	team.SchemeId = mockID

	updatedTeam, err := th.Suite.UpdateTeamScheme(th.BasicTeam)
	require.Nil(t, err)
	require.Equal(t, mockID, updatedTeam.SchemeId, "Wrong Team SchemeId")
}

func TestGetTeamMembers(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	var users []model.User
	users = append(users, *th.BasicUser)
	users = append(users, *th.BasicUser2)

	for i := 0; i < 8; i++ {
		user := model.User{
			Email:    strings.ToLower(model.NewId()) + "success+test@example.com",
			Username: fmt.Sprintf("user%v", i),
			Password: "passwd1",
			DeleteAt: int64(rand.Intn(2)),
		}
		ruser, err := th.Suite.CreateUser(th.Context, &user)
		require.Nil(t, err)
		require.NotNil(t, ruser)
		defer th.Suite.PermanentDeleteUser(th.Context, &user)

		_, _, err = th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		// Store the users for comparison later
		users = append(users, *ruser)
	}

	t.Run("Ensure Sorted By Username when TeamMemberGet options is passed", func(t *testing.T) {
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 100, &model.TeamMembersGetOptions{Sort: model.USERNAME})
		require.Nil(t, err)

		// Sort the users array by username
		sort.Slice(users, func(i, j int) bool {
			return users[i].Username < users[j].Username
		})

		// We should have the same number of users in both users and members array as we have not excluded any deleted members
		require.Equal(t, len(users), len(members))
		for i, member := range members {
			assert.Equal(t, users[i].Id, member.UserId)
		}
	})

	t.Run("Ensure ExcludedDeletedUsers when TeamMemberGetOptions is passed", func(t *testing.T) {
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 100, &model.TeamMembersGetOptions{ExcludeDeletedUsers: true})
		require.Nil(t, err)

		// Choose all users who aren't deleted from our users array
		var usersNotDeletedIDs []string
		var membersIDs []string
		for _, u := range users {
			if u.DeleteAt == 0 {
				usersNotDeletedIDs = append(usersNotDeletedIDs, u.Id)
			}
		}

		for _, m := range members {
			membersIDs = append(membersIDs, m.UserId)
		}

		require.Equal(t, len(usersNotDeletedIDs), len(membersIDs))
		require.ElementsMatch(t, usersNotDeletedIDs, membersIDs)
	})

	t.Run("Ensure Sorted By Username and ExcludedDeletedUsers when TeamMemberGetOptions is passed", func(t *testing.T) {
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 100, &model.TeamMembersGetOptions{Sort: model.USERNAME, ExcludeDeletedUsers: true})
		require.Nil(t, err)

		var usersNotDeleted []model.User
		for _, u := range users {
			if u.DeleteAt == 0 {
				usersNotDeleted = append(usersNotDeleted, u)
			}
		}

		// Sort our non deleted members by username
		sort.Slice(usersNotDeleted, func(i, j int) bool {
			return usersNotDeleted[i].Username < usersNotDeleted[j].Username
		})

		require.Equal(t, len(usersNotDeleted), len(members))
		for i, member := range members {
			assert.Equal(t, usersNotDeleted[i].Id, member.UserId)
		}
	})

	t.Run("Ensure Sorted By User ID when no TeamMemberGetOptions is passed", func(t *testing.T) {

		// Sort them by UserID because the result of GetTeamMembers() is also sorted
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})

		// Fetch team members multiple times
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 5, nil)
		require.Nil(t, err)

		// This should return 5 members
		members2, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 5, 6, nil)
		require.Nil(t, err)
		members = append(members, members2...)

		require.Equal(t, len(users), len(members))
		for i, member := range members {
			assert.Equal(t, users[i].Id, member.UserId)
		}
	})
}

func TestGetTeamStats(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	t.Run("without view restrictions", func(t *testing.T) {
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, nil)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 5, nil)
		require.Nil(t, err)
		assert.Equal(t, int64(len(members)), teamStats.TotalMemberCount)
		assert.Equal(t, int64(len(members)), teamStats.ActiveMemberCount)
	})

	t.Run("with view restrictions by this team", func(t *testing.T) {
		restrictions := &model.ViewUsersRestrictions{Teams: []string{th.BasicTeam.Id}}
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, restrictions)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		members, err := th.Suite.GetTeamMembers(th.BasicTeam.Id, 0, 5, nil)
		require.Nil(t, err)
		assert.Equal(t, int64(len(members)), teamStats.TotalMemberCount)
		assert.Equal(t, int64(len(members)), teamStats.ActiveMemberCount)
	})

	t.Run("with view restrictions by valid channel", func(t *testing.T) {
		restrictions := &model.ViewUsersRestrictions{Teams: []string{}, Channels: []string{th.BasicChannel.Id}}
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, restrictions)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		members, err := th.Suite.channels.GetChannelMembersPage(th.Context, th.BasicChannel.Id, 0, 5)
		require.Nil(t, err)
		assert.Equal(t, int64(len(members)), teamStats.TotalMemberCount)
		assert.Equal(t, int64(len(members)), teamStats.ActiveMemberCount)
	})

	t.Run("with view restrictions to not see anything", func(t *testing.T) {
		restrictions := &model.ViewUsersRestrictions{Teams: []string{}, Channels: []string{}}
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, restrictions)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		assert.Equal(t, int64(0), teamStats.TotalMemberCount)
		assert.Equal(t, int64(0), teamStats.ActiveMemberCount)
	})

	t.Run("with view restrictions by other team", func(t *testing.T) {
		restrictions := &model.ViewUsersRestrictions{Teams: []string{"other-team-id"}}
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, restrictions)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		assert.Equal(t, int64(0), teamStats.TotalMemberCount)
		assert.Equal(t, int64(0), teamStats.ActiveMemberCount)
	})

	t.Run("with view restrictions by not-existing channel", func(t *testing.T) {
		restrictions := &model.ViewUsersRestrictions{Teams: []string{}, Channels: []string{"test"}}
		teamStats, err := th.Suite.GetTeamStats(th.BasicTeam.Id, restrictions)
		require.Nil(t, err)
		require.NotNil(t, teamStats)
		assert.Equal(t, int64(0), teamStats.TotalMemberCount)
		assert.Equal(t, int64(0), teamStats.ActiveMemberCount)
	})
}

func TestUpdateTeamMemberRolesChangingGuest(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	t.Run("from guest to user", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateGuest(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		_, err = th.Suite.UpdateTeamMemberRoles(th.BasicTeam.Id, ruser.Id, "team_user")
		require.NotNil(t, err, "Should fail when try to modify the guest role")
	})

	t.Run("from user to guest", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		_, err = th.Suite.UpdateTeamMemberRoles(th.BasicTeam.Id, ruser.Id, "team_guest")
		require.NotNil(t, err, "Should fail when try to modify the guest role")
	})

	t.Run("from user to admin", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateUser(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		_, err = th.Suite.UpdateTeamMemberRoles(th.BasicTeam.Id, ruser.Id, "team_user team_admin")
		require.Nil(t, err, "Should work when you not modify guest role")
	})

	t.Run("from guest to guest plus custom", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateGuest(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		_, err = th.Suite.CreateRole(&model.Role{Name: "custom", DisplayName: "custom", Description: "custom"})
		require.Nil(t, err)

		_, err = th.Suite.UpdateTeamMemberRoles(th.BasicTeam.Id, ruser.Id, "team_guest custom")
		require.Nil(t, err, "Should work when you not modify guest role")
	})

	t.Run("a guest cant have user role", func(t *testing.T) {
		user := model.User{Email: strings.ToLower(model.NewId()) + "success+test@example.com", Nickname: "Darth Vader", Username: "vader" + model.NewId(), Password: "passwd1", AuthService: ""}
		ruser, _ := th.Suite.CreateGuest(th.Context, &user)

		_, _, err := th.Suite.AddUserToTeam(th.Context, th.BasicTeam.Id, ruser.Id, "")
		require.Nil(t, err)

		_, err = th.Suite.UpdateTeamMemberRoles(th.BasicTeam.Id, ruser.Id, "team_guest team_user")
		require.NotNil(t, err, "Should work when you not modify guest role")
	})
}

func TestInvalidateAllResendInviteEmailJobs(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	job, err := th.Suite.platform.Jobs.CreateJob(model.JobTypeResendInvitationEmail, map[string]string{})
	require.Nil(t, err)

	sysVar := &model.System{Name: job.Id, Value: "0"}
	e := th.Suite.platform.Store.System().SaveOrUpdate(sysVar)
	require.NoError(t, e)

	appErr := th.Suite.InvalidateAllResendInviteEmailJobs()
	require.Nil(t, appErr)

	j, e := th.Suite.platform.Store.Job().Get(job.Id)
	require.NoError(t, e)
	require.Equal(t, j.Status, model.JobStatusCanceled)

	_, sysValErr := th.Suite.platform.Store.System().GetByName(job.Id)
	var errNotFound *store.ErrNotFound
	require.ErrorAs(t, sysValErr, &errNotFound)
}

func TestInvalidateAllEmailInvites(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	t1 := model.Token{
		Token:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		CreateAt: model.GetMillis(),
		Type:     TokenTypeGuestInvitation,
		Extra:    "",
	}
	err := th.Suite.platform.Store.Token().Save(&t1)
	require.NoError(t, err)

	t2 := model.Token{
		Token:    "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
		CreateAt: model.GetMillis(),
		Type:     TokenTypeTeamInvitation,
		Extra:    "",
	}
	err = th.Suite.platform.Store.Token().Save(&t2)
	require.NoError(t, err)

	t3 := model.Token{
		Token:    "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
		CreateAt: model.GetMillis(),
		Type:     "other",
		Extra:    "",
	}
	err = th.Suite.platform.Store.Token().Save(&t3)
	require.NoError(t, err)

	appErr := th.Suite.InvalidateAllEmailInvites()
	require.Nil(t, appErr)

	_, err = th.Suite.platform.Store.Token().GetByToken(t1.Token)
	require.Error(t, err)

	_, err = th.Suite.platform.Store.Token().GetByToken(t2.Token)
	require.Error(t, err)

	_, err = th.Suite.platform.Store.Token().GetByToken(t3.Token)
	require.NoError(t, err)
}

func TestClearTeamMembersCache(t *testing.T) {
	th := SetupWithStoreMock(t)
	defer th.TearDown()

	mockStore := th.Suite.platform.Store.(*mocks.Store)
	mockTeamStore := mocks.TeamStore{}
	tms := []*model.TeamMember{}
	for i := 0; i < 200; i++ {
		tms = append(tms, &model.TeamMember{
			TeamId: "1",
		})
	}
	mockTeamStore.On("GetMembers", "teamID", 0, 100, mock.Anything).Return(tms, nil)
	mockTeamStore.On("GetMembers", "teamID", 100, 100, mock.Anything).Return([]*model.TeamMember{{
		TeamId: "1",
	}}, nil)
	mockStore.On("Team").Return(&mockTeamStore)
	mockStore.On("GetDBSchemaVersion").Return(1, nil)

	require.NoError(t, th.Suite.ClearTeamMembersCache("teamID"))
}

func TestInviteNewUsersToTeamGracefully(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	th.Suite.platform.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.EnableEmailInvitations = true
	})

	t.Run("it return list of email with no error on success", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		memberInvite := &model.MemberInvite{
			Emails: []string{"idontexist@mattermost.com"},
		}
		emailServiceMock.On("SendInviteEmails",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			memberInvite.Emails,
			"",
			mock.Anything,
			true,
		).Once().Return(nil)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteNewUsersToTeamGracefully(memberInvite, th.BasicTeam.Id, th.BasicUser.Id, "")
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Nil(t, res[0].Error)
	})

	t.Run("it should assign errors to emails when failing to send", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		memberInvite := &model.MemberInvite{
			Emails: []string{"idontexist@mattermost.com"},
		}
		emailServiceMock.On("SendInviteEmails",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			memberInvite.Emails,
			"",
			mock.Anything,
			true,
		).Once().Return(email.SendMailError)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteNewUsersToTeamGracefully(memberInvite, th.BasicTeam.Id, th.BasicUser.Id, "")
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.NotNil(t, res[0].Error)
	})

	t.Run("it return list of email with no error when inviting to team and channels using memberInvite struct", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		memberInvite := &model.MemberInvite{
			Emails:     []string{"idontexist@mattermost.com"},
			ChannelIds: []string{th.BasicChannel.Id},
		}
		emailServiceMock.On("SendInviteEmailsToTeamAndChannels",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("[]*model.Channel"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uint8"),
			memberInvite.Emails,
			"",
			mock.Anything,
			mock.AnythingOfType("string"),
			true,
		).Once().Return([]*model.EmailInviteWithError{}, nil)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteNewUsersToTeamGracefully(memberInvite, th.BasicTeam.Id, th.BasicUser.Id, "")
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Nil(t, res[0].Error)
	})

	t.Run("it return list of email with no error when inviting to team and channels using plain emails array", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		memberInvite := &model.MemberInvite{
			Emails: []string{"idontexist@mattermost.com"},
		}
		emailServiceMock.On("SendInviteEmails",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			[]string{"idontexist@mattermost.com"},
			"",
			mock.Anything,
			true,
		).Once().Return(nil)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteNewUsersToTeamGracefully(memberInvite, th.BasicTeam.Id, th.BasicUser.Id, "")
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Nil(t, res[0].Error)
	})
}

func TestInviteGuestsToChannelsGracefully(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	th.Suite.platform.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.EnableEmailInvitations = true
	})

	t.Run("it return list of email with no error on success", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		emailServiceMock.On("SendGuestInviteEmails",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("[]*model.Channel"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uint8"),
			[]string{"idontexist@mattermost.com"},
			"",
			"",
			true,
		).Once().Return(nil)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteGuestsToChannelsGracefully(th.BasicTeam.Id, &model.GuestsInvite{
			Emails:   []string{"idontexist@mattermost.com"},
			Channels: []string{th.BasicChannel.Id},
		}, th.BasicUser.Id)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Nil(t, res[0].Error)
	})

	t.Run("it should assign errors to emails when failing to send", func(t *testing.T) {
		emailServiceMock := emailmocks.ServiceInterface{}
		emailServiceMock.On("SendGuestInviteEmails",
			mock.AnythingOfType("*model.Team"),
			mock.AnythingOfType("[]*model.Channel"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uint8"),
			[]string{"idontexist@mattermost.com"},
			"",
			"",
			true,
		).Once().Return(email.SendMailError)
		th.Suite.email = &emailServiceMock

		res, err := th.Suite.InviteGuestsToChannelsGracefully(th.BasicTeam.Id, &model.GuestsInvite{
			Emails:   []string{"idontexist@mattermost.com"},
			Channels: []string{th.BasicChannel.Id},
		}, th.BasicUser.Id)

		require.Nil(t, err)
		require.Len(t, res, 1)
		require.NotNil(t, res[0].Error)
	})
}

func TestGetNewTeamMembersSince(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	team := th.CreateTeam()

	t.Run("counts team members", func(t *testing.T) {
		var originalExpectedCount int64
		var newTeamMemberJoinTime int64
		var anotherUser *model.User

		t.Run("since time 0", func(t *testing.T) {
			teamMembers, err := th.Suite.platform.Store.Team().GetMembers(team.Id, 0, 1000, nil)
			require.NoError(t, err)
			originalExpectedCount = int64(len(teamMembers))
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, originalExpectedCount, actualCount)
		})

		t.Run("after a new team member was added", func(t *testing.T) {
			anotherUser = th.CreateUser()
			newTeamMember, appErr := th.Suite.JoinUserToTeam(th.Context, team, anotherUser, "")
			newTeamMemberJoinTime = newTeamMember.CreateAt
			require.Nil(t, appErr)
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, originalExpectedCount+1, actualCount)
		})

		t.Run("after a team member was added to a different team, ensuring the wrong team's member count isn't incremented", func(t *testing.T) {
			anotherUser2 := th.CreateUser()
			anotherTeam := th.CreateTeam()
			_, appErr := th.Suite.JoinUserToTeam(th.Context, anotherTeam, anotherUser2, "")
			require.Nil(t, appErr)
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, originalExpectedCount+1, actualCount)
		})

		t.Run("since a given time", func(t *testing.T) {
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: newTeamMemberJoinTime, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, int64(1), actualCount)
		})

		t.Run("after a team member was removed", func(t *testing.T) {
			th.RemoveUserFromTeam(anotherUser, team)
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, originalExpectedCount, actualCount)
		})

		t.Run("after a user was deactivated", func(t *testing.T) {
			_, appErr := th.Suite.JoinUserToTeam(th.Context, team, anotherUser, "")
			require.Nil(t, appErr)
			_, beforeCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			_, appErr = th.Suite.UpdateActive(th.Context, anotherUser, false)
			defer th.Suite.UpdateActive(th.Context, anotherUser, true)
			require.Nil(t, appErr)
			_, afterCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, beforeCount-1, afterCount)
		})

		t.Run("after a user was permanently deleted", func(t *testing.T) {
			_, beforeCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			appErr = th.Suite.PermanentDeleteUser(th.Context, anotherUser)
			require.Nil(t, appErr)
			_, afterCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, beforeCount-1, afterCount)
		})

		t.Run("exclude bots", func(t *testing.T) {
			user := th.CreateUser()
			_, appErr := th.Suite.ConvertUserToBot(user)
			require.Nil(t, appErr)
			_, appErr = th.Suite.JoinUserToTeam(th.Context, team, user, "")
			require.Nil(t, appErr)
			_, actualCount, appErr := th.Suite.GetNewTeamMembersSince(th.Context, team.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Equal(t, originalExpectedCount, actualCount)
		})
	})

	t.Run("returns the correct team members", func(t *testing.T) {
		var originalExpectedMembers []*model.TeamMember
		var newTeamMemberJoinTime int64
		var anotherUser *model.User

		uIDs := func(members []*model.TeamMember) []string {
			ids := []string{}
			for _, member := range members {
				ids = append(ids, member.UserId)
			}
			return ids
		}

		nUIDs := func(members []*model.NewTeamMember) []string {
			ids := []string{}
			for _, member := range members {
				ids = append(ids, member.Id)
			}
			return ids
		}

		t.Run("since time 0", func(t *testing.T) {
			var err error
			originalExpectedMembers, err = th.Suite.platform.Store.Team().GetMembers(th.BasicTeam.Id, 0, 1000, nil)
			require.NoError(t, err)
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.ElementsMatch(t, uIDs(originalExpectedMembers), nUIDs(actualMembersList.Items))
		})

		t.Run("after a new team member was added", func(t *testing.T) {
			anotherUser = th.CreateUser()
			newTeamMember, appErr := th.Suite.JoinUserToTeam(th.Context, th.BasicTeam, anotherUser, "")
			newTeamMemberJoinTime = newTeamMember.CreateAt
			require.Nil(t, appErr)
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.ElementsMatch(t, append(uIDs(originalExpectedMembers), anotherUser.Id), nUIDs(actualMembersList.Items))
		})

		t.Run("after a team member was added to a different team, ensuring the wrong team's member count isn't incremented", func(t *testing.T) {
			anotherUser2 := th.CreateUser()
			anotherTeam := th.CreateTeam()
			_, appErr := th.Suite.JoinUserToTeam(th.Context, anotherTeam, anotherUser2, "")
			require.Nil(t, appErr)
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.ElementsMatch(t, append(uIDs(originalExpectedMembers), anotherUser.Id), nUIDs(actualMembersList.Items))
		})

		t.Run("since a given time", func(t *testing.T) {
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: newTeamMemberJoinTime, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Len(t, actualMembersList.Items, 1)
			require.Equal(t, anotherUser.Id, actualMembersList.Items[0].Id)
		})

		t.Run("after a team member was removed", func(t *testing.T) {
			th.RemoveUserFromTeam(anotherUser, th.BasicTeam)
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.ElementsMatch(t, uIDs(originalExpectedMembers), nUIDs(actualMembersList.Items))
		})

		t.Run("after a user was deactivated", func(t *testing.T) {
			_, appErr := th.Suite.JoinUserToTeam(th.Context, th.BasicTeam, anotherUser, "")
			require.Nil(t, appErr)
			beforeMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Contains(t, nUIDs(beforeMembersList.Items), anotherUser.Id)
			_, appErr = th.Suite.UpdateActive(th.Context, anotherUser, false)
			defer th.Suite.UpdateActive(th.Context, anotherUser, true)
			require.Nil(t, appErr)
			afterMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.NotContains(t, nUIDs(afterMembersList.Items), anotherUser.Id)
		})

		t.Run("after a user was permanently deleted", func(t *testing.T) {
			beforeMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.Contains(t, nUIDs(beforeMembersList.Items), anotherUser.Id)
			appErr = th.Suite.PermanentDeleteUser(th.Context, anotherUser)
			require.Nil(t, appErr)
			afterMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.NotContains(t, nUIDs(afterMembersList.Items), anotherUser.Id)
		})

		t.Run("exclude bots", func(t *testing.T) {
			user := th.CreateUser()
			_, appErr := th.Suite.ConvertUserToBot(user)
			require.Nil(t, appErr)
			_, appErr = th.Suite.JoinUserToTeam(th.Context, th.BasicTeam, user, "")
			require.Nil(t, appErr)
			actualMembersList, _, appErr := th.Suite.GetNewTeamMembersSince(th.Context, th.BasicTeam.Id, &model.InsightsOpts{StartUnixMilli: 0, Page: 0, PerPage: 1000})
			require.Nil(t, appErr)
			require.ElementsMatch(t, uIDs(originalExpectedMembers), nUIDs(actualMembersList.Items))
		})
	})
}
