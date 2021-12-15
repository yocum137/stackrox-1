package regocompile

import (
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	v1 "k8s.io/api/apps/v1"
)

type RegoCompilerForType struct {
}

func CompileRego(query *query.Query, objMeta *pathutil.AugmentedObjMeta) (string, error) {
	fieldToMetaPathMap, err := objMeta.MapSearchTagsToPaths()
	if err != nil {
		return "", err
	}
	_ = fieldToMetaPathMap
	return "", err
}
