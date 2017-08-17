package hash

type TreeNode struct {
	index int
	hash  *Hash
	value Node
	item  IHashable
	parent *TreeNode
	nodes []*TreeNode
}

func CreateTreeNode(item IHashable, nodes []*TreeNode) *TreeNode {
	tree := &TreeNode{
		item: item,
		nodes: nodes,
	}
	for _, node := range nodes {
		node.parent = tree
	}
	return tree
}

func (tree *TreeNode) Setup() {
	var index int
	tree.setup(&Hash{}, &index)
}

func (tree *TreeNode) setup(hash *Hash, index *int) {
	tree.hash = hash
	tree.index = *index
	*index++
	for _, node := range tree.nodes {
		node.setup(hash, index)
	}
}

func (tree *TreeNode) Compress() {
	compressNodes(tree.nodes)
}

func (tree *TreeNode) CalcAllHashes() {
	nodes := make([]Node, len(tree.nodes))
	for i, node := range tree.nodes {
		node.CalcAllHashes()
		nodes[i] = node.value
	}
	tree.value = joinHashes(tree.hash, tree.item, nodes)
}

func compressNodes(nodes []*TreeNode) {
	visited := make([]bool, len(nodes))
	for i := 1; i < len(visited); i++ {
		if nodes[i].value == nodes[i-1].value {
			visited[i] = true
			nodes[i].item = nodes[i-1].item
		}
	}
	for i, flag := range visited {
		if !flag {
			nodes[i].Compress()
		}
	}
}

func joinHashes(hash *Hash, item IHashable, nodes []Node) Node {
	result := hash.Calc(item.ToString() + "(")
	for _, node := range nodes {
		result = hash.Join(result, node)
		result = hash.Join(result, hash.Calc(","))
	}
	return hash.Join(result, hash.Calc(")"))
}

func (tree *TreeNode) last() int {
	if len(tree.nodes) == 0 {
		return tree.index
	}
	return tree.nodes[len(tree.nodes)-1].last()
}

func isEmpty(array []interface{}) bool {
	for _, value := range array {
		if value != nil {
			return false
		}
	}
	return true
}

func powers(x int) []int {
	result := []int{}
	current := 0
	for x > 0 {
		if x % 2 == 1 {
			result = append(result, current)
		}
		current++
		x /= 2
	}
	return result
}

func (tree *TreeNode) getAncestors(level int) {
	length := tree.last()
	current := make([]*TreeNode, length)
	for i, node := range allNodes {
		current[i] = node.parent
	}
	levels := [][]*TreeNode{}
	levels = append(levels, current)
	for {
		new := make([]*TreeNode, length)
		for i, node := range current {
			if node != nil {
				new[i] = current[node.index]
			}
		}
	}
}