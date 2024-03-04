// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockServiceClienter is a mock of ServiceClienter interface.
type MockServiceClienter struct {
	ctrl     *gomock.Controller
	recorder *MockServiceClienterMockRecorder
}

// MockServiceClienterMockRecorder is the mock recorder for MockServiceClienter.
type MockServiceClienterMockRecorder struct {
	mock *MockServiceClienter
}

// NewMockServiceClienter creates a new mock instance.
func NewMockServiceClienter(ctrl *gomock.Controller) *MockServiceClienter {
	mock := &MockServiceClienter{ctrl: ctrl}
	mock.recorder = &MockServiceClienterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceClienter) EXPECT() *MockServiceClienterMockRecorder {
	return m.recorder
}

// GetPublicKey mocks base method.
func (m *MockServiceClienter) GetPublicKey(ctx context.Context) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicKey", ctx)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicKey indicates an expected call of GetPublicKey.
func (mr *MockServiceClienterMockRecorder) GetPublicKey(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicKey", reflect.TypeOf((*MockServiceClienter)(nil).GetPublicKey), ctx)
}

// Register mocks base method.
func (m *MockServiceClienter) Register(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", ctx)
}

// Register indicates an expected call of Register.
func (mr *MockServiceClienterMockRecorder) Register(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockServiceClienter)(nil).Register), ctx)
}

// ResendEmailVerification mocks base method.
func (m *MockServiceClienter) ResendEmailVerification(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResendEmailVerification", ctx)
}

// ResendEmailVerification indicates an expected call of ResendEmailVerification.
func (mr *MockServiceClienterMockRecorder) ResendEmailVerification(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResendEmailVerification", reflect.TypeOf((*MockServiceClienter)(nil).ResendEmailVerification), ctx)
}

// VerifyEmail mocks base method.
func (m *MockServiceClienter) VerifyEmail(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VerifyEmail", ctx)
}

// VerifyEmail indicates an expected call of VerifyEmail.
func (mr *MockServiceClienterMockRecorder) VerifyEmail(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyEmail", reflect.TypeOf((*MockServiceClienter)(nil).VerifyEmail), ctx)
}
