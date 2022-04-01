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

type NamespacesStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestNamespacesStore(t *testing.T) {
	suite.Run(t, new(NamespacesStoreSuite))
}

func (s *NamespacesStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *NamespacesStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *NamespacesStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	namespaceMetadata := &storage.NamespaceMetadata{}
	s.NoError(testutils.FullInit(namespaceMetadata, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	foundNamespaceMetadata, exists, err := store.Get(ctx, namespaceMetadata.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNamespaceMetadata)

	s.NoError(store.Upsert(ctx, namespaceMetadata))
	foundNamespaceMetadata, exists, err = store.Get(ctx, namespaceMetadata.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(namespaceMetadata, foundNamespaceMetadata)

	namespaceMetadataCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(namespaceMetadataCount, 1)
	namespaceMetadataExists, err := store.Exists(ctx, namespaceMetadata.GetId())
	s.NoError(err)
	s.True(namespaceMetadataExists)
	s.NoError(store.Upsert(ctx, namespaceMetadata))
	foundNamespaceMetadata, exists, err = store.Get(ctx, namespaceMetadata.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(namespaceMetadata, foundNamespaceMetadata)
	s.NoError(store.Delete(ctx, namespaceMetadata.GetId()))
	foundNamespaceMetadata, exists, err = store.Get(ctx, namespaceMetadata.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNamespaceMetadata)

	var namespaceMetadatas []*storage.NamespaceMetadata
	for i := 0; i < 200; i++ {
		namespaceMetadata := &storage.NamespaceMetadata{}
		s.NoError(testutils.FullInit(namespaceMetadata, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		namespaceMetadatas = append(namespaceMetadatas, namespaceMetadata)
	}

	s.NoError(store.UpsertMany(ctx, namespaceMetadatas))

	namespaceMetadataCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(namespaceMetadataCount, 200)
}
