package main

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

type Array struct {
	item Item
}

func (item *Array) Type(suffix string) string {
	return "[]" + item.item.Type(suffix)
}

func (item *Array) IsObject() bool {
	return false
}

func (item *Array) AddProperties(set set.Set, safe bool) error {
	return fmt.Errorf("cannot add properties to an array")
}

func (item *Array) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	next, ok := object["items"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"array %s does not have items",
			prefix,
		)
	}
	objectType, ok := next["type"]
	if !ok {
		return fmt.Errorf(
			"items of array %s do not have a type",
			prefix,
		)
	}
	item.item, err = CreateItem(objectType)
	if err != nil {
		return fmt.Errorf("array %s: %v", prefix, err)
	}
	return item.item.Parse(prefix, next)
}

func (item *Array) CollectObjects(limit, offset int) (set.Set, error) {
	return item.item.CollectObjects(limit, offset)
}

func (item *Array) CollectProperties(limit, offset int) (set.Set, error) {
	return item.item.CollectProperties(limit, offset)
}