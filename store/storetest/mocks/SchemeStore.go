// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	model "github.com/mattermost/mattermost-server/v5/model"
)

// SchemeStore is an autogenerated mock type for the SchemeStore type
type SchemeStore struct {
	mock.Mock
}

// CountByScope provides a mock function with given fields: scope
func (_m *SchemeStore) CountByScope(scope string) (int64, error) {
	ret := _m.Called(scope)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(scope)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountWithoutPermission provides a mock function with given fields: scope, permissionID, roleScope, roleType
func (_m *SchemeStore) CountWithoutPermission(scope string, permissionID string, roleScope model.RoleScope, roleType model.RoleType) (int64, error) {
	ret := _m.Called(scope, permissionID, roleScope, roleType)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, string, model.RoleScope, model.RoleType) int64); ok {
		r0 = rf(scope, permissionID, roleScope, roleType)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, model.RoleScope, model.RoleType) error); ok {
		r1 = rf(scope, permissionID, roleScope, roleType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: schemeId
func (_m *SchemeStore) Delete(schemeId string) (*model.Scheme, error) {
	ret := _m.Called(schemeId)

	var r0 *model.Scheme
	if rf, ok := ret.Get(0).(func(string) *model.Scheme); ok {
		r0 = rf(schemeId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(schemeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: schemeId
func (_m *SchemeStore) Get(schemeId string) (*model.Scheme, error) {
	ret := _m.Called(schemeId)

	var r0 *model.Scheme
	if rf, ok := ret.Get(0).(func(string) *model.Scheme); ok {
		r0 = rf(schemeId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(schemeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPage provides a mock function with given fields: scope, offset, limit
func (_m *SchemeStore) GetAllPage(scope string, offset int, limit int) ([]*model.Scheme, error) {
	ret := _m.Called(scope, offset, limit)

	var r0 []*model.Scheme
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.Scheme); ok {
		r0 = rf(scope, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(scope, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: schemeName
func (_m *SchemeStore) GetByName(schemeName string) (*model.Scheme, error) {
	ret := _m.Called(schemeName)

	var r0 *model.Scheme
	if rf, ok := ret.Get(0).(func(string) *model.Scheme); ok {
		r0 = rf(schemeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(schemeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteAll provides a mock function with given fields:
func (_m *SchemeStore) PermanentDeleteAll() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: scheme
func (_m *SchemeStore) Save(scheme *model.Scheme) (*model.Scheme, error) {
	ret := _m.Called(scheme)

	var r0 *model.Scheme
	if rf, ok := ret.Get(0).(func(*model.Scheme) *model.Scheme); ok {
		r0 = rf(scheme)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Scheme) error); ok {
		r1 = rf(scheme)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
