package signature

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(dhaus): Make this a little bit easier to test / worthwhile to test.
// i.e. PUBLIC docker IMAGE that has been signed + the associated PUBLIC key.
func TestVerifierFactory_DeploymentVerifier(t *testing.T) {
	dep := &storage.Deployment{
		Name:                  "testing",
		Namespace:             "default",
		OrchestratorComponent: false,
		Containers: []*storage.Container{
			{
				Id: "test",
				Image: &storage.ContainerImage{
					Id: "testing",
					Name: &storage.ImageName{
						Registry: "ttl.sh",
						Remote:   "408ab9c1-2a84-4ac0-acb5-f969ced03f9e",
						Tag:      "72h",
						FullName: "ttl.sh/408ab9c1-2a84-4ac0-acb5-f969ced03f9e:72h",
					},
				},
			},
		},
		ServiceAccount: "default",
	}
	publicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFSVA2cm53MWRTS2pFQzRaOVhsQ2ZQN0lkNkNRVAp3aWpxZTZvNlpSZ25KeFV1Y3JIODIyQmdUYS9mWGpreE5ZUDR1ODBmQVFrNUFVZVhMc3kyeHRndUFBPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

	factory := NewVerifierFactory(WithBase64EncodedKeys([]string{publicKey + "=PUBLIC KEY"}), WithCustomKeyChain(authn.DefaultKeychain))
	verifier, err := factory.DeploymentVerifier(dep)
	require.NoError(t, err)
	res, verified, err := verifier.VerifyImages()
	require.NoError(t, err)
	assert.True(t, verified)
	assert.Equal(t, map[string]string{publicKey: "PUBLIC KEY"}, res)
}
