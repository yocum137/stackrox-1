package sac

import (
	"github.com/stackrox/rox/central/dackbox"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/sac/helpers"
	"github.com/stackrox/rox/pkg/search/filtered"
)

var (
	nsSAC = helpers.ForResource(resources.Namespace)

	nsSACFilter = filtered.MustCreateNewSACFilter(
		filtered.WithResourceHelper(nsSAC),
		filtered.WithScopeTransform(dackbox.NamespaceSACTransform),
		filtered.WithReadAccess())
)

// GetSACFilter returns the sac filter for image ids.
func GetSACFilter() filtered.Filter {
	return nsSACFilter
}
