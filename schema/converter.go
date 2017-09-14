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
	rawSuffix,
	interfaceSuffix,
	packageName string,
) (*Generated, error) {
	otherSet, err := parseAll(other)
	if err != nil {
		return nil, err
	}

	toConvertSet, err := parseAll(toConvert)
	if err != nil {
		return nil, err
	}

	if err = collectSchemas(toConvertSet, otherSet); err != nil {
		return nil, err
	}

	dbObjects := set.New()
	jsonObjects := set.New()
	for _, toConvertSchema := range toConvertSet {
		objectFromSchema, _ := toConvertSchema.(*Schema).collectObjects(1, 0)
		dbObjects.InsertAll(objectFromSchema)
		var object set.Set
		object, err = toConvertSchema.(*Schema).collectObjects(-1, 1)
		if err != nil {
			return nil, err
		}
		jsonObjects.InsertAll(object)
	}

	result := &Generated{}
	for _, object := range dbObjects.ToArray() {
		item := object.(*item.Object)
		result.RawCrud = append(
			result.RawCrud,
			item.GenerateFetch(packageName, rawSuffix, false, true),
			item.GenerateFetch(packageName, rawSuffix, true, true),
			item.GenerateList(packageName, rawSuffix, false, true),
			item.GenerateList(packageName, rawSuffix, true, true),
		)

		result.Crud = append(
			result.Crud,
			item.GenerateFetch(packageName, rawSuffix, false, false),
			item.GenerateFetch(packageName, rawSuffix, true, false),
			item.GenerateList(packageName, rawSuffix, false, false),
			item.GenerateList(packageName, rawSuffix, true, false),
		)
	}
	dbObjects.InsertAll(jsonObjects)
	for _, object := range dbObjects.ToArray() {
		item := object.(*item.Object)
		result.RawInterfaces = append(
			result.RawInterfaces,
			item.GenerateInterface(interfaceSuffix),
		)
		result.Interfaces = append(
			result.Interfaces,
			item.GenerateMutableInterface(interfaceSuffix, rawSuffix),
		)
		result.Structs = append(
			result.Structs,
			item.GenerateStruct(rawSuffix),
		)
		result.Implementations = append(
			result.Implementations,
			item.GenerateImplementation(interfaceSuffix, rawSuffix),
		)
		result.Constructors = append(
			result.Constructors,
			item.GenerateConstructor(rawSuffix),
		)
	}

	return result, nil
}
