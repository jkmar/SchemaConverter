package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

type Property struct {
	name string
	item Item
}

func CreateProperty(name string) *Property {
	return &Property{name: name}
}

func CreatePropertyWithType(name, itemType string) *Property {
	return &Property{name, CreateItemFromString(itemType)}
}

func (item *Property) Name() string {
	return item.name
}

func (item *Property) IsObject() bool {
	return item.item.IsObject()
}

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

func (item *Property) AddProperties(set set.Set, safe bool) error {
	return item.item.AddProperties(set, safe)
}

func (item *Property) CollectObjects(limit, offset int) (set.Set, error) {
	return item.item.CollectObjects(limit, offset)
}

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

func (item *Property) GenerateProperty(suffix, annotation string) string {
	return fmt.Sprintf(
		"\t%s %s `%s:\"%s\"`\n",
		util.ToGoName(item.name, ""),
		item.item.Type(suffix),
		annotation,
		item.name,
	)
}
