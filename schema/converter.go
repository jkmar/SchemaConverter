package schema

import (
	"fmt"
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
//   1. list of go structs as strings
//   2. error during execution
func Convert(
	other,
	toConvert []map[interface{}]interface{},
	annotationDB,
	annotationObject,
	suffix string,
) ([]string, error) {
	otherSet, err := parseAll(other)
	if err != nil {
		return nil, err
	}
	toConvertSet, err := parseAll(toConvert)
	if err != nil {
		return nil, err
	}
	if err := collectSchemas(toConvertSet, otherSet); err != nil {
		return nil, err
	}
	dbObjects := set.New()
	jsonObjects := set.New()
	for _, schema := range toConvertSet {
		schemaObject, _ := schema.(*Schema).collectObjects(1, 0)
		if err := dbObjects.SafeInsertAll(schemaObject); err != nil {
			return nil, fmt.Errorf(
				"multiple schemas with the same name: %s",
				schemaObject.Any().Name(),
			)
		}
		object, err := schema.(*Schema).collectObjects(-1, 1)
		if err != nil {
			return nil, err
		}
		jsonObjects.InsertAll(object)
	}
	result := []string{}
	for _, object := range dbObjects {
		result = append(result, object.(*item.Object).GenerateStruct(suffix, annotationDB))
	}
	for _, object := range jsonObjects {
		result = append(result, object.(*item.Object).GenerateStruct(suffix, annotationObject))
	}
	return result, nil
}
