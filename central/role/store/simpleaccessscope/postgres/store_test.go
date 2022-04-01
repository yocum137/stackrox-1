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
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type SimpleaccessscopesStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestSimpleaccessscopesStore(t *testing.T) {
	suite.Run(t, new(SimpleaccessscopesStoreSuite))
}

func (s *SimpleaccessscopesStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *SimpleaccessscopesStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *SimpleaccessscopesStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	simpleAccessScope := &storage.SimpleAccessScope{}
	s.NoError(testutils.FullInit(simpleAccessScope, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	foundSimpleAccessScope, exists, err := store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSimpleAccessScope)

	s.NoError(store.Upsert(ctx, simpleAccessScope))
	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(simpleAccessScope, foundSimpleAccessScope)

	simpleAccessScopeCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(simpleAccessScopeCount, 1)
	simpleAccessScopeExists, err := store.Exists(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(simpleAccessScopeExists)
	s.NoError(store.Upsert(ctx, simpleAccessScope))
	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(simpleAccessScope, foundSimpleAccessScope)
	s.NoError(store.Delete(ctx, simpleAccessScope.GetId()))
	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSimpleAccessScope)

	var simpleAccessScopes []*storage.SimpleAccessScope
	for i := 0; i < 200; i++ {
		simpleAccessScope := &storage.SimpleAccessScope{}
		s.NoError(testutils.FullInit(simpleAccessScope, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		simpleAccessScopes = append(simpleAccessScopes, simpleAccessScope)
	}

	s.NoError(store.UpsertMany(ctx, simpleAccessScopes))

	simpleAccessScopeCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(simpleAccessScopeCount, 200)
}
