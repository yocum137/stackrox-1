package object

import (
	"fmt"
	"testing"
)

type SubSubObject struct {
	Key, Value string
}

type SubObject struct {
	Value string
	ManySubObjects []*SubSubObject
}

type TestObject struct {
	SubObject *SubObject
	ManySubObjects []*SubObject
	Map map[string]string
}

func TestJSONPath(t *testing.T) {
	str := WalkObject((*TestObject)(nil))
	str.Walk(func(p *PathContext, f Field) bool {
		fmt.Println(p.JSONPath())
		return true
	})

	/*
	func JSONPath(p *PathContext, field Field) string {
		if p.Prev != nil && (IsSlice(p.Prev.Field) || IsMap(p.Prev.Field)) {
			return p.Prev.JSONPath()
		}
		ctxPath := p.JSONPath()
		jsonTag := field.Tags().Get("json").First()
		return stringutils.JoinNonEmpty(".", ctxPath, stringutils.FirstNonEmpty(jsonTag, p.Name()))
	}
	 */



}