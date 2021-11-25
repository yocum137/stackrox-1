package signature

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

type VerifierFactory interface {
	DeploymentVerifier(deployment *storage.Deployment) (DeploymentVerifier, error)
}

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

type VerifierFactoryOption func(factory *verifierFactory)

func WithCustomKeyChain(kc authn.Keychain) VerifierFactoryOption {
	return func(factory *verifierFactory) {
		factory.kc = kc
	}
}

func WithBase64EncodedKeys(keys []string) VerifierFactoryOption {
	return func(factory *verifierFactory) {
		factory.base64EncKeys = sanitizeBase64Keys(keys)
	}
}
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

func (v *verifierFactory) DeploymentVerifier(deployment *storage.Deployment) (DeploymentVerifier, error) {
	if v.kc == nil {
		kc, err := createK8Chain(deployment)
		if err != nil {
			log.Errorf("Error creating k8chain: %v", err)
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

func createK8Chain(deployment *storage.Deployment) (authn.Keychain, error) {
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
//	<base64 enc pub key>=<verifier type>
// This function will take care of splitting the verifier type and base64 enc public key.
func sanitizeBase64Keys(inputs []string) []string {
	base64EncKeys := make([]string, 0, len(inputs))
	for _, input := range inputs {
		base64EncKeys = append(base64EncKeys, strings.TrimPrefix(input, "PUBLIC KEY="))
	}
	return base64EncKeys
}
