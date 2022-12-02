// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datastore "github.com/stackrox/rox/central/processlisteningonport/datastore"
	storage "github.com/stackrox/rox/generated/storage"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// AddProcessListeningOnPort mocks base method.
func (m *MockDataStore) AddProcessListeningOnPort(arg0 context.Context, arg1 ...*storage.ProcessListeningOnPort) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddProcessListeningOnPort", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProcessListeningOnPort indicates an expected call of AddProcessListeningOnPort.
func (mr *MockDataStoreMockRecorder) AddProcessListeningOnPort(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessListeningOnPort", reflect.TypeOf((*MockDataStore)(nil).AddProcessListeningOnPort), varargs...)
}

// GetProcessListeningOnPort mocks base method.
func (m *MockDataStore) GetProcessListeningOnPort(ctx context.Context, opts datastore.GetOptions) (map[string][]*storage.ProcessListeningOnPort, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcessListeningOnPort", ctx, opts)
	ret0, _ := ret[0].(map[string][]*storage.ProcessListeningOnPort)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProcessListeningOnPort indicates an expected call of GetProcessListeningOnPort.
func (mr *MockDataStoreMockRecorder) GetProcessListeningOnPort(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcessListeningOnPort", reflect.TypeOf((*MockDataStore)(nil).GetProcessListeningOnPort), ctx, opts)
}
