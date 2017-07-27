package main

import "fmt"

type Array struct {
	Item Item
}

func (item *Array) Type() string {
	return "[]" + item.Item.Type()
}

func (item *Array) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	next, ok := object["items"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"invalid schema: array %s does not have items",
			prefix,
		)
	}
	objectType, ok := next["type"]
	if !ok {
		return fmt.Errorf(
			"invalid schema: items of array %s do not have a type",
			prefix,
		)
	}
	item.Item, err = CreateItem(objectType)
	if err != nil {
		return fmt.Errorf("invalid schema: array %s - %v", prefix, err)
	}
	return item.Item.Parse(prefix, next)
}
