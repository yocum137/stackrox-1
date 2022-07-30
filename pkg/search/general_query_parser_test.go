package search

import (
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stretchr/testify/assert"
)

func TestParseRawQuery(t *testing.T) {
	query := fmt.Sprintf("%s:field1,field12+%s:field2", DeploymentName, Category)
	expectedQuery := &aux.Query{
		Query: &aux.Query_Conjunction{Conjunction: &aux.ConjunctionQuery{
			Queries: []*aux.Query{
				{Query: &aux.Query_BaseQuery{
					BaseQuery: &aux.BaseQuery{
						Query: &aux.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &aux.MatchFieldQuery{Field: Category.String(), Value: "field2"},
						},
					},
				}},
				{Query: &aux.Query_Disjunction{Disjunction: &aux.DisjunctionQuery{
					Queries: []*aux.Query{
						{Query: &aux.Query_BaseQuery{
							BaseQuery: &aux.BaseQuery{
								Query: &aux.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &aux.MatchFieldQuery{Field: DeploymentName.String(), Value: "field1"},
								},
							},
						}},
						{Query: &aux.Query_BaseQuery{
							BaseQuery: &aux.BaseQuery{
								Query: &aux.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &aux.MatchFieldQuery{Field: DeploymentName.String(), Value: "field12"},
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

	expectedQuery = &aux.Query{
		Query: &aux.Query_Conjunction{Conjunction: &aux.ConjunctionQuery{
			Queries: []*aux.Query{
				{Query: &aux.Query_BaseQuery{
					BaseQuery: &aux.BaseQuery{
						Query: &aux.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &aux.MatchFieldQuery{Field: Category.String(), Value: "field2"},
						},
					},
				}},
				{Query: &aux.Query_Disjunction{Disjunction: &aux.DisjunctionQuery{
					Queries: []*aux.Query{
						{Query: &aux.Query_BaseQuery{
							BaseQuery: &aux.BaseQuery{
								Query: &aux.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &aux.MatchFieldQuery{Field: DeploymentName.String(), Value: "field1"},
								},
							},
						}},
						{Query: &aux.Query_BaseQuery{
							BaseQuery: &aux.BaseQuery{
								Query: &aux.BaseQuery_MatchFieldQuery{
									MatchFieldQuery: &aux.MatchFieldQuery{Field: DeploymentName.String(), Value: "field12"},
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
