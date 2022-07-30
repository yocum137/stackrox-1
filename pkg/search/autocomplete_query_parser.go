package search

import (
	"errors"
	"strings"

	"github.com/stackrox/rox/generated/auxpb"
)

// autocompleteQueryParser provides an autocomplete specific query parser.
type autocompleteQueryParser struct{}

// parse parses the input query.
func (pi autocompleteQueryParser) parse(input string) (*auxpb.Query, string, error) {
	// Handle empty input query case.
	if len(input) == 0 {
		return nil, "", errors.New("parser not configured to handle empty queries")
	}
	// Have a filled query, parse it.
	return pi.parseInternal(input)
}

func (pi autocompleteQueryParser) parseInternal(query string) (*auxpb.Query, string, error) {
	pairs := strings.Split(query, "+")

	queries := make([]*auxpb.Query, 0, len(pairs))
	var autocompleteKey string
	for i, pair := range pairs {
		key, commaSeparatedValues, valid := parsePair(pair, true)
		if !valid {
			continue
		}
		if i == len(pairs)-1 {
			queries = append(queries, queryFromFieldValues(key, strings.Split(commaSeparatedValues, ","), true))
			autocompleteKey = key
		} else {
			queries = append(queries, queryFromFieldValues(key, strings.Split(commaSeparatedValues, ","), false))
		}
	}

	// We always want to return an error here, because it means that the query is ill-defined.
	if len(queries) == 0 {
		return nil, "", errors.New("after parsing, query is empty")
	}

	return ConjunctionQuery(queries...), autocompleteKey, nil
}
