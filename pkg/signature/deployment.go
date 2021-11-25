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

func (d deploymentVerifier) VerifyImages() (map[string]string, bool, error) {
	log.Info("Starting validation")
	verifiedBase64EncKeysSet := set.NewStringSet()
	log.Infof("Verifier: %+v", d.iv)
	for _, container := range d.deployment.GetContainers() {
		res := d.iv.VerifySignature(types.ToImage(container.GetImage()))

		if res.Err != nil {
			log.Infof("Encountered an error while trying to verify image within deployment %q: %v",
				d.deployment.GetId(), res.Err)
			return nil, false, res.Err
		}

		if !res.Verified {
			log.Infof("Failed to verify image within deployment %q", d.deployment.GetId())
			return nil, false, nil
		}

		verifiedBase64EncKeysSet.Add(res.VerifiedKey)
	}
	return createResultMap(verifiedBase64EncKeysSet.AsSlice()), true, nil
}

func createResultMap(keys []string) map[string]string {
	resultMap := make(map[string]string, len(keys))
	for _, key := range keys {
		// TODO(dhaus): This should have the type of the verifier used.
		resultMap["PUBLIC KEY"] = key
	}
	return resultMap
}
