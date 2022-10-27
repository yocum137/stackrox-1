// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PodsStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	testDB      *pgtest.TestPostgres
}

func TestPodsStore(t *testing.T) {
	suite.Run(t, new(PodsStoreSuite))
}

func (s *PodsStoreSuite) SetupSuite() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")

	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.Pool)
}

func (s *PodsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE pods CASCADE")
	s.T().Log("pods", tag)
	s.NoError(err)
}

func (s *PodsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
	s.envIsolator.RestoreAll()
}

func (s *PodsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	pod := &storage.Pod{}
	s.NoError(testutils.FullInit(pod, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundPod, exists, err := store.Get(ctx, pod.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPod)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, pod))
	foundPod, exists, err = store.Get(ctx, pod.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(pod, foundPod)

	podCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, podCount)
	podCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(podCount)

	podExists, err := store.Exists(ctx, pod.GetId())
	s.NoError(err)
	s.True(podExists)
	s.NoError(store.Upsert(ctx, pod))
	s.ErrorIs(store.Upsert(withNoAccessCtx, pod), sac.ErrResourceAccessDenied)

	foundPod, exists, err = store.Get(ctx, pod.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(pod, foundPod)

	s.NoError(store.Delete(ctx, pod.GetId()))
	foundPod, exists, err = store.Get(ctx, pod.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPod)
	s.NoError(store.Delete(withNoAccessCtx, pod.GetId()))

	var pods []*storage.Pod
	var podIDs []string
	for i := 0; i < 200; i++ {
		pod := &storage.Pod{}
		s.NoError(testutils.FullInit(pod, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		pods = append(pods, pod)
		podIDs = append(podIDs, pod.GetId())
	}

	s.NoError(store.UpsertMany(ctx, pods))

	podCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, podCount)

	s.NoError(store.DeleteMany(ctx, podIDs))

	podCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, podCount)
}

func (s *PodsStoreSuite) TestSACUpsert() {
	obj := &storage.Pod{}
	s.NoError(testutils.FullInit(obj, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	ctxs := getSACContexts(obj, storage.Access_READ_WRITE_ACCESS)
	for name, expectedErr := range map[string]error{
		withAllAccess:           nil,
		withNoAccess:            sac.ErrResourceAccessDenied,
		withNoAccessToCluster:   sac.ErrResourceAccessDenied,
		withAccessToDifferentNs: sac.ErrResourceAccessDenied,
		withAccess:              nil,
		withAccessToCluster:     nil,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.Upsert(ctxs[name], obj), expectedErr)
		})
	}
}

func (s *PodsStoreSuite) TestSACUpsertMany() {
	obj := &storage.Pod{}
	s.NoError(testutils.FullInit(obj, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	ctxs := getSACContexts(obj, storage.Access_READ_WRITE_ACCESS)
	for name, expectedErr := range map[string]error{
		withAllAccess:           nil,
		withNoAccess:            sac.ErrResourceAccessDenied,
		withNoAccessToCluster:   sac.ErrResourceAccessDenied,
		withAccessToDifferentNs: sac.ErrResourceAccessDenied,
		withAccess:              nil,
		withAccessToCluster:     nil,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.UpsertMany(ctxs[name], []*storage.Pod{obj}), expectedErr)
		})
	}
}

func (s *PodsStoreSuite) TestSACCount() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           2,
		withNoAccess:            0,
		withNoAccessToCluster:   0,
		withAccessToDifferentNs: 0,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			count, err := s.store.Count(ctxs[name])
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *PodsStoreSuite) TestSACWalk() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIds := range map[string][]string{
		withAllAccess:           []string{objA.GetId(), objB.GetId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetId()},
		withAccessToCluster:     []string{objA.GetId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			ids := []string{}
			getIds := func(obj *storage.Pod) error {
				ids = append(ids, obj.GetId())
				return nil
			}
			err := s.store.Walk(ctxs[name], getIds)
			assert.NoError(t, err)
			assert.ElementsMatch(t, expectedIds, ids)
		})
	}
}

