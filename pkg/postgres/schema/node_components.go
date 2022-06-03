// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableNodeComponentsStmt holds the create statement for table `node_components`.
	CreateTableNodeComponentsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists node_components (
                   Id varchar,
                   Name varchar,
                   Version varchar,
                   Source integer,
                   RiskScore numeric,
                   TopCvss numeric,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		GormModel: (*NodeComponents)(nil),
		Indexes:   []string{},
		Children:  []*postgres.CreateStmts{},
	}

	// NodeComponentsSchema is the go schema for table `node_components`.
	NodeComponentsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("node_components")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ImageComponent)(nil)), "node_components")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_IMAGE_COMPONENTS, "imagecomponent", (*storage.ImageComponent)(nil)))
		RegisterTable(schema, CreateTableNodeComponentsStmt)
		return schema
	}()
)

const (
	NodeComponentsTableName = "node_components"
)

// NodeComponents holds the Gorm model for Postgres table `node_components`.
type NodeComponents struct {
	Id         string             `gorm:"column:id;type:varchar;primaryKey"`
	Name       string             `gorm:"column:name;type:varchar"`
	Version    string             `gorm:"column:version;type:varchar"`
	Source     storage.SourceType `gorm:"column:source;type:integer"`
	RiskScore  float32            `gorm:"column:riskscore;type:numeric"`
	TopCvss    float32            `gorm:"column:topcvss;type:numeric"`
	Serialized []byte             `gorm:"column:serialized;type:bytea"`
}
