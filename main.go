package main

import (
	"fmt"
	"io/ioutil"

	"strings"

	"flag"

	Log "github.com/Sirupsen/logrus"
	"github.com/serenize/snaker"
	"gopkg.in/yaml.v2"
)

type Array struct {
	ItemType string
	Object   Object
}

type Property struct {
	Name   string
	Type   string
	Object Object
	Array  Array
}

type Object struct {
	Name       string
	Properties []Property
}

type SchemaRoot struct {
	Prefix  string
	Extends []string
	Schema  Object
}

var objectStore = map[string]Object{}

func parseArray(name string, node map[interface{}]interface{}) Array {
	a := Array{}

	itemsNode := node["items"].(map[interface{}]interface{})
	a.ItemType = parseType(itemsNode["type"])

	if a.ItemType == "object" {
		a.Object = parseObject(name, itemsNode)
	}

	return a
}

func toGoName(suffix, name string) string {
	name = strings.Replace(name+suffix, "-", "_", -1)
	return snaker.SnakeToCamel(name)
}

func parseObject(name string, obj map[interface{}]interface{}) Object {
	o := Object{Name: name}
	yamlProperties := obj["properties"].(map[interface{}]interface{})
	o.Properties = parseProperties(yamlProperties)
	for i, objProperty := range o.Properties {
		node := yamlProperties[objProperty.Name].(map[interface{}]interface{})
		if objProperty.Type == "object" {
			objProperty.Object = parseObject(name+"_"+objProperty.Name, node)
		}
		if objProperty.Type == "array" {
			objProperty.Array = parseArray(name+"_"+objProperty.Name, node)
		}
		o.Properties[i] = objProperty
	}
	objectStore[name] = o
	return o
}

func parseType(t interface{}) string {
	switch v := t.(type) {
	case string:
		return v
	case []interface{}:
		for _, item := range v {
			strItem := item.(string)
			if strItem != "null" {
				return strItem
			}
		}

	}
	panic(fmt.Sprintf("unsupported type: %T", t))
}

func parseProperties(obj map[interface{}]interface{}) []Property {
	properties := []Property{}
	for property, definition := range obj {
		p := Property{
			Name: property.(string),
			Type: parseType(definition.(map[interface{}]interface{})["type"]),
		}
		properties = append(properties, p)
	}
	return properties
}

func parseSchemaRoot(s interface{}) SchemaRoot {
	root := s.(map[interface{}]interface{})
	yamlSchema := root["schema"].(map[interface{}]interface{})
	schema := SchemaRoot{
		Schema: parseObject(root["id"].(string), yamlSchema),
	}
	return schema
}

func generateStruct(suffix, annotation string, o Object) string {
	code := "type " + toGoName(suffix, o.Name) + " struct {\n"
	for _, property := range o.Properties {
		code += "    " + toGoName("", property.Name) + " "
		if property.Type == "array" {
			code += "[]"
			if property.Array.ItemType == "object" {
				code += toGoName(suffix, property.Array.Object.Name)
			} else {
				code += property.Array.ItemType
			}
		} else if property.Type == "object" {
			code += toGoName(suffix, property.Object.Name)
		} else {
			code += property.Type
		}
		code += " `" + annotation + ":\"" + property.Name + "\"`\n"
	}
	code += "}"
	return code
}

func main() {
	inputSchema := flag.String("schema", "", "path to yaml with schema")
	suffix := flag.String("suffix", "", "suffix appended to struct names")
	flag.Parse()

	if *inputSchema == "" {
		Log.Info("Missing input schema")
		return
	}

	inputContent, err := ioutil.ReadFile(*inputSchema)
	if err != nil {
		panic(fmt.Sprintf("Failed to open %s file", *inputSchema))
	}
	schema := map[string]interface{}{}

	err = yaml.Unmarshal(inputContent, &schema)
	if err != nil {
		panic("Cannot parse given schema")
	}

	for _, schema := range schema["schemas"].([]interface{}) {
		root := parseSchemaRoot(schema)
		fmt.Println(generateStruct(*suffix, "db", root.Schema))
	}

	for _, obj := range objectStore {
		fmt.Println(generateStruct(*suffix, "json", obj))
	}

}
