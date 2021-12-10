package object

import (
	"reflect"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// StructType distinguishes between different struct types
type StructType int

// Different types of structs
const (
	MESSAGE StructType = iota
	TIME
	ONEOF
)

// Field is the interface that exposes fields in the objects
type Field interface {
	Name() string
	Type() reflect.Kind
	Tags() Tags
}

func newField(field reflect.StructField, kind reflect.Kind) Field {
	return baseField{
		name: field.Name,
		kind: kind,
		tags: Tags{
			StructTag: field.Tag,
		},
	}
}

type baseField struct {
	name string
	kind reflect.Kind
	tags Tags
}

func (f baseField) Name() string {
	return f.name
}

func (f baseField) Type() reflect.Kind {
	return f.kind
}

func (f baseField) Tags() Tags {
	return f.tags
}

// Map defines a map type
type Map struct {
	Field
	Value Field
}

// Slice defines a repeated field
type Slice struct {
	Field
	Value Field
}

// Enum wraps an int32 with the proto descriptor
type Enum struct {
	Field
	Descriptor *descriptor.EnumDescriptorProto
}

// Tags is a utility around the StructTag object
type Tags struct {
	reflect.StructTag
}

// TagValues is a wrapper around the struct tags that provides utility functions
type TagValues string

// First returns the first element in the tag. Useful for things like search
func (t TagValues) First() string {
	return strings.Split(string(t), ",")[0]
}

// Lookup returns the value of the key and a bool if the key exists in the tag
func (t TagValues) Lookup(key string) (string, bool) {
	for _, value := range strings.Split(string(t), ",") {
		spl := strings.Split(value, "=")
		if spl[0] == key {
			if len(spl) > 1 {
				return spl[1], true
			}
			return "", true
		}
	}
	return "", false
}

// Exists checks the existence of the key in that tag
func (t TagValues) Exists(key string) bool {
	_, ok := t.Lookup(key)
	return ok
}

// Get returns the value of the key in the tag
func (t TagValues) Get(key string) string {
	val, _ := t.Lookup(key)
	return val
}

// Get returns the tag string for the specified key
func (t Tags) Get(key string) TagValues {
	values, _ := t.Lookup(key)
	return values
}

// Lookup returns the tag value for the specified key and a bool that signifies if the key exists in the tag
func (t Tags) Lookup(key string) (TagValues, bool) {
	tagStr, ok := t.StructTag.Lookup(key)
	return TagValues(tagStr), ok
}

// Struct defines a struct in the object
type Struct struct {
	Field
	Fields     []Field
	StructType StructType
}

type contextFunc func(ctx PathContext, obj Field) bool

type PathContext struct {

}

func (s Struct) Walk(fn contextFunc) {



})

func (s Struct) walkHelper(path PathContext, fn contextFunc) {



	fn(path)
	for field := range s.Fields {

	}
})