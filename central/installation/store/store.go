package store

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
)

// internalStore provides access to the data layer
type internalStore interface {
	GetAll(ctx context.Context) ([]*storage.InstallationInfo, error)
	Upsert(ctx context.Context, installationinfo *storage.InstallationInfo) error
}

// Store provides access to the data layer
type Store interface {
	GetInstallationInfo(ctx context.Context) (*storage.InstallationInfo, error)
	AddInstallationInfo(ctx context.Context, installationinfo *storage.InstallationInfo) error
}

// NewStore creates a new installation store from an internal store
func NewStore(internal internalStore) Store {
	return &storeImpl{
		internal: internal,
	}
}

type storeImpl struct {
	internal internalStore
}

func (s *storeImpl) GetInstallationInfo(ctx context.Context) (*storage.InstallationInfo, error) {
	installations, err := s.internal.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(installations) == 0 {
		return nil, nil
	}
	return installations[0], nil
}

func (s *storeImpl) AddInstallationInfo(ctx context.Context, installationinfo *storage.InstallationInfo) error {
	return s.internal.Upsert(ctx, installationinfo)
}
