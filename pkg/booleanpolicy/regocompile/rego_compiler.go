package regocompile

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
)

type regoCompilerForType struct {
	fieldToMetaPathMap *pathutil.FieldToMetaPathMap
}

type RegoCompilerForType interface {
	CompileRegoBasedEvaluator(query *query.Query) (evaluator.Evaluator, error)
}

type regoBasedEvaluator struct {
}

func (r *regoBasedEvaluator) Evaluate(obj pathutil.AugmentedValue) (*evaluator.Result, bool) {
}

func NewRegoCompilerForType(objMeta *pathutil.AugmentedObjMeta) (RegoCompilerForType, error) {
	fieldToMetaPathMap, err := objMeta.MapSearchTagsToPaths()
	if err != nil {
		return nil, err
	}
	return &regoCompilerForType{fieldToMetaPathMap: fieldToMetaPathMap}, nil
}

func pathToKey(path []string) string {
	return strings.Join(path, ".")
}

func (r *regoCompilerForType) CompileRegoBasedEvaluator(query *query.Query) (evaluator.Evaluator, error) {
	regoModule, err := r.compileRego(query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile rego: %w", err)
	}
	regoObj, err := rego.New(
		rego.Query("out = data.policy.main.violations"),
		rego.Module("main.policy", regoModule),
	).PrepareForEval(context.Background())
}

func (r *regoCompilerForType) compileRego(query *query.Query) (string, error) {
	pathsToArrayIndexes := make(map[string]int)

	for _, fieldQuery := range query.FieldQueries {
		field := fieldQuery.Field
		metaPathToField, found := r.fieldToMetaPathMap.Get(field)
		if !found {
			return "", fmt.Errorf("field %v not in object", field)
		}
		var path []string
		for i, elem := range metaPathToField {
			path = append(path, elem.JSONTag)
			if i == len(metaPathToField)-1 {
				continue
			}
			if elem.Type.Kind() == reflect.Slice || elem.Type.Kind() == reflect.Array {
				pathKey := pathToKey(path)
				idx, ok := pathsToArrayIndexes[pathKey]
				if !ok {
					idx = len(pathsToArrayIndexes)
					pathsToArrayIndexes[pathKey] = idx
				}
			}
		}
	}

	return "", nil
}
