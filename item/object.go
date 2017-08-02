package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

type Object struct {
	objectType string
	properties set.Set
}

func (item *Object) Name() string {
	return item.objectType
}

func (item *Object) Type(suffix string) string {
	return util.ToGoName(item.Name(), suffix)
}

func (item *Object) IsObject() bool {
	return true
}

func (item *Object) AddProperties(properties set.Set, safe bool) error {
	if properties.Empty() {
		return nil
	}
	if item.properties == nil {
		item.properties = set.New()
	}
	if safe {
		if err := item.properties.SafeInsertAll(properties); err != nil {
			return fmt.Errorf(
				"object %s: multiple properties have the same name",
				item.Name(),
			)
		}
	} else {
		properties.InsertAll(item.properties)
		item.properties = properties
	}
	return nil
}

func (item *Object) Parse(prefix string, object map[interface{}]interface{}) error {
	item.objectType = prefix
	item.properties = set.New()
	properties, ok := object["properties"]
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
				item.Name(),
			)
		}
		newProperty := CreateProperty(strProperty)
		item.properties.Insert(newProperty)
		definitionMap, ok := definition.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf(
				"object %s has invalid property %s",
				item.Name(),
				strProperty,
			)
		}
		if err := newProperty.Parse(prefix, definitionMap); err != nil {
			return err
		}
	}
	return nil
}

func (item *Object) CollectObjects(limit, offset int) (set.Set, error) {
	if limit == 0 {
		return nil, nil
	}
	result := set.New()
	if offset <= 0 {
		result.Insert(item)
	}
	for _, property := range item.properties {
		other, err := property.(*Property).CollectObjects(limit-1, offset-1)
		if err != nil {
			return nil, err
		}
		if err = result.SafeInsertAll(other); err != nil {
			return nil, fmt.Errorf(
				"multiple objects with the same type at object %s",
				item.Name(),
			)
		}
	}
	return result, nil
}

func (item *Object) CollectProperties(limit, offset int) (set.Set, error) {
	result := set.New()
	for _, property := range item.properties {
		other, err := property.(*Property).CollectProperties(limit, offset)
		if err != nil {
			return nil, err
		}
		err = result.SafeInsertAll(other)
		if err != nil {
			return nil, fmt.Errorf(
				"multiple properties with the same name at object %s",
				item.Name(),
			)
		}
	}
	return result, nil
}

func (item *Object) GenerateStruct(suffix, annotation string) string {
	code := "type " + item.Type(suffix) + " struct {\n"
	properties := item.properties.ToArray()
	for _, property := range properties {
		code += property.(*Property).GenerateProperty(suffix, annotation)
	}
	return code + "}\n"
}
