/*
package main

import (
	"strings"
	"github.com/serenize/snaker"
	"fmt"
)

func CreateProperty(name string) *Property {
	return &Property{Name: name}
}

func CreateItem(itemType interface{}) Item {
	return CreateItemFromString(parseType(itemType))
}

func CreateItemFromString(itemType string) Item {
	switch itemType {
	case "item":
		return &Array{}
	case "object":
		return &Object{}
	default:
		return &PlainItem{}
	}
}

type Item interface {
	Type() string
	Parse(map[interface{}]interface{})
	//Collect(string, map[string]*Property)
	GenerateStruct(string) string
	GenerateProperty(string, string) string
}

type Property struct {
	FullName string
	Name string
	Item Item
}

func (item *Property) Type() string {
	return item.FullName
}

func (item *Property) Parse(object map[interface{}]interface{}) {
	item.Item = CreateItem(object["type"])
	item.Parse(object)
}

func (item *Property) GenerateStruct(annotation,  string) string {
	if item.Item.Type() != "object" {
		return ""
	}
	code := "type" + toGoName(item.Type()) + " struct {\n"

}

func (item *Property) GenerateProperty(fullName, annotation string) string {
	return fmt.Sprintf(
		"\t%s %s `%s:\"%s\"`\n",
		item.Name,
		item.Item.GenerateProperty(fullName + item.Type(), annotation),
		annotation,
		item.Name,
	)
}

func (item *PlainItem) GenerateProperty(fullName, annotation string) string {
	return item.Type()
}

func (item *Array) GenerateProperty(fullName, annotation string) string {
	return "[]" + item.Item.GenerateProperty(fullName, annotation)
}

func (item *Object) GenerateProperty(fullName, annotation string) string {
	return fullName
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

func (object *Object) Type() string {
	return "object"
}

/*
func (item Object) Parse(object map[interface{}]interface{}) {
	properties := object["properties"].(map[interface{}]interface{})
	item.Properties = parseProperties(object["properties"].(map[interface{}]interface{}))
	for _, property := range item.Properties {
		property.Item.Parse(properties[property.Name].(map[interface{}]interface{}))
	}
}
*/

/*
func (item *Object) Parse(object map[interface{}]interface{}) {
	next := object["properties"].(map[interface{}]interface{})
	item.Properties = []*Property{}
	for property, definition := range next {
		item.Properties = append(item.Properties, CreateProperty(property.(string)))
		item.Properties[len(item.Properties)-1].Parse(definition.(map[interface{}]interface{}))
	}
}

type Array struct {
	Item Item
}

func (array *Array) Type() string {
	return "item"
}

func (array *Array) Parse(object map[interface{}]interface{}) {
	next := object["items"].(map[interface{}]interface{})
	array.Item = CreateItem(next["type"])
	array.Item.Parse(next)
}

func toGoName(suffix, name string) string {
	name = strings.Replace(name+suffix, "-", "_", -1)
	return snaker.SnakeToCamel(name)
}

