package declarativeconfig

import (
	"context"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/gogo/protobuf/proto"
	authProviderDatastore "github.com/stackrox/rox/central/authprovider/datastore"
	groupDatastore "github.com/stackrox/rox/central/group/datastore"
	roleDatastore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/declarativeconfig"
	"github.com/stackrox/rox/pkg/declarativeconfig/transform"
	"github.com/stackrox/rox/pkg/k8scfgwatch"
	"github.com/stackrox/rox/pkg/maputil"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	declarativeConfigDir = "/run/stackrox.io/declarative-configuration"
)

type protoMessagesByType = map[reflect.Type][]proto.Message

type ReconciliationErrorReporter interface {
	ProcessError(protoValue proto.Message, err error)
}

type managerImpl struct {
	once sync.Once

	universalTransformer transform.Transformer

	transformedMessagesByHandler map[string]protoMessagesByType
	transformedMessagesMutex     sync.RWMutex

	reconciliationTickerDuration time.Duration
	watchIntervalDuration        time.Duration

	reconciliationTicker *time.Ticker
	shortCircuitSignal   concurrency.Signal

	roleDS                      roleDatastore.DataStore
	groupDS                     groupDatastore.DataStore
	authProviderDS              authproviders.Store
	authProviderRegistry        authproviders.Registry
	reconciliationCtx           context.Context
	reconciliationErrorReporter ReconciliationErrorReporter
}

var (
	authProviderType  = reflect.TypeOf((*storage.AuthProvider)(nil))
	accessScopeType   = reflect.TypeOf((*storage.SimpleAccessScope)(nil))
	groupType         = reflect.TypeOf((*storage.Group)(nil))
	permissionSetType = reflect.TypeOf((*storage.PermissionSet)(nil))
	roleType          = reflect.TypeOf((*storage.Role)(nil))
)

// New creates a new instance of Manager.
// Note that it will not watch the declarative configuration directories when created, only after
// ReconcileDeclarativeConfigurations has been called.
func New(reconciliationTickerDuration, watchIntervalDuration time.Duration, roleDS roleDatastore.DataStore, groupDS groupDatastore.DataStore, authProviderDS authproviders.Store, registry authproviders.Registry, reconciliationErrorReporter ReconciliationErrorReporter) Manager {
	writeDeclarativeRoleCtx := declarativeconfig.WithModifyDeclarativeResource(context.Background())
	writeDeclarativeRoleCtx = sac.WithGlobalAccessScopeChecker(writeDeclarativeRoleCtx,
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			// TODO: ROX-14398 Replace Role with Access
			sac.ResourceScopeKeys(resources.Role, resources.Access)))
	return &managerImpl{
		universalTransformer:         transform.New(),
		transformedMessagesByHandler: map[string]protoMessagesByType{},
		reconciliationTickerDuration: reconciliationTickerDuration,
		watchIntervalDuration:        watchIntervalDuration,
		roleDS:                       roleDS,
		groupDS:                      groupDS,
		authProviderDS:               authProviderDS,
		reconciliationCtx:            writeDeclarativeRoleCtx,
		reconciliationErrorReporter:  reconciliationErrorReporter,
		authProviderRegistry:         registry,
	}
}

func (m *managerImpl) ReconcileDeclarativeConfigurations() {
	m.once.Do(func() {
		// For each directory within the declarative configuration path, create a watch handler.
		// The reason we need multiple watch handlers and cannot simply watch the root directory is that
		// changes to directories are ignored within the watch handler.
		entries, err := os.ReadDir(declarativeConfigDir)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("Declarative configuration directory does not exist, no reconciliation will be done")
				return
			}
			utils.Should(err)
			return
		}

		var startedWatchHandler bool
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			log.Infof("Start watch handler for declarative configuration for path %s",
				path.Join(declarativeConfigDir, entry.Name()))
			wh := newWatchHandler(entry.Name(), m)
			// Set Force to true, so we explicitly retry watching the files within the directory and not stop on the first
			// error occurred.
			watchOpts := k8scfgwatch.Options{Interval: m.watchIntervalDuration, Force: true}
			_ = k8scfgwatch.WatchConfigMountDir(context.Background(), declarativeConfigDir,
				k8scfgwatch.DeduplicateWatchErrors(wh), watchOpts)
			startedWatchHandler = true
		}

		// Only start the reconciliation loop if at least one watch handler has been started.
		if startedWatchHandler {
			log.Info("Start the reconciliation loop for declarative configurations")
			m.startReconciliationLoop()
		}
	})
}

// UpdateDeclarativeConfigContents will take the file contents and transform these to declarative configurations.
func (m *managerImpl) UpdateDeclarativeConfigContents(handlerID string, contents [][]byte) {
	configurations, err := declarativeconfig.ConfigurationFromRawBytes(contents...)
	if err != nil {
		log.Errorf("Error during unmarshalling of declarative configuration files: %+v", err)
		return
	}
	transformedConfigurations := make(map[reflect.Type][]proto.Message, len(configurations))
	for _, configuration := range configurations {
		transformedConfig, err := m.universalTransformer.Transform(configuration)
		if err != nil {
			log.Errorf("Error during transforming declarative configuration %+v: %+v", configuration, err)
			continue
		}
		for protoType, protoMessages := range transformedConfig {
			transformedConfigurations[protoType] = append(transformedConfigurations[protoType], protoMessages...)
		}
	}

	m.transformedMessagesMutex.Lock()
	m.transformedMessagesByHandler[handlerID] = transformedConfigurations
	m.transformedMessagesMutex.Unlock()
	m.shortCircuitReconciliationLoop()
}

