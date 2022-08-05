package store

import (
	"context"

	"github.com/stackrox/rox/central/cve/converter"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

// Store provides storage functionality for cluster-cve edges.
//go:generate mockgen-wrapper
type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storage.ClusterCVEEdge, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.ClusterCVEEdge, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.ClusterCVEEdge, []int, error)

	Upsert(ctx context.Context, cveParts ...converter.ClusterCVEParts) error
	Delete(ctx context.Context, ids ...string) error
}
