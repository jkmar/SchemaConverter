package main

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/item"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

type Schema struct {
	parent  string
	extends []string
	schema  *item.Property
}

func (schema *Schema) Name() string {
	return schema.schema.Name()
}

func (schema *Schema) Bases() []string {
	return schema.extends
}

func (schema *Schema) getName(object map[interface{}]interface{}) error {
	id, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("schema does not have an id")
	}
	schema.schema = item.CreateProperty(id)
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
		item.CreatePropertyWithType(util.AddName(schema.parent, "id"),
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
	result, err := schema.schema.CollectObjects(limit, offset)
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
	result, err := schema.schema.CollectProperties(limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid schema %s: %v",
			schema.Name(),
			err,
		)
	}
	return result, nil
}

func (schema *Schema) Join(edges []*node) error {
	properties := set.New()
	for _, node := range edges {
		newProperties, err := node.schema.CollectProperties(2, 1)
		if err != nil {
			return fmt.Errorf(
				"error on joining schema %s: %v",
				schema.Name(),
				err,
			)
		}
		if err = properties.SafeInsertAll(newProperties); err != nil {
			return fmt.Errorf(
				"multiple properties with the same name in bases of schema %s",
				schema.Name(),
			)
		}
	}
	err := schema.schema.AddProperties(properties, false)
	if err != nil {
		return fmt.Errorf(
			"schema %s should be a an object",
			schema.Name(),
		)
	}
	return nil
}

func ParseAll(objects []map[interface{}]interface{}) (set.Set, error) {
	set := set.New()
	for _, object := range objects {
		schema := &Schema{}
		if err := schema.Parse(object); err != nil {
			return nil, err
		}
		if err := set.SafeInsert(schema); err != nil {
			return nil, fmt.Errorf(
				"multiple schemas with the same name: %s",
				schema.Name(),
			)
		}
	}
	return set, nil
}