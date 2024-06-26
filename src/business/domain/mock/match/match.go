// Code generated by MockGen. DO NOT EDIT.
// Source: matchs/match.go
//
// Generated by this command:
//
//	mockgen -source=matchs/match.go -destination=mock/match/match.go
//
// Package mock_match is a generated GoMock package.
package mock_match

import (
	context "context"
	entity "loverly/src/business/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockInterface) Create(ctx context.Context, param entity.Match) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, param)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockInterfaceMockRecorder) Create(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInterface)(nil).Create), ctx, param)
}

// GetByUserId mocks base method.
func (m *MockInterface) GetByUserId(ctx context.Context, userId int64) ([]entity.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", ctx, userId)
	ret0, _ := ret[0].([]entity.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockInterfaceMockRecorder) GetByUserId(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockInterface)(nil).GetByUserId), ctx, userId)
}
