// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type TestMultiKeyStructsStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestTestMultiKeyStructsStore(t *testing.T) {
	suite.Run(t, new(TestMultiKeyStructsStoreSuite))
}

func (s *TestMultiKeyStructsStoreSuite) SetupSuite() {
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

func (s *TestMultiKeyStructsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.pool.Exec(ctx, "TRUNCATE test_multi_key_structs CASCADE")
	s.T().Log("test_multi_key_structs", tag)
	s.NoError(err)
}

func (s *TestMultiKeyStructsStoreSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *TestMultiKeyStructsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	testMultiKeyStruct := &storage.TestMultiKeyStruct{}
	s.NoError(testutils.FullInit(testMultiKeyStruct, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundTestMultiKeyStruct, exists, err := store.Get(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestMultiKeyStruct)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, testMultiKeyStruct))
	foundTestMultiKeyStruct, exists, err = store.Get(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2())
	s.NoError(err)
	s.True(exists)
	s.Equal(testMultiKeyStruct, foundTestMultiKeyStruct)

	testMultiKeyStructCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, testMultiKeyStructCount)
	testMultiKeyStructCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(testMultiKeyStructCount)

	testMultiKeyStructExists, err := store.Exists(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2())
	s.NoError(err)
	s.True(testMultiKeyStructExists)
	s.NoError(store.Upsert(ctx, testMultiKeyStruct))
	s.ErrorIs(store.Upsert(withNoAccessCtx, testMultiKeyStruct), sac.ErrResourceAccessDenied)

	foundTestMultiKeyStruct, exists, err = store.Get(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2())
	s.NoError(err)
	s.True(exists)
	s.Equal(testMultiKeyStruct, foundTestMultiKeyStruct)

	s.NoError(store.Delete(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2()))
	foundTestMultiKeyStruct, exists, err = store.Get(ctx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestMultiKeyStruct)
	s.NoError(store.Delete(withNoAccessCtx, testMultiKeyStruct.GetKey1(), testMultiKeyStruct.GetKey2()))

	var testMultiKeyStructs []*storage.TestMultiKeyStruct
	for i := 0; i < 200; i++ {
		testMultiKeyStruct := &storage.TestMultiKeyStruct{}
		s.NoError(testutils.FullInit(testMultiKeyStruct, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		testMultiKeyStructs = append(testMultiKeyStructs, testMultiKeyStruct)
	}

	s.NoError(store.UpsertMany(ctx, testMultiKeyStructs))

	testMultiKeyStructCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, testMultiKeyStructCount)
}
