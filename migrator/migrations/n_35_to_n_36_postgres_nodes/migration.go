// Code originally generated by pg-bindings generator. DO NOT EDIT.

package n35ton36

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/dackboxhelper"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	"github.com/stackrox/rox/migrator/migrations/n_35_to_n_36_postgres_nodes/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_35_to_n_36_postgres_nodes/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	nodeConverter "github.com/stackrox/rox/pkg/nodes/converter"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 35 // 146

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum + 1)}, // 147
		Run: func(databases *types.Databases) error {
			legacyStore := legacy.New(dackboxhelper.GetMigrationDackBox(), dackboxhelper.GetMigrationKeyFence(), true)
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving nodes from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 500
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB, true)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNodesStmt)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNodeCvesStmt)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNodeComponentsStmt)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNodeComponentEdgesStmt)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNodeComponentsCvesEdgesStmt)
	return walk(ctx, legacyStore, func(obj *storage.Node) error {
		nodeConverter.FillV2NodeVulnerabilities(obj)
		if err := store.Upsert(ctx, obj); err != nil {
			log.WriteToStderrf("failed to persist nodes to store %v", err)
			return err
		}
		return nil
	})
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.Node) error) error {
	return store_walk(ctx, s, fn)
}

func store_walk(ctx context.Context, s legacy.Store, fn func(obj *storage.Node) error) error {
	ids, err := s.GetIDs(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize

		if end > len(ids) {
			end = len(ids)
		}
		objs, _, err := s.GetMany(ctx, ids[i:end])
		if err != nil {
			return err
		}
		for _, obj := range objs {
			if err = fn(obj); err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
