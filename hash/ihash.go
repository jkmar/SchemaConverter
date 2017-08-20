package hash

// IHashable represent tree node can be hashed
type IHashable interface {
	// ToString should return string representation of a value of a node
	// return:
	//   string identifying a value of a node
	ToString() string

	// Compress should compress child of a node
	// args:
	//   1. IHashable - source (node to compress the other node to)
	//   2. IHashable - destination (node to compress)
	Compress(IHashable, IHashable)

	// GetChildren should return children of a node
	// return:
	//   array of children of a node
	GetChildren() []IHashable
}
