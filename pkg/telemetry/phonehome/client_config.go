package phonehome

import "context"

// interceptor is a function which will be called on every API call if none of
// the previous interceptors in the chain returned false.
// An interceptor function may add custom properties to the props map so that
// they appear in the event.
type interceptor func(rp *RequestParams, props map[string]any) bool

// Config represents a telemetry client instance configuration.
type Config struct {
	// ClientID identifies an entity that reports telemetry data.
	ClientID string
	// GroupID identifies the main group to which the client belongs.
	GroupID string
	// Client static properties reported as client identity.
	Properties map[string]any
	// Random client configuration properties.
	Config map[string]any

	// Map of event name to the list of interceptors, that gather properties for
	// the event.
	interceptors map[string][]interceptor
}

// GatherFunc returns properties gathered by a data source.
type GatherFunc func(context.Context) (map[string]any, error)
