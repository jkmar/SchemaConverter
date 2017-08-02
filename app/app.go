package app

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/reader"
	"github.com/zimnx/YamlSchemaToGoStruct/schema"
	"io/ioutil"
	"strings"
)

func readConfig(config, input string) ([]map[interface{}]interface{}, error) {
	if len(config) == 0 {
		return nil, nil
	}
	return reader.ReadAll(config, input)
}

func writeResult(structs []string, input, output string) error {
	result := "package " + strings.TrimSuffix(input, ".yaml")
	for _, goStruct := range structs {
		result = fmt.Sprintf("%s\n\n%s", result, goStruct)
	}
	if len(output) == 0 {
		fmt.Print(result)
		return nil
	}
	return ioutil.WriteFile(output, []byte(result), 0644)
}

func Run(input, output, config, db, json, suffix string) error {
	other, err := readConfig(config, input)
	if err != nil {
		return err
	}
	objects, err := reader.ReadSingle(input)
	if err != nil {
		return err
	}
	result, err := schema.Convert(other, objects, db, json, suffix)
	if err != nil {
		return err
	}
	return writeResult(result, input, output)
}
