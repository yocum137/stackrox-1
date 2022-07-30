package edgefields

import (
	"context"
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stretchr/testify/assert"
)

func TestGetCVEEdgeQuery(t *testing.T) {
	query := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Fixable.String(), Value: "true"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.ClusterID.String(), Value: "cluster1"},
						},
					},
				}},
			},
		}},
	}

	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_Disjunction{
					Disjunction: &auxpb.DisjunctionQuery{
						Queries: []*auxpb.Query{
							{Query: &auxpb.Query_BaseQuery{
								BaseQuery: &auxpb.BaseQuery{
									Query: &auxpb.BaseQuery_MatchFieldQuery{
										MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Fixable.String(), Value: "true"},
									},
								},
							}},
							{Query: &auxpb.Query_BaseQuery{
								BaseQuery: &auxpb.BaseQuery{
									Query: &auxpb.BaseQuery_MatchFieldQuery{
										MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.ClusterCVEFixable.String(), Value: "true"},
									},
								},
							}},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.ClusterID.String(), Value: "cluster1"},
						},
					},
				}},
			},
		}},
	}

	getCVEEdgeQuery(query)
	assert.Equal(t, expectedQuery, query)
}

func TestSnoozedQueryHandler(t *testing.T) {
	scopedCtx := scoped.Context(context.Background(), scoped.Scope{
		ID:    "img1",
		Level: v1.SearchCategory_IMAGES,
	})
	snoozedCVEsQuery := search.NewQueryBuilder().AddBools(search.CVESuppressed, true).ProtoQuery()
	observedCVEsQuery := search.NewQueryBuilder().AddBools(search.CVESuppressed, false).ProtoQuery()
	cveStateQuery := search.NewQueryBuilder().AddExactMatches(search.VulnerabilityState, storage.VulnerabilityState_DEFERRED.String(), storage.VulnerabilityState_FALSE_POSITIVE.String()).ProtoQuery()
	conjunction := search.ConjunctionQuery(snoozedCVEsQuery, cveStateQuery)

	for _, c := range []struct {
		desc     string
		incoming *auxpb.Query
		expected *auxpb.Query
		ctx      context.Context
	}{
		{
			desc:     "query is not in image scope; nothing to do",
			incoming: snoozedCVEsQuery.Clone(),
			expected: snoozedCVEsQuery,
			ctx:      context.Background(),
		},
		{
			desc:     "query is in image scope; should be updated",
			incoming: snoozedCVEsQuery.Clone(),
			expected: conjunction,
			ctx:      scopedCtx,
		},
		{
			desc:     "not querying snoozed cves; should not be updated",
			incoming: observedCVEsQuery.Clone(),
			expected: observedCVEsQuery,
			ctx:      scopedCtx,
		},
		{
			desc:     "nothing to do",
			incoming: conjunction.Clone(),
			expected: conjunction,
			ctx:      scopedCtx,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			assert.EqualValues(t, c.expected, handleSnoozedCVEQuery(c.ctx, c.incoming))
		})
	}
}
