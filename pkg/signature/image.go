package signature

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/stackrox/rox/generated/storage"
)

// ImageVerifier takes an image and verifies whether the image is signed against a set of public keys.
// It returns a result which includes whether the signature could be verified and optionally if the
// verification was successful the public key that verified the signature as a base64 encoded string.
type ImageVerifier interface {
	VerifySignature(image *storage.Image) VerificationResult
}

// VerificationResult will be returned by all SignatureVerifier's containing information
// about whether the image signature was verified. If it was verified, a list of public keys which matched
// the signature is returned.
type VerificationResult struct {
	VerifiedKey string
	Verified    bool
	Err         error
}

func DefaultImageVerifier(base64EncPubKeys []string, kc authn.Keychain) (ImageVerifier, error) {
	return newPublicKeyVerifier(base64EncPubKeys, kc)
}
