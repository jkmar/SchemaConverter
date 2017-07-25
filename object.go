package main

import (
	"strings"
	"github.com/serenize/snaker"
)

func CreateProperty(name, itemType string) *Property {
	return &Property{
		name,
		CreateItem(itemType),
	}
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
		return PlainItem{}
	}
}

type Item interface {
	Type() string
	Parse(map[interface{}]interface{})
	Collect(string, map[string]*Property)
}

type Property struct {
	Name string
	Item Item
}

type PlainItem struct {
	ItemType string
}

func (item *PlainItem) Type() string {
	return item.ItemType
}

func (item *PlainItem) Parse(object map[interface{}]interface{}) {
	item.ItemType = object["type"].(string)
}

func (item *PlainItem)  Collect(prefix string, objects map[string]*Property) {
}

type Object struct {
	Properties []*Property
}

func (object Object) Type() string {
	return "object"
}

func (item Object) Parse(object map[interface{}]interface{}) {
	properties := object["properties"].(map[interface{}]interface{})
	item.Properties = parseProperties(object["properties"].(map[interface{}]interface{}))
	for _, property := range item.Properties {
		property.Item.Parse(properties[property.Name].(map[interface{}]interface{}))
	}
}

func (object Object) Collect(prefix string, objects map[string]*Property) {
	for _, property := range object.Properties {

	}
}

func parseProperties(object map[interface{}]interface{}) []*Property {
	properties := []*Property{}
	for property, definition := range object {
		itemType := parseType(definition.(map[interface{}]interface{})["type"])
		properties = append(properties, CreateProperty(property.(string), itemType))
	}
	return properties
}

type Array struct {
	Item Item
}

func (array Array) Type() string {
	return "array"
}

func (array Array) Parse(object map[interface{}]interface{}) {
	items := object["items"].(map[interface{}]interface{})
	array.Item = CreateItem(parseType(items["type"]))
	array.Item.Parse(items)
}

func (property *Property) Collect(prefix string, objects map[string]interface{}) {
	objects[prefix+property.Name] = property.Item
	property.Item.Collect(prefix+property.Name, objects)
}

func toGoName(suffix, name string) string {
	name = strings.Replace(name+suffix, "-", "_", -1)
	return snaker.SnakeToCamel(name)
}

type ByName []*Property

func (array ByName) Len() int {
	return len(array)
}

func (array ByName) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func (array ByName) Less(i, j int) bool {
	return array[i].Name < array[j].Name
}