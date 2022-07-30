package compound

import (
	"github.com/stackrox/rox/generated/auxpb"
)

type searchRequestSpec struct {
	or                     []*searchRequestSpec
	and                    []*searchRequestSpec
	boolean                *booleanRequestSpec
	leftJoinWithRightOrder *joinRequestSpec
	base                   *baseRequestSpec
}

type booleanRequestSpec struct {
	must    *searchRequestSpec
	mustNot *searchRequestSpec
}

type joinRequestSpec struct {
	left  *searchRequestSpec
	right *searchRequestSpec
}

type baseRequestSpec struct {
	Spec  *SearcherSpec
	Query *auxpb.Query
}
