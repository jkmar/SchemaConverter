package set

import (
	"fmt"
	"sort"
)

type Set map[string]Element

func New() Set {
	return make(map[string]Element)
}

func (set Set) Empty() bool {
	return set == nil || len(set) == 0
}

func (set Set) Size() int {
	if set == nil {
		return 0
	}
	return len(set)
}

func (set Set) Contains(element Element) bool {
	if set == nil {
		return false
	}
	_, ok := set[element.Name()]
	return ok
}

func (set Set) Delete(element Element) {
	if set != nil {
		delete(set, element.Name())
	}
}

func (set Set) Insert(element Element) {
	if set != nil {
		set[element.Name()] = element
	}
}

func (set Set) InsertAll(other Set) {
	if !other.Empty() {
		for _, value := range other {
			set.Insert(value)
		}
	}
}

func (set *Set) SafeInsert(element Element) error {
	if set.Contains(element) {
		return fmt.Errorf(
			"the element with the name %s already in the set",
			element.Name(),
		)
	}
	set.Insert(element)
	return nil
}

func (set Set) SafeInsertAll(other Set) error {
	for _, value := range other {
		if set.Contains(value) {
			return fmt.Errorf(
				"the element with the name %s already in the set",
				value.Name(),
			)
		}
	}
	set.InsertAll(other)
	return nil
}

func (set Set) ToArray() []Element {
	if set == nil {
		return nil
	}
	result := make([]Element, len(set))
	i := 0
	for _, value := range set {
		result[i] = value
		i++
	}
	sort.Sort(byName(result))
	return result
}