package store

import (
	globaldb "github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/pkg/storecache"
	sync "github.com/stackrox/rox/pkg/sync"
	utils "github.com/stackrox/rox/pkg/utils"
)

var (
	singleton     Store
	singletonInit sync.Once
)

// Singleton returns a singleton of the Store class
func Singleton() Store {
	singletonInit.Do(func() {
		store, err := New(globaldb.GetGlobalDB(), storecache.NewMapBackedCache())
		utils.Must(err)
		singleton = store
	})
	return singleton
}
