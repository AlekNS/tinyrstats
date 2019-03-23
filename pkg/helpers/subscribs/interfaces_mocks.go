// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package subscribs is a generated GoMock package.
package subscribs

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEventHandler is a mock of EventHandler interface
type MockEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockEventHandlerMockRecorder
}

// MockEventHandlerMockRecorder is the mock recorder for MockEventHandler
type MockEventHandlerMockRecorder struct {
	mock *MockEventHandler
}

// NewMockEventHandler creates a new mock instance
func NewMockEventHandler(ctrl *gomock.Controller) *MockEventHandler {
	mock := &MockEventHandler{ctrl: ctrl}
	mock.recorder = &MockEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventHandler) EXPECT() *MockEventHandlerMockRecorder {
	return m.recorder
}

// Emit mocks base method
func (m *MockEventHandler) Emit(args ...interface{}) error {
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Emit", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Emit indicates an expected call of Emit
func (mr *MockEventHandlerMockRecorder) Emit(args ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Emit", reflect.TypeOf((*MockEventHandler)(nil).Emit), args...)
}

// On mocks base method
func (m *MockEventHandler) On(arg0 *HandlerOnFunc) (HandlerOffFunc, error) {
	ret := m.ctrl.Call(m, "On", arg0)
	ret0, _ := ret[0].(HandlerOffFunc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// On indicates an expected call of On
func (mr *MockEventHandlerMockRecorder) On(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "On", reflect.TypeOf((*MockEventHandler)(nil).On), arg0)
}

// Off mocks base method
func (m *MockEventHandler) Off(arg0 *HandlerOnFunc) error {
	ret := m.ctrl.Call(m, "Off", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Off indicates an expected call of Off
func (mr *MockEventHandlerMockRecorder) Off(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Off", reflect.TypeOf((*MockEventHandler)(nil).Off), arg0)
}

// OffAll mocks base method
func (m *MockEventHandler) OffAll() error {
	ret := m.ctrl.Call(m, "OffAll")
	ret0, _ := ret[0].(error)
	return ret0
}

// OffAll indicates an expected call of OffAll
func (mr *MockEventHandlerMockRecorder) OffAll() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OffAll", reflect.TypeOf((*MockEventHandler)(nil).OffAll))
}