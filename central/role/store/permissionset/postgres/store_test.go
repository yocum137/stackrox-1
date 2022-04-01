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

type PermissionsetsStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestPermissionsetsStore(t *testing.T) {
	suite.Run(t, new(PermissionsetsStoreSuite))
}

func (s *PermissionsetsStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *PermissionsetsStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *PermissionsetsStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	permissionSet := &storage.PermissionSet{}
	s.NoError(testutils.FullInit(permissionSet, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	foundPermissionSet, exists, err := store.Get(ctx, permissionSet.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPermissionSet)

	s.NoError(store.Upsert(ctx, permissionSet))
	foundPermissionSet, exists, err = store.Get(ctx, permissionSet.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(permissionSet, foundPermissionSet)

	permissionSetCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(permissionSetCount, 1)
	permissionSetExists, err := store.Exists(ctx, permissionSet.GetId())
	s.NoError(err)
	s.True(permissionSetExists)
	s.NoError(store.Upsert(ctx, permissionSet))
	foundPermissionSet, exists, err = store.Get(ctx, permissionSet.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(permissionSet, foundPermissionSet)
	s.NoError(store.Delete(ctx, permissionSet.GetId()))
	foundPermissionSet, exists, err = store.Get(ctx, permissionSet.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPermissionSet)

	var permissionSets []*storage.PermissionSet
	for i := 0; i < 200; i++ {
		permissionSet := &storage.PermissionSet{}
		s.NoError(testutils.FullInit(permissionSet, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		permissionSets = append(permissionSets, permissionSet)
	}

	s.NoError(store.UpsertMany(ctx, permissionSets))

	permissionSetCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(permissionSetCount, 200)
}
