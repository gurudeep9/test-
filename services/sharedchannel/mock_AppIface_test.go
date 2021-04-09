// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make sharedchannel-mocks`.

package sharedchannel

import (
	filestore "github.com/mattermost/mattermost-server/v5/shared/filestore"
	mock "github.com/stretchr/testify/mock"

	model "github.com/mattermost/mattermost-server/v5/model"
)

// MockAppIface is an autogenerated mock type for the AppIface type
type MockAppIface struct {
	mock.Mock
}

// AddUserToChannel provides a mock function with given fields: user, channel, skipTeamMemberIntegrityCheck
func (_m *MockAppIface) AddUserToChannel(user *model.User, channel *model.Channel, skipTeamMemberIntegrityCheck bool) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(user, channel, skipTeamMemberIntegrityCheck)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(*model.User, *model.Channel, bool) *model.ChannelMember); ok {
		r0 = rf(user, channel, skipTeamMemberIntegrityCheck)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.User, *model.Channel, bool) *model.AppError); ok {
		r1 = rf(user, channel, skipTeamMemberIntegrityCheck)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// AddUserToTeamByTeamId provides a mock function with given fields: teamId, user
func (_m *MockAppIface) AddUserToTeamByTeamId(teamId string, user *model.User) *model.AppError {
	ret := _m.Called(teamId, user)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, *model.User) *model.AppError); ok {
		r0 = rf(teamId, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// CreateChannelWithUser provides a mock function with given fields: channel, userId
func (_m *MockAppIface) CreateChannelWithUser(channel *model.Channel, userId string) (*model.Channel, *model.AppError) {
	ret := _m.Called(channel, userId)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(*model.Channel, string) *model.Channel); ok {
		r0 = rf(channel, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Channel, string) *model.AppError); ok {
		r1 = rf(channel, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreatePost provides a mock function with given fields: post, channel, triggerWebhooks, setOnline
func (_m *MockAppIface) CreatePost(post *model.Post, channel *model.Channel, triggerWebhooks bool, setOnline bool) (*model.Post, *model.AppError) {
	ret := _m.Called(post, channel, triggerWebhooks, setOnline)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*model.Post, *model.Channel, bool, bool) *model.Post); ok {
		r0 = rf(post, channel, triggerWebhooks, setOnline)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Post, *model.Channel, bool, bool) *model.AppError); ok {
		r1 = rf(post, channel, triggerWebhooks, setOnline)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateUploadSession provides a mock function with given fields: us
func (_m *MockAppIface) CreateUploadSession(us *model.UploadSession) (*model.UploadSession, *model.AppError) {
	ret := _m.Called(us)

	var r0 *model.UploadSession
	if rf, ok := ret.Get(0).(func(*model.UploadSession) *model.UploadSession); ok {
		r0 = rf(us)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UploadSession)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.UploadSession) *model.AppError); ok {
		r1 = rf(us)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// DeletePost provides a mock function with given fields: postID, deleteByID
func (_m *MockAppIface) DeletePost(postID string, deleteByID string) (*model.Post, *model.AppError) {
	ret := _m.Called(postID, deleteByID)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string, string) *model.Post); ok {
		r0 = rf(postID, deleteByID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(postID, deleteByID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// DeleteReactionForPost provides a mock function with given fields: reaction
func (_m *MockAppIface) DeleteReactionForPost(reaction *model.Reaction) *model.AppError {
	ret := _m.Called(reaction)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.AppError); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// FileReader provides a mock function with given fields: path
func (_m *MockAppIface) FileReader(path string) (filestore.ReadCloseSeeker, *model.AppError) {
	ret := _m.Called(path)

	var r0 filestore.ReadCloseSeeker
	if rf, ok := ret.Get(0).(func(string) filestore.ReadCloseSeeker); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(filestore.ReadCloseSeeker)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(path)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetOrCreateDirectChannel provides a mock function with given fields: userId, otherUserId, channelOptions
func (_m *MockAppIface) GetOrCreateDirectChannel(userId string, otherUserId string, channelOptions ...model.ChannelOption) (*model.Channel, *model.AppError) {
	_va := make([]interface{}, len(channelOptions))
	for _i := range channelOptions {
		_va[_i] = channelOptions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, userId, otherUserId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string, string, ...model.ChannelOption) *model.Channel); ok {
		r0 = rf(userId, otherUserId, channelOptions...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, ...model.ChannelOption) *model.AppError); ok {
		r1 = rf(userId, otherUserId, channelOptions...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// PatchChannelModerationsForChannel provides a mock function with given fields: channel, channelModerationsPatch
func (_m *MockAppIface) PatchChannelModerationsForChannel(channel *model.Channel, channelModerationsPatch []*model.ChannelModerationPatch) ([]*model.ChannelModeration, *model.AppError) {
	ret := _m.Called(channel, channelModerationsPatch)

	var r0 []*model.ChannelModeration
	if rf, ok := ret.Get(0).(func(*model.Channel, []*model.ChannelModerationPatch) []*model.ChannelModeration); ok {
		r0 = rf(channel, channelModerationsPatch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ChannelModeration)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Channel, []*model.ChannelModerationPatch) *model.AppError); ok {
		r1 = rf(channel, channelModerationsPatch)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// PermanentDeleteChannel provides a mock function with given fields: channel
func (_m *MockAppIface) PermanentDeleteChannel(channel *model.Channel) *model.AppError {
	ret := _m.Called(channel)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Channel) *model.AppError); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// SaveReactionForPost provides a mock function with given fields: reaction
func (_m *MockAppIface) SaveReactionForPost(reaction *model.Reaction) (*model.Reaction, *model.AppError) {
	ret := _m.Called(reaction)

	var r0 *model.Reaction
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.Reaction); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reaction)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Reaction) *model.AppError); ok {
		r1 = rf(reaction)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SendEphemeralPost provides a mock function with given fields: userId, post
func (_m *MockAppIface) SendEphemeralPost(userId string, post *model.Post) *model.Post {
	ret := _m.Called(userId, post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string, *model.Post) *model.Post); ok {
		r0 = rf(userId, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	return r0
}

// UpdatePost provides a mock function with given fields: post, safeUpdate
func (_m *MockAppIface) UpdatePost(post *model.Post, safeUpdate bool) (*model.Post, *model.AppError) {
	ret := _m.Called(post, safeUpdate)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*model.Post, bool) *model.Post); ok {
		r0 = rf(post, safeUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Post, bool) *model.AppError); ok {
		r1 = rf(post, safeUpdate)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}
