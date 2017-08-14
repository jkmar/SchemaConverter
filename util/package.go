package util

import (
	"fmt"
	"strings"
)

func packageName(input string) string {
	array := strings.Split(input, "/")
	return "package " + strings.TrimSuffix(array[len(array)-1], ".yaml")
}

func CollectData(name string, data []string) string {
	prefix := packageName(name) + "\n"
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
