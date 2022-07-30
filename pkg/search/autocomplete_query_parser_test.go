package search

import (
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stretchr/testify/assert"
)

func TestParseAutocompleteQuery(t *testing.T) {
	query := fmt.Sprintf("%s:field1,field12+%s:field2", DeploymentName, Category)
	expectedQuery := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
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
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: Category.String(), Value: "field2", Highlight: true},
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
