// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost-server/server/v8/public/model"
	mock "github.com/stretchr/testify/mock"
)

// LicenseStore is an autogenerated mock type for the LicenseStore type
type LicenseStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *LicenseStore) Get(id string) (*model.LicenseRecord, error) {
	ret := _m.Called(id)

	var r0 *model.LicenseRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.LicenseRecord, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.LicenseRecord); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LicenseRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *LicenseStore) GetAll() ([]*model.LicenseRecord, error) {
	ret := _m.Called()

	var r0 []*model.LicenseRecord
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*model.LicenseRecord, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*model.LicenseRecord); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.LicenseRecord)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: license
func (_m *LicenseStore) Save(license *model.LicenseRecord) (*model.LicenseRecord, error) {
	ret := _m.Called(license)

	var r0 *model.LicenseRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.LicenseRecord) (*model.LicenseRecord, error)); ok {
		return rf(license)
	}
	if rf, ok := ret.Get(0).(func(*model.LicenseRecord) *model.LicenseRecord); ok {
		r0 = rf(license)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LicenseRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.LicenseRecord) error); ok {
		r1 = rf(license)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewLicenseStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewLicenseStore creates a new instance of LicenseStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLicenseStore(t mockConstructorTestingTNewLicenseStore) *LicenseStore {
	mock := &LicenseStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
