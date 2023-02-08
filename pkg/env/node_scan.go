package env

import "time"

var (
	// NodeRescanInterval will set the duration for when to fetch node inventory to be scanned for vulnerabilities
	NodeRescanInterval = registerDurationSetting("ROX_NODE_RESCAN_INTERVAL", 4*time.Hour)

	// NodeInventoryCacheDuration defines the time after which a cached inventory is considered outdated
	NodeInventoryCacheDuration = registerDurationSetting("ROX_NODE_INVENTORY_CACHE_TIME", 10*time.Minute)

	// NodeInventoryInitialBackoff defines the initial time in seconds a Node Inventory will be delayed if a backoff file is found
	NodeInventoryInitialBackoff = RegisterIntegerSetting("ROX_NODE_INVENTORY_INITIAL_BACKOFF", 30)

	// NodeInventoryBackoffIncrement sets the seconds that are added on each interrupted run
	NodeInventoryBackoffIncrement = RegisterIntegerSetting("ROX_NODE_INVENTORY_BACKOFF_INCREMENT", 5)
)
