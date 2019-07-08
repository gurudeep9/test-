// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"
import store "github.com/mattermost/mattermost-server/store"

// TeamStore is an autogenerated mock type for the TeamStore type
type TeamStore struct {
	mock.Mock
}

// AnalyticsGetTeamCountForScheme provides a mock function with given fields: schemeId
func (_m *TeamStore) AnalyticsGetTeamCountForScheme(schemeId string) store.StoreChannel {
	ret := _m.Called(schemeId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(schemeId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// AnalyticsTeamCount provides a mock function with given fields:
func (_m *TeamStore) AnalyticsTeamCount() (int64, *model.AppError) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func() *model.AppError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// ClearAllCustomRoleAssignments provides a mock function with given fields:
func (_m *TeamStore) ClearAllCustomRoleAssignments() store.StoreChannel {
	ret := _m.Called()

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func() store.StoreChannel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// ClearCaches provides a mock function with given fields:
func (_m *TeamStore) ClearCaches() {
	_m.Called()
}

// Get provides a mock function with given fields: id
func (_m *TeamStore) Get(id string) (*model.Team, *model.AppError) {
	ret := _m.Called(id)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetActiveMemberCount provides a mock function with given fields: teamId
func (_m *TeamStore) GetActiveMemberCount(teamId string) (int64, *model.AppError) {
	ret := _m.Called(teamId)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(teamId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(teamId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *TeamStore) GetAll() ([]*model.Team, *model.AppError) {
	ret := _m.Called()

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func() []*model.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func() *model.AppError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAllForExportAfter provides a mock function with given fields: limit, afterId
func (_m *TeamStore) GetAllForExportAfter(limit int, afterId string) ([]*model.TeamForExport, *model.AppError) {
	ret := _m.Called(limit, afterId)

	var r0 []*model.TeamForExport
	if rf, ok := ret.Get(0).(func(int, string) []*model.TeamForExport); ok {
		r0 = rf(limit, afterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamForExport)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, string) *model.AppError); ok {
		r1 = rf(limit, afterId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAllPage provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllPage(offset int, limit int) ([]*model.Team, *model.AppError) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, int) *model.AppError); ok {
		r1 = rf(offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAllPrivateTeamListing provides a mock function with given fields:
func (_m *TeamStore) GetAllPrivateTeamListing() store.StoreChannel {
	ret := _m.Called()

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func() store.StoreChannel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// GetAllPrivateTeamPageListing provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllPrivateTeamPageListing(offset int, limit int) ([]*model.Team, *model.AppError) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, int) *model.AppError); ok {
		r1 = rf(offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAllTeamListing provides a mock function with given fields:
func (_m *TeamStore) GetAllTeamListing() store.StoreChannel {
	ret := _m.Called()

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func() store.StoreChannel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// GetAllTeamPageListing provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllTeamPageListing(offset int, limit int) ([]*model.Team, *model.AppError) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, int) *model.AppError); ok {
		r1 = rf(offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetByInviteId provides a mock function with given fields: inviteId
func (_m *TeamStore) GetByInviteId(inviteId string) (*model.Team, *model.AppError) {
	ret := _m.Called(inviteId)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(inviteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(inviteId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *TeamStore) GetByName(name string) (*model.Team, *model.AppError) {
	ret := _m.Called(name)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelUnreadsForAllTeams provides a mock function with given fields: excludeTeamId, userId
func (_m *TeamStore) GetChannelUnreadsForAllTeams(excludeTeamId string, userId string) ([]*model.ChannelUnread, *model.AppError) {
	ret := _m.Called(excludeTeamId, userId)

	var r0 []*model.ChannelUnread
	if rf, ok := ret.Get(0).(func(string, string) []*model.ChannelUnread); ok {
		r0 = rf(excludeTeamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ChannelUnread)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(excludeTeamId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelUnreadsForTeam provides a mock function with given fields: teamId, userId
func (_m *TeamStore) GetChannelUnreadsForTeam(teamId string, userId string) ([]*model.ChannelUnread, *model.AppError) {
	ret := _m.Called(teamId, userId)

	var r0 []*model.ChannelUnread
	if rf, ok := ret.Get(0).(func(string, string) []*model.ChannelUnread); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ChannelUnread)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(teamId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetMember provides a mock function with given fields: teamId, userId
func (_m *TeamStore) GetMember(teamId string, userId string) (*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userId)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(string, string) *model.TeamMember); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(teamId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetMembers provides a mock function with given fields: teamId, offset, limit, restrictions
func (_m *TeamStore) GetMembers(teamId string, offset int, limit int, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, offset, limit, restrictions)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, int, int, *model.ViewUsersRestrictions) []*model.TeamMember); ok {
		r0 = rf(teamId, offset, limit, restrictions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int, *model.ViewUsersRestrictions) *model.AppError); ok {
		r1 = rf(teamId, offset, limit, restrictions)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetMembersByIds provides a mock function with given fields: teamId, userIds, restrictions
func (_m *TeamStore) GetMembersByIds(teamId string, userIds []string, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userIds, restrictions)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, []string, *model.ViewUsersRestrictions) []*model.TeamMember); ok {
		r0 = rf(teamId, userIds, restrictions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, []string, *model.ViewUsersRestrictions) *model.AppError); ok {
		r1 = rf(teamId, userIds, restrictions)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamMembersForExport provides a mock function with given fields: userId
func (_m *TeamStore) GetTeamMembersForExport(userId string) ([]*model.TeamMemberForExport, *model.AppError) {
	ret := _m.Called(userId)

	var r0 []*model.TeamMemberForExport
	if rf, ok := ret.Get(0).(func(string) []*model.TeamMemberForExport); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMemberForExport)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamsByScheme provides a mock function with given fields: schemeId, offset, limit
func (_m *TeamStore) GetTeamsByScheme(schemeId string, offset int, limit int) store.StoreChannel {
	ret := _m.Called(schemeId, offset, limit)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, int, int) store.StoreChannel); ok {
		r0 = rf(schemeId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// GetTeamsByUserId provides a mock function with given fields: userId
func (_m *TeamStore) GetTeamsByUserId(userId string) store.StoreChannel {
	ret := _m.Called(userId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// GetTeamsForUser provides a mock function with given fields: userId
func (_m *TeamStore) GetTeamsForUser(userId string) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(userId)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string) []*model.TeamMember); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamsForUserWithPagination provides a mock function with given fields: userId, page, perPage
func (_m *TeamStore) GetTeamsForUserWithPagination(userId string, page int, perPage int) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(userId, page, perPage)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.TeamMember); ok {
		r0 = rf(userId, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(userId, page, perPage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTotalMemberCount provides a mock function with given fields: teamId
func (_m *TeamStore) GetTotalMemberCount(teamId string) (int64, *model.AppError) {
	ret := _m.Called(teamId)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(teamId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(teamId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUserTeamIds provides a mock function with given fields: userId, allowFromCache
func (_m *TeamStore) GetUserTeamIds(userId string, allowFromCache bool) store.StoreChannel {
	ret := _m.Called(userId, allowFromCache)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, bool) store.StoreChannel); ok {
		r0 = rf(userId, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// InvalidateAllTeamIdsForUser provides a mock function with given fields: userId
func (_m *TeamStore) InvalidateAllTeamIdsForUser(userId string) {
	_m.Called(userId)
}

// MigrateTeamMembers provides a mock function with given fields: fromTeamId, fromUserId
func (_m *TeamStore) MigrateTeamMembers(fromTeamId string, fromUserId string) store.StoreChannel {
	ret := _m.Called(fromTeamId, fromUserId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, string) store.StoreChannel); ok {
		r0 = rf(fromTeamId, fromUserId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// PermanentDelete provides a mock function with given fields: teamId
func (_m *TeamStore) PermanentDelete(teamId string) *model.AppError {
	ret := _m.Called(teamId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// RemoveAllMembersByTeam provides a mock function with given fields: teamId
func (_m *TeamStore) RemoveAllMembersByTeam(teamId string) store.StoreChannel {
	ret := _m.Called(teamId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// RemoveAllMembersByUser provides a mock function with given fields: userId
func (_m *TeamStore) RemoveAllMembersByUser(userId string) *model.AppError {
	ret := _m.Called(userId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// RemoveMember provides a mock function with given fields: teamId, userId
func (_m *TeamStore) RemoveMember(teamId string, userId string) store.StoreChannel {
	ret := _m.Called(teamId, userId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, string) store.StoreChannel); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// ResetAllTeamSchemes provides a mock function with given fields:
func (_m *TeamStore) ResetAllTeamSchemes() store.StoreChannel {
	ret := _m.Called()

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func() store.StoreChannel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Save provides a mock function with given fields: team
func (_m *TeamStore) Save(team *model.Team) (*model.Team, *model.AppError) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Team) *model.AppError); ok {
		r1 = rf(team)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SaveMember provides a mock function with given fields: member, maxUsersPerTeam
func (_m *TeamStore) SaveMember(member *model.TeamMember, maxUsersPerTeam int) store.StoreChannel {
	ret := _m.Called(member, maxUsersPerTeam)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.TeamMember, int) store.StoreChannel); ok {
		r0 = rf(member, maxUsersPerTeam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// SearchAll provides a mock function with given fields: term
func (_m *TeamStore) SearchAll(term string) ([]*model.Team, *model.AppError) {
	ret := _m.Called(term)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(term)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(term)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SearchByName provides a mock function with given fields: name
func (_m *TeamStore) SearchByName(name string) ([]*model.Team, *model.AppError) {
	ret := _m.Called(name)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SearchOpen provides a mock function with given fields: term
func (_m *TeamStore) SearchOpen(term string) ([]*model.Team, *model.AppError) {
	ret := _m.Called(term)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(term)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(term)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SearchPrivate provides a mock function with given fields: term
func (_m *TeamStore) SearchPrivate(term string) ([]*model.Team, *model.AppError) {
	ret := _m.Called(term)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(term)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(term)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Update provides a mock function with given fields: team
func (_m *TeamStore) Update(team *model.Team) (*model.Team, *model.AppError) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Team) *model.AppError); ok {
		r1 = rf(team)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateDisplayName provides a mock function with given fields: name, teamId
func (_m *TeamStore) UpdateDisplayName(name string, teamId string) *model.AppError {
	ret := _m.Called(name, teamId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, string) *model.AppError); ok {
		r0 = rf(name, teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// UpdateLastTeamIconUpdate provides a mock function with given fields: teamId, curTime
func (_m *TeamStore) UpdateLastTeamIconUpdate(teamId string, curTime int64) *model.AppError {
	ret := _m.Called(teamId, curTime)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, int64) *model.AppError); ok {
		r0 = rf(teamId, curTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// UpdateMember provides a mock function with given fields: member
func (_m *TeamStore) UpdateMember(member *model.TeamMember) (*model.TeamMember, *model.AppError) {
	ret := _m.Called(member)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(*model.TeamMember) *model.TeamMember); ok {
		r0 = rf(member)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.TeamMember) *model.AppError); ok {
		r1 = rf(member)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UserBelongsToTeams provides a mock function with given fields: userId, teamIds
func (_m *TeamStore) UserBelongsToTeams(userId string, teamIds []string) store.StoreChannel {
	ret := _m.Called(userId, teamIds)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, []string) store.StoreChannel); ok {
		r0 = rf(userId, teamIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}
