package main

import "fmt"

type Property struct {
	Name string
	Item Item
}

func CreateProperty(name string) *Property {
	return &Property{Name: name}
}

func (item *Property) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	objectType, ok := object["type"]
	if !ok {
		return fmt.Errorf(
			"invalid schema: property %s does not have a type",
			addName(prefix, item.Name),
		)
	}
	item.Item, err = CreateItem(objectType)
	if err != nil {
		err = fmt.Errorf(
			"invalid schema: property %s - %v",
			addName(prefix, item.Name),
			err,
		)
		return
	}
	return item.Item.Parse(addName(prefix, item.Name), object)
}

func (item *Property) Collect(depth int) []*Object {
	return item.Item.Collect(depth)
}

func (item *Property) GenerateProperty(suffix, annotation string) string {
	return fmt.Sprintf(
		"\t%s %s `%s:\"%s\"`\n",
		toGoName(item.Name, ""),
		item.Item.Type(suffix),
		annotation,
		item.Name,
	)
}
