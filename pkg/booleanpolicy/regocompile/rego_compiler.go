package regocompile

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/utils"
)

type regoCompilerForType struct {
	fieldToMetaPathMap *pathutil.FieldToMetaPathMap
}

type RegoCompilerForType interface {
	CompileRegoBasedEvaluator(query *query.Query) (evaluator.Evaluator, error)
}

type regoBasedEvaluator struct {
	q rego.PreparedEvalQuery
}

func convertBindingToResult(binding interface{}) (m map[string][]string, err error) {
	panicked := true
	defer func() {
		if r := recover(); r != nil || panicked {
			err = fmt.Errorf("panic running evaluator: %v", r)
		}
	}()
	m = make(map[string][]string)
	for k, v := range binding.(map[string]interface{}) {
		vAsInterfaceSlice := v.([]interface{})
		vAsString := make([]string, 0, len(vAsInterfaceSlice))
		for _, val := range vAsInterfaceSlice {
			vAsString = append(vAsString, fmt.Sprintf("%s", val))
		}
		m[k] = vAsString
	}
	panicked = false
	return m, nil
}

func (r *regoBasedEvaluator) Evaluate(obj pathutil.AugmentedValue) (*evaluator.Result, bool) {
	inMemVal, err := obj.GetFullValue()
	// If there is an error here, it is a programming error. Let's not panic in prod over it.
	if err != nil {
		utils.Should(err)
		return nil, false
	}
	resultSet, err := r.q.Eval(context.Background(), rego.EvalInput(inMemVal))
	// If there is an error here, it is a programming error. Let's not panic in prod over it.
	if err != nil {
		utils.Should(err)
		return nil, false
	}
	if len(resultSet) != 1 {
		utils.Should(fmt.Errorf("invalid resultSet: %+v", resultSet))
		return nil, false
	}
	result := resultSet[0]
	outBindings, found := result.Bindings["out"].([]interface{})
	if !found {
		utils.Should(errors.New("resultSet didn't contain the expected bindings"))
		return nil, false
	}

	// This means it didn't match.
	if len(outBindings) == 0 {
		return nil, false
	}

	res := &evaluator.Result{}
	// Our queries are constructed so that each binding will be a map[string][]interface{}.
	// rego, however will store this as a map[string]interface{}, with each value being an []interface{}
	for _, binding := range outBindings {
		match, err := convertBindingToResult(binding)
		if err != nil {
			utils.Should(fmt.Errorf("failed to convert binding %+v: %w", binding, err))
			return nil, false
		}
		res.Matches = append(res.Matches, match)
	}
	return res, true
}

func MustCreateRegoCompiler(objMeta *pathutil.AugmentedObjMeta) RegoCompilerForType {
	r, err := CreateRegoCompiler(objMeta)
	utils.Must(err)
	return r
}

func CreateRegoCompiler(objMeta *pathutil.AugmentedObjMeta) (RegoCompilerForType, error) {
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
	q, err := rego.New(
		rego.Query("out = data.policy.main.violations"),
		rego.Module("main.policy", regoModule),
	).PrepareForEval(context.Background())
	return &regoBasedEvaluator{q: q}, nil
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
