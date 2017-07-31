package main

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

type Schema struct {
	parent  string
	extends []string
	schema  *Property
}

func (schema *Schema) Name() string {
	return schema.schema.Name()
}

func (schema *Schema) getName(object map[interface{}]interface{}) error {
	id, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("schema does not have an id")
	}
	schema.schema = CreateProperty(id)
	return nil
}

func (schema *Schema) getParent(object map[interface{}]interface{}) {
	schema.parent, _ = object["parent"].(string)
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
	schema.extends = bases
	return nil
}

func (schema *Schema) addParent() error {
	if schema.parent == "" {
		return nil
	}
	set := set.New()
	set.Insert(
		CreatePropertyWithType(addName(schema.parent, "id"),
		"string"),
	)
	return schema.schema.AddProperties(set, true)
}

func (schema *Schema) Parse(object map[interface{}]interface{}) error {
	if err := schema.getName(object); err != nil {
		return err
	}
	schema.getParent(object)
	if err := schema.getBaseSchemas(object); err != nil {
		return fmt.Errorf(
			"invalid schema %s: %v",
			schema.schema.Name(),
			err,
		)
	}
	next, ok := object["schema"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"invalid schema %s: schema does not have a \"schema\"",
			schema.Name(),
		)
	}
	if err := schema.schema.Parse("", next); err != nil {
		return fmt.Errorf(
			"invalid schema %s: %v",
			schema.Name(),
			err,
		)
	}
	if !schema.schema.IsObject() {
		return fmt.Errorf(
			"invalid schema %s: schema should be an object",
			schema.Name(),
		)
	}
	err := schema.addParent()
	if err != nil {
		return fmt.Errorf("invalid schema %s: v",
			schema.Name(),
			err,
		)
	}
	return nil
}

func (schema *Schema) CollectObjects(limit, offset int) (set.Set, error) {
	result, err := schema.CollectObjects(limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid schema %s: %v",
			schema.Name(),
			err,
		)
	}
	return result, nil
}

func (schema *Schema) CollectProperties(limit, offset int) (set.Set, error) {
	result, err := schema.CollectProperties(limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid schema %s: %v",
			schema.Name(),
			err,
		)
	}
	return result, nil
}
