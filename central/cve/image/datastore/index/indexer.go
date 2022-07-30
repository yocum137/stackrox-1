package index

import (
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// Indexer provides functionality to index node cves.
//go:generate mockgen-wrapper
type Indexer interface {
	AddImageCVE(cve *storage.ImageCVE) error
	AddImageCVEs(cves []*storage.ImageCVE) error
	Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteImageCVE(id string) error
	DeleteImageCVEs(ids []string) error
	Search(q *aux.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}
