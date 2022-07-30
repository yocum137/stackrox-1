package index

import (
	bleve "github.com/blevesearch/bleve"
	"github.com/stackrox/rox/generated/auxpb"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// Indexer is the process indicator indexer.
//go:generate mockgen-wrapper
type Indexer interface {
	AddProcessIndicator(processindicator *storage.ProcessIndicator) error
	AddProcessIndicators(processindicators []*storage.ProcessIndicator) error
	DeleteProcessIndicator(id string) error
	DeleteProcessIndicators(ids []string) error
	MarkInitialIndexingComplete() error
	NeedsInitialIndexing() (bool, error)
	Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
	Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error)
}

// New returns a new process indicator indexer.
func New(index bleve.Index) Indexer {
	return &indexerImpl{index: index}
}
