package item

import (
	"fmt"
	"strings"

	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

// DBKind is an implementation of Kind interface
type DBKind struct {
}

// Type implementation
func (dbKind *DBKind) Type(suffix string, item Item) string {
	if item.IsNull() {
		return "goext." + getNullType(suffix, item)
	}
	return item.Type(suffix)
}

// InterfaceType implementation
func (dbKind *DBKind) InterfaceType(suffix string, item Item) string {
	if item.IsNull() {
		return dbKind.Type(suffix, item)
	}
	return item.InterfaceType(suffix)
}

func dbAnnotation(name string, item Item) string {
	return fmt.Sprintf("db:\"%s\"", name)
}

// Annotation implementation
func (dbKind *DBKind) Annotation(name string, item Item) string {
	return fmt.Sprintf(
		"`%s %s`",
		dbAnnotation(name, item),
		jsonAnnotation(name, item),
	)
}

// Default implementation
func (dbKind *DBKind) Default(suffix string, item Item) string {
	if item.IsNull() {
		return fmt.Sprintf(
			"goext.Make%s(%s)",
			getNullType(suffix, item),
			item.Default(suffix),
		)
	}
	return item.Default(suffix)
}

func getNullType(suffix string, item Item) string {
	return util.ToGoName(
		"Null",
		strings.TrimSuffix(item.Type(suffix), "64"),
	)
}
