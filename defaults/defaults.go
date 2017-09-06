package defaults

// PlainDefaults is an interface for default values
type PlainDefaults interface {
	Write() string
}

// CreatePlainDefaults is a PlainDefaults factory
func CreatePlainDefaults(value interface{}) PlainDefaults {
	switch value.(type) {
	case string:
		return &StringDefault{value: value.(string)}
	case int:
		return &IntDefault{value: value.(int)}
	case bool:
		return &BoolDefault{value: value.(bool)}
	case float64:
		return &FloatDefault{value: value.(float64)}
	}
	return nil
}
