package search

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/object"
	"github.com/stackrox/rox/pkg/search/enumregistry"
)

// Walk iterates over the obj and creates a search.Map object from the found struct tags
func WalkVisitor(category v1.SearchCategory, prefix string, obj interface{}) OptionsMap {
	objStruct := object.WalkObject(obj)

	fields := make(map[FieldLabel]*Field)

	objStruct.Walk(func(path *object.PathContext, field object.Field) bool {
		searchTag := field.Tags().Get("search")
		if searchTag.Ignored() {
			return false
		}
		if object.IsStruct(field) || object.IsSlice(field) || searchTag.Empty() {
			return true
		}

		fieldName := searchTag.First()
		if !FieldLabelSet.Contains(fieldName) {
			log.Panicf("Field %q is not a valid FieldLabel. You may need to add it to pkg/search/options.go", fieldName)
		}
		var dataType v1.SearchDataType
		switch {
		case object.IsStruct(field), object.IsSlice(field):
			return true
		case object.IsTime(field):
			dataType = v1.SearchDataType_SEARCH_DATETIME
		case object.IsEnum(field):
			enumregistry.Add(prefix, field.(object.Enum).Descriptor)
			dataType = v1.SearchDataType_SEARCH_ENUM
		case object.IsMap(field):
			dataType = v1.SearchDataType_SEARCH_MAP
		default:
			switch field.Type() {
			case reflect.String:
				dataType = v1.SearchDataType_SEARCH_STRING
			case reflect.Bool:
				dataType = v1.SearchDataType_SEARCH_BOOL
			case reflect.Int32, reflect.Uint32, reflect.Uint64, reflect.Int64, reflect.Float32, reflect.Float64:
				dataType = v1.SearchDataType_SEARCH_NUMERIC
			default:
				panic(fmt.Sprintf("Type %s for field %s is not currently handled", field.Type(), prefix))
			}
		}

		fields[FieldLabel(fieldName)] = &Field{
			FieldPath: prefix + "." + object.JSONPath(path, field),
			Store:     searchTag.Exists("store"),
			Hidden:    searchTag.Exists("hidden"),
			Category:  category,
			Analyzer:  searchTag.Get("analyzer"),
			Type:      dataType,
		}

		fmt.Println(prefix + "." + object.JSONPath(path, field), field.Name())
		return true
	})
	return OptionsMapFromMap(category, fields)
}
