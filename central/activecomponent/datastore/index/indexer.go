package index

import (
	"github.com/blevesearch/bleve"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

//go:generate mockgen-wrapper
// Indexer is the interface for indexing active component
type Indexer interface {
	Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error)
	Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}

// New creates a new active component indexer
func New(index bleve.Index) Indexer {
	return &indexerImpl{index: index}
}
