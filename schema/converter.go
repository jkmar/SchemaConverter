package schema

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/item"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

func Convert(
	other,
	toConvert []map[interface{}]interface{},
	prefixDB,
	prefixObject,
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
		result = append(result, object.(*item.Object).GenerateStruct(suffix, prefixDB))
	}
	for _, object := range jsonObjects {
		result = append(result, object.(*item.Object).GenerateStruct(suffix, prefixObject))
	}
	return result, nil
}
