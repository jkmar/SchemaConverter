package schema

import (
	"github.com/zimnx/YamlSchemaToGoStruct/item"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

// Convert converts given maps describing schemas to go structs
// args:
//   other []map[interface{}]interface{} - maps describing schemas than
//                                         should not be converted to go structs
//   toConvert []map[interface{}]interface{} - maps describing schemas that
//                                             should be converted to go structs
//   annotationDB string - annotation added to each field in schemas
//   annotationObject string - annotation added to each field in objects
//   suffix string - suffix added to each type name
// return:
//   1. list of go interfaces as strings
//   2. list of go structs as strings
//   3. list of implementations of interfaces as strings
//   4. error during execution
func Convert(
	other,
	toConvert []map[interface{}]interface{},
	suffix string,
) (
	interfaces,
	structs,
	implementations []string,
	err error,
) {
	var otherSet set.Set
	otherSet, err = parseAll(other)
	if err != nil {
		return
	}
	var toConvertSet set.Set
	toConvertSet, err = parseAll(toConvert)
	if err != nil {
		return
	}
	if err = collectSchemas(toConvertSet, otherSet); err != nil {
		return
	}
	dbObjects := set.New()
	jsonObjects := set.New()
	for _, toConvertSchema := range toConvertSet {
		objectFromSchema, _ := toConvertSchema.(*Schema).collectObjects(1, 0)
		dbObjects.InsertAll(objectFromSchema)
		var object set.Set
		object, err = toConvertSchema.(*Schema).collectObjects(-1, 1)
		if err != nil {
			return
		}
		jsonObjects.InsertAll(object)
	}
	for _, object := range dbObjects {
		dbObject := object.(*item.Object)
		interfaces = append(interfaces, dbObject.GenerateInterface(suffix))
		structs = append(structs, dbObject.GenerateStruct(suffix))
		implementations = append(implementations, dbObject.GenerateImplementation(suffix))
	}
	for _, object := range jsonObjects {
		jsonObject := object.(*item.Object)
		interfaces = append(interfaces, jsonObject.GenerateInterface(suffix))
		structs = append(structs, jsonObject.GenerateStruct(suffix))
		implementations = append(implementations, jsonObject.GenerateImplementation(suffix))
	}
	return
}
