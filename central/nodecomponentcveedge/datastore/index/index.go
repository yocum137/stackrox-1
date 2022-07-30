package index

import (
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// Indexer provides indexing functionality for storage.NodeComponentCVEEdge objects.
type Indexer interface {
	AddNodeComponentCVEEdge(componentcveedge *storage.NodeComponentCVEEdge) error
	AddNodeComponentCVEEdges(componentcveedges []*storage.NodeComponentCVEEdge) error
	Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteNodeComponentCVEEdge(id string) error
	DeleteNodeComponentCVEEdges(ids []string) error
	Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}
