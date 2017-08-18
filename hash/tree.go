package hash

import (
	"sort"
	"fmt"
)

// TreeNode is node in hash tree
type TreeNode struct {
	index  int
	hash   *Hash
	value  Node
	item   IHashable
	parent *TreeNode
	nodes  []*TreeNode
}

// CreateTreeNode creates node from hashable item
func CreateTreeNode(item IHashable, nodes []*TreeNode) *TreeNode {
	tree := &TreeNode{
		item:  item,
		nodes: nodes,
	}
	for _, node := range nodes {
		node.parent = tree
	}
	return tree
}

// Run compresses items that have the same hash and level th parent
func (tree *TreeNode) Run(level int)  {
	tree.setup()
	tree.calcAllHashes()
	tree.compress(level)
}

func (tree *TreeNode) setup() {
	tree.setupHash(&Hash{})
}

func (tree *TreeNode) setupHash(hash *Hash) {
	tree.hash = hash
	for _, node := range tree.nodes {
		node.setupHash(hash)
	}
}

func (tree *TreeNode) getAllNodes() []*TreeNode {
	var (
		result []*TreeNode
		index  int
	)
	tree.allNodes(&index, &result)
	return result
}

func (tree *TreeNode) allNodes(index *int, result *[]*TreeNode) {
	*result = append(*result, tree)
	tree.index = *index
	*index++
	for _, node := range tree.nodes {
		node.allNodes(index, result)
	}
}

func (tree *TreeNode) calcAllHashes() {
	nodes := make([]Node, len(tree.nodes))
	for i, node := range tree.nodes {
		node.calcAllHashes()
		nodes[i] = node.value
	}
	tree.value = joinHashes(tree.hash, *tree.item, nodes)
}

func joinHashes(hash *Hash, item IHashable, nodes []Node) Node {
	result := hash.Calc(item.ToString() + "(")
	for _, node := range nodes {
		result = hash.Join(result, node)
		result = hash.Join(result, hash.Calc(","))
	}
	return hash.Join(result, hash.Calc(")"))
}

func (tree *TreeNode) compress(level int) {
	tree.getAncestors(level)
	nodes := tree.getAllNodes()
	sort.Sort(byHash(nodes))

	for _, node := range nodes {
		fmt.Printf("%p\n", *node.item)
	}

	var index int
	for index = 0; index < len(nodes) && nodes[index].parent == nil; index++ {
	}
	index++
	for ; index < len(nodes); index++ {
		if nodes[index].value == nodes[index-1].value &&
			nodes[index].parent == nodes[index-1].parent {
			fmt.Println("equal")
			*nodes[index].item = *nodes[index-1].item
		}
	}

	for _, node := range nodes {
		fmt.Printf("%p\n", *node.item)
	}
}

func (tree *TreeNode) getAncestors(level int) {
	powers := powers(level)
	last := level
	if level >= 1 {
		last = powers[len(powers)-1]
	}

	allNodes := tree.getAllNodes()
	length := len(allNodes)

	current := make([]*TreeNode, length)
	for i, node := range allNodes {
		current[i] = node.parent
	}

	levels := [][]*TreeNode{current}
	for limit := 0; limit < last && !isEmpty(levels); limit++ {
		new := make([]*TreeNode, length)
		for i, node := range current {
			if node != nil {
				new[i] = current[node.index]
			}
		}
		levels = append(levels, new)
		current = new
	}

	if last >= len(levels) {
		for _, node := range allNodes {
			node.parent = nil
		}
		return
	}

	for _, node := range allNodes {
		current := node
		for _, power := range powers {
			if current == nil {
				break
			}
			current = levels[power][current.index]
		}
		node.parent = current
	}
}

func isEmpty(array [][]*TreeNode) bool {
	for _, value := range array[len(array)-1] {
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
		if x%2 == 1 {
			result = append(result, current)
		}
		current++
		x /= 2
	}
	return result
}

// Sorting nodes by hash
type byHash []*TreeNode

func (array byHash) Len() int {
	return len(array)
}

func (array byHash) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func (array byHash) Less(i, j int) bool {
	if array[i].parent == nil {
		return true
	}
	if array[j].parent == nil {
		return false
	}
	if array[i].value.value == array[j].value.value {
		if array[i].value.length == array[j].value.length {
			return array[i].parent.index < array[j].parent.index
		}
		return array[i].value.length < array[j].value.length
	}
	return array[i].value.value < array[j].value.value

}
