// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quadev-ltd/qd-qpi-gateway/internal/middleware (interfaces: AutheticationMiddlewarer)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockAutheticationMiddlewarer is a mock of AutheticationMiddlewarer interface.
type MockAutheticationMiddlewarer struct {
	ctrl     *gomock.Controller
	recorder *MockAutheticationMiddlewarerMockRecorder
}

// MockAutheticationMiddlewarerMockRecorder is the mock recorder for MockAutheticationMiddlewarer.
type MockAutheticationMiddlewarerMockRecorder struct {
	mock *MockAutheticationMiddlewarer
}

// NewMockAutheticationMiddlewarer creates a new mock instance.
func NewMockAutheticationMiddlewarer(ctrl *gomock.Controller) *MockAutheticationMiddlewarer {
	mock := &MockAutheticationMiddlewarer{ctrl: ctrl}
	mock.recorder = &MockAutheticationMiddlewarerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutheticationMiddlewarer) EXPECT() *MockAutheticationMiddlewarerMockRecorder {
	return m.recorder
}

// RefreshAuthentication mocks base method.
func (m *MockAutheticationMiddlewarer) RefreshAuthentication(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RefreshAuthentication", arg0)
}

// RefreshAuthentication indicates an expected call of RefreshAuthentication.
func (mr *MockAutheticationMiddlewarerMockRecorder) RefreshAuthentication(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAuthentication", reflect.TypeOf((*MockAutheticationMiddlewarer)(nil).RefreshAuthentication), arg0)
}

// RequireAuthentication mocks base method.
func (m *MockAutheticationMiddlewarer) RequireAuthentication(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequireAuthentication", arg0)
}

// RequireAuthentication indicates an expected call of RequireAuthentication.
func (mr *MockAutheticationMiddlewarerMockRecorder) RequireAuthentication(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequireAuthentication", reflect.TypeOf((*MockAutheticationMiddlewarer)(nil).RequireAuthentication), arg0)
}
