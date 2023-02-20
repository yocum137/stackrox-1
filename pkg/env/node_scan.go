package env

import "time"

var (
	// NodeRescanInterval sets the time between periodic node rescans for changed components. Results are cached in case they need to be resubmitted for NodeInventoryCacheDuration.
	NodeRescanInterval = registerDurationSetting("ROX_NODE_RESCAN_INTERVAL", 4*time.Hour)

	// NodeInventoryCacheDuration defines the time after which a cached inventory is considered outdated. Defaults to 90% of rescan interval.
	NodeInventoryCacheDuration = registerDurationSetting("ROX_NODE_INVENTORY_CACHE_TIME", time.Duration(NodeRescanInterval.DurationSetting().Nanoseconds()-NodeRescanInterval.DurationSetting().Nanoseconds()/10))

	// NodeInventoryInitialBackoff defines the initial time in seconds a Node Inventory will be delayed if a backoff file is found
	NodeInventoryInitialBackoff = registerDurationSetting("ROX_NODE_INVENTORY_INITIAL_BACKOFF", 30*time.Second)

	// NodeInventoryBackoffIncrement sets the seconds that are added on each interrupted run
	NodeInventoryBackoffIncrement = registerDurationSetting("ROX_NODE_INVENTORY_BACKOFF_INCREMENT", 5*time.Second)

	// NodeInventoryMaxBackoff is the upper boundary of backoff. Defaults to 5m in seconds, being 50% of Kubernetes restart policy stability timer.
	NodeInventoryMaxBackoff = registerDurationSetting("ROX_NODE_INVENTORY_MAX_BACKOFF", 300*time.Second)
)
