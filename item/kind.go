package item

// Kind is an interface for a type of property
// it can be either db or json
type Kind interface {
	// Type should return a go type of a property
	// args:
	//   1. string - a suffix added to a type
	//   2. item - an item of a property
	// return:
	//   go type of a property
	Type(string, Item) string

	// Type should return an interface type of a property
	// args:
	//   1. string - a suffix added to a type
	//   2. item - an item of a property
	// return:
	//   interface type of a property
	InterfaceType(string, Item) string

	// Annotation should return an annotation for a property is a go struct
	// args:
	//   1. string - a name of a property
	//   2. item - an item of a property
	// return:
	//   go annotation of a property
	Annotation(string, Item) string
}
