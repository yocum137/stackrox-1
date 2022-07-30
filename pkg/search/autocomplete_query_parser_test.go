package search

import (
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stretchr/testify/assert"
)

func TestParseAutocompleteQuery(t *testing.T) {
	query := fmt.Sprintf("%s:field1,field12+%s:field2", DeploymentName, Category)
	expectedQuery := &aux.Query{
		Query: &aux.Query_Conjunction{Conjunction: &aux.ConjunctionQuery{
			Queries: []*aux.Query{
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
				{Query: &aux.Query_BaseQuery{
					BaseQuery: &aux.BaseQuery{
						Query: &aux.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &aux.MatchFieldQuery{Field: Category.String(), Value: "field2", Highlight: true},
						},
					},
				}},
			},
		}},
	}
	expectedKey := Category.String()

	var actualKey string
	actualQuery, actualKey, err := autocompleteQueryParser{}.parse(query)
	assert.NoError(t, err)
	assert.Equal(t, expectedKey, actualKey)
	assert.Equal(t, expectedQuery, actualQuery)

	_, _, err = autocompleteQueryParser{}.parse("")
	assert.Error(t, err)

	// An invalid query should always return an error.
	query = "INVALIDQUERY"
	_, _, err = autocompleteQueryParser{}.parse(query)
	assert.Error(t, err)
}
