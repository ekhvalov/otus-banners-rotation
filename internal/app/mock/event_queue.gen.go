// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ekhvalov/otus-banners-rotation/internal/app (interfaces: EventQueue)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	app "github.com/ekhvalov/otus-banners-rotation/internal/app"
	gomock "github.com/golang/mock/gomock"
)

// MockEventQueue is a mock of EventQueue interface.
type MockEventQueue struct {
	ctrl     *gomock.Controller
	recorder *MockEventQueueMockRecorder
}

// MockEventQueueMockRecorder is the mock recorder for MockEventQueue.
type MockEventQueueMockRecorder struct {
	mock *MockEventQueue
}

// NewMockEventQueue creates a new mock instance.
func NewMockEventQueue(ctrl *gomock.Controller) *MockEventQueue {
	mock := &MockEventQueue{ctrl: ctrl}
	mock.recorder = &MockEventQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventQueue) EXPECT() *MockEventQueueMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockEventQueue) Put(arg0 context.Context, arg1 app.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockEventQueueMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockEventQueue)(nil).Put), arg0, arg1)
}
