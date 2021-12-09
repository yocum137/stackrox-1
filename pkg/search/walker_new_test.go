package search

import (
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	newOptions := Walk(v1.SearchCategory_DEPLOYMENTS, "deployment", (*storage.Deployment)(nil)).Original()
	legacyOptions := walkLegacy(v1.SearchCategory_DEPLOYMENTS, "deployment", (*storage.Deployment)(nil)).Original()

	assert.Equal(t, len(newOptions), len(legacyOptions))
	for k, v := range newOptions {
		t.Run(k.String(), func(t *testing.T) {
			assert.Equal(t, v, legacyOptions[k])
		})
	}
}
