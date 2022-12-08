package phonehome

import (
	"context"
	"time"
)

// Interceptor is a function which will be called on every API call if none of
// the previous interceptors in the chain returned false.
// An Interceptor function may add custom properties to the props map so that
// they appear in the event.
type Interceptor func(rp *RequestParams, props map[string]any) bool

// Config represents a telemetry client instance configuration.
type Config struct {
	// ClientID identifies an entity that reports telemetry data.
	ClientID string
	// ClientName tells what kind of client is sending data.
	ClientName string
	// GroupID identifies the main group to which the client belongs.
	GroupID string
	// The period of identity gathering. Default is 1 hour.
	GatherPeriod time.Duration

	telemeter Telemeter
	gatherer  Gatherer

	// Map of event name to the list of interceptors, that gather properties for
	// the event.
	interceptors map[string][]Interceptor
}

// GatherFunc returns properties gathered by a data source.
type GatherFunc func(context.Context) (map[string]any, error)
