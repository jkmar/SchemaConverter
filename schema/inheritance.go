package schema

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

type node struct {
	value *Schema
	edges []*node
	mark  int
}

type graph struct {
	allSchemas map[string]*Schema
	schemas    set.Set
}

func createGraph(toConvert, other set.Set) (*graph, error) {
	allSchemas := map[string]*Schema{}
	if err := other.SafeInsertAll(toConvert); err != nil {
		return nil, fmt.Errorf("multiple schemas with the same name")
	}
	for name, otherSchema := range other {
		allSchemas[name] = otherSchema.(*Schema)
	}
	return &graph{allSchemas, toConvert}, nil
}

func (graph *graph) sort() ([]*node, error) {
	nodes := map[string]*node{}

	getNode := func(id string) (*node, error) {
		if node, ok := nodes[id]; ok {
			return node, nil
		}
		value, ok := graph.allSchemas[id]
		if !ok {
			return nil, fmt.Errorf(
				"schema with id %s does not exist",
				id,
			)
		}
		node := &node{value: value}
		nodes[id] = node
		return node, nil
	}

	getNeighbours := func(schema *Schema) ([]*node, error) {
		result := make([]*node, len(schema.bases()))
		for i, name := range schema.extends {
			var err error
			if result[i], err = getNode(name); err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	result := []*node{}

	var visit func(*node) error
	visit = func(node *node) error {
		if node.mark == 2 {
			return nil
		} else if node.mark == 1 {
			return fmt.Errorf("cyclic dependencies detected")
		}
		var err error
		node.mark = 1
		if node.edges, err = getNeighbours(node.value); err != nil {
			return err
		}
		for _, neighbour := range node.edges {
			err := visit(neighbour)
			if err != nil {
				return err
			}
		}
		node.mark = 2
		result = append(result, node)
		return nil
	}

	for _, vertex := range graph.schemas {
		node, _ := getNode(vertex.Name())
		if err := visit(node); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (node *node) join() error {
	return node.value.join(node.edges)
}

func updateSchemas(schemas []*node) error {
	for i := len(schemas) - 1; i >= 0; i-- {
		if err := schemas[i].join(); err != nil {
			return err
		}
	}
	return nil
}

func collectSchemas(toConvert, other set.Set) error {
	graph, err := createGraph(toConvert, other)
	if err != nil {
		return err
	}
	nodes, err := graph.sort()
	if err != nil {
		return err
	}
	return updateSchemas(nodes)
}
