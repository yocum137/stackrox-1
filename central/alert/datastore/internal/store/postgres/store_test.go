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
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type AlertsStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestAlertsStore(t *testing.T) {
	suite.Run(t, new(AlertsStoreSuite))
}

func (s *AlertsStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *AlertsStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *AlertsStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	alert := &storage.Alert{}
	s.NoError(testutils.FullInit(alert, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	alert.Entity = &storage.Alert_Deployment_{Deployment: &storage.Alert_Deployment{Id: "orig"}}

	foundAlert, exists, err := store.Get(ctx, alert.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundAlert)

	s.NoError(store.Upsert(ctx, alert))
	foundAlert, exists, err = store.Get(ctx, alert.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(alert, foundAlert)

	alertCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(alertCount, 1)

	alertExists, err := store.Exists(ctx, alert.GetId())
	s.NoError(err)
	s.True(alertExists)
	s.NoError(store.Upsert(ctx, alert))

	foundAlert, exists, err = store.Get(ctx, alert.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(alert, foundAlert)

	alert.GetDeployment().Id = "other"
	alert.Id = "new"
	s.NoError(store.Upsert(ctx, alert))

	indexer := NewIndexer(pool)
	out, _ := indexer.Search(search.NewQueryBuilder().AddStrings(search.DeploymentID, "orig").ProtoQuery())
	_ = out

	s.NoError(store.Delete(ctx, alert.GetId()))
	foundAlert, exists, err = store.Get(ctx, alert.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundAlert)
}
