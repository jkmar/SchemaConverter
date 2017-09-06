package item

// ParseContext represents context used in parsing
type ParseContext struct {
	prefix string
	level int
	required bool
	defaults interface{}
	data map[interface{}]interface{}
}