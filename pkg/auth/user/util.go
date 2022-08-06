package user

import (
	"sort"

	"github.com/stackrox/rox/generated/storage"
)

// ConvertAttributes converts a map of user attributes to v1.UserAttribute
func ConvertAttributes(attrMap map[string][]string) []*storage.AuthStatus_UserAttribute {
	if attrMap == nil {
		return nil
	}

	result := make([]*storage.AuthStatus_UserAttribute, 0, len(attrMap))
	for k, vs := range attrMap {
		attr := &storage.AuthStatus_UserAttribute{
			Key:    k,
			Values: vs,
		}
		result = append(result, attr)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
}
