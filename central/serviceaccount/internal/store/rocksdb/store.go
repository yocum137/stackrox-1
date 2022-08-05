// Code generated by rocksdb-bindings generator. DO NOT EDIT.

package rocksdb

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/db"
	"github.com/stackrox/rox/pkg/rocksdb"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
)

var (
	log = logging.LoggerForModule()

	bucket = []byte("service_accounts")
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	GetIDs(ctx context.Context) ([]string, error)
	Get(ctx context.Context, id string) (*storage.ServiceAccount, bool, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.ServiceAccount, []int, error)
	Upsert(ctx context.Context, obj *storage.ServiceAccount) error
	UpsertMany(ctx context.Context, objs []*storage.ServiceAccount) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
	Walk(ctx context.Context, fn func(obj *storage.ServiceAccount) error) error
	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)

	// Unused and only exists to satisfy interfaces used for Postgres
	GetByQuery(ctx context.Context, q *v1.Query) ([]*storage.ServiceAccount, error)
}

type storeImpl struct {
	crud db.Crud
}

func alloc() proto.Message {
	return &storage.ServiceAccount{}
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(*storage.ServiceAccount).GetId())
}

// New returns a new Store instance using the provided rocksdb instance.
func New(db *rocksdb.RocksDB) Store {
	globaldb.RegisterBucket(bucket, "ServiceAccount")
	return &storeImpl{
		crud: generic.NewCRUD(db, bucket, keyFunc, alloc, false),
	}
}

// Count returns the number of objects in the store
func (b *storeImpl) Count(_ context.Context) (int, error) {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.Count, "ServiceAccount")

	return b.crud.Count()
}

// Exists returns if the id exists in the store
func (b *storeImpl) Exists(_ context.Context, id string) (bool, error) {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.Exists, "ServiceAccount")

	return b.crud.Exists(id)
}

// GetIDs returns all the IDs for the store
func (b *storeImpl) GetIDs(_ context.Context) ([]string, error) {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.GetAll, "ServiceAccountIDs")

	return b.crud.GetKeys()
}

// Get returns the object, if it exists from the store
func (b *storeImpl) Get(_ context.Context, id string) (*storage.ServiceAccount, bool, error) {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.Get, "ServiceAccount")

	msg, exists, err := b.crud.Get(id)
	if err != nil || !exists {
		return nil, false, err
	}
	return msg.(*storage.ServiceAccount), true, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (b *storeImpl) GetMany(_ context.Context, ids []string) ([]*storage.ServiceAccount, []int, error) {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.GetMany, "ServiceAccount")

	msgs, missingIndices, err := b.crud.GetMany(ids)
	if err != nil {
		return nil, nil, err
	}
	objs := make([]*storage.ServiceAccount, 0, len(msgs))
	for _, m := range msgs {
		objs = append(objs, m.(*storage.ServiceAccount))
	}
	return objs, missingIndices, nil
}

// Upsert inserts the object into the DB
func (b *storeImpl) Upsert(_ context.Context, obj *storage.ServiceAccount) error {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.Add, "ServiceAccount")

	return b.crud.Upsert(obj)
}

// UpsertMany batches objects into the DB
func (b *storeImpl) UpsertMany(_ context.Context, objs []*storage.ServiceAccount) error {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.AddMany, "ServiceAccount")

	msgs := make([]proto.Message, 0, len(objs))
	for _, o := range objs {
		msgs = append(msgs, o)
    }

	return b.crud.UpsertMany(msgs)
}

// Delete removes the specified ID from the store
func (b *storeImpl) Delete(_ context.Context, id string) error {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.Remove, "ServiceAccount")

	return b.crud.Delete(id)
}

// Delete removes the specified IDs from the store
func (b *storeImpl) DeleteMany(_ context.Context, ids []string) error {
	defer metrics.SetRocksDBOperationDurationTime(time.Now(), ops.RemoveMany, "ServiceAccount")

	return b.crud.DeleteMany(ids)
}

// Walk iterates over all of the objects in the store and applies the closure
func (b *storeImpl) Walk(_ context.Context, fn func(obj *storage.ServiceAccount) error) error {
	return b.crud.Walk(func(msg proto.Message) error {
		return fn(msg.(*storage.ServiceAccount))
	})
}

// AckKeysIndexed acknowledges the passed keys were indexed
func (b *storeImpl) AckKeysIndexed(_ context.Context, keys ...string) error {
	return b.crud.AckKeysIndexed(keys...)
}

// GetKeysToIndex returns the keys that need to be indexed
func (b *storeImpl) GetKeysToIndex(_ context.Context) ([]string, error) {
	return b.crud.GetKeysToIndex()
}

// GetByQuery is unused and only exists to satisfy interfaces used for Postgres
func (b * storeImpl) GetByQuery(ctx context.Context, q *v1.Query) ([]*storage.ServiceAccount, error) {
	panic("unimplemented")
}
