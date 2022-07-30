package search

import (
	"testing"

	"github.com/stackrox/rox/generated/auxpb"
	"github.com/stretchr/testify/assert"
)

func TestEmptyQuery(t *testing.T) {
	assert.Equal(t, &auxpb.Query{}, NewQueryBuilder().ProtoQuery())
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
		expected := &auxpb.Query{
			Query: &auxpb.Query_BaseQuery{
				BaseQuery: &auxpb.BaseQuery{
					Query: &auxpb.BaseQuery_DocIdQuery{
						DocIdQuery: &auxpb.DocIDQuery{
							Ids: c.docIDs,
						},
					},
				},
			},
		}
		assert.Equal(t, expected, q)
	}
}
