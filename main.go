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
//	schema  Object
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
//	a.itemType = ParseType(itemsNode["type"])
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
//			Type: ParseType(definition.(map[interface{}]interface{})["type"]),
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
//		schema: parseObject(root["id"].(string), yamlSchema),
//	}
//	return schema
//}
//
//func generateStruct(suffix, annotation string, o Object) string {
//	code := "type " + ToGoName(suffix, o.name) + " struct {\n"
//	for _, property := range o.properties {
//		code += "    " + ToGoName("", property.name) + " "
//		if property.Type == "item" {
//			code += "[]"
//			if property.Array.itemType == "object" {
//				code += ToGoName(suffix, property.Array.Object.name)
//			} else {
//				code += mapType(property.Array.itemType)
//			}
//		} else if property.Type == "object" {
//			code += ToGoName(suffix, property.Object.name)
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

	tmp, err := Convert(nil, objects, "db", "json", *suffix)
	if err != nil {
		panic(err)
	}
	for _, x := range tmp {
		fmt.Println(x)
	}
}
