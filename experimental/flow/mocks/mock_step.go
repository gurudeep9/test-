// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-api/experimental/flow/steps (interfaces: Step)

// Package mock_flow is a generated GoMock package.
package mock_flow

import (
	gomock "github.com/golang/mock/gomock"
	freetext_fetcher "github.com/mattermost/mattermost-plugin-api/experimental/freetext_fetcher"
	model "github.com/mattermost/mattermost-server/v5/model"
	reflect "reflect"
)

// MockStep is a mock of Step interface
type MockStep struct {
	ctrl     *gomock.Controller
	recorder *MockStepMockRecorder
}

// MockStepMockRecorder is the mock recorder for MockStep
type MockStepMockRecorder struct {
	mock *MockStep
}

// NewMockStep creates a new mock instance
func NewMockStep(ctrl *gomock.Controller) *MockStep {
	mock := &MockStep{ctrl: ctrl}
	mock.recorder = &MockStepMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStep) EXPECT() *MockStepMockRecorder {
	return m.recorder
}

// GetFreetextFetcher mocks base method
func (m *MockStep) GetFreetextFetcher() freetext_fetcher.FreetextFetcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFreetextFetcher")
	ret0, _ := ret[0].(freetext_fetcher.FreetextFetcher)
	return ret0
}

// GetFreetextFetcher indicates an expected call of GetFreetextFetcher
func (mr *MockStepMockRecorder) GetFreetextFetcher() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFreetextFetcher", reflect.TypeOf((*MockStep)(nil).GetFreetextFetcher))
}

// GetPropertyName mocks base method
func (m *MockStep) GetPropertyName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertyName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPropertyName indicates an expected call of GetPropertyName
func (mr *MockStepMockRecorder) GetPropertyName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertyName", reflect.TypeOf((*MockStep)(nil).GetPropertyName))
}

// IsEmpty mocks base method
func (m *MockStep) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty
func (mr *MockStepMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockStep)(nil).IsEmpty))
}

// PostSlackAttachment mocks base method
func (m *MockStep) PostSlackAttachment(arg0 string, arg1 int) *model.SlackAttachment {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostSlackAttachment", arg0, arg1)
	ret0, _ := ret[0].(*model.SlackAttachment)
	return ret0
}

// PostSlackAttachment indicates an expected call of PostSlackAttachment
func (mr *MockStepMockRecorder) PostSlackAttachment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSlackAttachment", reflect.TypeOf((*MockStep)(nil).PostSlackAttachment), arg0, arg1)
}

// ResponseSlackAttachment mocks base method
func (m *MockStep) ResponseSlackAttachment(arg0 interface{}) *model.SlackAttachment {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponseSlackAttachment", arg0)
	ret0, _ := ret[0].(*model.SlackAttachment)
	return ret0
}

// ResponseSlackAttachment indicates an expected call of ResponseSlackAttachment
func (mr *MockStepMockRecorder) ResponseSlackAttachment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseSlackAttachment", reflect.TypeOf((*MockStep)(nil).ResponseSlackAttachment), arg0)
}

// ShouldSkip mocks base method
func (m *MockStep) ShouldSkip(arg0 interface{}) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldSkip", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// ShouldSkip indicates an expected call of ShouldSkip
func (mr *MockStepMockRecorder) ShouldSkip(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldSkip", reflect.TypeOf((*MockStep)(nil).ShouldSkip), arg0)
}
