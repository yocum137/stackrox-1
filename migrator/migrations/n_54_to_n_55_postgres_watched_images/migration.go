// Code generated by pg-bindings generator. DO NOT EDIT.
package n54ton55

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_54_to_n_55_postgres_watched_images/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_54_to_n_55_postgres_watched_images/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 54 // 165

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum+1)}, // 166
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving watched_images from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = frozenSchema.WatchedImagesSchema
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableWatchedImagesStmt)
	var watchedImages []*storage.WatchedImage
	err := walk(ctx, legacyStore, func(obj *storage.WatchedImage) error {
		watchedImages = append(watchedImages, obj)
		if len(watchedImages) == batchSize {
			if err := store.UpsertMany(ctx, watchedImages); err != nil {
				log.WriteToStderrf("failed to persist watched_images to store %v", err)
				return err
			}
			watchedImages = watchedImages[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(watchedImages) > 0 {
		if err = store.UpsertMany(ctx, watchedImages); err != nil {
			log.WriteToStderrf("failed to persist watched_images to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.WatchedImage) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
