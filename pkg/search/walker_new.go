package search

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/object"
	"github.com/stackrox/rox/pkg/search/enumregistry"
	"github.com/stackrox/rox/pkg/stringutils"
)

type searchWalker struct {
	category v1.SearchCategory
	fields   map[FieldLabel]*Field
}

// Walk iterates over the obj and creates a search.Map object from the found struct tags
func Walk(category v1.SearchCategory, prefix string, obj interface{}) OptionsMap {
	objStruct := object.WalkObject(obj)

	sw := searchWalker{
		category: category,
		fields:   make(map[FieldLabel]*Field),
	}
	sw.walkRecursive(prefix, objStruct)
	return OptionsMapFromMap(category, sw.fields)
}

func (s *searchWalker) addSearchField(path string, field object.Field, dataType v1.SearchDataType) {
	searchTag, ok := field.Tags().Lookup("search")
	if !ok {
		return
	}

	fieldName := searchTag.First()
	if !FieldLabelSet.Contains(fieldName) {
		log.Panicf("Field %q is not a valid FieldLabel. You may need to add it to pkg/search/options.go", fieldName)
	}

	s.fields[FieldLabel(fieldName)] = &Field{
		FieldPath: path,
		Store:     searchTag.Exists("store"),
		Hidden:    searchTag.Exists("hidden"),
		Category:  s.category,
		Analyzer:  searchTag.Get("analyzer"),
		Type:      dataType,
	}
}

func (s *searchWalker) handlePrimitive(prefix string, obj object.Field) {
	switch obj.Type() {
	case reflect.String:
		s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_STRING)
	case reflect.Bool:
		s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_BOOL)
	case reflect.Int32:
		if enum, ok := obj.(object.Enum); ok {
			enumregistry.Add(prefix, enum.Descriptor)
			s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_ENUM)
			return
		}
		s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_NUMERIC)
	case reflect.Uint32, reflect.Uint64, reflect.Int64, reflect.Float32, reflect.Float64:
		s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_NUMERIC)
	default:
		panic(fmt.Sprintf("Type %s for field %s is not currently handled", obj.Type(), prefix))
	}
}

func (s *searchWalker) walkRecursive(parentPrefix string, obj object.Field) {
	if obj.Tags().Get("search") == "-" {
		return
	}

	var jsonName string
	if values, ok := obj.Tags().Lookup("json"); ok {
		jsonName = values.First()
	} else {
		jsonName = obj.Name()
	}

	prefix := stringutils.JoinNonEmpty(".", parentPrefix, jsonName)
	switch obj.Type() {
	case reflect.Slice:
		slice := obj.(object.Slice)
		if slice.Value.Type() == reflect.Struct {
			s.walkRecursive(parentPrefix, slice.Value)
			return
		}
		s.handlePrimitive(prefix, slice.Value)
	case reflect.Struct:
		structObj := obj.(object.Struct)
		switch structObj.StructType {
		case object.TIME:
			s.addSearchField(prefix+".seconds", obj, v1.SearchDataType_SEARCH_DATETIME)
		case object.ONEOF, object.MESSAGE:
			for _, field := range structObj.Fields {
				s.walkRecursive(prefix, field)
			}
		}
	case reflect.Map:
		s.addSearchField(prefix, obj, v1.SearchDataType_SEARCH_MAP)
	default:
		s.handlePrimitive(prefix, obj)
	case reflect.Interface:
	}
}
