package search

import (
	"context"

	"github.com/stackrox/rox/central/clustercveedge/index"
	"github.com/stackrox/rox/central/clustercveedge/store"
	cveIndex "github.com/stackrox/rox/central/cve/index"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/search"
)

// Searcher provides search functionality on existing cves.
//go:generate mockgen-wrapper
type Searcher interface {
	Search(ctx context.Context, query *auxpb.Query) ([]search.Result, error)
	Count(ctx context.Context, query *auxpb.Query) (int, error)
	SearchEdges(context.Context, *auxpb.Query) ([]*v1.SearchResult, error)
	SearchRawEdges(ctx context.Context, query *auxpb.Query) ([]*storage.ClusterCVEEdge, error)
}

// New returns a new instance of Searcher for the given storage and index.
func New(storage store.Store, clusterCVEEdgeIndexer index.Indexer, cveIndexer cveIndex.Indexer, graphProvider graph.Provider) Searcher {
	return &searcherImpl{
		storage:       storage,
		indexer:       clusterCVEEdgeIndexer,
		searcher:      formatSearcher(clusterCVEEdgeIndexer, cveIndexer),
		graphProvider: graphProvider,
	}
}
