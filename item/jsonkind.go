package item

import "fmt"

// JSONKind is an implementation of Kind interface
type JSONKind struct {
}

// Type implementation
func (jsonKind *JSONKind) Type(suffix string, item Item) string {
	return item.Type(suffix)
}

// Annotation implementation
func (jsonKind *JSONKind) Annotation(name string, item Item) string {
	var annotation string
	if item.IsNull() {
		annotation = ",omitempty"
	}
	return fmt.Sprintf(
		"`json:\"%s%s\"`",
		name,
		annotation,
	)
}
