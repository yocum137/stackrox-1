// Package operator contains "operational logic" that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.

package operator

import (
	"context"

	"github.com/stackrox/rox/pkg/logging"
	"k8s.io/client-go/kubernetes"
)

var (
	log = logging.LoggerForModule()
)

type operatorImpl struct {
	k8sClient     kubernetes.Interface
}

// Operator performs some operational logic that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.
type Operator interface {
	Start(ctx context.Context) error
}

// New creates a new operator
func New(k8sClient kubernetes.Interface) Operator {
	return &operatorImpl{k8sClient: k8sClient}
}

// Start launches the processes that implement the "operational logic" of Sensor.
func (o *operatorImpl) Start(ctx context.Context) error {
	log.Info("Starting embedded operator.")

	log.Info("Embedded operator started.")

	return nil // FIXME
}