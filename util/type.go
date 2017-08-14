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

func TryToAddName(prefix, suffix string) string {
	if prefix == "" {
		return ""
	}
	return AddName(prefix, suffix)
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
func ParseType(itemType interface{}) (string, bool, error) {
	switch goType := itemType.(type) {
	case string:
		return mapType(goType), false, nil
	case []interface{}:
		var (
			result string
			null   bool
		)
		for _, singleType := range goType {
			if strType, ok := singleType.(string); ok {
				if strType == "null" {
					null = true
				} else if result == "" {
					result = mapType(strType)
				}
			}
		}
		if result != "" {
			return result, null, nil
		}

	}
	return "", false, fmt.Errorf("unsupported type: %T", itemType)
}
