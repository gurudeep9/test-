// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make sharedchannel-mocks`.

package sharedchannel

import (
	mlog "github.com/mattermost/mattermost/server/public/shared/mlog"
	mock "github.com/stretchr/testify/mock"

	model "github.com/mattermost/mattermost/server/public/model"

	remotecluster "github.com/mattermost/mattermost/server/v8/platform/services/remotecluster"

	store "github.com/mattermost/mattermost/server/v8/channels/store"
)

// MockServerIface is an autogenerated mock type for the ServerIface type
type MockServerIface struct {
	mock.Mock
}

// AddClusterLeaderChangedListener provides a mock function with given fields: listener
func (_m *MockServerIface) AddClusterLeaderChangedListener(listener func()) string {
	ret := _m.Called(listener)

	if len(ret) == 0 {
		panic("no return value specified for AddClusterLeaderChangedListener")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(func()) string); ok {
		r0 = rf(listener)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Config provides a mock function with given fields:
func (_m *MockServerIface) Config() *model.Config {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Config")
	}

	var r0 *model.Config
	if rf, ok := ret.Get(0).(func() *model.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Config)
		}
	}

	return r0
}

// GetRemoteClusterService provides a mock function with given fields:
func (_m *MockServerIface) GetRemoteClusterService() remotecluster.RemoteClusterServiceIFace {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRemoteClusterService")
	}

	var r0 remotecluster.RemoteClusterServiceIFace
	if rf, ok := ret.Get(0).(func() remotecluster.RemoteClusterServiceIFace); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(remotecluster.RemoteClusterServiceIFace)
		}
	}

	return r0
}

// GetStore provides a mock function with given fields:
func (_m *MockServerIface) GetStore() store.Store {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStore")
	}

	var r0 store.Store
	if rf, ok := ret.Get(0).(func() store.Store); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.Store)
		}
	}

	return r0
}

// IsLeader provides a mock function with given fields:
func (_m *MockServerIface) IsLeader() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsLeader")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Log provides a mock function with given fields:
func (_m *MockServerIface) Log() *mlog.Logger {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Log")
	}

	var r0 *mlog.Logger
	if rf, ok := ret.Get(0).(func() *mlog.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mlog.Logger)
		}
	}

	return r0
}

// RemoveClusterLeaderChangedListener provides a mock function with given fields: id
func (_m *MockServerIface) RemoveClusterLeaderChangedListener(id string) {
	_m.Called(id)
}

// NewMockServerIface creates a new instance of MockServerIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServerIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServerIface {
	mock := &MockServerIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
