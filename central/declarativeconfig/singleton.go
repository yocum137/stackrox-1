package declarativeconfig

import (
	"github.com/gogo/protobuf/proto"
	authProviderDatastore "github.com/stackrox/rox/central/authprovider/datastore"
	roleDatastore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once     sync.Once
	instance Manager
)

type noOpErrorReporter struct{}

func (n noOpErrorReporter) ProcessError(_ proto.Message, _ error) {}

// ManagerSingleton provides the instance of Manager to use.
func ManagerSingleton(registry authproviders.Registry) Manager {
	once.Do(func() {
		instance = New(
			env.DeclarativeConfigReconcileInterval.DurationSetting(),
			env.DeclarativeConfigWatchInterval.DurationSetting(),
			roleDatastore.Singleton(),
			authProviderDatastore.Singleton(),
			registry,
			// TODO(ROX-15088): replace with actual health reporter
			noOpErrorReporter{})
	})
	return instance
}
