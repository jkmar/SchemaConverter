package item

// ParseContext represents context used in parsing
type ParseContext struct {
	Prefix   string
	Level    int
	Required bool
	Defaults interface{}
	Data     map[interface{}]interface{}
}
