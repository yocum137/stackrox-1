// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	storage "github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type PolicyStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestPolicyStore(t *testing.T) {
	suite.Run(t, new(PolicyStoreSuite))
}

func (s *PolicyStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *PolicyStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *PolicyStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	policy := &storage.Policy{}
	s.NoError(testutils.FullInit(policy, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	foundPolicy, exists, err := store.Get(ctx, policy.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPolicy)

	s.NoError(store.Upsert(ctx, policy))
	foundPolicy, exists, err = store.Get(ctx, policy.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(policy, foundPolicy)

	policyCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(policyCount, 1)
	policyExists, err := store.Exists(ctx, policy.GetId())
	s.NoError(err)
	s.True(policyExists)
	s.NoError(store.Upsert(ctx, policy))
	foundPolicy, exists, err = store.Get(ctx, policy.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(policy, foundPolicy)
	s.NoError(store.Delete(ctx, policy.GetId()))
	foundPolicy, exists, err = store.Get(ctx, policy.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPolicy)

	var policys []*storage.Policy
	for i := 0; i < 200; i++ {
		policy := &storage.Policy{}
		s.NoError(testutils.FullInit(policy, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		policys = append(policys, policy)
	}

	s.NoError(store.UpsertMany(ctx, policys))

	policyCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(policyCount, 200)
}
