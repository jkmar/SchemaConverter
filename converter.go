package main

import (
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/item"
)

func Convert(
	other,
	toConvert []map[interface{}]interface{},
	prefixDB,
	prefixObject,
	suffix string,
) ([]string, error) {
	otherSet, err := ParseAll(other)
	if err != nil {
		return nil, err
	}
	toConvertSet, err := ParseAll(toConvert)
	if err != nil {
		return nil, err
	}
	if err := collectSchemas(toConvertSet, otherSet); err != nil {
		return nil, err
	}
	dbObjects := set.New()
	jsonObjects := set.New()
	for _, schema := range toConvertSet {
		schemaObject, _ := schema.(*Schema).CollectObjects(1, 0)
		if err := dbObjects.SafeInsertAll(schemaObject); err != nil {
			return nil, fmt.Errorf(
				"multiple schemas with the same name: %s",
				schemaObject.Any().Name(),
			)
		}
		object, err := schema.(*Schema).CollectObjects(-1, 1)
		if err != nil {
			return nil, err
		}
		// TODO check if objects are equal
		jsonObjects.InsertAll(object);
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
