package env

import "time"

var (
	// NodeRescanInterval will set the duration for when to fetch node inventory to be scanned for vulnerabilities
	NodeRescanInterval = registerDurationSetting("ROX_NODE_RESCAN_INTERVAL", 4*time.Hour)

	// NodeInventoryCacheDuration defines the time after which a cached inventory is considered outdated
	NodeInventoryCacheDuration = registerDurationSetting("ROX_NODE_INVENTORY_CACHE", 10*time.Minute)
)
