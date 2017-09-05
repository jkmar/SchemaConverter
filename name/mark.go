package name

// Mark is a type used to change names
type Mark interface {
	// Change should change string
	// args:
	//   1. *string - string; a string to be changed
	// return:
	//   true iff. string was changed
	Change(*string) bool

	// Update should update mark based on other mark
	// args:
	//   1. Mark - mark; a mark on which update is based
	Update(Mark)

	// lengthDifference should return a difference of
	// lengths between string before and after the change
	// return:
	//   difference of lengths
	lengthDifference() int
}

// CreateMark is a factory for Mark interface
func CreateMark(prefix string) Mark {
	return &CommonMark{
		used:  false,
		begin: len(prefix),
	}
}
