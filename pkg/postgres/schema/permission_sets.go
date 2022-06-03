// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTablePermissionSetsStmt holds the create statement for table `permission_sets`.
	CreateTablePermissionSetsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists permission_sets (
                   Id varchar,
                   Name varchar UNIQUE,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		GormModel: (*PermissionSets)(nil),
		Indexes:   []string{},
		Children:  []*postgres.CreateStmts{},
	}

	// PermissionSetsSchema is the go schema for table `permission_sets`.
	PermissionSetsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("permission_sets")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.PermissionSet)(nil)), "permission_sets")
		RegisterTable(schema, CreateTablePermissionSetsStmt)
		return schema
	}()
)

const (
	PermissionSetsTableName = "permission_sets"
)

// PermissionSets holds the Gorm model for Postgres table `permission_sets`.
type PermissionSets struct {
	Id         string `gorm:"column:id;type:varchar;primaryKey"`
	Name       string `gorm:"column:name;type:varchar;unique"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
