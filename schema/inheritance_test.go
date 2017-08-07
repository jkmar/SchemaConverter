package schema

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

var _ = Describe("inheritance tests", func() {
	var (
		getObject = func(base, name, itemType string) map[interface{}]interface{} {
			return map[interface{}]interface{}{
				"id":     name,
				"parent": "p" + name,
				"schema": map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						base: map[interface{}]interface{}{
							"type": itemType,
						},
						name: map[interface{}]interface{}{
							"type": "string",
						},
					},
				},
			}
		}

		createFromObject = func(object map[interface{}]interface{}) *Schema {
			schema := &Schema{}
			err := schema.parse(object)
			Expect(err).ToNot(HaveOccurred())
			return schema
		}
	)

	Describe("collect tests", func() {
		It("Should return error for multiple schemas with the same name", func() {
			schema := createFromObject(getObject("a", "b", "string"))
			other := createFromObject(getObject("a", "b", "number"))
			Expect(schema).ToNot(Equal(other))
			toConvert := set.New()
			toConvert.Insert(schema)
			otherSet := set.New()
			otherSet.Insert(other)
			expected := fmt.Errorf("multiple schemas with the same name")
			err := collectSchemas(toConvert, otherSet)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for non existing base schema", func() {
			id := "b"
			object := getObject("a", "a", "string")
			object["extends"] = []interface{}{id}
			toConvert := set.New()
			toConvert.Insert(createFromObject(object))
			expected := fmt.Errorf(
				"schema with id %s does not exist",
				id,
			)
			err := collectSchemas(toConvert, set.New())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return cyclic dependency error", func() {
			objects := make([]map[interface{}]interface{}, 3)
			objects[0] = getObject("a", "a", "string")
			objects[0]["extends"] = []interface{}{"b"}
			objects[1] = getObject("b", "b", "string")
			objects[1]["extends"] = []interface{}{"c"}
			objects[2] = getObject("c", "c", "string")
			objects[2]["extends"] = []interface{}{"a"}
			toConvert := set.New()
			for _, object := range objects {
				toConvert.Insert(createFromObject(object))
			}
			expected := fmt.Errorf("cyclic dependencies detected")
			err := collectSchemas(toConvert, set.New())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return join error", func() {
			id := "a"
			base := "base"
			objects := make([]map[interface{}]interface{}, 3)
			objects[0] = getObject(id, id, "string")
			objects[0]["extends"] = []interface{}{"b", "c"}
			objects[1] = getObject(base, "b", "string")
			objects[2] = getObject(base, "c", "string")
			toConvert := set.New()
			for _, object := range objects {
				toConvert.Insert(createFromObject(object))
			}
			expected := fmt.Errorf(
				"multiple properties with the same name in bases of schema %s",
				id,
			)
			err := collectSchemas(toConvert, set.New())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should correctly update schemas", func() {
			objects := make([]map[interface{}]interface{}, 3)
			objects[0] = getObject("a", "a", "string")
			objects[0]["extends"] = []interface{}{"b", "c"}
			objects[1] = getObject("b", "b", "string")
			objects[1]["extends"] = []interface{}{"c"}
			objects[2] = getObject("c", "c", "string")
			toConvert := set.New()
			for _, object := range objects {
				toConvert.Insert(createFromObject(object))
			}
			err := collectSchemas(toConvert, set.New())
			Expect(err).ToNot(HaveOccurred())
			array := toConvert.ToArray()
			for i, schema := range array {
				newSet, err := schema.(*Schema).collectProperties(-1, 1)
				Expect(err).ToNot(HaveOccurred())
				Expect(newSet.Size()).To(Equal(2 * (3 - i)))
			}
		})
	})
})
