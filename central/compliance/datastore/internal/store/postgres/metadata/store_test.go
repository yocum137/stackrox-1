// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ComplianceRunMetadataStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestComplianceRunMetadataStore(t *testing.T) {
	suite.Run(t, new(ComplianceRunMetadataStoreSuite))
}

func (s *ComplianceRunMetadataStoreSuite) SetupSuite() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	ctx := sac.WithAllAccess(context.Background())

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.Require().NoError(err)

	Destroy(ctx, pool)

	s.pool = pool
	gormDB := pgtest.OpenGormDB(s.T(), source)
	defer pgtest.CloseGormDB(s.T(), gormDB)
	s.store = CreateTableAndNewStore(ctx, pool, gormDB)
}

func (s *ComplianceRunMetadataStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.pool.Exec(ctx, "TRUNCATE compliance_run_metadata CASCADE")
	s.T().Log("compliance_run_metadata", tag)
	s.NoError(err)
}

func (s *ComplianceRunMetadataStoreSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *ComplianceRunMetadataStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	complianceRunMetadata := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(complianceRunMetadata, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundComplianceRunMetadata, exists, err := store.Get(ctx, complianceRunMetadata.GetRunId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceRunMetadata)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, complianceRunMetadata))
	foundComplianceRunMetadata, exists, err = store.Get(ctx, complianceRunMetadata.GetRunId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceRunMetadata, foundComplianceRunMetadata)

	complianceRunMetadataCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, complianceRunMetadataCount)
	complianceRunMetadataCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(complianceRunMetadataCount)

	complianceRunMetadataExists, err := store.Exists(ctx, complianceRunMetadata.GetRunId())
	s.NoError(err)
	s.True(complianceRunMetadataExists)
	s.NoError(store.Upsert(ctx, complianceRunMetadata))
	s.ErrorIs(store.Upsert(withNoAccessCtx, complianceRunMetadata), sac.ErrResourceAccessDenied)

	foundComplianceRunMetadata, exists, err = store.Get(ctx, complianceRunMetadata.GetRunId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceRunMetadata, foundComplianceRunMetadata)

	s.NoError(store.Delete(ctx, complianceRunMetadata.GetRunId()))
	foundComplianceRunMetadata, exists, err = store.Get(ctx, complianceRunMetadata.GetRunId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceRunMetadata)
	s.NoError(store.Delete(withNoAccessCtx, complianceRunMetadata.GetRunId()))

	var complianceRunMetadatas []*storage.ComplianceRunMetadata
	for i := 0; i < 200; i++ {
		complianceRunMetadata := &storage.ComplianceRunMetadata{}
		s.NoError(testutils.FullInit(complianceRunMetadata, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceRunMetadatas = append(complianceRunMetadatas, complianceRunMetadata)
	}

	s.NoError(store.UpsertMany(ctx, complianceRunMetadatas))

	complianceRunMetadataCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, complianceRunMetadataCount)
}

func (s *ComplianceRunMetadataStoreSuite) TestSACUpsert() {
	obj := &storage.ComplianceRunMetadata{}
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

func (s *ComplianceRunMetadataStoreSuite) TestSACUpsertMany() {
	obj := &storage.ComplianceRunMetadata{}
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
			assert.ErrorIs(t, s.store.UpsertMany(ctxs[name], []*storage.ComplianceRunMetadata{obj}), expectedErr)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACCount() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
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

func (s *ComplianceRunMetadataStoreSuite) TestSACWalk() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIds := range map[string][]string{
		withAllAccess:           []string{objA.GetRunId(), objB.GetRunId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetRunId()},
		withAccessToCluster:     []string{objA.GetRunId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			ids := []string{}
			getIds := func(obj *storage.ComplianceRunMetadata) error {
				ids = append(ids, obj.GetRunId())
				return nil
			}
			err := s.store.Walk(ctxs[name], getIds)
			assert.NoError(t, err)
			assert.ElementsMatch(t, expectedIds, ids)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACGetIDs() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIds := range map[string][]string{
		withAllAccess:           []string{objA.GetRunId(), objB.GetRunId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetRunId()},
		withAccessToCluster:     []string{objA.GetRunId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			ids, err := s.store.GetIDs(ctxs[name])
			assert.NoError(t, err)
			assert.EqualValues(t, expectedIds, ids)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACExists() {
	objA := &storage.ComplianceRunMetadata{}
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
			exists, err := s.store.Exists(ctxs[name], objA.GetRunId())
			assert.NoError(t, err)
			assert.Equal(t, expected, exists)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACGet() {
	objA := &storage.ComplianceRunMetadata{}
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
			actual, exists, err := s.store.Get(ctxs[name], objA.GetRunId())
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

func (s *ComplianceRunMetadataStoreSuite) TestSACDelete() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
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

			assert.NoError(t, s.store.Delete(ctxs[name], objA.GetRunId()))
			assert.NoError(t, s.store.Delete(ctxs[name], objB.GetRunId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACDeleteMany() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
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
				objA.GetRunId(),
				objB.GetRunId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *ComplianceRunMetadataStoreSuite) TestSACGetMany() {
	objA := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunMetadata{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.store.Upsert(withAllAccessCtx, objA)
	s.store.Upsert(withAllAccessCtx, objB)

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]struct {
		elems          []*storage.ComplianceRunMetadata
		missingIndices []int
	}{
		withAllAccess:           {elems: []*storage.ComplianceRunMetadata{objA, objB}, missingIndices: []int{}},
		withNoAccess:            {elems: []*storage.ComplianceRunMetadata{}, missingIndices: []int{0, 1}},
		withNoAccessToCluster:   {elems: []*storage.ComplianceRunMetadata{}, missingIndices: []int{0, 1}},
		withAccessToDifferentNs: {elems: []*storage.ComplianceRunMetadata{}, missingIndices: []int{0, 1}},
		withAccess:              {elems: []*storage.ComplianceRunMetadata{objA}, missingIndices: []int{1}},
		withAccessToCluster:     {elems: []*storage.ComplianceRunMetadata{objA}, missingIndices: []int{1}},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, missingIndices, err := s.store.GetMany(ctxs[name], []string{objA.GetRunId(), objB.GetRunId()})
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

func getSACContexts(obj *storage.ComplianceRunMetadata, access storage.Access) map[string]context.Context {
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
