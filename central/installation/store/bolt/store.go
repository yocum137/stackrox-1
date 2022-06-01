package bolt

import (
	"context"
	"time"

	proto "github.com/gogo/protobuf/proto"
	metrics "github.com/stackrox/rox/central/metrics"
	storage "github.com/stackrox/rox/generated/storage"
	singletonstore "github.com/stackrox/rox/pkg/bolthelper/singletonstore"
	ops "github.com/stackrox/rox/pkg/metrics"
	bbolt "go.etcd.io/bbolt"
)

var (
	bucketName = []byte("installationInfo")
)

// New creates a new bolt store
func New(db *bbolt.DB) *store {
	return &store{underlying: singletonstore.New(db, bucketName, func() proto.Message {
		return new(storage.InstallationInfo)
	}, "InstallationInfo")}
}

type store struct {
	underlying singletonstore.SingletonStore
}

func (s *store) Upsert(ctx context.Context, installationinfo *storage.InstallationInfo) error {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Add, "InstallationInfo")
	return s.underlying.Create(installationinfo)
}

func (s *store) GetAll(ctx context.Context) ([]*storage.InstallationInfo, error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.GetAll, "InstallationInfo")
	msg, err := s.underlying.Get()
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, nil
	}
	return []*storage.InstallationInfo{msg.(*storage.InstallationInfo)}, nil
}
