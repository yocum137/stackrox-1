package search

import (
	"testing"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stretchr/testify/assert"
)

func TestEmptyQuery(t *testing.T) {
	assert.Equal(t, &aux.Query{}, NewQueryBuilder().ProtoQuery())
}

func TestDocIDs(t *testing.T) {
	cases := []struct {
		desc   string
		docIDs []string
	}{
		{
			desc:   "no doc ids",
			docIDs: []string{},
		},
		{
			desc:   "one doc id",
			docIDs: []string{"1"},
		},
		{
			desc:   "two doc ids",
			docIDs: []string{"1", "2"},
		},
	}
	for _, c := range cases {
		q := NewQueryBuilder().AddDocIDs(c.docIDs...).ProtoQuery()
		expected := &aux.Query{
			Query: &aux.Query_BaseQuery{
				BaseQuery: &aux.BaseQuery{
					Query: &aux.BaseQuery_DocIdQuery{
						DocIdQuery: &aux.DocIDQuery{
							Ids: c.docIDs,
						},
					},
				},
			},
		}
		assert.Equal(t, expected, q)
	}
}
