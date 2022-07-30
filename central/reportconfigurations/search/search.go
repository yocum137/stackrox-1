package search

import (
	"context"

	"github.com/stackrox/rox/central/reportconfigurations/index"
	"github.com/stackrox/rox/central/reportconfigurations/store"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/search"
)

var (
	log = logging.LoggerForModule()
)

// Searcher provides search functionality on existing report configurations.
//go:generate mockgen-wrapper
type Searcher interface {
	Search(ctx context.Context, query *auxpb.Query) ([]search.Result, error)
	SearchReportConfigurations(ctx context.Context, query *auxpb.Query) ([]*storage.ReportConfiguration, error)
	Count(ctx context.Context, query *auxpb.Query) (int, error)
}

// New returns a new instance of Searcher for the given storage and index.
func New(storage store.Store, indexer index.Indexer) *searcherImpl {
	return &searcherImpl{
		storage:  storage,
		indexer:  indexer,
		searcher: formatSearcher(indexer),
	}
}
