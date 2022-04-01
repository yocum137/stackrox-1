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

type NetworkentityStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestNetworkentityStore(t *testing.T) {
	suite.Run(t, new(NetworkentityStoreSuite))
}

func (s *NetworkentityStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *NetworkentityStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *NetworkentityStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	networkEntity := &storage.NetworkEntity{}
	s.NoError(testutils.FullInit(networkEntity, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	foundNetworkEntity, exists, err := store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkEntity)

	s.NoError(store.Upsert(ctx, networkEntity))
	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkEntity, foundNetworkEntity)

	networkEntityCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(networkEntityCount, 1)
	networkEntityExists, err := store.Exists(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(networkEntityExists)
	s.NoError(store.Upsert(ctx, networkEntity))
	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkEntity, foundNetworkEntity)
	s.NoError(store.Delete(ctx, networkEntity.GetInfo().GetId()))
	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkEntity)

	var networkEntitys []*storage.NetworkEntity
	for i := 0; i < 200; i++ {
		networkEntity := &storage.NetworkEntity{}
		s.NoError(testutils.FullInit(networkEntity, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkEntitys = append(networkEntitys, networkEntity)
	}

	s.NoError(store.UpsertMany(ctx, networkEntitys))

	networkEntityCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(networkEntityCount, 200)
}
