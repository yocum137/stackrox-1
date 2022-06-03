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

type NetworkpoliciesStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestNetworkpoliciesStore(t *testing.T) {
	suite.Run(t, new(NetworkpoliciesStoreSuite))
}

func (s *NetworkpoliciesStoreSuite) SetupTest() {
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

func (s *NetworkpoliciesStoreSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *NetworkpoliciesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	networkPolicy := &storage.NetworkPolicy{}
	s.NoError(testutils.FullInit(networkPolicy, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundNetworkPolicy, exists, err := store.Get(ctx, networkPolicy.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkPolicy)

	s.NoError(store.Upsert(ctx, networkPolicy))
	foundNetworkPolicy, exists, err = store.Get(ctx, networkPolicy.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkPolicy, foundNetworkPolicy)

	networkPolicyCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, networkPolicyCount)

	networkPolicyExists, err := store.Exists(ctx, networkPolicy.GetId())
	s.NoError(err)
	s.True(networkPolicyExists)
	s.NoError(store.Upsert(ctx, networkPolicy))

	foundNetworkPolicy, exists, err = store.Get(ctx, networkPolicy.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkPolicy, foundNetworkPolicy)

	s.NoError(store.Delete(ctx, networkPolicy.GetId()))
	foundNetworkPolicy, exists, err = store.Get(ctx, networkPolicy.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkPolicy)

	var networkPolicys []*storage.NetworkPolicy
	for i := 0; i < 200; i++ {
		networkPolicy := &storage.NetworkPolicy{}
		s.NoError(testutils.FullInit(networkPolicy, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkPolicys = append(networkPolicys, networkPolicy)
	}

	s.NoError(store.UpsertMany(ctx, networkPolicys))

	networkPolicyCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, networkPolicyCount)
}
