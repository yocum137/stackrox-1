package requestinfo

import (
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestRefererStrings(t *testing.T) {
	// The direct dependency on generated/api/v1 of the non-test package needed to
	// be removed to prevent issues around upgrading scanner code. Make sure the constants
	// introduced for that purpose have the same values as the previously used proto enum names.
	assert.Equal(t, v1.Audit_UI.String(), refererUI)
	assert.Equal(t, v1.Audit_API.String(), refererAPI)
}
