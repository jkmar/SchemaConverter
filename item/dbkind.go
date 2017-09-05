package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
	"strings"
)

// DBKind is an implementation of Kind interface
type DBKind struct {
}

// Type implementation
func (dbKind *DBKind) Type(suffix string, item Item) string {
	if item.IsNull() {
		return "goext." + util.ToGoName(
			"Null",
			strings.TrimSuffix(item.Type(suffix), "64"),
		)
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

// Annotation implementation
func (dbKind *DBKind) Annotation(name string, item Item) string {
	return fmt.Sprintf(
		"`db:\"%s\"`",
		name,
	)
}
