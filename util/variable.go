package util

import (
	"unicode"
	"unicode/utf8"
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

// ToLowerFirst makes the first character of a given string lowercase
func ToLowerFirst(text string) string {
	if text == "" {
		return ""
	}
	rune, pos := utf8.DecodeRuneInString(text)
	return string(unicode.ToLower(rune)) + text[pos:]
}

// VariableName gets a variable name from its type
func VariableName(name string) string {
	result := ToLowerFirst(name)
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
