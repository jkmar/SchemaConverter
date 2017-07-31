package main

import "fmt"

type Node struct {
	Schema *Schema
	Edges []*Node
	mark int
}

var allSchemas map[string]*Schema

func sort(schemas []*Schema) ([]*Node, error) {
	nodes := map[string]*Node{}

	getNode := func(id string) *Node {
		if node, ok := nodes[id]; ok {
			return node
		}
		node := &Node{Schema: allSchemas[id]}
		nodes[id] = node
		return node
	}

	getNeighbours := func(schema *Schema) []*Node {
		result := make([]*Node, len(schema.extends))
		for i, name := range schema.extends {
			result[i] = getNode(name)
		}
		return result
	}

	result := []*Node{}

	var visit func(*Node) error
	visit = func(node *Node) error {
		if node.mark == 2 {
			return nil
		} else if node.mark == 1 {
			return fmt.Errorf("Cyclic dependencies detected")
		}
		node.mark = 1
		node.Edges = getNeighbours(node.Schema)
		for _, neighbour := range node.Edges {
			err := visit(neighbour)
			if err != nil {
				return err
			}
		}
		node.mark = 2
		result = append(result, node)
		return nil
	}

	for _, schema := range schemas {
		err := visit(getNode(schema.Name()))
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (node *Node) Join() error {
	return node.Schema.Join(node.Edges)
}

func updateSchemas(schemas []*Node) error {
	for i := len(schemas)-1; i >= 0; i-- {
		if err := schemas[i].Join(); err != nil {
			return err
		}
	}
	return nil
}
