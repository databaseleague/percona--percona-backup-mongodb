// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/percona/mongodb-backup/proto/api (interfaces: ApiServer)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/percona/mongodb-backup/proto/api"
	reflect "reflect"
)

// MockApiServer is a mock of ApiServer interface
type MockApiServer struct {
	ctrl     *gomock.Controller
	recorder *MockApiServerMockRecorder
}

// MockApiServerMockRecorder is the mock recorder for MockApiServer
type MockApiServerMockRecorder struct {
	mock *MockApiServer
}

// NewMockApiServer creates a new mock instance
func NewMockApiServer(ctrl *gomock.Controller) *MockApiServer {
	mock := &MockApiServer{ctrl: ctrl}
	mock.recorder = &MockApiServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApiServer) EXPECT() *MockApiServerMockRecorder {
	return m.recorder
}

// GetClients mocks base method
func (m *MockApiServer) GetClients(arg0 *api.Empty, arg1 api.Api_GetClientsServer) error {
	ret := m.ctrl.Call(m, "GetClients", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetClients indicates an expected call of GetClients
func (mr *MockApiServerMockRecorder) GetClients(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClients", reflect.TypeOf((*MockApiServer)(nil).GetClients), arg0, arg1)
}
