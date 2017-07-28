package main

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func getSchemas(filename string, object []interface{}) ([]map[interface{}]interface{}, error) {
	result := make([]map[interface{}]interface{}, len(object))
	for i, item := range object {
		var ok bool
		result[i], ok = item.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"error in file %s: schema should have type map[interface{}]interface{}",
				filename,
			)
		}
	}
	return result, nil
}

func ReadSchemasFromFile(filename string) ([]map[interface{}]interface{}, error) {
	inputContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open file %s",
			filename,
		)
	}
	object := map[interface{}]interface{}{}
	err = yaml.Unmarshal(inputContent, &object)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot parse given schema from file %s",
			filename,
		)
	}
	schemas, ok := object["schemas"].([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"no schemas found in file %s",
			filename,
		)
	}
	return getSchemas(filename, schemas)
}
