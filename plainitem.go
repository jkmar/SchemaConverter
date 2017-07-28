package main

import "fmt"

type PlainItem struct {
	ItemType string
}

func (item *PlainItem) Type(suffix string) string {
	return item.ItemType
}

func (item *PlainItem) IsObject() bool {
	return false
}

func (item *PlainItem) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	objectType, ok := object["type"]
	if !ok {
		return fmt.Errorf(
			"invalid schema: item %s does not have a type",
			prefix,
		)
	}
	item.ItemType, err = parseType(objectType)
	if err != nil {
		err = fmt.Errorf(
			"invalid schema: item %s - %v",
			prefix,
			err,
		)
	}
	return
}

func (item *PlainItem) Collect(depth int) []*Object {
	return nil
}