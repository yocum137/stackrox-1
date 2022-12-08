package phonehome

import (
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/phonehome/segment"
)

var (
	onceTelemeter sync.Once
)

// Enabled returns true if telemetry data collection is enabled.
func Enabled() bool {
	return segment.Enabled()
}

// Telemeter returns the instance of the telemeter.
func (cfg *Config) Telemeter() Telemeter {
	onceTelemeter.Do(func() {
		cfg.telemeter = segment.NewTelemeter(cfg.ClientID, cfg.Properties)
		// Central adds itself to the tenant group, adding its properties to the
		// group properties:
		cfg.telemeter.Group(cfg.GroupID, cfg.ClientID, cfg.Properties)
		// Add the local admin user as well:
		cfg.telemeter.Group(cfg.GroupID, "local:"+cfg.ClientID+":admin", cfg.Properties)
	})
	return cfg.telemeter
}
