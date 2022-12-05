package service

import (
	"context"
	"encoding/base64"

	"github.com/golang/protobuf/proto"
	cTLS "github.com/google/certificate-transparency-go/tls"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/tlsconfig"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/cryptoutils"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz/allow"
	"github.com/stackrox/rox/pkg/mtls"
	pkgPH "github.com/stackrox/rox/pkg/telemetry/phonehome"
	"github.com/stackrox/rox/pkg/version"
	"google.golang.org/grpc"
)

// Service is the struct that manages the Metadata API
type serviceImpl struct {
	v1.UnimplementedMetadataServiceServer
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterMetadataServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterMetadataServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, allow.Anonymous().Authorized(ctx, fullMethodName)
}

// GetMetadata returns the metadata for Rox.
func (s *serviceImpl) GetMetadata(ctx context.Context, _ *v1.Empty) (*v1.Metadata, error) {
	metadata := &v1.Metadata{
		BuildFlavor:   buildinfo.BuildFlavor,
		ReleaseBuild:  buildinfo.ReleaseBuild,
		LicenseStatus: v1.Metadata_VALID,
	}
	id, _ := authn.IdentityFromContext(ctx)
	if pkgPH.Enabled() {
		metadata.StorageKeyV1 = env.TelemetryStorageKey.Setting()
		metadata.TelemetryEndpoint = env.TelemetryEndpoint.Setting()
		metadata.UserId = pkgPH.HashUserID(id)
	}
	// Only return the version to logged in users, not anonymous users.
	if id != nil {
		metadata.Version = version.GetMainVersion()
	}
	return metadata, nil
}

// TLSChallenge returns all trusted CAs (i.e. secret/additional-ca) and centrals cert chain. This is necessary if
// central is running behind load balancer with self-signed certificates.
//
// To validate that the list of trust roots comes directly from central and have not been tampered with,
// Central will cryptographically sign it with the private key of its service certificate.
//
// 1. External challenge token, generated by the external service
// 2. Central challenge token, generated by central itself
// 3. Payload (i.e. certificates)
func (s *serviceImpl) TLSChallenge(ctx context.Context, req *v1.TLSChallengeRequest) (*v1.TLSChallengeResponse, error) {
	sensorChallenge := req.GetChallengeToken()
	sensorChallengeBytes, err := base64.URLEncoding.DecodeString(sensorChallenge)
	if err != nil {
		return nil, errors.Wrapf(errox.InvalidArgs, "challenge token must be a valid base64 string: %s", err)
	}
	if len(sensorChallengeBytes) != centralsensor.ChallengeTokenLength {
		return nil, errors.Wrapf(errox.InvalidArgs, "base64 decoded challenge token must be %d bytes long, got %s", centralsensor.ChallengeTokenLength, sensorChallenge)
	}

	// Create central challenge token
	nonceGenerator := cryptoutils.NewNonceGenerator(centralsensor.ChallengeTokenLength, nil)
	centralChallenge, err := nonceGenerator.Nonce()
	if err != nil {
		return nil, errors.Errorf("Could not create central challenge: %s", err)
	}

	_, caCertDERBytes, err := mtls.CACert()
	if err != nil {
		return nil, errors.Errorf("Could not read CA cert and private key: %s", err)
	}

	leafCert, err := mtls.LeafCertificateFromFile()
	if err != nil {
		return nil, errors.Errorf("Could not load leaf certificate: %s", err)
	}

	additionalCAs, err := tlsconfig.GetAdditionalCAs()
	if err != nil {
		return nil, errors.Errorf("reading additional CAs: %s", err)
	}

	// add default leaf cert to additional CAs
	defaultCertChain, err := tlsconfig.GetDefaultCertChain()
	if err != nil {
		return nil, errors.Errorf("could not read default CA cert: %s", err)
	}
	if len(defaultCertChain) > 0 {
		additionalCAs = append(additionalCAs, defaultCertChain[0])
	}

	// Write trust info to proto struct
	trustInfo := &v1.TrustInfo{
		CentralChallenge: centralChallenge,
		SensorChallenge:  sensorChallenge,
		CertChain: [][]byte{
			leafCert.Certificate[0],
			caCertDERBytes,
		},
		AdditionalCas: additionalCAs,
	}
	trustInfoBytes, err := proto.Marshal(trustInfo)
	if err != nil {
		return nil, errors.Errorf("Could not marshal trust info: %s", err)
	}

	// Create signature from CA key
	sign, err := cTLS.CreateSignature(cryptoutils.DerefPrivateKey(leafCert.PrivateKey), cTLS.SHA256, trustInfoBytes)
	if err != nil {
		return nil, errors.Errorf("Could not sign trust info: %s", err)
	}

	resp := &v1.TLSChallengeResponse{
		Signature:           sign.Signature,
		TrustInfoSerialized: trustInfoBytes,
	}

	return resp, nil
}
