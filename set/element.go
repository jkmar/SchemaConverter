package set

type Element interface {
	Name() string
}

type byName []Element

func (array byName) Len() int {
	return len(array)
}

func (array byName) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func (array byName) Less(i, j int) bool {
	return array[i].Name() < array[j].Name()
}
