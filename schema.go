package main

import "fmt"

type Schema struct {
	Parent string
	Extends []string
	Schema *Property
}

func (schema *Schema) getName(object map[interface{}]interface{}) error {
	id, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("schema does not have an id")
	}
	schema.Schema = CreateProperty(id)
	return nil
}

func (schema *Schema) getParent(object map[interface{}]interface{}) {
	schema.Parent, _ = object["parent"].(string)
}

func (schema *Schema) getBaseSchemas(object map[interface{}]interface{}) error {
	extends, ok := object["extends"].([]interface{})
	if !ok {
		return nil
	}
	bases := make([]string, len(extends))
	for i, base := range extends {
		bases[i], ok = base.(string)
		if !ok {
			return fmt.Errorf("one of the base schemas is not a string")
		}
	}
	schema.Extends = bases
	return nil
}

func (schema *Schema) Parse(object map[interface{}]interface{}) error {
	err := schema.getName(object)
	if err != nil {
		return err
	}
	schema.getParent(object)
	err = schema.getBaseSchemas(object)
	if err != nil {
		return fmt.Errorf(
			"invalid schema %s: %v",
			schema.Schema.Name,
			err,
		)
	}
	err = schema.Schema.Parse("", object)
	if err != nil {
		return fmt.Errorf("%s - %v", schema.Schema.Name, err)
	}
	if !schema.Schema.Item.IsObject() {
		return fmt.Errorf(
			"invalid schema %s: schema should be an object",
			schema.Schema.Name,
		)
	}
	return nil
}

func (schema *Schema) Collect(depth int) []*Object {
	return schema.Schema.Collect(depth)
}