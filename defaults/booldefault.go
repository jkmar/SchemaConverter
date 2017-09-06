package defaults

import "strconv"

// BoolDefault represents bool default value
type BoolDefault struct {
	value bool
}

// Write implementation
func (boolDefault *BoolDefault) Write() string {
	return strconv.FormatBool(boolDefault.value)
}