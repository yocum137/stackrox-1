package policyutils

import (
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// DeploymentExclusionToQuery returns the proto query to get all excluded deployments
func DeploymentExclusionToQuery(exclusions []*storage.Exclusion) *auxpb.Query {
	var queries []*auxpb.Query
	for _, exclusion := range exclusions {
		subqueries := make([]*auxpb.Query, 0, 2)
		if exclusion.GetDeployment() != nil {
			if exclusion.GetDeployment().GetName() != "" {
				subqueries = append(subqueries, search.NewQueryBuilder().AddExactMatches(search.DeploymentName,
					exclusion.GetDeployment().GetName()).ProtoQuery())
			}
			if exclusion.GetDeployment().GetScope() != nil {
				subqueries = append(subqueries, ScopeToQuery([]*storage.Scope{exclusion.GetDeployment().GetScope()}))
			}

			if len(subqueries) == 0 {
				continue
			}

			queries = append(queries, search.ConjunctionQuery(subqueries...))
		}
	}

	if len(queries) == 0 {
		return search.MatchNoneQuery()
	}

	return search.DisjunctionQuery(queries...)
}
