package defaults

import "strconv"

// FloatDefault represents float default value
type FloatDefault struct {
	value float64
}

// Write implementation
func (floatDefault *FloatDefault) Write() string {
	return strconv.FormatFloat(floatDefault.value, 'E', -1, 64)
}
