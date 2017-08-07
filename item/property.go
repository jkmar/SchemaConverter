package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

// Property is a type for an item with name
type Property struct {
	name string
	item Item
}

// CreateProperty is a constructor
func CreateProperty(name string) *Property {
	return &Property{name: name}
}

// CreatePropertyWithType creates property with an item of given type
func CreatePropertyWithType(name, itemType string) *Property {
	return &Property{name, createItemFromString(itemType)}
}

// Name gets a property name
func (item *Property) Name() string {
	return item.name
}

// IsObject checks if an item in property is an object
func (item *Property) IsObject() bool {
	_, ok := item.item.(*Object)
	return ok
}

// Parse creates property from given map
// args:
//   prefix string - a prefix added to items type
//   object map[interface{}]interface{} - map from which a property is created
// return:
//   1. error during execution
func (item *Property) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	objectType, ok := object["type"]
	if !ok {
		return fmt.Errorf(
			"property %s does not have a type",
			util.AddName(prefix, item.name),
		)
	}
	item.item, err = CreateItem(objectType)
	if err != nil {
		return fmt.Errorf(
			"property %s: %v",
			util.AddName(prefix, item.name),
			err,
		)
	}
	return item.item.Parse(util.AddName(prefix, item.name), object)
}

// AddProperties adds properties to items of given property
// args:
//   set set.Set [Property] - a set of properties
//   safe bool - flag; if in the set exists a property with the same type
//               as one of the items properties, then if flag is set
//               an error should be returned,
//               otherwise that property should be ignored
// return:
//   1. error during execution
func (item *Property) AddProperties(set set.Set, safe bool) error {
	return item.item.AddProperties(set, safe)
}

// CollectObjects should return a set of objects contained within a property
// args:
//   1. int - limit; how deep to search for an object; starting from 0;
//            if limit is negative this parameter is ignored.
//   2. int - offset; from which level gathering objects should begin;
// return:
//   1. set of collected objects
//   2. error during execution
func (item *Property) CollectObjects(limit, offset int) (set.Set, error) {
	return item.item.CollectObjects(limit, offset)
}

// CollectProperties should return a set properties contained within a property
// args:
//   1. int - limit; how deep to search for a property; starting from 0;
//            if limit is negative this parameter is ignored.
//   2. int - offset; from which level gathering properties should begin;
// return:
//   1. set of collected properties
//   2. error during execution
func (item *Property) CollectProperties(limit, offset int) (set.Set, error) {
	if limit == 0 {
		return nil, nil
	}
	result, err := item.item.CollectProperties(limit-1, offset-1)
	if err != nil {
		return nil, err
	}
	if offset <= 0 {
		if result == nil {
			result = set.New()
		}
		err = result.SafeInsert(item)
		if err != nil {
			return nil, fmt.Errorf(
				"multiple properties with the same name: %s",
				item.name,
			)
		}
	}
	return result, nil
}

// GenerateProperty creates a property of a go struct from given property
// with suffix added to type name
func (item *Property) GenerateProperty(suffix, annotation string) string {
	return fmt.Sprintf(
		"\t%s %s `%s:\"%s\"`\n",
		util.ToGoName(item.name, ""),
		item.item.Type(suffix),
		annotation,
		item.name,
	)
}
