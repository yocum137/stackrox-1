package parser

import (
	"net/url"
	"testing"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/assert"
)

func TestParseURLQuery(t *testing.T) {
	vals := url.Values{
		"query":                          []string{"Namespace:ABC"},
		"pagination.offset":              []string{"5"},
		"pagination.limit":               []string{"50"},
		"pagination.sortOption.field":    []string{"Deployment"},
		"pagination.sortOption.reversed": []string{"true"},
	}

	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Namespace.String(), Value: "ABC"},
				},
			},
		},
		Pagination: &auxpb.QueryPagination{
			Offset: 5,
			Limit:  50,
			SortOptions: []*auxpb.QuerySortOption{
				{
					Field:    search.DeploymentName.String(),
					Reversed: true,
				},
			},
		},
	}

	actual, _, err := ParseURLQuery(vals)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, actual)
}

func TestParseURLQueryWithExtraValues(t *testing.T) {
	vals := url.Values{
		"query":                          []string{"Namespace:ABC"},
		"pagination.offset":              []string{"5"},
		"pagination.limit":               []string{"50"},
		"pagination.sortOption.field":    []string{"Deployment"},
		"pagination.sortOption.reversed": []string{"true"},
		"blah":                           []string{"blah"},
	}

	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Namespace.String(), Value: "ABC"},
				},
			},
		},
		Pagination: &auxpb.QueryPagination{
			Offset: 5,
			Limit:  50,
			SortOptions: []*auxpb.QuerySortOption{
				{
					Field:    search.DeploymentName.String(),
					Reversed: true,
				},
			},
		},
	}

	actual, _, err := ParseURLQuery(vals)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, actual)
}

func TestParseURLQueryConjunctionQuery(t *testing.T) {
	vals := url.Values{
		"query":                          []string{"Namespace:ABC+Cluster:ABC"},
		"pagination.offset":              []string{"5"},
		"pagination.limit":               []string{"50"},
		"pagination.sortOption.field":    []string{"Deployment"},
		"pagination.sortOption.reversed": []string{"true"},
	}

	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{
					Query: &auxpb.Query_BaseQuery{
						BaseQuery: &auxpb.BaseQuery{
							Query: &auxpb.BaseQuery_MatchFieldQuery{
								MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Cluster.String(), Value: "ABC"},
							},
						},
					},
				},
				{
					Query: &auxpb.Query_BaseQuery{
						BaseQuery: &auxpb.BaseQuery{
							Query: &auxpb.BaseQuery_MatchFieldQuery{
								MatchFieldQuery: &auxpb.MatchFieldQuery{Field: search.Namespace.String(), Value: "ABC"},
							},
						},
					},
				},
			},
		}},
		Pagination: &auxpb.QueryPagination{
			Offset: 5,
			Limit:  50,
			SortOptions: []*auxpb.QuerySortOption{
				{
					Field:    search.DeploymentName.String(),
					Reversed: true,
				},
			},
		},
	}

	actual, _, err := ParseURLQuery(vals)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedQuery, actual)
}
