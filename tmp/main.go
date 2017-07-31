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
//	itemType string
//	Object   Object
//}
//
//type Property struct {
//	name   string
//	Type   string
//	Object Object
//	Array  Array
//}
//
//type Object struct {
//	name       string
//	properties []Property
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
//	a.itemType = parseType(itemsNode["type"])
//
//	if a.itemType == "object" {
//		a.Object = parseObject(name, itemsNode)
//	}
//
//	return a
//}
//

//func parseObject(name string, obj map[interface{}]interface{}) Object {
//	o := Object{name: name}
//	yamlProperties := obj["properties"].(map[interface{}]interface{})
//	o.properties = parseProperties(yamlProperties)
//	for i, objProperty := range o.properties {
//		node := yamlProperties[objProperty.name].(map[interface{}]interface{})
//		if objProperty.Type == "object" {
//			objProperty.Object = parseObject(name+"_"+objProperty.name, node)
//		}
//		if objProperty.Type == "item" {
//			objProperty.Array = parseArray(name+"_"+objProperty.name, node)
//		}
//		o.properties[i] = objProperty
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
//			name: property.(string),
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
//	code := "type " + toGoName(suffix, o.name) + " struct {\n"
//	for _, property := range o.properties {
//		code += "    " + toGoName("", property.name) + " "
//		if property.Type == "item" {
//			code += "[]"
//			if property.Array.itemType == "object" {
//				code += toGoName(suffix, property.Array.Object.name)
//			} else {
//				code += mapType(property.Array.itemType)
//			}
//		} else if property.Type == "object" {
//			code += toGoName(suffix, property.Object.name)
//		} else {
//			code += mapType(property.Type)
//		}
//		code += " `" + annotation + ":\"" + property.name + "\"`\n"
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
