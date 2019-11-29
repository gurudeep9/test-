// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/v5/model"

// DataRetentionInterface is an autogenerated mock type for the DataRetentionInterface type
type DataRetentionInterface struct {
	mock.Mock
}

// GetPolicy provides a mock function with given fields:
func (_m *DataRetentionInterface) GetPolicy() (*model.DataRetentionPolicy, *model.AppError) {
	ret := _m.Called()

	var r0 *model.DataRetentionPolicy
	if rf, ok := ret.Get(0).(func() *model.DataRetentionPolicy); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DataRetentionPolicy)
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
