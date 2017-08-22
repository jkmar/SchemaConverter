package name

type Mark interface {
	Change(*string) bool
	Update(Mark)
	lengthDifference() int
}

func CreateMark(prefix string) Mark {
	return &CommonMark{
		used: false,
		begin: len(prefix),
	}
}
