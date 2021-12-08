package signature

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
)

// VerifierFactory creates a factory that can instantiate verifiers for Deployments.
type VerifierFactory interface {
	DeploymentVerifier(deployment *storage.Deployment) (DeploymentVerifier, error)
}

// MockVerifierFactory creates a mock VerifierFactory for testing purposes
func MockVerifierFactory(keys map[string]string, fail bool, validate bool) VerifierFactory {
	return &mockVerifierFactory{d: &mockDeploymentVerifier{
		keysMap:  keys,
		fail:     fail,
		validate: validate,
	}}
}

type mockVerifierFactory struct {
	d DeploymentVerifier
}

func (m *mockVerifierFactory) DeploymentVerifier(deployment *storage.Deployment) (DeploymentVerifier, error) {
	return m.d, nil
}

// VerifierFactoryOption provides option to inject values within the factory.
type VerifierFactoryOption func(factory *verifierFactory)

// WithCustomKeyChain provides a VerifierFactoryOption to inject a custom authn.Keychain to use when authenticating with
// OCI compatible container registries. By default, the k8schain will be used.
// NOTE: This is especially useful during testing to use i.e. the authn.DefaultKeychain instead.
func WithCustomKeyChain(kc authn.Keychain) VerifierFactoryOption {
	return func(factory *verifierFactory) {
		factory.kc = kc
	}
}

// WithBase64EncodedPublicKeys provides a VerifierFactoryOption to inject a list of base64 encoded public keys. The
// signature will be verified against these keys.
func WithBase64EncodedPublicKeys(keys []string) VerifierFactoryOption {
	return func(factory *verifierFactory) {
		factory.base64EncKeys = sanitizeBase64Keys(keys)
	}
}

// NewVerifierFactory creates a factory capable of creating a verifier.
func NewVerifierFactory(opts ...VerifierFactoryOption) VerifierFactory {
	f := &verifierFactory{}
	for _, option := range opts {
		option(f)
	}
	return f
}

type verifierFactory struct {
	base64EncKeys []string
	kc            authn.Keychain
}

// DeploymentVerifier will create a verifier that handles verification of all images within a storage.Deployment.
func (v *verifierFactory) DeploymentVerifier(deployment *storage.Deployment) (DeploymentVerifier, error) {
	if v.kc == nil {
		kc, err := createK8SChain(deployment)
		if err != nil {
			return nil, err
		}
		v.kc = kc
	}

	iv, err := DefaultImageVerifier(v.base64EncKeys, v.kc)
	if err != nil {
		return nil, err
	}

	return &deploymentVerifier{iv: iv, deployment: deployment}, nil
}

// createK8SChain will create a authn.Keychain that will emulate the authentication mechanism the kublet uses.
// It consumes the namespace, service account and pull secrets associated with the storage.Deployment.
// The creation will fail, if ANY of the specified pull secrets is not found or any other unexpected error occured.
// When called in a non-kubernetes context, the function will fail at it relies on access to the kubernetes API.
func createK8SChain(deployment *storage.Deployment) (authn.Keychain, error) {
	k8c, err := k8schain.NewInCluster(context.Background(), k8schain.Options{
		Namespace:          deployment.GetNamespace(),
		ServiceAccountName: deployment.GetServiceAccount(),
		ImagePullSecrets:   deployment.GetImagePullSecrets(),
	})

	if err != nil {
		return nil, errors.Wrap(err, "creating k8schain for retrieving images")
	}
	return k8c, nil
}

// sanitizeBase64Keys sanitizes the input of the base64 keys. Since these are coming from the policy itself and the
// values are stored within a map, they will have the following structure:
//	<verifier type>=<base64 enc pub key>
// This function will take care of splitting the verifier type and base64 enc public key IF it exists.
func sanitizeBase64Keys(inputs []string) []string {
	base64EncKeys := make([]string, 0, len(inputs))
	for _, input := range inputs {
		base64EncKeys = append(base64EncKeys, strings.TrimPrefix(input, "PUBLIC KEY="))
	}
	return base64EncKeys
}