func (s *PodsStoreSuite) TestSACGetIDs() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIds := range map[string][]string{
		withAllAccess:           []string{objA.GetId(), objB.GetId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetId()},
		withAccessToCluster:     []string{objA.GetId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			ids, err := s.store.GetIDs(ctxs[name])
			assert.NoError(t, err)
			assert.EqualValues(t, expectedIds, ids)
		})
	}
}

func (s *PodsStoreSuite) TestSACExists() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]bool{
		withAllAccess:           true,
		withNoAccess:            false,
		withNoAccessToCluster:   false,
		withAccessToDifferentNs: false,
		withAccess:              true,
		withAccessToCluster:     true,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			exists, err := s.store.Exists(ctxs[name], objA.GetId())
			assert.NoError(t, err)
			assert.Equal(t, expected, exists)
		})
	}
}

func (s *PodsStoreSuite) TestSACGet() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]bool{
		withAllAccess:           true,
		withNoAccess:            false,
		withNoAccessToCluster:   false,
		withAccessToDifferentNs: false,
		withAccess:              true,
		withAccessToCluster:     true,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, exists, err := s.store.Get(ctxs[name], objA.GetId())
			assert.NoError(t, err)
			assert.Equal(t, expected, exists)
			if expected == true {
				assert.Equal(t, objA, actual)
			} else {
				assert.Nil(t, actual)
			}
		})
	}
}

func (s *PodsStoreSuite) TestSACDelete() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	withAllAccessCtx := sac.WithAllAccess(context.Background())

	ctxs := getSACContexts(objA, storage.Access_READ_WRITE_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           0,
		withNoAccess:            2,
		withNoAccessToCluster:   2,
		withAccessToDifferentNs: 2,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.Delete(ctxs[name], objA.GetId()))
			assert.NoError(t, s.store.Delete(ctxs[name], objB.GetId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *PodsStoreSuite) TestSACDeleteMany() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	withAllAccessCtx := sac.WithAllAccess(context.Background())

	ctxs := getSACContexts(objA, storage.Access_READ_WRITE_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           0,
		withNoAccess:            2,
		withNoAccessToCluster:   2,
		withAccessToDifferentNs: 2,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.DeleteMany(ctxs[name], []string{
				objA.GetId(),
				objB.GetId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *PodsStoreSuite) TestSACGetMany() {
	objA := &storage.Pod{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Pod{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]struct {
		elems          []*storage.Pod
		missingIndices []int
	}{
		withAllAccess:           {elems: []*storage.Pod{objA, objB}, missingIndices: []int{}},
		withNoAccess:            {elems: []*storage.Pod{}, missingIndices: []int{0, 1}},
		withNoAccessToCluster:   {elems: []*storage.Pod{}, missingIndices: []int{0, 1}},
		withAccessToDifferentNs: {elems: []*storage.Pod{}, missingIndices: []int{0, 1}},
		withAccess:              {elems: []*storage.Pod{objA}, missingIndices: []int{1}},
		withAccessToCluster:     {elems: []*storage.Pod{objA}, missingIndices: []int{1}},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, missingIndices, err := s.store.GetMany(ctxs[name], []string{objA.GetId(), objB.GetId()})
			assert.NoError(t, err)
			assert.Equal(t, expected.elems, actual)
			assert.Equal(t, expected.missingIndices, missingIndices)
		})
	}

	s.T().Run("with no ids", func(t *testing.T) {
		actual, missingIndices, err := s.store.GetMany(withAllAccessCtx, []string{})
		assert.Nil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, missingIndices)
	})
}

const (
	withAllAccess           = "AllAccess"
	withNoAccess            = "NoAccess"
	withAccessToDifferentNs = "AccessToDifferentNs"
	withAccess              = "Access"
	withAccessToCluster     = "AccessToCluster"
	withNoAccessToCluster   = "NoAccessToCluster"
)

func getSACContexts(obj *storage.Pod, access storage.Access) map[string]context.Context {
	return map[string]context.Context{
		withAllAccess: sac.WithAllAccess(context.Background()),
		withNoAccess:  sac.WithNoAccess(context.Background()),
		withAccessToDifferentNs: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
				sac.NamespaceScopeKeys("unknown ns"),
			)),
		withAccess: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
				sac.NamespaceScopeKeys(obj.GetNamespace()),
			)),
		withAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
			)),
		withNoAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys("unknown cluster"),
			)),
	}
}
