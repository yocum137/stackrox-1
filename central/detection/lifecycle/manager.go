package lifecycle

import (
	"time"

	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/detection/alertmanager"
	"github.com/stackrox/rox/central/detection/deploytime"
	"github.com/stackrox/rox/central/detection/runtime"
	"github.com/stackrox/rox/central/enrichment"
	processDatastore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/sensorevent/service/pipeline"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"golang.org/x/time/rate"
)

const (
	rateLimitDuration = 20 * time.Second
	tickerDuration    = 1 * time.Minute
)

var (
	logger = logging.LoggerForModule()
)

// A Manager manages deployment/policy lifecycle updates.
type Manager interface {
	IndicatorAdded(indicator *storage.ProcessIndicator, injector pipeline.EnforcementInjector) error
	// DeploymentUpdated processes a new or updated deployment, generating and updating alerts in the store and returning
	// enforcement action.
	DeploymentUpdated(deployment *storage.Deployment) (string, storage.EnforcementAction, error)
	UpsertPolicy(policy *storage.Policy) error

	DeploymentRemoved(deployment *storage.Deployment) error
	RemovePolicy(policyID string) error
}

// NewManager returns a new manager with the injected dependencies.
func NewManager(enricher enrichment.Enricher, deploytimeDetector deploytime.Detector, runtimeDetector runtime.Detector,
	deploymentDatastore deploymentDatastore.DataStore, processesDataStore processDatastore.DataStore, alertManager alertmanager.AlertManager) Manager {
	m := &managerImpl{
		enricher:            enricher,
		deploytimeDetector:  deploytimeDetector,
		runtimeDetector:     runtimeDetector,
		alertManager:        alertManager,
		deploymentDataStore: deploymentDatastore,
		processesDataStore:  processesDataStore,

		queuedIndicators: make(map[string]indicatorWithInjector),

		limiter: rate.NewLimiter(rate.Every(rateLimitDuration), 5),
		ticker:  time.NewTicker(tickerDuration),
	}
	go m.flushQueuePeriodically()
	return m
}
