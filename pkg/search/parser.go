package search

import (
	"strings"

	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/set"
)

var (
	log = logging.LoggerForModule()
)

// ParseQueryForAutocomplete parses the input string specific for autocomplete requests.
func ParseQueryForAutocomplete(query string) (*auxpb.Query, string, error) {
	return autocompleteQueryParser{}.parse(query)
}

// ParseQuery parses the input query with the supplied options.
func ParseQuery(query string, opts ...ParseQueryOption) (*auxpb.Query, error) {
	parser := generalQueryParser{}
	for _, opt := range opts {
		opt(&parser)
	}
	return parser.parse(query)
}

// ParseQueryOption represents an option to use when parsing queries.
type ParseQueryOption func(parser *generalQueryParser)

// MatchAllIfEmpty will cause an empty query to be returned if the input query is empty (as opposed to an error).
func MatchAllIfEmpty() ParseQueryOption {
	return func(parser *generalQueryParser) {
		parser.MatchAllIfEmpty = true
	}
}

// ExcludeFieldLabel removes a specific options key from the search if it exists
func ExcludeFieldLabel(k FieldLabel) ParseQueryOption {
	return func(parser *generalQueryParser) {
		if parser.ExcludedFieldLabels == nil {
			parser.ExcludedFieldLabels = set.NewStringSet()
		}
		parser.ExcludedFieldLabels.Add(k.String())
	}
}

// FilterFields uses a predicate to filter our fields from a raw query based on the field key.
func FilterFields(query string, pred func(field string) bool) string {
	if query == "" {
		return query
	}
	pairs := splitQuery(query)
	pairsToKeep := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		key, _, valid := parsePair(pair, false)
		if !valid {
			continue
		}
		if !pred(key) {
			continue
		}
		pairsToKeep = append(pairsToKeep, pair)
	}
	return strings.Join(pairsToKeep, "+")
}

// Extracts "key", "value1,value2" from a string in the format key:value1,value2
func parsePair(pair string, allowEmpty bool) (key string, values string, valid bool) {
	pair = strings.TrimSpace(pair)
	if len(pair) == 0 {
		return
	}

	spl := strings.SplitN(pair, ":", 2)
	// len < 2 implies there isn't a colon and the second check verifies that the : wasn't the last char
	if len(spl) < 2 || (spl[1] == "" && !allowEmpty) {
		return
	}
	// If empty strings are allowed, it means we're treating them as wildcards.
	if allowEmpty {
		if spl[1] == "" {
			spl[1] = WildcardString
		} else if string(spl[1][len(spl[1])-1]) == "," {
			spl[1] = spl[1] + WildcardString
		}
	}
	return spl[0], spl[1], true
}

func queryFromFieldValues(field string, values []string, highlight bool) *auxpb.Query {
	queries := make([]*auxpb.Query, 0, len(values))
	for _, value := range values {
		queries = append(queries, MatchFieldQuery(field, value, highlight))
	}

	return DisjunctionQuery(queries...)
}

// DisjunctionQuery returns a disjunction query of the provided queries.
func DisjunctionQuery(queries ...*auxpb.Query) *auxpb.Query {
	return disjunctOrConjunctQueries(false, queries...)
}

// ConjunctionQuery returns a conjunction query of the provided queries.
func ConjunctionQuery(queries ...*auxpb.Query) *auxpb.Query {
	return disjunctOrConjunctQueries(true, queries...)
}

// Helper function that DisjunctionQuery and ConjunctionQuery proxy to.
// Do NOT call this directly.
func disjunctOrConjunctQueries(isConjunct bool, queries ...*auxpb.Query) *auxpb.Query {
	if len(queries) == 0 {
		return &auxpb.Query{}
	}

	if len(queries) == 1 {
		return queries[0]
	}
	if isConjunct {
		return &auxpb.Query{
			Query: &auxpb.Query_Conjunction{Conjunction: &auxpb.ConjunctionQuery{Queries: queries}},
		}
	}

	return &auxpb.Query{
		Query: &auxpb.Query_Disjunction{Disjunction: &auxpb.DisjunctionQuery{Queries: queries}},
	}
}

