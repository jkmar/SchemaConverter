package item

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

// DBKind is an implementation of Kind interface
type DBKind struct {
}

// Type implementation
func (dbKind *DBKind) Type(suffix string, item Item) string {
	if item.IsNull() {
		return "sql." + util.ToGoName("Null", item.Type(suffix))
	}
	return item.Type(suffix)
}

// Annotation implementation
func (dbKind *DBKind) Annotation(name string, item Item) string {
	return fmt.Sprintf(
		"`db:\"%s\"`",
		name,
	)
}
