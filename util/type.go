package util

import (
	"fmt"
	"github.com/serenize/snaker"
	"strings"
)

var typeMapping = map[string]string{
	"integer":  "int64",
	"number":   "int64",
	"boolean":  "bool",
	"abstract": "object",
}

// AddName creates a snake case name from prefix and suffix
func AddName(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return prefix + "_" + suffix
}

// ToGoName creates a camel case name from prefix and suffix
func ToGoName(prefix, suffix string) string {
	name := strings.Replace(AddName(prefix, suffix), "-", "_", -1)
	return snaker.SnakeToCamel(name)
}

func mapType(typeName string) string {
	if mappedName, ok := typeMapping[typeName]; ok {
		return mappedName
	}
	return typeName
}

// ParseType converts an interface to a name of a go type
func ParseType(itemType interface{}) (string, error) {
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
