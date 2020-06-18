// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-api/experimental/freetext_fetcher (interfaces: FreetextStore)

// Package mock_freetext_fetcher is a generated GoMock package.
package mock_freetext_fetcher

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFreetextStore is a mock of FreetextStore interface
type MockFreetextStore struct {
	ctrl     *gomock.Controller
	recorder *MockFreetextStoreMockRecorder
}

// MockFreetextStoreMockRecorder is the mock recorder for MockFreetextStore
type MockFreetextStoreMockRecorder struct {
	mock *MockFreetextStore
}

// NewMockFreetextStore creates a new mock instance
func NewMockFreetextStore(ctrl *gomock.Controller) *MockFreetextStore {
	mock := &MockFreetextStore{ctrl: ctrl}
	mock.recorder = &MockFreetextStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFreetextStore) EXPECT() *MockFreetextStoreMockRecorder {
	return m.recorder
}

// ShouldProcessFreetext mocks base method
func (m *MockFreetextStore) ShouldProcessFreetext(arg0, arg1 string) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldProcessFreetext", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ShouldProcessFreetext indicates an expected call of ShouldProcessFreetext
func (mr *MockFreetextStoreMockRecorder) ShouldProcessFreetext(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldProcessFreetext", reflect.TypeOf((*MockFreetextStore)(nil).ShouldProcessFreetext), arg0, arg1)
}

// StartFetching mocks base method
func (m *MockFreetextStore) StartFetching(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartFetching", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartFetching indicates an expected call of StartFetching
func (mr *MockFreetextStoreMockRecorder) StartFetching(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartFetching", reflect.TypeOf((*MockFreetextStore)(nil).StartFetching), arg0, arg1, arg2)
}

// StopFetching mocks base method
func (m *MockFreetextStore) StopFetching(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopFetching", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopFetching indicates an expected call of StopFetching
func (mr *MockFreetextStoreMockRecorder) StopFetching(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopFetching", reflect.TypeOf((*MockFreetextStore)(nil).StopFetching), arg0)
}
