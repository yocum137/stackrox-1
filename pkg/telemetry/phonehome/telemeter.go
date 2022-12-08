package phonehome

import (
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/phonehome/segment"
)

var (
	onceTelemeter sync.Once
)

// Telemeter defines a common interface for telemetry gatherers.
//go:generate mockgen-wrapper
type Telemeter interface {
	Start()
	Stop()
	Identify(userID string, props map[string]any)
	Track(event, userID string, props map[string]any)
	Group(groupID, userID string, props map[string]any)
}

// Enabled returns true if telemetry data collection is enabled.
func Enabled() bool {
	return segment.Enabled()
}

// Telemeter returns the instance of the telemeter.
func (cfg *Config) Telemeter() Telemeter {
	onceTelemeter.Do(func() {
		cfg.telemeter = segment.NewTelemeter(cfg.ClientID, cfg.ClientName)
	})
	return cfg.telemeter
}
