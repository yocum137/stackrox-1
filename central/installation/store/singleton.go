package store

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/installation/store/bolt"
	"github.com/stackrox/rox/central/installation/store/postgres"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/uuid"
)

var (
	storeSingleton Store
	singletonInit  sync.Once
)

// Singleton returns a singleton of the InstallationInfo store
func Singleton() Store {
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.InstallationInfo)))

	singletonInit.Do(func() {
		var internalStore internalStore
		if features.PostgresDatastore.Enabled() {
			internalStore = postgres.New(context.TODO(), globaldb.GetPostgres())
		} else {
			internalStore = bolt.New(globaldb.GetGlobalDB())
		}
		storeSingleton = NewStore(internalStore)

		info, err := storeSingleton.GetInstallationInfo(ctx)
		if err != nil {
			panic(err)
		}

		if info == nil {
			info := &storage.InstallationInfo{
				Id:      uuid.NewV4().String(),
				Created: types.TimestampNow(),
			}
			err = storeSingleton.AddInstallationInfo(ctx, info)
			if err != nil {
				panic(err)
			}
		}
	})
	return storeSingleton
}
