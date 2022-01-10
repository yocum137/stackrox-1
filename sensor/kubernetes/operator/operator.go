// Package operator contains "operational logic" that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.

package operator

import (
	"context"

	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/kubernetes/sensor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"k8s.io/client-go/kubernetes"

	"github.com/stackrox/rox/generated/internalapi/central"
	grpcUtil "github.com/stackrox/rox/pkg/grpc/util"
	"github.com/stackrox/rox/sensor/common/clusterid"
)

var (
	log = logging.LoggerForModule()
	grpcCallOptions = []grpc.CallOption{grpc.UseCompressor(gzip.Name)}
)

type operatorImpl struct {
	k8sClient     kubernetes.Interface
	centralConnection *grpcUtil.LazyClientConn
	ctx context.Context
	localScannerServiceClient central.LocalScannerServiceClient
}

// Operator performs some operational logic that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.
type Operator interface {
	Start(ctx context.Context) error
}

// New creates a new operator
func New(k8sClient kubernetes.Interface, centralConnection *grpcUtil.LazyClientConn) Operator {
	return &operatorImpl{
		k8sClient: k8sClient,
		// Operations on this connection will block until centralConnection.Set is called
		centralConnection: centralConnection,
	}
}

// Start launches the processes that implement the "operational logic" of Sensor.
func (o *operatorImpl) Start(ctx context.Context) error {
	log.Info("Starting embedded operator.")

	o.ctx = ctx
	// TODO log
	o.localScannerServiceClient = central.NewLocalScannerServiceClient(o.centralConnection)

	log.Info("Embedded operator started.")

	// FIXME
	_, err := o.issueScannerCertificates()
	if err != nil {
		return err
	}

	return nil // FIXME
}

func (o *operatorImpl) issueScannerCertificates() (*central.IssueLocalScannerCertsResponse, error) {
	// Blocks until CentralCommunication sets the cluster id after receiving it from centralHello.
	clusterID := clusterid.Get()
	// Current requirements expect sensor and local scanner to run on the same namespace
	localScannerNamespace := sensor.GetSensorNamespace()
	request := central.IssueLocalScannerCertsRequest{
		ClusterId: clusterID,
		Namespace: localScannerNamespace,
	}

	// FIXME: timeout with context
	return o.localScannerServiceClient.IssueLocalScannerCerts(o.ctx, &request, grpcCallOptions...)
}
