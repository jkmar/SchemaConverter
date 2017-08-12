package util

import (
	"github.com/serenize/snaker"
	"strings"
)

var keywords = []string{
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"goto",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",
}

// VariableName gets a variable name from its type
func VariableName(name string) string {
	result := snaker.SnakeToCamelLower(name)
	for _, keyword := range keywords {
		if result == keyword {
			return result + "Object"
		}
	}
	return result
}

// IndexVariable returns a name of variable used in for loop
func IndexVariable(depth int) rune {
	return rune('i' + depth - 1)
}

// Indent returns an indent with given width
func Indent(width int) string {
	return strings.Repeat("\t", width)
}
