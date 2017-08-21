package name

import "strings"

const new = "common"

type CommonMark struct {
	used  bool
	begin int
	end   int
}

func (commonMark *CommonMark) Change(string *string) bool {
	if strings.HasPrefix((*string)[commonMark.begin:], new) {
		return false
	}
	result := (*string)[:commonMark.begin] + new
	if commonMark.used {
		result += (*string)[commonMark.end:]
	} else {
		commonMark.used = true
		commonMark.end = len(*string)
	}
	*string = result
	return true
}

func (commonMark *CommonMark) Update(mark Mark) {
	difference := mark.lengthDifference()
	commonMark.begin += difference
	commonMark.end += difference
}

func (commonMark *CommonMark) lengthDifference() int {
	if commonMark.used {
		return len(new) - commonMark.end + commonMark.begin
	}
	return 0
}