func queryFromBaseQuery(baseQuery *auxpb.BaseQuery) *auxpb.Query {
	return &auxpb.Query{
		Query: &auxpb.Query_BaseQuery{BaseQuery: baseQuery},
	}
}

// MatchFieldQuery returns a match field query.
// It's a simple convenience wrapper around initializing the struct.
func MatchFieldQuery(field, value string, highlight bool) *auxpb.Query {
	return queryFromBaseQuery(&auxpb.BaseQuery{
		Query: &auxpb.BaseQuery_MatchFieldQuery{MatchFieldQuery: &auxpb.MatchFieldQuery{Field: field, Value: value, Highlight: highlight}},
	})
}

// matchLinkedFieldsQuery returns a query that matches
func matchLinkedFieldsQuery(fieldValues []fieldValue) *auxpb.Query {
	mfqs := make([]*auxpb.MatchFieldQuery, len(fieldValues))
	for i, fv := range fieldValues {
		mfqs[i] = &auxpb.MatchFieldQuery{Field: fv.l.String(), Value: fv.v, Highlight: fv.highlighted}
	}

	return queryFromBaseQuery(&auxpb.BaseQuery{
		Query: &auxpb.BaseQuery_MatchLinkedFieldsQuery{MatchLinkedFieldsQuery: &auxpb.MatchLinkedFieldsQuery{
			Query: mfqs,
		}},
	})
}

func docIDQuery(ids []string) *auxpb.Query {
	return queryFromBaseQuery(&auxpb.BaseQuery{
		Query: &auxpb.BaseQuery_DocIdQuery{DocIdQuery: &auxpb.DocIDQuery{Ids: ids}},
	})
}

//go:generate stringer -type=QueryModifier
// QueryModifier describes the query modifiers for a specific individual query
type QueryModifier int

// These are the currently supported modifiers
const (
	AtLeastOne QueryModifier = iota
	Negation
	Regex
	Equality
)

// GetValueAndModifiersFromString parses the raw value string into its value and modifiers
func GetValueAndModifiersFromString(value string) (string, []QueryModifier) {
	var queryModifiers []QueryModifier
	trimmedValue := value
	// We only allow at most one modifier from the set {atleastone, negation}.
	// Anything more, we treat as part of the string to query for.
	var negationOrAtLeastOneFound bool
forloop:
	for {
		switch {
		// AtLeastOnePrefix is !! so it must come before negation prefix
		case !negationOrAtLeastOneFound && strings.HasPrefix(trimmedValue, AtLeastOnePrefix) && len(trimmedValue) > len(AtLeastOnePrefix):
			trimmedValue = trimmedValue[len(AtLeastOnePrefix):]
			queryModifiers = append(queryModifiers, AtLeastOne)
			negationOrAtLeastOneFound = true
		case !negationOrAtLeastOneFound && strings.HasPrefix(trimmedValue, NegationPrefix) && len(trimmedValue) > len(NegationPrefix):
			trimmedValue = trimmedValue[len(NegationPrefix):]
			queryModifiers = append(queryModifiers, Negation)
			negationOrAtLeastOneFound = true
		case strings.HasPrefix(trimmedValue, RegexPrefix) && len(trimmedValue) > len(RegexPrefix):
			trimmedValue = strings.ToLower(trimmedValue[len(RegexPrefix):])
			queryModifiers = append(queryModifiers, Regex)
			break forloop // Once we see that it's a regex, we don't check for special-characters in the rest of the string.
		case strings.HasPrefix(trimmedValue, EqualityPrefixSuffix) && strings.HasSuffix(trimmedValue, EqualityPrefixSuffix) && len(trimmedValue) > 2*len(EqualityPrefixSuffix):
			trimmedValue = trimmedValue[len(EqualityPrefixSuffix) : len(trimmedValue)-len(EqualityPrefixSuffix)]
			queryModifiers = append(queryModifiers, Equality)
			break forloop // Once it's within quotes, we take the value inside as is, and don't try to extract modifiers.
		default:
			break forloop
		}
	}
	return trimmedValue, queryModifiers
}
