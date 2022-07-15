package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gogo/protobuf/types"
	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	configDatastore "github.com/stackrox/rox/central/config/datastore"
	"github.com/stackrox/rox/central/cve/datastore"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/globaldb/export"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/central/option"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/sac"
)

func main() {
	if features.PostgresDatastore.Enabled() {
		return
	}

	migrations.DBMountPathInt = "local/database-restore/full"
	option.CentralOptions.DBPathBase = migrations.CurrentPath()

	blevePath := filepath.Join(option.CentralOptions.DBPathBase, "bleve")
	globalindex.DefaultBlevePath = filepath.Join(blevePath, "default")
	globalindex.DefaultTmpBlevePath = filepath.Join(blevePath, "tmp")
	globalindex.SeparateIndexPath = filepath.Join(blevePath, "separate")

	// Can start accessing _most_ singletons. Some singletons that access certificates will fail
	datastore.Singleton()

	ctx := sac.WithAllAccess(context.Background())
	configDS := configDatastore.Singleton()
	config, err := getConfig(configDS, ctx)
	if err != nil {
		return
	}

	zapLevel, _ := logging.LevelForLabel("Warn")
	logging.SetGlobalLogLevel(zapLevel)

	fmt.Println("=== Initial Values ===")

	if err := printValues(config, ctx); err != nil {
		return
	}

	zapLevel, _ = logging.LevelForLabel("Info")
	logging.SetGlobalLogLevel(zapLevel)

	fmt.Println("=== Saving New Values ===")
	config.GetPrivateConfig().GetDecommissionedClusterRetention().RetentionDurationDays = 15
	config.GetPrivateConfig().GetDecommissionedClusterRetention().CreatedAt = updateTime(-50 * 24 * time.Hour)
	config.GetPrivateConfig().GetDecommissionedClusterRetention().IgnoreClusterLabels = map[string]string{"shouldDecomm": "false", "k2": "v2"}
	if err := configDS.UpsertConfig(ctx, config); err != nil {
		fmt.Printf("Failed to save new config: %+v\n", err)
		return
	}

	// do it again because other changes will force the last update to get reset
	config.GetPrivateConfig().GetDecommissionedClusterRetention().LastUpdated = updateTime(-50 * 24 * time.Hour)
	if err := configDS.UpsertConfig(ctx, config); err != nil {
		fmt.Printf("Failed to save new config: %+v\n", err)
		return
	}

	clusterDS := clusterDatastore.Singleton()
	clusters, _ := clusterDS.GetClusters(ctx)
	for _, c := range clusters {

		if c.GetName() == "production" {
			c.GetHealthStatus().LastContact = updateTime(-100 * 24 * time.Hour)
		} else {
			c.GetHealthStatus().LastContact = updateTime(-100 * 24 * time.Hour)
		}
		if err := clusterDS.UpdateClusterHealth(ctx, c.GetId(), c.GetHealthStatus()); err != nil {
			fmt.Printf("Failed to update cluster health for cluster %s : %q", c.GetName(), err)
			return
		}

		if c.GetName() == "production" {
			c.Labels = map[string]string{
				"shouldDecomm":"false",
				"otherKey": "ja",
			}

			if err := clusterDS.UpdateCluster(ctx, c); err != nil {
				fmt.Printf("Failed to update cluster labels for cluster %s : %q", c.GetName(), err)
				return
			}
		}
	}

	fmt.Println("=== Values After Setting ===")

	config, err = getConfig(configDS, ctx)
	if err != nil {
		return
	}
	if err := printValues(config, ctx); err != nil {
		return
	}

	// contract DB
	f, err := os.OpenFile("local/fixed_db_export.zip", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	boltDB, rocksDB := globaldb.GetGlobalDB(), globaldb.GetRocksDB()
	if err := export.Backup(ctx, boltDB, rocksDB, true, f); err != nil {
		fmt.Printf("Failed to write contracted db: %q\n", err)
	}
}

func updateTime(duration time.Duration) *types.Timestamp {
	newTime, _ := types.TimestampProto(time.Now().Add(duration))
	return newTime
}

func getConfig(configDS configDatastore.DataStore, ctx context.Context) (*storage.Config, error) {
	config, err := configDS.GetConfig(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if config == nil {
		fmt.Println("UNEXPECTED: Got nil config")
		return nil, err
	}
	return config, nil
}

func printValues(config *storage.Config, ctx context.Context) error {
	pvtConfig := config.GetPrivateConfig()
	fmt.Printf("Cluster Retention: %d\n", pvtConfig.GetDecommissionedClusterRetention().GetRetentionDurationDays())

	configCreatedAt, _ := types.TimestampFromProto(pvtConfig.GetDecommissionedClusterRetention().GetCreatedAt())
	fmt.Printf("Cluster config created: %s\n", configCreatedAt.Format(time.RFC1123Z))

	lastUpdateTime, _ := types.TimestampFromProto(pvtConfig.GetDecommissionedClusterRetention().GetLastUpdated())
	fmt.Printf("Cluster config last updated: %s\n", lastUpdateTime.Format(time.RFC1123Z))

	fmt.Printf("Cluster config ignore labels:\n")
	for k, v := range pvtConfig.GetDecommissionedClusterRetention().GetIgnoreClusterLabels() {
		fmt.Printf("    %s => %v\n", k, v)
	}

	clusterDS := clusterDatastore.Singleton()
	clusters, err := clusterDS.GetClusters(ctx)
	if err != nil {
		fmt.Println("UNEXPECTED: Got no clusters")
		return err
	}
	for _, c := range clusters {
		contactTime, _ := types.TimestampFromProto(c.GetHealthStatus().GetLastContact())
		fmt.Printf("[Cluster %s] Status: %q -- Last Contact Time: %s\n", c.GetName(), c.GetHealthStatus().GetSensorHealthStatus(), contactTime.Format(time.RFC1123Z))
		fmt.Printf("[Cluster %s] Labels:\n", c.GetName())

		for k, v := range c.GetLabels() {
			fmt.Printf("    %s => %v\n", k, v)
		}
	}

	return nil
}
