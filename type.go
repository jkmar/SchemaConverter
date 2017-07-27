package main

import (
	"fmt"
	"strings"
	"github.com/serenize/snaker"
)

var typeMapping = map[string]string{
	"integer": "int64",
	"number":  "int64",
	"boolean": "bool",
}

func addName(prefix, sufix string) string {
	if prefix == "" {
		return sufix
	}
	return prefix + "_" + sufix
}

func toGoName(prefix, suffix string) string {
	name := strings.Replace(prefix+suffix, "-", "_", -1)
	return snaker.SnakeToCamel(name)
}

func mapType(typeName string) string {
	if mappedName, ok := typeMapping[typeName]; ok {
		return mappedName
	}
	return typeName
}

func parseType(itemType interface{}) (string, error) {
	switch goType := itemType.(type) {
	case string:
		return mapType(goType), nil
	case []interface{}:
		for _, item := range goType {
			if strItem, ok := item.(string); ok && strItem != "null" {
				return mapType(strItem), nil
			}
		}

	}
	return "", fmt.Errorf("unsupported type: %T", itemType)
}
