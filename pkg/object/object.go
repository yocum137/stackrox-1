package object

import (
	"reflect"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/stackrox/rox/pkg/stringutils"
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
	Children() []Field
}

func IsEnum(f Field) bool {
	_, ok := f.(Enum)
	return ok
}

func IsSlice(f Field) bool {
	_, ok := f.(Slice)
	return ok
}

func IsTime(f Field) bool {
	_, ok := f.(Time)
	return ok
}

func IsStruct(f Field) bool {
	_, ok := f.(Struct)
	return ok
}

func IsMap(f Field) bool {
	_, ok := f.(Map)
	return ok
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

func (f baseField) Children() []Field {
	return nil
}

// Map defines a map type
type Map struct {
	Field
	Value Field
}

func (m Map) Children() []Field {
	return []Field {m.Value}
}

type Time struct {
	Field
}

// Slice defines a repeated field
type Slice struct {
	Field
	Value Field
}

func (s Slice) Children() []Field {
	return []Field{s.Value}
}

// Enum wraps an int32 with the proto descriptor
type Enum struct {
	Field
	Descriptor *descriptor.EnumDescriptorProto
}

// Struct defines a struct in the object
type Struct struct {
	Field
	Fields     []Field
	StructType StructType
}

func (s Struct) Children() []Field {
	return s.Fields
}

type contextFunc func(ctx *PathContext, obj Field) bool

func JSONPath(p *PathContext, field Field) string {
	if p.Prev != nil && (IsSlice(p.Prev.Field) || IsMap(p.Prev.Field)) {
		return p.Prev.JSONPath()
	}
	ctxPath := p.JSONPath()
	jsonTag := field.Tags().Get("json").First()
	return stringutils.JoinNonEmpty(".", ctxPath, stringutils.FirstNonEmpty(jsonTag, p.Name()))
}

type PathContext struct {
	Prev *PathContext
	Field
}

func (p *PathContext) JSONPath() string {
	if IsSlice(p.Field) || IsMap(p.Field) {
		return p.Prev.JSONPath()
	}
	jsonTag := p.Field.Tags().Get("json").First()
	if p.Prev == nil {
		return stringutils.FirstNonEmpty(jsonTag, p.Name())
	}
	return stringutils.JoinNonEmpty(".", p.Prev.JSONPath(), stringutils.FirstNonEmpty(jsonTag, p.Name()))
}

func (s Struct) Walk(fn contextFunc) {
	ctx := &PathContext{
		Prev:  nil,
		Field: s,
	}
	walkHelper(ctx, fn)
}

func walkHelper(path *PathContext, fn contextFunc) {
	for _, field := range path.Field.Children() {
		if ok := fn(path, field); !ok {
			continue
		}
		childCtx := &PathContext{
			Prev:  path,
			Field: field,
		}
		walkHelper(childCtx, fn)
	}
}
