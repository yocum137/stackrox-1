package main

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=DataType
type DataType int

const (
	BOOL         DataType = 0
	NUMERIC      DataType = 1
	STRING       DataType = 2
	DATETIME     DataType = 3
	MAP          DataType = 4
	ENUM         DataType = 5
	ARRAY        DataType = 6
	STRING_ARRAY DataType = 7
	INT_ARRAY    DataType = 8
)

func dataTypeToSQLType(dataType DataType) string {
	var sqlType string
	switch dataType {
	case BOOL:
		sqlType = "bool"
	case NUMERIC:
		sqlType = "numeric"
	case STRING:
		sqlType = "varchar"
	case DATETIME:
		sqlType = "timestamp"
	case MAP:
		sqlType = "jsonb"
	case ENUM:
		sqlType = "integer"
	case STRING_ARRAY:
		sqlType = "text[]"
	case INT_ARRAY:
		sqlType = "intarray"
	default:
		panic(dataType.String())
	}
	return sqlType
}

func fieldsFromPath(b *strings.Builder, table *Table) {
	for i, elem := range table.Elems {
		if !elem.IsSearchable {
			continue
		}
		if !(table.Parent == nil && i == 0) {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%s %s", elem.SQLPath(), dataTypeToSQLType(elem.DataType))
	}
	for _, child := range table.Embedded {
		fieldsFromPath(b, child)
	}
}
