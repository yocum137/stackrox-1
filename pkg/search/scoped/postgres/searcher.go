package postgres

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stackrox/rox/pkg/utils"
)

// WithScoping allows the input searcher to be scoped.
func WithScoping(searcher search.Searcher) search.Searcher {
	return search.FuncSearcher{
		SearchFunc: func(ctx context.Context, q *aux.Query) ([]search.Result, error) {
			scopes, hasScope := scoped.GetAllScopes(ctx)
			if hasScope {
				var err error
				q, err = scopeQuery(q, scopes)
				if err != nil || q == nil {
					return nil, err
				}
			}
			return searcher.Search(ctx, q)
		},
		CountFunc: func(ctx context.Context, q *aux.Query) (int, error) {
			scopes, hasScope := scoped.GetAllScopes(ctx)
			if hasScope {
				var err error
				q, err = scopeQuery(q, scopes)
				if err != nil || q == nil {
					return 0, err
				}
			}
			return searcher.Count(ctx, q)
		},
	}
}

func scopeQuery(q *aux.Query, scopes []scoped.Scope) (*aux.Query, error) {
	pagination := q.GetPagination()
	q.Pagination = nil
	conjuncts := []*aux.Query{q}
	for _, scope := range scopes {
		schema := mapping.GetTableFromCategory(scope.Level)
		if schema == nil {
			utils.Should(errors.Errorf("no schema registered for search category %s", scope.Level))
			return q, nil
		}
		idField := schema.ID()
		conjuncts = append(conjuncts, search.NewQueryBuilder().AddExactMatches(search.FieldLabel(idField.Search.FieldName), scope.ID).ProtoQuery())
	}
	ret := search.ConjunctionQuery(conjuncts...)
	ret.Pagination = pagination
	return ret, nil
}
