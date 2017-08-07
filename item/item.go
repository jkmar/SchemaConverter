package item

import (
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

// Item is an interface for a type of a variable
type Item interface {
	// Type should return go type of item
	// args:
	//   1. string - a suffix added to type
	// return:
	//   type of item with suffix appended
	Type(string) string

	// AddProperties should add properties to an item
	// args:
	//   1. set.Set [Property] - a set of properties
	//   2. bool - flag; if in the set exists a property with the same type
	//             as one of the items properties, then if flag is set
	//             an error should be returned,
	//             otherwise that property should be ignored
	// return:
	//   1. error during execution
	AddProperties(set.Set, bool) error

	// Parse should create an item from given map
	// args:
	//   1. string - a prefix added to items type
	//   2. map[interface{}]interface{} - map from which an item is created
	// return:
	//   1. error during execution
	Parse(string, map[interface{}]interface{}) error

	// CollectObjects should return a set of objects contained within an item
	// args:
	//   1. int - limit; how deep to search for an object; starting from 0;
	//            if limit is negative this parameter is ignored.
	//   2. int - offset; from which level gathering objects should begin;
	// return:
	//   1. set of collected objects
	//   2. error during execution
	// example:
	//   let objects be denoted by o and other items by i
	//   suppose we have the following tree:
	//             o1
	//            / \
	//           o2  o3
	//          /  \   \
	//        o4   o5   o6
	//        / \   \    \
	//       o7  i1  i2   i3
	//
	// CollectObjects(3, 1) should return a set of o2, o3, o4, o5, o6
	// CollectObjects(2, 2) should return an empty set
	// CollectObjects(-1, 4) should return a set of o7
	CollectObjects(int, int) (set.Set, error)

	// CollectProperties should return a set properties contained within an item
	// args:
	//   1. int - limit; how deep to search for a property; starting from 0;
	//            if limit is negative this parameter is ignored.
	//   2. int - offset; from which level gathering properties should begin;
	// return:
	//   1. set of collected properties
	//   2. error during execution
	CollectProperties(int, int) (set.Set, error)
}

// CreateItem is a factory for items
func CreateItem(itemType interface{}) (Item, error) {
	strType, err := util.ParseType(itemType)
	if err != nil {
		return nil, err
	}
	return createItemFromString(strType), nil
}

func createItemFromString(itemType string) Item {
	switch itemType {
	case "array":
		return &Array{}
	case "object":
		return &Object{}
	default:
		return &PlainItem{}
	}
}
