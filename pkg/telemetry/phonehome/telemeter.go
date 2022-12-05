package phonehome

// Telemeter defines a common interface for telemetry gatherers.
//go:generate mockgen-wrapper
type Telemeter interface {
	Start()
	Stop()
	GetID() string
	Identify(props map[string]any)
	Track(event, userID string, props map[string]any)
	Group(groupID, userID string, props map[string]any)
}
