package main

import (
	"fmt"
)

type Property struct {
	Name string
	Item Item
}

func CreateProperty(name string) *Property {
	return &Property{Name: name}
}

func CreateItem(itemType interface{}) Item {
	return CreateItemFromString(parseType(itemType))
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

type Schema struct {
	Parent string
	Extends []string
	Schema *Property
}

func (schema *Schema) Name() string {
	return schema.Schema.Name
}







func (item *Property) Parse(prefix string, object map[interface{}]interface{}) {
	item.Item = CreateItem(object["type"])
	item.Parse(prefix+item.Name, object)
}

func (item *Property) GenerateProperty(sufix, annotation string) string {
	return fmt.Sprintf(
		"\t%s %s `%s:\"%s\"`\n",
		toGoName(item.Name, ""),
		toGoName(item.Item.Type(), sufix),
		sufix,
		annotation,
		item.Name,
	)
}



func (schema *Schema) Parse(object map[interface{}]interface{}) {
	parent, ok := object["parent"].(string)
	if (ok) {
		schema.Parent = parent
	}
	extends, ok := object["extends"].([]string)
	if (ok) {
		schema.Extends = extends
	}
	schema.Schema = CreateProperty(object["id"].(string))
	schema.Schema.Parse("", object["schema"].(map[interface{}]interface{}))
}