// shortCircuitReconciliationLoop will short circuit the reconciliation loop and trigger a reconciliation loop run.
// Note that the reconciliation loop will not be run if:
//   - the short circuit loop signal has not been reset yet and is de-duped.
func (m *managerImpl) shortCircuitReconciliationLoop() {
	// In case the signal is already triggered, the current call (and the Signal() call) will be effectively de-duped.
	m.shortCircuitSignal.Signal()
}

func (m *managerImpl) startReconciliationLoop() {
	m.reconciliationTicker = time.NewTicker(m.reconciliationTickerDuration)

	go m.reconciliationLoop()
}

func (m *managerImpl) reconciliationLoop() {
	// While we currently do not have an exit in the form of "stopping" the reconciliation, still, ensure that
	// the ticker is stopped when we stop running the reconciliation.
	defer m.reconciliationTicker.Stop()
	for {
		select {
		case <-m.shortCircuitSignal.Done():
			m.shortCircuitSignal.Reset()
			m.runReconciliation()
		case <-m.reconciliationTicker.C:
			m.runReconciliation()
		}
	}
}

func (m *managerImpl) runReconciliation() {
	m.transformedMessagesMutex.RLock()
	transformedMessagesByHandler := maputil.ShallowClone(m.transformedMessagesByHandler)
	m.transformedMessagesMutex.RUnlock()
	m.reconcileTransformedMessages(transformedMessagesByHandler)
}

func (m *managerImpl) reconcileTransformedMessages(transformedMessagesByHandler map[string]protoMessagesByType) {
	log.Debugf("Run reconciliation for the next handlers: %v", maputil.Keys(transformedMessagesByHandler))
	transformedMessages := map[reflect.Type][]proto.Message{}
	for _, protoMessagesByType := range transformedMessagesByHandler {
		for protoType, protoMessages := range protoMessagesByType {
			transformedMessages[protoType] = append(transformedMessages[protoType], protoMessages...)
		}

	}
	m.reconcileUpsertAccessScopes(transformedMessages)
	m.reconcileUpsertPermissionSets(transformedMessages)
	m.reconcileUpsertRoles(transformedMessages)
	m.reconcileUpsertAuthProviders(transformedMessages)
	m.reconcileUpsertGroups(transformedMessages)
	// TODO(ROX-14694): Add deletion of resources.
	log.Debugf("Deleting all proto messages that have traits.Origin==DECLARATIVE but are not contained"+
		" within the current list of transformed messages: %+v", transformedMessagesByHandler)
}

func (m *managerImpl) reconcileUpsertAccessScopes(transformedMessages map[reflect.Type][]proto.Message) {
	accessScopes, ok := transformedMessages[accessScopeType]
	// No access scopes to reconcile.
	if !ok {
		return
	}
	for _, accessScope := range accessScopes {
		err := m.roleDS.UpsertAccessScope(m.reconciliationCtx, accessScope.(*storage.SimpleAccessScope))
		if err != nil {
			m.reconciliationErrorReporter.ProcessError(accessScope, err)
		}
	}
}

func (m *managerImpl) reconcileUpsertPermissionSets(transformedMessages map[reflect.Type][]proto.Message) {
	permissionSets, ok := transformedMessages[permissionSetType]
	// No permission sets to reconcile.
	if !ok {
		return
	}
	for _, permissionSet := range permissionSets {
		err := m.roleDS.UpsertPermissionSet(m.reconciliationCtx, permissionSet.(*storage.PermissionSet))
		if err != nil {
			m.reconciliationErrorReporter.ProcessError(permissionSet, err)
		}
	}
}

func (m *managerImpl) reconcileUpsertRoles(transformedMessages map[reflect.Type][]proto.Message) {
	roles, ok := transformedMessages[roleType]
	// No roles to reconcile.
	if !ok {
		return
	}
	for _, role := range roles {
		err := m.roleDS.UpsertRole(m.reconciliationCtx, role.(*storage.Role))
		if err != nil {
			m.reconciliationErrorReporter.ProcessError(role, err)
		}
	}
}

func (m *managerImpl) reconcileUpsertAuthProviders(transformedMessages map[reflect.Type][]proto.Message) {
	authProviders, ok := transformedMessages[authProviderType]
	// No auth providers to reconcile.
	if !ok {
		return
	}
	for _, protoMessage := range authProviders {
		authProvider := protoMessage.(*storage.AuthProvider)
		if err := m.authProviderRegistry.DeleteProvider(m.reconciliationCtx, authProvider.GetId(), true, true); err != nil {
			m.reconciliationErrorReporter.ProcessError(authProvider, err)
			continue
		}

		if _, err := m.authProviderRegistry.CreateProvider(m.reconciliationCtx, authproviders.WithStorageView(authProvider),
			authproviders.WithAttributeVerifier(authProvider),
			authproviders.WithValidateCallback(authProviderDatastore.Singleton())); err != nil {
			m.reconciliationErrorReporter.ProcessError(authProvider, err)
		}
	}
}

func (m *managerImpl) reconcileUpsertGroups(transformedMessages map[reflect.Type][]proto.Message) {
	groups, ok := transformedMessages[groupType]
	// No access scopes to reconcile.
	if !ok {
		return
	}
	for _, group := range groups {
		err := m.groupDS.Upsert(m.reconciliationCtx, group.(*storage.Group))
		if err != nil {
			m.reconciliationErrorReporter.ProcessError(group, err)
		}
	}
}
