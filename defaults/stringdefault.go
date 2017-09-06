package defaults

// StringDefault represents string default value
type StringDefault struct {
	value string
}

// Write implementation
func (stringDefault *StringDefault) Write() string {
	return `"` + stringDefault.value + `"`
}
