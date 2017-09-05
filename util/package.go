package util

import (
	"fmt"
)

// CollectData creates string representing go file
func CollectData(name string, data []string) string {
	prefix := "package " + name + "\n"
	result := prefix
	for _, element := range data {
		if element != "" {
			result = fmt.Sprintf("%s\n%s", result, element)
		}
	}
	if result == prefix {
		return ""
	}
	return result
}
