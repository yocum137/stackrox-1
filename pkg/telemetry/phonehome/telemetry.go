package phonehome

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authn/basic"
)

// HashUserID anonymizes user ID so that it can be sent to the external
// telemetry storage for product data analysis.
func (cfg *Config) HashUserID(id authn.Identity) string {
	if id == nil {
		return "local:" + cfg.ClientID + ":unauthenticated"
	}
	if basic.IsBasicIdentity(id) {
		return "local:" + cfg.ClientID + ":" + id.FullName()
	}
	sha := sha256.New()
	_, _ = sha.Write([]byte(id.UID()))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}
