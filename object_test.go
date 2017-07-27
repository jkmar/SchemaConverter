package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("object tests", func() {
	Describe("type tests", func() {
		It("Should return correct object type", func() {
			itemType := "ab"
			object := Object{objectType: itemType}
			expected := itemType
			result := object.Type()
			Expect(result).To(Equal(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item   *Object
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Object{}
		})

		It("Should return error for object with no properties", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"invalid schema: object %s does not have properties",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for object with invalid property definition", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": 1,
				},
			}
			expected := fmt.Errorf(
				"invalid schema: object %s has invalid property a",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for object with invalid property", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"a": 1,
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema: property %s does not have a type",
				addName(prefix, "a"),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid object", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"type": "string",
					},
					"b": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "string",
						},
					},
					"c": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"x": map[interface{}]interface{}{
								"type": "boolean",
							},
							"y": map[interface{}]interface{}{
								"type": "string",
							},
						},
					},
				},
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.objectType).To(Equal(prefix))
			Expect(len(item.Properties)).To(Equal(len(object["properties"].(map[interface{}]interface{}))))
			names := make([]string, len(item.Properties))
			for i, property := range item.Properties {
				names[i] = property.Name
				switch names[i] {
				case "a":
					Expect(property.Item.Type()).To(Equal("string"))
				case "b":
					Expect(property.Item.Type()).To(Equal("[]string"))
				case "c":
					Expect(property.Item.Type()).To(Equal(addName(prefix, "c")))
				}
			}
			sort.Strings(names)
			Expect(names).To(Equal([]string{"a", "b", "c"}))
		})
	})
})
