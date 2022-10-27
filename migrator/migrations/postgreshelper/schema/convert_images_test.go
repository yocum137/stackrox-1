// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestImageSerialization(t *testing.T) {
	obj := &storage.Image{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertImageFromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertImageToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}
