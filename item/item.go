package item

import (
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

type Item interface {
	Type(string) string
	IsObject() bool
	AddProperties(set.Set, bool) error
	Parse(string, map[interface{}]interface{}) error
	CollectObjects(int, int) (set.Set, error)
	CollectProperties(int, int) (set.Set, error)
}

func CreateItem(itemType interface{}) (Item, error) {
	strType, err := util.ParseType(itemType)
	if err != nil {
		return nil, err
	}
	return CreateItemFromString(strType), nil
}

func CreateItemFromString(itemType string) Item {
	switch itemType {
	case "array":
		return &Array{}
	case "object":
		return &Object{}
	default:
		return &PlainItem{}
	}
}
