// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/deployment/datastore (interfaces: DataStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// CountDeployments mocks base method
func (m *MockDataStore) CountDeployments(arg0 context.Context) (int, error) {
	ret := m.ctrl.Call(m, "CountDeployments", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountDeployments indicates an expected call of CountDeployments
func (mr *MockDataStoreMockRecorder) CountDeployments(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountDeployments", reflect.TypeOf((*MockDataStore)(nil).CountDeployments), arg0)
}

// GetAllDeployments mocks base method
func (m *MockDataStore) GetAllDeployments(arg0 context.Context) ([]*storage.Deployment, error) {
	ret := m.ctrl.Call(m, "GetAllDeployments", arg0)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDeployments indicates an expected call of GetAllDeployments
func (mr *MockDataStoreMockRecorder) GetAllDeployments(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDeployments", reflect.TypeOf((*MockDataStore)(nil).GetAllDeployments), arg0)
}

// GetDeployment mocks base method
func (m *MockDataStore) GetDeployment(arg0 context.Context, arg1 string) (*storage.Deployment, bool, error) {
	ret := m.ctrl.Call(m, "GetDeployment", arg0, arg1)
	ret0, _ := ret[0].(*storage.Deployment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeployment indicates an expected call of GetDeployment
func (mr *MockDataStoreMockRecorder) GetDeployment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployment", reflect.TypeOf((*MockDataStore)(nil).GetDeployment), arg0, arg1)
}

// GetDeployments mocks base method
func (m *MockDataStore) GetDeployments(arg0 context.Context, arg1 []string) ([]*storage.Deployment, error) {
	ret := m.ctrl.Call(m, "GetDeployments", arg0, arg1)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeployments indicates an expected call of GetDeployments
func (mr *MockDataStoreMockRecorder) GetDeployments(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployments", reflect.TypeOf((*MockDataStore)(nil).GetDeployments), arg0, arg1)
}

// GetImagesForDeployment mocks base method
func (m *MockDataStore) GetImagesForDeployment(arg0 context.Context, arg1 *storage.Deployment) ([]*storage.Image, error) {
	ret := m.ctrl.Call(m, "GetImagesForDeployment", arg0, arg1)
	ret0, _ := ret[0].([]*storage.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImagesForDeployment indicates an expected call of GetImagesForDeployment
func (mr *MockDataStoreMockRecorder) GetImagesForDeployment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImagesForDeployment", reflect.TypeOf((*MockDataStore)(nil).GetImagesForDeployment), arg0, arg1)
}

// ListDeployment mocks base method
func (m *MockDataStore) ListDeployment(arg0 context.Context, arg1 string) (*storage.ListDeployment, bool, error) {
	ret := m.ctrl.Call(m, "ListDeployment", arg0, arg1)
	ret0, _ := ret[0].(*storage.ListDeployment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListDeployment indicates an expected call of ListDeployment
func (mr *MockDataStoreMockRecorder) ListDeployment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeployment", reflect.TypeOf((*MockDataStore)(nil).ListDeployment), arg0, arg1)
}

// ListDeployments mocks base method
func (m *MockDataStore) ListDeployments(arg0 context.Context) ([]*storage.ListDeployment, error) {
	ret := m.ctrl.Call(m, "ListDeployments", arg0)
	ret0, _ := ret[0].([]*storage.ListDeployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeployments indicates an expected call of ListDeployments
func (mr *MockDataStoreMockRecorder) ListDeployments(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeployments", reflect.TypeOf((*MockDataStore)(nil).ListDeployments), arg0)
}

// RemoveDeployment mocks base method
func (m *MockDataStore) RemoveDeployment(arg0 context.Context, arg1, arg2 string) error {
	ret := m.ctrl.Call(m, "RemoveDeployment", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDeployment indicates an expected call of RemoveDeployment
func (mr *MockDataStoreMockRecorder) RemoveDeployment(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDeployment", reflect.TypeOf((*MockDataStore)(nil).RemoveDeployment), arg0, arg1, arg2)
}

// Search mocks base method
func (m *MockDataStore) Search(arg0 context.Context, arg1 *v1.Query) ([]search.Result, error) {
	ret := m.ctrl.Call(m, "Search", arg0, arg1)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockDataStoreMockRecorder) Search(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), arg0, arg1)
}

// SearchDeployments mocks base method
func (m *MockDataStore) SearchDeployments(arg0 context.Context, arg1 *v1.Query) ([]*v1.SearchResult, error) {
	ret := m.ctrl.Call(m, "SearchDeployments", arg0, arg1)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchDeployments indicates an expected call of SearchDeployments
func (mr *MockDataStoreMockRecorder) SearchDeployments(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchDeployments), arg0, arg1)
}

// SearchListDeployments mocks base method
func (m *MockDataStore) SearchListDeployments(arg0 context.Context, arg1 *v1.Query) ([]*storage.ListDeployment, error) {
	ret := m.ctrl.Call(m, "SearchListDeployments", arg0, arg1)
	ret0, _ := ret[0].([]*storage.ListDeployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListDeployments indicates an expected call of SearchListDeployments
func (mr *MockDataStoreMockRecorder) SearchListDeployments(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchListDeployments), arg0, arg1)
}

// SearchRawDeployments mocks base method
func (m *MockDataStore) SearchRawDeployments(arg0 context.Context, arg1 *v1.Query) ([]*storage.Deployment, error) {
	ret := m.ctrl.Call(m, "SearchRawDeployments", arg0, arg1)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawDeployments indicates an expected call of SearchRawDeployments
func (mr *MockDataStoreMockRecorder) SearchRawDeployments(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchRawDeployments), arg0, arg1)
}

// UpdateDeployment mocks base method
func (m *MockDataStore) UpdateDeployment(arg0 context.Context, arg1 *storage.Deployment) error {
	ret := m.ctrl.Call(m, "UpdateDeployment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeployment indicates an expected call of UpdateDeployment
func (mr *MockDataStoreMockRecorder) UpdateDeployment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeployment", reflect.TypeOf((*MockDataStore)(nil).UpdateDeployment), arg0, arg1)
}

// UpsertDeployment mocks base method
func (m *MockDataStore) UpsertDeployment(arg0 context.Context, arg1 *storage.Deployment) error {
	ret := m.ctrl.Call(m, "UpsertDeployment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertDeployment indicates an expected call of UpsertDeployment
func (mr *MockDataStoreMockRecorder) UpsertDeployment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertDeployment", reflect.TypeOf((*MockDataStore)(nil).UpsertDeployment), arg0, arg1)
}
