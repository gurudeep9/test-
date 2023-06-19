// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost/server/public/plugin/experimental/bot/poster (interfaces: Poster)

// Package mock_bot is a generated GoMock package.
package mock_bot

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mattermost/mattermost/server/public/model"
)

// MockPoster is a mock of Poster interface.
type MockPoster struct {
	ctrl     *gomock.Controller
	recorder *MockPosterMockRecorder
}

// MockPosterMockRecorder is the mock recorder for MockPoster.
type MockPosterMockRecorder struct {
	mock *MockPoster
}

// NewMockPoster creates a new mock instance.
func NewMockPoster(ctrl *gomock.Controller) *MockPoster {
	mock := &MockPoster{ctrl: ctrl}
	mock.recorder = &MockPosterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoster) EXPECT() *MockPosterMockRecorder {
	return m.recorder
}

// DM mocks base method.
func (m *MockPoster) DM(arg0, arg1 string, arg2 ...interface{}) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DM", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DM indicates an expected call of DM.
func (mr *MockPosterMockRecorder) DM(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DM", reflect.TypeOf((*MockPoster)(nil).DM), varargs...)
}

// DMWithAttachments mocks base method.
func (m *MockPoster) DMWithAttachments(arg0 string, arg1 ...*model.SlackAttachment) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DMWithAttachments", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DMWithAttachments indicates an expected call of DMWithAttachments.
func (mr *MockPosterMockRecorder) DMWithAttachments(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DMWithAttachments", reflect.TypeOf((*MockPoster)(nil).DMWithAttachments), varargs...)
}

// DeletePost mocks base method.
func (m *MockPoster) DeletePost(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockPosterMockRecorder) DeletePost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPoster)(nil).DeletePost), arg0)
}

// Ephemeral mocks base method.
func (m *MockPoster) Ephemeral(arg0, arg1, arg2 string, arg3 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Ephemeral", varargs...)
}

// Ephemeral indicates an expected call of Ephemeral.
func (mr *MockPosterMockRecorder) Ephemeral(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ephemeral", reflect.TypeOf((*MockPoster)(nil).Ephemeral), varargs...)
}

// UpdatePost mocks base method.
func (m *MockPoster) UpdatePost(arg0 *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePost", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePost indicates an expected call of UpdatePost.
func (mr *MockPosterMockRecorder) UpdatePost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePost", reflect.TypeOf((*MockPoster)(nil).UpdatePost), arg0)
}

// UpdatePostByID mocks base method.
func (m *MockPoster) UpdatePostByID(arg0, arg1 string, arg2 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePostByID", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePostByID indicates an expected call of UpdatePostByID.
func (mr *MockPosterMockRecorder) UpdatePostByID(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePostByID", reflect.TypeOf((*MockPoster)(nil).UpdatePostByID), varargs...)
}

// UpdatePosterID mocks base method.
func (m *MockPoster) UpdatePosterID(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdatePosterID", arg0)
}

// UpdatePosterID indicates an expected call of UpdatePosterID.
func (mr *MockPosterMockRecorder) UpdatePosterID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePosterID", reflect.TypeOf((*MockPoster)(nil).UpdatePosterID), arg0)
}
