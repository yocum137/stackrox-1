package object

import (
	"reflect"
	"strings"
)

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

func (t TagValues) Empty() bool {
	return t == ""
}

// Get returns the value of the key in the tag
func (t TagValues) Get(key string) string {
	val, _ := t.Lookup(key)
	return val
}

func (t TagValues) Ignored() bool {
	return t == "-"
}

func (t TagValues) String() string {
	return string(t)
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
