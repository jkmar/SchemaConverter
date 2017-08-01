package reader

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func getSchemas(filename string, object []interface{}) ([]map[interface{}]interface{}, error) {
	result := make([]map[interface{}]interface{}, len(object))
	for i, item := range object {
		var ok bool
		if result[i], ok = item.(map[interface{}]interface{}); !ok {
			return nil, fmt.Errorf(
				"error in file %s: schema should have type map[interface{}]interface{}",
				filename,
			)
		}
	}
	return result, nil
}

func getSchemasFromFile(filename string) ([]interface{}, error) {
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
	return schemas, nil
}

func ReadSingle(filename string) ([]map[interface{}]interface{}, error) {
	schemas, err := getSchemasFromFile(filename)
	if err != nil {
		return nil, err
	}
	return getSchemas(filename, schemas)
}

func ReadAll(filename string) ([]map[interface{}]interface{}, error) {
	schemas, err := getSchemasFromFile(filename)
	if err != nil {
		return nil, err
	}
	result := []map[interface{}]interface{}{}
	for _, schema := range schemas {
		newFilename, ok := schema.(string)
		if !ok {
			return nil, fmt.Errorf(
				"in config file %s schemas should be filenames",
				filename,
			)
		}
		if !strings.HasPrefix(newFilename, "embed") {
			schemas, err := ReadSingle(newFilename)
			if err != nil {
				return nil, err
			}
			result = append(result, schemas...)
		}
	}
	return result, nil
}
