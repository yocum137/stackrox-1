package m172tom173

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/m_172_to_m_173_network_flows_partition/schema"
	"github.com/stackrox/rox/migrator/types"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/uuid"
	"gorm.io/gorm"
)

var (
	startSeqNum = 172

	migration = types.Migration{
		StartingSeqNum: startSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startSeqNum + 1)}, // 173
		Run: func(databases *types.Databases) error {
			err := MigrateToPartitions(databases.GormDB, databases.PostgresDB)
			if err != nil {
				return errors.Wrap(err, "updating network_flows to partitions")
			}
			return nil
		},
	}

	log = loghelper.LogWrapper{}
)

// MigrateToPartitions updates the btree network flow indexes to be hash
func MigrateToPartitions(gormDB *gorm.DB, db *postgres.DB) error {
	// First get the distinct clusters in the network_flows table
	clusters, err := getClusters(db)
	if err != nil {
		log.WriteToStderrf("unable to retrieve clusters from network_flows, %v", err)
		return err
	}

	// Rename the original table
	err = gormDB.Migrator().RenameTable("network_flows", "network_flows_no_partition")
	if err != nil {
		log.WriteToStderrf("unable to rename network_flows to network_flows_no_partition, %v", err)
		return err
	}

	// Now apply the updated schema to create a partition table with updated index types
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNetworkFlowsStmt)

	// Check what is up with indexes
	indexes, err := gormDB.Migrator().GetIndexes("network_flows")
	log.WriteToStderrf("Indexes right after create -- %v", indexes)

	// Create the partition and move the data
	for _, cluster := range clusters {
		err = createPartition(db, cluster)
		if err != nil {
			log.WriteToStderrf("unable to create partition for cluster %q, %v", cluster, err)
			return err
		}

		err = migrateData(db, cluster)
		if err != nil {
			log.WriteToStderrf("unable to move data for cluster %q, %v", cluster, err)
			return err
		}
	}

	// Drop the old table
	err = gormDB.Migrator().DropTable("network_flows_no_partition")
	if err != nil {
		log.WriteToStderrf("unable to drop table network_flows_no_partition, %v", err)
		return err
	}

	indexes, err = gormDB.Migrator().GetIndexes("network_flows")
	log.WriteToStderrf("Indexes at end of the migration -- %v", indexes)

	return nil
}

func getClusters(db *postgres.DB) ([]string, error) {
	var clusters []string
	getClustersStmt := "select distinct clusterid from network_flows;"

	row, err := db.Query(context.Background(), getClustersStmt)
	if err != nil {
		return nil, err
	}

	defer row.Close()
	for row.Next() {
		var cluster string
		if err := row.Scan(&cluster); err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func createPartition(db *postgres.DB, cluster string) error {
	partitionPostFix := strings.ReplaceAll(cluster, "-", "_")
	partitionCreate := `create table if not exists network_flows_%s partition of network_flows 
		for values in ('%s')`

	_, err := db.Exec(context.Background(), fmt.Sprintf(partitionCreate, partitionPostFix, cluster))
	if err != nil {
		return err
	}

	return nil
}

func migrateData(db *postgres.DB, cluster string) error {
	partitionPostFix := strings.ReplaceAll(cluster, "-", "_")
	moveDataStmt := fmt.Sprintf("INSERT INTO network_flows_%s SELECT * FROM network_flows_no_partition WHERE clusterid = $1", partitionPostFix)

	clusterUUID, err := uuid.FromString(cluster)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), moveDataStmt, clusterUUID)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
