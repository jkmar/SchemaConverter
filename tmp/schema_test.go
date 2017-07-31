package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("schema tests", func() {
	Describe("get name tests", func() {
		var schema *Schema

		BeforeEach(func() {
			schema = &Schema{}
		})

		It("Should return error for schema with no name", func() {
			object := map[interface{}]interface{}{}
			expected := fmt.Errorf("schema does not have an id")
			err := schema.getName(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for schema with invalid name", func() {
			object := map[interface{}]interface{}{
				"id": 1,
			}
			expected := fmt.Errorf("schema does not have an id")
			err := schema.getName(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should get correct schema name", func() {
			name := "name"
			object := map[interface{}]interface{}{
				"id": name,
			}
			err := schema.getName(object)
			Expect(err).ToNot(HaveOccurred())
			Expect(schema.Schema.name).To(Equal(name))
		})
	})

	Describe("get parent tests", func() {
		var schema *Schema

		BeforeEach(func() {
			schema = &Schema{}
		})

		It("Should get empty parent for schema with invalid parent", func() {
			object := map[interface{}]interface{}{
				"parent": 1,
			}
			schema.getParent(object)
			Expect(schema.Parent).To(BeEmpty())
		})

		It("Should get correct parent", func() {
			name := "name"
			object := map[interface{}]interface{}{
				"parent": name,
			}
			schema.getParent(object)
			Expect(schema.Parent).To(Equal(name))
		})
	})

	Describe("get base schemas tests", func() {
		var schema *Schema

		BeforeEach(func() {
			schema = &Schema{}
		})

		It("Should return error for schema with invalid base", func() {
			object := map[interface{}]interface{}{
				"extends": []interface{}{
					"a",
					1,
				},
			}
			expected := fmt.Errorf("one of the base schemas is not a string")
			err := schema.getBaseSchemas(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should get empty bases for schema with invalid base", func() {
			object := map[interface{}]interface{}{
				"extends": 1,
			}
			err := schema.getBaseSchemas(object)
			Expect(err).ToNot(HaveOccurred())
			Expect(schema.Extends).To(BeNil())
		})

		It("Should get correct base schemas", func() {
			object := map[interface{}]interface{}{
				"extends": []interface{}{"a", "b", "c"},
			}
			err := schema.getBaseSchemas(object)
			Expect(err).To(BeNil())
			Expect(schema.Extends).To(Equal([]string{"a", "b", "c"}))
		})
	})

	Describe("parse tests", func() {
		var schema *Schema

		BeforeEach(func() {
			schema = &Schema{}
		})

		It("Should return error for schema with no name", func() {
			object := map[interface{}]interface{}{}
			expected := fmt.Errorf("schema does not have an id")
			err := schema.Parse(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for schema with invalid base schema", func() {
			name := "test"
			object := map[interface{}]interface{}{
				"id": name,
				"extends": []interface{}{
					1,
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: one of the base schemas is not a string",
				name,
			)
			err := schema.Parse(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for schema with no schema", func() {
			name := "test"
			object := map[interface{}]interface{}{
				"id":         name,
			}
			expected := fmt.Errorf(
				"invalid schema %s: schema does not have a \"schema\"",
				name,
			)
			err := schema.Parse(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for invalid schema", func() {
			name := "test"
			object := map[interface{}]interface{}{
				"id":         name,
				"schema": map[interface{}]interface{}{
					"properties": 1,
				},
			}
			expected := fmt.Errorf(
				"%s - invalid schema: property %s does not have a type",
				name,
				name,
			)
			err := schema.Parse(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for schema which is not an object", func() {
			name := "test"
			object := map[interface{}]interface{}{
				"id":   name,
				"schema": map[interface{}]interface{}{
					"type": "string",
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: schema should be an object",
				name,
			)
			err := schema.Parse(object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse correct schema", func() {
			name := "test"
			object := map[interface{}]interface{}{
				"id":   name,
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"a": map[interface{}]interface{}{
							"type": "string",
						},
					},
				},
			}
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			Expect(schema.Schema.name).To(Equal(name))
			Expect(schema.Schema.item.IsObject()).To(BeTrue())
			Expect(schema.Schema.item.Type("")).To(Equal("Test"))
		})
	})

	Describe("collect tests", func() {
		It("Should collect all object", func() {
			schema := &Schema{}
			names := []string{"A", "B", "C"}
			object := map[interface{}]interface{}{
				"id":   names[0],
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						names[1]: map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"a": map[interface{}]interface{}{
									"type": "string",
								},
								names[2]: map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
								names[0]: map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
							},
						},
						names[0]: map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"a": map[interface{}]interface{}{
									"type": "number",
								},
							},
						},
					},
				},
			}
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			result := schema.Collect(-1)
			Expect(len(result)).To(Equal(5))
			expected := []string{
				names[0],
				names[0]+names[0],
				names[0]+names[1],
				names[0]+names[1]+names[0],
				names[0]+names[1]+names[2],
			}
			types := make([]string, len(expected))
			for i, value := range result {
				types[i] = value.Type("")
			}
			sort.Strings(types)
			Expect(types).To(Equal(expected))
		})
	})
})
