package main

import (
	"flag"
	"fmt"
)

//import (
//	"fmt"
//	"io/ioutil"
//
//	"strings"
//
//	"flag"
//
//	Log "github.com/Sirupsen/logrus"
//	"github.com/serenize/snaker"
//	"gopkg.in/yaml.v2"
//)
//
//type Array struct {
//	ItemType string
//	Object   Object
//}
//
//type Property struct {
//	Name   string
//	Type   string
//	Object Object
//	Array  Array
//}
//
//type Object struct {
//	Name       string
//	Properties []Property
//}
//
//type SchemaRoot struct {
//	Prefix  string
//	Extends []string
//	Schema  Object
//}
//
//var objectStore = map[string]Object{}
//
//
//
//func parseArray(name string, node map[interface{}]interface{}) Array {
//	a := Array{}
//
//	itemsNode := node["items"].(map[interface{}]interface{})
//	a.ItemType = parseType(itemsNode["type"])
//
//	if a.ItemType == "object" {
//		a.Object = parseObject(name, itemsNode)
//	}
//
//	return a
//}
//

//func parseObject(name string, obj map[interface{}]interface{}) Object {
//	o := Object{Name: name}
//	yamlProperties := obj["properties"].(map[interface{}]interface{})
//	o.Properties = parseProperties(yamlProperties)
//	for i, objProperty := range o.Properties {
//		node := yamlProperties[objProperty.Name].(map[interface{}]interface{})
//		if objProperty.Type == "object" {
//			objProperty.Object = parseObject(name+"_"+objProperty.Name, node)
//		}
//		if objProperty.Type == "item" {
//			objProperty.Array = parseArray(name+"_"+objProperty.Name, node)
//		}
//		o.Properties[i] = objProperty
//	}
//	objectStore[name] = o
//	return o
//}
//
//
//func parseProperties(obj map[interface{}]interface{}) []Property {
//	properties := []Property{}
//	for property, definition := range obj {
//		p := Property{
//			Name: property.(string),
//			Type: parseType(definition.(map[interface{}]interface{})["type"]),
//		}
//		properties = append(properties, p)
//	}
//	return properties
//}
//
//func parseSchemaRoot(s interface{}) SchemaRoot {
//	root := s.(map[interface{}]interface{})
//	yamlSchema := root["schema"].(map[interface{}]interface{})
//	schema := SchemaRoot{
//		Schema: parseObject(root["id"].(string), yamlSchema),
//	}
//	return schema
//}
//
//func generateStruct(suffix, annotation string, o Object) string {
//	code := "type " + toGoName(suffix, o.Name) + " struct {\n"
//	for _, property := range o.Properties {
//		code += "    " + toGoName("", property.Name) + " "
//		if property.Type == "item" {
//			code += "[]"
//			if property.Array.ItemType == "object" {
//				code += toGoName(suffix, property.Array.Object.Name)
//			} else {
//				code += mapType(property.Array.ItemType)
//			}
//		} else if property.Type == "object" {
//			code += toGoName(suffix, property.Object.Name)
//		} else {
//			code += mapType(property.Type)
//		}
//		code += " `" + annotation + ":\"" + property.Name + "\"`\n"
//	}
//	code += "}"
//	return code
//}
//
func main() {
	inputSchema := flag.String("schema", "", "path to yaml with schema")
	suffix := flag.String("suffix", "", "suffix appended to struct names")
	flag.Parse()

	objects, err := ReadSchemasFromFile(*inputSchema)
	if err != nil {
		panic(err)
	}


	schemas := make([]*Schema, len(objects))
	for i, object := range objects {
		schemas[i] = &Schema{}
		err := schemas[i].Parse(object)
		if err != nil {
			panic(err)
		}
	}

	toOutput := []*Object{}
	for _, schema := range schemas {
		toOutput = append(toOutput, schema.Collect(-1)...)
	}

	for _, output := range toOutput {
		fmt.Println(output.GenerateStruct(*suffix, "db"))
	}
}
