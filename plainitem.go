package main

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

type PlainItem struct {
	itemType string
}

func (item *PlainItem) Type(suffix string) string {
	return item.itemType
}

func (item *PlainItem) IsObject() bool {
	return false
}

func (item *PlainItem) AddProperties(set set.Set, safe bool) error {
	return fmt.Errorf("cannot add properties to a plain item")
}

func (item *PlainItem) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	objectType, ok := object["type"]
	if !ok {
		return fmt.Errorf(
			"item %s does not have a type",
			prefix,
		)
	}
	item.itemType, err = parseType(objectType)
	if err != nil {
		err = fmt.Errorf(
			"item %s: %v",
			prefix,
			err,
		)
	}
	return
}

func (item *PlainItem) CollectObjects(limit, offset int) (set.Set, error) {
	return nil, nil
}

func (item *PlainItem) CollectProperties(limit, offset int) (set.Set, error) {
	return nil, nil
}