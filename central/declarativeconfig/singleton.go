package declarativeconfig

import (
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

// ManagerSingleton provides the instance of Manager to use.
func ManagerSingleton(registry authproviders.Registry) Manager {
	once.Do(func() {
		instance = New(
			env.DeclarativeConfigReconcileInterval.DurationSetting(),
			env.DeclarativeConfigWatchInterval.DurationSetting(),
			roleDatastore.Singleton(),
			authProviderDatastore.Singleton(),
			registry)
	})
	return instance
}
