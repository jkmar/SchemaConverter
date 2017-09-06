package defaults

import "strconv"

// IntDefault represents int default value
type IntDefault struct {
	value int
}

// Write implementation
func (intDefault *IntDefault) Write() string {
	return strconv.Itoa(intDefault.value)
}
