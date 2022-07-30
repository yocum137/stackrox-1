package search

import (
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stretchr/testify/assert"
)

func TestParseRawQuery(t *testing.T) {
	query := fmt.Sprintf("%s:field1,field12+%s:field2", DeploymentName, Category)
	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: Category.String(), Value: "field2"},
						},
					},
				}},
				{Query: &auxpb.Query_Disjunction{Disjunction: &auxpb.DisjunctionQuery{
					Queries: []*auxpb.Query{
						{Query: &auxpb.Query_BaseQuery{
							BaseQuery: &auxpb.BaseQuery{
								Query: &auxpb.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "field1"},
								},
							},
						}},
						{Query: &auxpb.Query_BaseQuery{
							BaseQuery: &auxpb.BaseQuery{
								Query: &auxpb.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "field12"},
								},
							},
						}},
					},
				}}},
			},
		}},
	}
	actualQuery, err := generalQueryParser{}.parse(query)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, actualQuery)

	query = fmt.Sprintf("%s:field1,field12 + %s:field2", DeploymentName, Category)

	expectedQuery = &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: Category.String(), Value: "field2"},
						},
					},
				}},
				{Query: &auxpb.Query_Disjunction{Disjunction: &auxpb.DisjunctionQuery{
					Queries: []*auxpb.Query{
						{Query: &auxpb.Query_BaseQuery{
							BaseQuery: &auxpb.BaseQuery{
								Query: &auxpb.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "field1"},
								},
							},
						}},
						{Query: &auxpb.Query_BaseQuery{
							BaseQuery: &auxpb.BaseQuery{
								Query: &auxpb.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "field12"},
								},
							},
						}},
					},
				}}},
			},
		}},
	}
	actualQuery, err = generalQueryParser{}.parse(query)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, actualQuery)

	_, err = generalQueryParser{}.parse("")
	assert.Error(t, err)
	actualQuery, err = generalQueryParser{MatchAllIfEmpty: true}.parse("")
	assert.NoError(t, err)
	assert.Equal(t, EmptyQuery(), actualQuery)

	// An invalid query should return an error.
	query = "INVALIDQUERY"
	_, err = generalQueryParser{}.parse(query)
	assert.Error(t, err)
}
