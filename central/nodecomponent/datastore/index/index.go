package index

import (
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// Indexer provides funtionality to index node components.
type Indexer interface {
	AddNodeComponent(components *storage.NodeComponent) error
	AddNodeComponents(components []*storage.NodeComponent) error
	Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteNodeComponent(id string) error
	DeleteNodeComponents(ids []string) error
	Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}
