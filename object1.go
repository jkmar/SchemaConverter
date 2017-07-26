package main

type Object struct {
	objectType string
	Properties []*Property
}

func (item *Object) Type() string {
	return item.objectType
}

func (item *Object) Parse(prefix string, object map[interface{}]interface{}) {
	next := object["properties"].(map[interface{}]interface{})
	item.objectType = prefix
	item.Properties = []*Property{}
	for property, definition := range next {
		item.Properties = append(item.Properties, CreateProperty(property.(string)))
		item.Properties[len(item.Properties)-1].Parse(prefix, definition.(map[interface{}]interface{}))
	}
}

func (item *Object) GenerateStruct(sufix, annotation string) string {
	code := "type " + toGoName(item.Type(), sufix) + " struct {\n"
	for _, propery := range item.Properties {
		code += propery.GenerateProperty(sufix, annotation)
	}
	return code + "}\n"
}