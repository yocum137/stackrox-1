package getters

import (
	"context"

	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
)

// ProcessBaselines encapsulates the sub-interface of the process baselines datastore required for risk.
type ProcessBaselines interface {
	GetProcessBaseline(context.Context, *storage.ProcessBaselineKey) (*storage.ProcessBaseline, error)
}

// ProcessIndicators encapulates the sub-interface of the process indicator datastore required for risk.
type ProcessIndicators interface {
	SearchRawProcessIndicators(ctx context.Context, q *auxpb.Query) ([]*storage.ProcessIndicator, error)
}
