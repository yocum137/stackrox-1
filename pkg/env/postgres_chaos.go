package env

var (
	// PostgresChaosPercent defines the percentage of queries that will fail with transient error codes
	PostgresChaosPercent = RegisterIntegerSetting("ROX_POSTGRES_CHAOS_PERCENT", 0)
)
