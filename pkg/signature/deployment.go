package signature

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/set"
)

// DeploymentVerifier takes all images of a storage.Deployment and verifies whether they are signed
// by against a set of public keys.
// It returns a boolean indicating whether the signature could be verified and the public key
// with which the verification was successful.
type DeploymentVerifier interface {
	VerifyImages() (map[string]string, bool, error)
}

type mockDeploymentVerifier struct {
	keysMap  map[string]string
	fail     bool
	validate bool
}

func (m *mockDeploymentVerifier) VerifyImages() (map[string]string, bool, error) {
	if m.fail {
		return nil, false, errors.New("sample failure")
	}

	if m.validate {
		return m.keysMap, true, nil
	}
	return nil, false, nil
}

type deploymentVerifier struct {
	iv         ImageVerifier
	deployment *storage.Deployment
}

// VerifyImages will use the injected storage.Deployment and the ImageVerifier to verify all images within a
// deployment and their signature. If the signature is successfully verified, it will return true as well as a
// map of a verifier type and the public key that successfully validated the signature.
// If the verification fails the function will return false and a nil map.
// If any error occurred during signature verification it will return an error and a false value.
func (d deploymentVerifier) VerifyImages() (map[string]string, bool, error) {
	verifiedBase64EncKeysSet := set.NewStringSet()
	for _, container := range d.deployment.GetContainers() {
		// TODO(dhaus): We can improve this here by handling all signature verification async.
		res := d.iv.VerifySignature(types.ToImage(container.GetImage()))

		if res.Err != nil {
			return nil, false, res.Err
		}

		if !res.Verified {
			return nil, false, nil
		}

		verifiedBase64EncKeysSet.Add(res.VerifiedKey)
	}
	return createResultMap(verifiedBase64EncKeysSet.AsSlice()), true, nil
}

func createResultMap(keys []string) map[string]string {
	resultMap := make(map[string]string, len(keys))
	for _, key := range keys {
		// TODO(dhaus): Needs chaning once its possible to have base64enc values as keys.
		resultMap["PUBLIC KEY"] = key
	}
	return resultMap
}
