// Code generated by MockGen. DO NOT EDIT.
// Source: gh.go

// Package garbanzo is a generated GoMock package.
package garbanzo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockclientI is a mock of clientI interface.
type MockclientI struct {
	ctrl     *gomock.Controller
	recorder *MockclientIMockRecorder
}

// MockclientIMockRecorder is the mock recorder for MockclientI.
type MockclientIMockRecorder struct {
	mock *MockclientI
}

// NewMockclientI creates a new mock instance.
func NewMockclientI(ctrl *gomock.Controller) *MockclientI {
	mock := &MockclientI{ctrl: ctrl}
	mock.recorder = &MockclientIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockclientI) EXPECT() *MockclientIMockRecorder {
	return m.recorder
}

// getNotifications mocks base method.
func (m *MockclientI) getNotifications() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getNotifications")
	ret0, _ := ret[0].(error)
	return ret0
}

// getNotifications indicates an expected call of getNotifications.
func (mr *MockclientIMockRecorder) getNotifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getNotifications", reflect.TypeOf((*MockclientI)(nil).getNotifications))
}
