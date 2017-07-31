package main

import (
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

type node struct {
	schema *Schema
	edges  []*node
	mark   int
}

type schemaGraph struct {
	allSchemas map[string]*Schema
	schemas set.Set
}

func createGraph(toConvert, other set.Set) (*schemaGraph, error) {
	allSchemas := map[string]*Schema{}
	if err := other.SafeInsertAll(toConvert); err != nil {
		fmt.Errorf("multiple schemas with the same name")
	}
	for name, schema := range toConvert {
		allSchemas[name] = schema.(*Schema)
	}
	return &schemaGraph{allSchemas, toConvert}, nil
}

func (graph *schemaGraph) sort() ([]*node, error) {
	nodes := map[string]*node{}

	getNode := func(id string) *node {
		if node, ok := nodes[id]; ok {
			return node
		}
		node := &node{schema: graph.allSchemas[id]}
		nodes[id] = node
		return node
	}

	getNeighbours := func(schema *Schema) []*node {
		result := make([]*node, len(schema.Bases()))
		for i, name := range schema.extends {
			result[i] = getNode(name)
		}
		return result
	}

	result := []*node{}

	var visit func(*node) error
	visit = func(node *node) error {
		if node.mark == 2 {
			return nil
		} else if node.mark == 1 {
			return fmt.Errorf("Cyclic dependencies detected")
		}
		node.mark = 1
		node.edges = getNeighbours(node.schema)
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

	for _, schema := range graph.schemas {
		err := visit(getNode(schema.Name()))
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (node *node) join() error {
	return node.schema.Join(node.edges)
}

func updateSchemas(schemas []*node) error {
	for i := len(schemas)-1; i >= 0; i-- {
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

