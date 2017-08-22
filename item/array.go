package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/hash"
	"github.com/zimnx/YamlSchemaToGoStruct/name"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

// Array is an implementation of Item interface
type Array struct {
	arrayItem Item
}

// Copy implementation
func (array *Array) Copy() Item {
	newArray := *array
	return &newArray
}

// ToString implementation
func (array *Array) ToString() string {
	return "#[]"
}

// Compress implementation
func (array *Array) Compress(source, destination hash.IHashable) {
	if sourceItem, ok := source.(Item); array.arrayItem == destination && ok {
		array.arrayItem = sourceItem
	}
}

// GetChildren implementation
func (array *Array) GetChildren() []hash.IHashable {
	return []hash.IHashable{array.arrayItem}
}

// ChangeName implementation
func (array *Array) ChangeName(mark name.Mark) {
	array.arrayItem.ChangeName(mark)
}

// ContainsObject implementation
func (array *Array) ContainsObject() bool {
	return array.arrayItem.ContainsObject()
}

// IsNull implementation
func (array *Array) IsNull() bool {
	return false
}

// MakeRequired implementation
func (array *Array) MakeRequired() {
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

// GenerateGetter implementation
func (array *Array) GenerateGetter(
	variable,
	argument,
	interfaceSuffix string,
	depth int,
) string {
	indent := util.Indent(depth)
	var resultSuffix string
	if depth == 1 {
		if !array.ContainsObject() {
			return fmt.Sprintf(
				"%sreturn %s",
				indent,
				variable,
			)
		}
		resultSuffix = fmt.Sprintf(
			"\n%sreturn %s",
			indent,
			argument,
		)
	}
	index := arrayIndex(depth)
	return fmt.Sprintf(
		"%s%s make(%s, len(%s))\n%sfor %c := range %s {\n%s\n%s}%s",
		indent,
		util.ResultPrefix(argument, depth, true),
		array.InterfaceType(interfaceSuffix),
		variable,
		indent,
		util.IndexVariable(depth),
		variable,
		array.arrayItem.GenerateGetter(
			variable+index,
			argument+index,
			interfaceSuffix,
			depth+1,
		),
		indent,
		resultSuffix,
	)
}

// GenerateSetter implementation
func (array *Array) GenerateSetter(
	variable,
	argument,
	typeSuffix string,
	depth int,
) string {
	indent := util.Indent(depth)
	if _, ok := array.arrayItem.(*PlainItem); ok {
		return fmt.Sprintf(
			"%s%s = %s",
			indent,
			variable,
			argument,
		)
	}
	index := arrayIndex(depth)
	return fmt.Sprintf(
		"%s%s = make(%s, len(%s))\n%sfor %c := range %s {\n%s\n%s}",
		indent,
		variable,
		array.Type(typeSuffix),
		argument,
		indent,
		util.IndexVariable(depth),
		argument,
		array.arrayItem.GenerateSetter(
			variable+index,
			argument+index,
			typeSuffix,
			depth+1,
		),
		indent,
	)
}

func arrayIndex(depth int) string {
	return fmt.Sprintf("[%c]", util.IndexVariable(depth))
}
