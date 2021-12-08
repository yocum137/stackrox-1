package signature

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	ggcrRemote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/oci/remote"
	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/stackrox/rox/generated/storage"
)

type publicKeyVerifier struct {
	base64EncPublicKeys []string
	parsedPublicKeys    []crypto.PublicKey
	kc                  authn.Keychain
}

func newPublicKeyVerifier(base64EncKeys []string, kc authn.Keychain) (*publicKeyVerifier, error) {
	parsedKeys := make([]crypto.PublicKey, 0, len(base64EncKeys))
	for _, base64EncKey := range base64EncKeys {
		decodedKey, err := base64.StdEncoding.DecodeString(base64EncKey)
		if err != nil {
			return nil, errors.Wrapf(err, "decoding base64 key %q", base64EncKey)
		}

		keyBlock, _ := pem.Decode(decodedKey)
		if keyBlock == nil || keyBlock.Type != "PUBLIC KEY" {
			return nil, fmt.Errorf("failed to decode PEM block containing public key: %q", base64EncKey)
		}

		parsedKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "parsing key as public key")
		}
		parsedKeys = append(parsedKeys, parsedKey)
	}

	return &publicKeyVerifier{
		base64EncPublicKeys: base64EncKeys,
		parsedPublicKeys:    parsedKeys,
		kc:                  kc,
	}, nil
}

// VerifySignature will verify whether an image is signed against a set of public keys.
// This function follows cosign approach and its workflow to validate public keys.
// 1. The image and repository references are created.
// 2. For each public key:
//		a. Try to download the signature file, inferred from the image / repo references.
//		b. Verify the signature against the public key.
// The image's signature will be marked as unverified when no matching public key is found, no signature
// is found or any error occurred.
func (p *publicKeyVerifier) VerifySignature(image *storage.Image) VerificationResult {
	imageFullName := image.GetName().GetFullName()
	ctx := context.Background()

	if image.GetNotPullable() {
		return VerificationResult{}
	}

	registryOpts := options.RegistryOptions{}
	remoteOpts, err := registryOpts.ClientOpts(ctx)
	if err != nil {
		return VerificationResult{Err: errors.Wrap(err, "initializing remote opts")}
	}

	checkOpts := &cosign.CheckOpts{
		RegistryClientOpts: remoteOpts,
		Annotations:        map[string]interface{}{},
	}

	// The repository reference expects the following:
	// 		<registry/remote>
	// Need to strip the tag value as well as the ":" from the image's full name.
	repoRef, err := name.NewRepository(strings.TrimSuffix(imageFullName, ":"+image.GetName().GetTag()))
	if err != nil {
		return VerificationResult{Err: errors.Wrapf(err,
			"getting repository reference for image %q", imageFullName)}
	}

	checkOpts.RegistryClientOpts = append(checkOpts.RegistryClientOpts, remote.WithTargetRepository(repoRef),
		remote.WithRemoteOptions(ggcrRemote.WithAuthFromKeychain(p.kc)))

	imageRef, err := name.ParseReference(imageFullName)
	if err != nil {
		return VerificationResult{Err: errors.Wrapf(err, "getting image reference for image %q", imageFullName)}
	}

	for pubKeyIndex, pubKey := range p.parsedPublicKeys {
		checkOpts.SigVerifier, err = signature.LoadVerifier(pubKey, crypto.SHA256)
		if err != nil {
			return VerificationResult{Err: errors.Wrap(err, "getting cosign verifier")}
		}

		// The signature is verified when no error is returned, since we do not additionally do any attestation with
		// cosign.
		_, _, err = cosign.Verify(ctx, imageRef, cosign.SignaturesAccessor, checkOpts)
		if err == nil {
			return VerificationResult{Verified: true, VerifiedKey: p.base64EncPublicKeys[pubKeyIndex]}
		}
	}

	return VerificationResult{}
}
