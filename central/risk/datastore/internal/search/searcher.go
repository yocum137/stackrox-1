package search

import (
	"context"

	"github.com/stackrox/rox/central/risk/datastore/internal/index"
	"github.com/stackrox/rox/central/risk/datastore/internal/store"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// Searcher provides search functionality on existing risks
//go:generate mockgen-wrapper
type Searcher interface {
	Search(ctx context.Context, q *auxpb.Query) ([]search.Result, error)
	Count(ctx context.Context, q *auxpb.Query) (int, error)
	SearchRawRisks(ctx context.Context, q *auxpb.Query) ([]*storage.Risk, error)
}

// New returns a new instance of Searcher for the given storage and indexer.
func New(storage store.Store, indexer index.Indexer) Searcher {
	return &searcherImpl{
		storage:  storage,
		indexer:  indexer,
		searcher: formatSearcher(indexer),
	}
}
