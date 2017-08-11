package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

// Array is an implementation of Item interface
type Array struct {
	arrayItem Item
}

// IsNull implementation
func (array *Array) IsNull() bool {
	return false
}

// Type implementation
func (array *Array) Type(suffix string) string {
	return "[]" + array.arrayItem.Type(suffix)
}

// InterfaceType implementation
func (array *Array) InterfaceType(suffix string) string {
	return "[]" + array.arrayItem.InterfaceType(suffix)
}

// AddProperties implementation
func (array *Array) AddProperties(set set.Set, safe bool) error {
	return fmt.Errorf("cannot add properties to an array")
}

// Parse implementation
func (array *Array) Parse(
	prefix string,
	level int,
	required bool,
	data map[interface{}]interface{},
) (err error) {
	next, ok := data["items"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"array %s does not have items",
			prefix,
		)
	}
	objectType, ok := next["type"]
	if !ok {
		return fmt.Errorf(
			"items of array %s do not have a type",
			prefix,
		)
	}
	array.arrayItem, err = CreateItem(objectType)
	if err != nil {
		return fmt.Errorf("array %s: %v", prefix, err)
	}
	return array.arrayItem.Parse(prefix, level, required, next)
}

// CollectObjects implementation
func (array *Array) CollectObjects(limit, offset int) (set.Set, error) {
	return array.arrayItem.CollectObjects(limit, offset)
}

// CollectProperties implementation
func (array *Array) CollectProperties(limit, offset int) (set.Set, error) {
	return array.arrayItem.CollectProperties(limit, offset)
}

//func (array *Array) GenerateSetter(prefix, arg string) string {
//	if _, ok := array.arrayItem.(*PlainItem); ok {
//		return fmt.Sprintf(
//			"%s = %s",
//			prefix,
//			arg,
//		)
//	}
//	return fmt.Sprintf(
//		"%s = make(%s, len(%s))\n"
//	)
//}
