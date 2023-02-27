// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
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

// Count mocks base method.
func (m *MockDataStore) Count(ctx context.Context, q *v1.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockDataStoreMockRecorder) Count(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDataStore)(nil).Count), ctx, q)
}

// CountNodes mocks base method.
func (m *MockDataStore) CountNodes(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountNodes", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountNodes indicates an expected call of CountNodes.
func (mr *MockDataStoreMockRecorder) CountNodes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountNodes", reflect.TypeOf((*MockDataStore)(nil).CountNodes), ctx)
}

// DeleteAllNodesForCluster mocks base method.
func (m *MockDataStore) DeleteAllNodesForCluster(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllNodesForCluster", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllNodesForCluster indicates an expected call of DeleteAllNodesForCluster.
func (mr *MockDataStoreMockRecorder) DeleteAllNodesForCluster(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllNodesForCluster", reflect.TypeOf((*MockDataStore)(nil).DeleteAllNodesForCluster), ctx, clusterID)
}

// DeleteNodes mocks base method.
func (m *MockDataStore) DeleteNodes(ctx context.Context, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteNodes", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNodes indicates an expected call of DeleteNodes.
func (mr *MockDataStoreMockRecorder) DeleteNodes(ctx interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodes", reflect.TypeOf((*MockDataStore)(nil).DeleteNodes), varargs...)
}

// Exists mocks base method.
func (m *MockDataStore) Exists(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockDataStoreMockRecorder) Exists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDataStore)(nil).Exists), ctx, id)
}

// GetManyNodeMetadata mocks base method.
func (m *MockDataStore) GetManyNodeMetadata(ctx context.Context, ids []string) ([]*storage.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManyNodeMetadata", ctx, ids)
	ret0, _ := ret[0].([]*storage.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManyNodeMetadata indicates an expected call of GetManyNodeMetadata.
func (mr *MockDataStoreMockRecorder) GetManyNodeMetadata(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManyNodeMetadata", reflect.TypeOf((*MockDataStore)(nil).GetManyNodeMetadata), ctx, ids)
}

// GetNode mocks base method.
func (m *MockDataStore) GetNode(ctx context.Context, id string) (*storage.Node, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", ctx, id)
	ret0, _ := ret[0].(*storage.Node)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNode indicates an expected call of GetNode.
func (mr *MockDataStoreMockRecorder) GetNode(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockDataStore)(nil).GetNode), ctx, id)
}

// GetNodesBatch mocks base method.
func (m *MockDataStore) GetNodesBatch(ctx context.Context, ids []string) ([]*storage.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodesBatch", ctx, ids)
	ret0, _ := ret[0].([]*storage.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodesBatch indicates an expected call of GetNodesBatch.
func (mr *MockDataStoreMockRecorder) GetNodesBatch(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodesBatch", reflect.TypeOf((*MockDataStore)(nil).GetNodesBatch), ctx, ids)
}

// Search mocks base method.
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchNodes mocks base method.
func (m *MockDataStore) SearchNodes(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchNodes", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchNodes indicates an expected call of SearchNodes.
func (mr *MockDataStoreMockRecorder) SearchNodes(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchNodes", reflect.TypeOf((*MockDataStore)(nil).SearchNodes), ctx, q)
}

// SearchRawNodes mocks base method.
func (m *MockDataStore) SearchRawNodes(ctx context.Context, q *v1.Query) ([]*storage.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawNodes", ctx, q)
	ret0, _ := ret[0].([]*storage.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawNodes indicates an expected call of SearchRawNodes.
func (mr *MockDataStoreMockRecorder) SearchRawNodes(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawNodes", reflect.TypeOf((*MockDataStore)(nil).SearchRawNodes), ctx, q)
}

// UpsertNode mocks base method.
func (m *MockDataStore) UpsertNode(ctx context.Context, node *storage.Node, ignoreScan bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertNodeNoScan", ctx, node)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertNode indicates an expected call of UpsertNode.
func (mr *MockDataStoreMockRecorder) UpsertNode(ctx, node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertNodeNoScan", reflect.TypeOf((*MockDataStore)(nil).UpsertNode), ctx, node)
}
