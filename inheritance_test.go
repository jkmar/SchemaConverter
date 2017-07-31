package main

/*
import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
	"github.com/zimnx/YamlSchemaToGoStruct/item"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

var _ = Describe("inheritance tests", func() {
	var (
		createA := func() *Schema {
			schema := &Schema{}
			object := map[interface{}]interface{}{
				"id": "a",
				"parent": "pa",
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"a": map[interface{}]interface{}{
							"type": "boolean",
						},
						"b": map[interface{}]interface{}{

						},
					},
				},
			}
		}

		getCorrect := func() set.Set {

		}
	)

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
			Expect(schema.schema.Name()).To(Equal(name))
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
			Expect(schema.parent).To(BeEmpty())
		})

		It("Should get correct parent", func() {
			name := "name"
			object := map[interface{}]interface{}{
				"parent": name,
			}
			schema.getParent(object)
			Expect(schema.parent).To(Equal(name))
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
			Expect(schema.extends).To(BeNil())
		})

		It("Should get correct base schemas", func() {
			object := map[interface{}]interface{}{
				"extends": []interface{}{"a", "b", "c"},
			}
			err := schema.getBaseSchemas(object)
			Expect(err).To(BeNil())
			Expect(schema.extends).To(Equal([]string{"a", "b", "c"}))
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
				"invalid schema %s: property %s does not have a type",
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

		It("Should return error for schema with invalid parent", func() {
			id := "test_schema"
			name := "test"
			object := map[interface{}]interface{}{
				"id": id,
				"parent": name,
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						util.AddName(name, "id"): map[interface{}]interface{}{
							"type": "boolean",
						},
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: object %s: multiple properties have the same name",
				id,
				id,
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
			Expect(schema.Name()).To(Equal(name))
			Expect(schema.schema.IsObject()).To(BeTrue())
		})
	})

	Describe("parse all tests", func() {
		It("Should return error for invalid schema", func() {
			name := "test"
			objects := []map[interface{}]interface{}{
				{
					"id": name + "1",
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"a": map[interface{}]interface{}{
								"type": "string",
							},
						},
					},
				},
				{
					"id": name,
					"schema": map[interface{}]interface{}{
						"type": "string",
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: schema should be an object",
				name,
			)
			_, err := ParseAll(objects)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for multiple schemas with the same name", func() {
			name := "test"
			objects := []map[interface{}]interface{}{
				{
					"id": name,
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"a": map[interface{}]interface{}{
								"type": "string",
							},
						},
					},
				},
				{
					"id": name,
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"b": map[interface{}]interface{}{
								"type": "number",
							},
						},
					},
				},
			}
			expected := fmt.Errorf(
				"multiple schemas with the same name: %s",
				name,
			)
			_, err := ParseAll(objects)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid schemas", func() {
			name := "test"
			objects := []map[interface{}]interface{}{
				{
					"id": name + "0",
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"a": map[interface{}]interface{}{
								"type": "string",
							},
						},
					},
				},
				{
					"id": name + "1",
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"b": map[interface{}]interface{}{
								"type": "number",
							},
						},
					},
				},
			}
			set, err := ParseAll(objects)
			Expect(err).ToNot(HaveOccurred())
			Expect(set.Size()).To(Equal(len(objects)))
		})
	})

	Describe("collect objects tests", func() {
		It("Should return error for schema with multiple objects with the same name", func() {
			name := "name"
			schema := &Schema{}
			object := map[interface{}]interface{}{
				"id": name,
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"A_": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"B": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
							},
						},
						"A": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"_B": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
							},
						},
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: multiple objects with the same type at object %s",
				name,
				name,
			)
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			_, err = schema.CollectObjects(-1, 0)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

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
			result, err := schema.CollectObjects(-1, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(result)).To(Equal(5))
			array := result.ToArray()
			Expect(util.ToGoName(array[0].Name(), "")).To(Equal(names[0]))
			Expect(util.ToGoName(array[1].Name(), "")).To(Equal(names[0]+names[0]))
			Expect(util.ToGoName(array[2].Name(), "")).To(Equal(names[0]+names[1]))
			Expect(util.ToGoName(array[3].Name(), "")).To(Equal(names[0]+names[1]+names[0]))
			Expect(util.ToGoName(array[4].Name(), "")).To(Equal(names[0]+names[1]+names[2]))
		})
	})

	Describe("collect properties tests", func() {
		It("Should return error for schema with multiple properties with the same name", func() {
			name := "name"
			schema := &Schema{}
			object := map[interface{}]interface{}{
				"id": name,
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"A_": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"B": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
							},
						},
						"A": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"_B": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"a": map[interface{}]interface{}{
											"type": "number",
										},
									},
								},
							},
						},
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema %s: multiple properties with the same name at object %s",
				name,
				name,
			)
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			_, err = schema.CollectProperties(-1, 0)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should collect all properties", func() {
			schema := &Schema{}
			names := []string{"A", "B", "C", "D"}
			object := map[interface{}]interface{}{
				"id": names[3],
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						names[1]: map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								names[0]: map[interface{}]interface{}{
									"type": "string",
								},
								names[2]: map[interface{}]interface{}{
									"type": "number",
								},
							},
						},
					},
				},
			}
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			result, err := schema.CollectProperties(-1, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(result)).To(Equal(len(names)))
			array := result.ToArray()
			Expect(util.ToGoName(array[0].Name(), "")).To(Equal(names[0]))
			Expect(util.ToGoName(array[1].Name(), "")).To(Equal(names[1]))
			Expect(util.ToGoName(array[2].Name(), "")).To(Equal(names[2]))
			Expect(util.ToGoName(array[3].Name(), "")).To(Equal(names[3]))
		})
	})

	Describe("join tests", func() {
		var (
			schema *Schema
			names = []string{"a", "b", "c", "d"}
		)

		BeforeEach(func() {
			schema = &Schema{}
			object := map[interface{}]interface{}{
				"id": names[0],
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						names[1]: map[interface{}]interface{}{
							"type": "string",
						},
						names[2]: map[interface{}]interface{}{
							"type": "number",
						},
					},
				},
			}
			err := schema.Parse(object)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should return error when joining to invalid schema", func() {
			other := &Schema{schema: item.CreatePropertyWithType(names[3], "string")}
			expected := fmt.Errorf(
				"schema %s should be an object",
				other.Name(),
			)
			err := other.Join([]*node{{schema: schema}})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error when schemas in nodes share properties", func() {
			nodes := []*node{{schema: schema}, {schema: schema}}
			expected := fmt.Errorf(
				"multiple properties with the same name in bases of schema %s",
				names[0],
			)
			err := schema.Join(nodes)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should join correctly", func() {
			other := &Schema{}
			object := map[interface{}]interface{}{
				"id": names[0],
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						names[3]: map[interface{}]interface{}{
							"type": "boolean",
						},
						names[2]: map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"a": map[interface{}]interface{}{
									"type": "string",
								},
								"b": map[interface{}]interface{}{
									"type": "boolean",
								},
							},
						},
					},
				},
			}
			err := other.Parse(object)
			Expect(err).ToNot(HaveOccurred())
			err = schema.Join([]*node{{schema: other}})
			Expect(err).ToNot(HaveOccurred())
			properties, err := schema.CollectProperties(-1, 1)
			Expect(err).ToNot(HaveOccurred())
			array := properties.ToArray()
			Expect(len(array)).To(Equal(3))
			Expect(array[2].(*item.Property).IsObject()).To(BeFalse())
		})
	})
})
*/