// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ProcessListeningOnPortsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestProcessListeningOnPortsStore(t *testing.T) {
	suite.Run(t, new(ProcessListeningOnPortsStoreSuite))
}

func (s *ProcessListeningOnPortsStoreSuite) SetupSuite() {
	s.T().Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")

	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.Pool)
}

func (s *ProcessListeningOnPortsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE process_listening_on_ports CASCADE")
	s.T().Log("process_listening_on_ports", tag)
	s.NoError(err)
}

func (s *ProcessListeningOnPortsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ProcessListeningOnPortsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	processListeningOnPortStorage := &storage.ProcessListeningOnPortStorage{}
	s.NoError(testutils.FullInit(processListeningOnPortStorage, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundProcessListeningOnPortStorage, exists, err := store.Get(ctx, processListeningOnPortStorage.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundProcessListeningOnPortStorage)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, processListeningOnPortStorage))
	foundProcessListeningOnPortStorage, exists, err = store.Get(ctx, processListeningOnPortStorage.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(processListeningOnPortStorage, foundProcessListeningOnPortStorage)

	processListeningOnPortStorageCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, processListeningOnPortStorageCount)
	processListeningOnPortStorageCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(processListeningOnPortStorageCount)

	processListeningOnPortStorageExists, err := store.Exists(ctx, processListeningOnPortStorage.GetId())
	s.NoError(err)
	s.True(processListeningOnPortStorageExists)
	s.NoError(store.Upsert(ctx, processListeningOnPortStorage))
	s.ErrorIs(store.Upsert(withNoAccessCtx, processListeningOnPortStorage), sac.ErrResourceAccessDenied)

	foundProcessListeningOnPortStorage, exists, err = store.Get(ctx, processListeningOnPortStorage.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(processListeningOnPortStorage, foundProcessListeningOnPortStorage)

	s.NoError(store.Delete(ctx, processListeningOnPortStorage.GetId()))
	foundProcessListeningOnPortStorage, exists, err = store.Get(ctx, processListeningOnPortStorage.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundProcessListeningOnPortStorage)
	s.NoError(store.Delete(withNoAccessCtx, processListeningOnPortStorage.GetId()))

	var processListeningOnPortStorages []*storage.ProcessListeningOnPortStorage
	var processListeningOnPortStorageIDs []string
	for i := 0; i < 200; i++ {
		processListeningOnPortStorage := &storage.ProcessListeningOnPortStorage{}
		s.NoError(testutils.FullInit(processListeningOnPortStorage, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		processListeningOnPortStorages = append(processListeningOnPortStorages, processListeningOnPortStorage)
		processListeningOnPortStorageIDs = append(processListeningOnPortStorageIDs, processListeningOnPortStorage.GetId())
	}

	s.NoError(store.UpsertMany(ctx, processListeningOnPortStorages))

	processListeningOnPortStorageCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, processListeningOnPortStorageCount)

	s.NoError(store.DeleteMany(ctx, processListeningOnPortStorageIDs))

	processListeningOnPortStorageCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, processListeningOnPortStorageCount)
}
