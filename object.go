package main

import "fmt"

type Object struct {
	objectType string
	Properties []*Property
}

func (item *Object) Type() string {
	return item.objectType
}

func (item *Object) Parse(prefix string, object map[interface{}]interface{}) (err error) {
	next, ok := object["properties"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf(
			"invalid schema: object %s does not have properties",
			prefix,
		)
	}
	item.objectType = prefix
	item.Properties = []*Property{}
	for property, definition := range next {
		item.Properties = append(item.Properties, CreateProperty(property.(string)))
		definitionMap, ok := definition.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf(
				"invalid schema: object %s has invalid property %s",
				item.objectType,
				property.(string),
			)
		}
		err = item.Properties[len(item.Properties)-1].Parse(prefix, definitionMap)
		if err != nil {
			return
		}
	}
	return
}

//func (item *Object) GenerateStruct(sufix, annotation string) string {
//	code := "type " + toGoName(item.Type(), sufix) + " struct {\n"
//	for _, propery := range item.Properties {
//		code += propery.GenerateProperty(sufix, annotation)
//	}
//	return code + "}\n"
//}
