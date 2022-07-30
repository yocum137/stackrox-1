package search

import (
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
)

func TestFilterQuery(t *testing.T) {
	optionsMap := Walk(v1.SearchCategory_IMAGES, "derp", &storage.Image{})

	query := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}

	newQuery, filtered := FilterQueryWithMap(query, optionsMap)
	assert.True(t, filtered)
	assert.Equal(t, &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
				},
			},
		},
	}, newQuery)

	var expected *auxpb.Query
	newQuery, filtered = FilterQueryWithMap(EmptyQuery(), optionsMap)
	assert.False(t, filtered)
	assert.Equal(t, expected, newQuery)

	q := &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{
						Field: ImageSHA.String(),
						Value: "blah",
					},
				},
			},
		},
	}
	newQuery, filtered = FilterQueryWithMap(q, optionsMap)
	assert.False(t, filtered)
	assert.Equal(t, q, newQuery)
}

func TestInverseFilterQuery(t *testing.T) {
	optionsMap := Walk(v1.SearchCategory_IMAGES, "derp", &storage.Image{})

	query := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}

	newQuery, filtered := InverseFilterQueryWithMap(query, optionsMap)
	assert.True(t, filtered)
	assert.Equal(t, &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}, newQuery)

	var expected *auxpb.Query
	newQuery, filtered = InverseFilterQueryWithMap(EmptyQuery(), optionsMap)
	assert.False(t, filtered)
	assert.Equal(t, expected, newQuery)

	q := &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{
						Field: ImageSHA.String(),
						Value: "blah",
					},
				},
			},
		},
	}
	newQuery, filtered = InverseFilterQueryWithMap(q, optionsMap)
	assert.False(t, filtered)
	assert.Equal(t, expected, newQuery)
}

func TestAddAsConjunction(t *testing.T) {
	toAdd := &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
				},
			},
		},
	}

	addTo := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}

	expected := &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
						},
					},
				}},
			},
		}},
	}

	added, err := AddAsConjunction(toAdd, addTo)
	assert.NoError(t, err)
	assert.Equal(t, expected, added)

	addTo = &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{
			BaseQuery: &auxpb.BaseQuery{
				Query: &auxpb.BaseQuery_MatchFieldQuery{
					MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
				},
			},
		},
	}

	expected = &auxpb.Query{
		Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}

	added, err = AddAsConjunction(toAdd, addTo)
	assert.NoError(t, err)
	assert.Equal(t, expected, added)

	addTo = &auxpb.Query{
		Query: &auxpb.Query_Disjunction{Disjunction: &auxpb.DisjunctionQuery{
			Queries: []*auxpb.Query{
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: CVE.String(), Value: "cveId"},
						},
					},
				}},
				{Query: &auxpb.Query_BaseQuery{
					BaseQuery: &auxpb.BaseQuery{
						Query: &auxpb.BaseQuery_MatchFieldQuery{
							MatchFieldQuery: &auxpb.MatchFieldQuery{Field: DeploymentName.String(), Value: "depname"},
						},
					},
				}},
			},
		}},
	}

	_, err = AddAsConjunction(toAdd, addTo)
	assert.NoError(t, err)
}
