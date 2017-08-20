package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/hash"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
	"strings"
)

// Object is an implementation of Item interface
type Object struct {
	objectType string
	properties set.Set
}

// ToString implementation
func (object *Object) ToString() string {
	return "#*"
}

// Compress implementation
func (object *Object) Compress(source, destination hash.IHashable) {
	if destinationProperty, ok := destination.(*Property); ok {
		if sourceProperty, ok := source.(*Property); ok {
			if object.properties.Contains(destinationProperty) {
				object.properties.Delete(destinationProperty)
				object.properties.Insert(sourceProperty)
			}
		}
	}
}

// GetChildren implementation
func (object *Object) GetChildren() []hash.IHashable {
	sorted := object.properties.ToArray()
	result := make([]hash.IHashable, len(sorted))
	for i, property := range sorted {
		result[i] = property.(hash.IHashable)
	}
	return result
}

// ContainsObject implementation
func (object *Object) ContainsObject() bool {
	return true
}

// IsNull implementation
func (object *Object) IsNull() bool {
	return false
}

// Name is a function that allows object to be used as a set element
func (object *Object) Name() string {
	return object.objectType
}

// Type implementation
func (object *Object) Type(suffix string) string {
	return "*" + object.getType(suffix)
}

// InterfaceType implementation
func (object *Object) InterfaceType(suffix string) string {
	return "I" + object.getType(suffix)
}

// AddProperties implementation
func (object *Object) AddProperties(properties set.Set, safe bool) error {
	if properties.Empty() {
		return nil
	}
	if object.properties == nil {
		object.properties = set.New()
	}
	if safe {
		if err := object.properties.SafeInsertAll(properties); err != nil {
			return fmt.Errorf(
				"object %s: multiple properties have the same name",
				object.Name(),
			)
		}
	} else {
		properties.InsertAll(object.properties)
		object.properties = properties
	}
	return nil
}

// Parse implementation
func (object *Object) Parse(
	prefix string,
	level int,
	required bool,
	data map[interface{}]interface{},
) error {
	object.objectType = prefix
	object.properties = set.New()
	requiredMap, err := parseRequired(data)
	if err != nil {
		return fmt.Errorf(
			"object %s: %v",
			prefix,
			err,
		)
	}
	properties, ok := data["properties"]
	if !ok {
		return nil
	}
	next, ok := properties.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"object %s has invalid properties",
			prefix,
		)
	}
	for property, definition := range next {
		strProperty, ok := property.(string)
		if !ok {
			return fmt.Errorf(
				"object %s has property which name is not a string",
				object.Name(),
			)
		}
		newProperty := CreateProperty(strProperty)
		object.properties.Insert(newProperty)
		definitionMap, ok := definition.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf(
				"object %s has invalid property %s",
				object.Name(),
				strProperty,
			)
		}
		if err := newProperty.Parse(
			prefix,
			level,
			requiredMap[strProperty],
			definitionMap,
		); err != nil {
			return err
		}
	}
	return nil
}

// CollectObjects implementation
func (object *Object) CollectObjects(limit, offset int) (set.Set, error) {
	if limit == 0 {
		return nil, nil
	}
	result := set.New()
	if offset <= 0 {
		result.Insert(object)
	}
	for _, property := range object.properties {
		other, err := property.(*Property).CollectObjects(limit-1, offset-1)
		if err != nil {
			return nil, err
		}
		if err = result.SafeInsertAll(other); err != nil {
			return nil, fmt.Errorf(
				"multiple objects with the same type at object %s",
				object.Name(),
			)
		}
	}
	return result, nil
}

// CollectProperties implementation
func (object *Object) CollectProperties(limit, offset int) (set.Set, error) {
	result := set.New()
	for _, property := range object.properties {
		other, err := property.(*Property).CollectProperties(limit, offset)
		if err != nil {
			return nil, err
		}
		err = result.SafeInsertAll(other)
		if err != nil {
			return nil, fmt.Errorf(
				"multiple properties with the same name at object %s",
				object.Name(),
			)
		}
	}
	return result, nil
}

// GenerateGetter implementation
func (object *Object) GenerateGetter(
	variable,
	argument,
	suffix string,
	depth int,
) string {
	return fmt.Sprintf(
		"%s%s %s",
		util.Indent(depth),
		util.ResultPrefix(argument, depth, false),
		variable,
	)
}

// GenerateSetter implementation
func (object *Object) GenerateSetter(
	variable,
	argument,
	suffix string,
	depth int,
) string {
	return fmt.Sprintf(
		"%s%s = %s.(%s)",
		util.Indent(depth),
		variable,
		argument,
		object.Type(suffix),
	)
}

// GenerateStruct creates a struct of an object
// with suffix added to type name of each field
func (object *Object) GenerateStruct(suffix string) string {
	code := "type " + object.getType(suffix) + " struct {\n"
	properties := object.properties.ToArray()
	for _, property := range properties {
		code += property.(*Property).GenerateProperty(suffix)
	}
	return code + "}\n"
}

// GenerateInterface creates an interface of an object
// with suffix added to objects type
func (object *Object) GenerateInterface(suffix string) string {
	code := "type " + object.InterfaceType(suffix) + " interface {\n"
	properties := object.properties.ToArray()
	for _, property := range properties {
		code += fmt.Sprintf(
			"\t%s\n\t%s\n",
			property.(*Property).GetterHeader(suffix),
			property.(*Property).SetterHeader(suffix, false),
		)
	}
	return code + "}\n"
}

// GenerateImplementation creates an implementation of an objects
// getter and setter methods
func (object *Object) GenerateImplementation(suffix string) string {
	variable := util.VariableName(util.AddName(object.objectType, suffix))
	prefix := fmt.Sprintf(
		"func (%s %s) ",
		variable,
		object.Type(suffix),
	)
	properties := object.properties.ToArray()
	code := ""
	for _, property := range properties {
		code += fmt.Sprintf(
			"%s%s\n\n%s%s\n\n",
			prefix,
			property.(*Property).GenerateGetter(variable, suffix),
			prefix,
			property.(*Property).GenerateSetter(variable, suffix),
		)
	}
	return strings.TrimSuffix(code, "\n")
}

func parseRequired(data map[interface{}]interface{}) (map[string]bool, error) {
	required, ok := data["required"]
	if !ok {
		return nil, nil
	}
	list, ok := required.([]interface{})
	if !ok {
		return nil, fmt.Errorf("required should be a list of strings")
	}
	result := map[string]bool{}
	for _, element := range list {
		elementString, ok := element.(string)
		if !ok {
			return nil, fmt.Errorf("required should be a list of strings")
		}
		result[elementString] = true
	}
	return result, nil
}

func (object *Object) getType(suffix string) string {
	return util.ToGoName(object.Name(), suffix)
}
