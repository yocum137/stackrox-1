package search

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/pkg/utils"
)

// ApplyFnToAllBaseQueries walks recursively over the query, applying fn to all the base queries.
func ApplyFnToAllBaseQueries(q *auxpb.Query, fn func(*auxpb.BaseQuery)) {
	if q.GetQuery() == nil {
		return
	}

	switch typedQ := q.GetQuery().(type) {
	case *auxpb.Query_Disjunction:
		for _, subQ := range typedQ.Disjunction.GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
	case *auxpb.Query_Conjunction:
		for _, subQ := range typedQ.Conjunction.GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
	case *auxpb.Query_BooleanQuery:
		for _, subQ := range typedQ.BooleanQuery.GetMust().GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
		for _, subQ := range typedQ.BooleanQuery.GetMustNot().GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
	case *auxpb.Query_BaseQuery:
		fn(typedQ.BaseQuery)
	default:
		utils.Should(fmt.Errorf("unhandled query type: %T; query was %s", q, proto.MarshalTextString(q)))
	}
}

// FilterQueryWithMap removes match fields portions of the query that are not in the input options map.
func FilterQueryWithMap(q *auxpb.Query, optionsMap OptionsMap) (*auxpb.Query, bool) {
	var areFieldsFiltered bool
	filtered, _ := FilterQuery(q, func(bq *auxpb.BaseQuery) bool {
		matchFieldQuery, ok := bq.GetQuery().(*auxpb.BaseQuery_MatchFieldQuery)
		if ok {
			if _, isValid := optionsMap.Get(matchFieldQuery.MatchFieldQuery.GetField()); isValid {
				return true
			}
		}
		areFieldsFiltered = true
		return false
	})
	return filtered, areFieldsFiltered
}

// InverseFilterQueryWithMap removes match fields portions of the query that are in the input options map.
func InverseFilterQueryWithMap(q *auxpb.Query, optionsMap OptionsMap) (*auxpb.Query, bool) {
	var areFieldsFiltered bool
	filtered, _ := FilterQuery(q, func(bq *auxpb.BaseQuery) bool {
		matchFieldQuery, ok := bq.GetQuery().(*auxpb.BaseQuery_MatchFieldQuery)
		if ok {
			if _, isValid := optionsMap.Get(matchFieldQuery.MatchFieldQuery.GetField()); !isValid {
				areFieldsFiltered = true
				return true
			}
		}
		return false
	})
	return filtered, areFieldsFiltered
}

// AddAsConjunction adds the input toAdd query to the input addTo query at the top level, either by appending it to the
// conjunction list, or, if it is a base query, by making it a conjunction. Explicity disallows nested queries, as the
// resulting query is expected to be either a base query, or a flat query.
func AddAsConjunction(toAdd *auxpb.Query, addTo *auxpb.Query) (*auxpb.Query, error) {
	if addTo.Query == nil {
		return toAdd, nil
	}
	switch typedQ := addTo.GetQuery().(type) {
	case *auxpb.Query_Conjunction:
		typedQ.Conjunction.Queries = append(typedQ.Conjunction.Queries, toAdd)
		return addTo, nil
	case *auxpb.Query_BaseQuery, *auxpb.Query_Disjunction:
		return ConjunctionQuery(toAdd, addTo), nil
	default:
		return nil, errors.New("cannot add to a non-nil, non-conjunction/disjunction, non-base query")
	}
}

// FilterQuery applies the given function on every base query, and returns a new
// query that has only the sub-queries that the function returns true for.
// It will NOT mutate q unless the function passed mutates its argument.
func FilterQuery(q *auxpb.Query, fn func(*auxpb.BaseQuery) bool) (*auxpb.Query, bool) {
	if q.GetQuery() == nil {
		return nil, false
	}
	switch typedQ := q.GetQuery().(type) {
	case *auxpb.Query_Disjunction:
		filteredQueries := filterQueriesByFunction(typedQ.Disjunction.GetQueries(), fn)
		if len(filteredQueries) == 0 {
			return nil, false
		}
		return DisjunctionQuery(filteredQueries...), true
	case *auxpb.Query_Conjunction:
		filteredQueries := filterQueriesByFunction(typedQ.Conjunction.GetQueries(), fn)
		if len(filteredQueries) == 0 {
			return nil, false
		}
		return ConjunctionQuery(filteredQueries...), true
	case *auxpb.Query_BaseQuery:
		if fn(typedQ.BaseQuery) {
			return q, true
		}
		return nil, false
	default:
		log.Errorf("Unhandled query type: %T; query was %s", q, proto.MarshalTextString(q))
		return nil, false
	}
}

// Helper function used by FilterQuery.
func filterQueriesByFunction(qs []*auxpb.Query, fn func(*auxpb.BaseQuery) bool) (filteredQueries []*auxpb.Query) {
	for _, q := range qs {
		filteredQuery, found := FilterQuery(q, fn)
		if found {
			filteredQueries = append(filteredQueries, filteredQuery)
		}
	}
	return
}

// AddRawQueriesAsConjunction adds the input toAdd raw query to the input addTo raw query
func AddRawQueriesAsConjunction(toAdd string, addTo string) string {
	if toAdd == "" && addTo == "" {
		return ""
	}

	if addTo == "" {
		return toAdd
	}

	if toAdd == "" {
		return addTo
	}

	return addTo + "+" + toAdd
}
