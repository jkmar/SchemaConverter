package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("property tests", func() {
	Describe("creation tests", func() {
		It("Should create property with correct name", func() {
			name := "name"
			property := CreateProperty(name)
			Expect(property.Name).To(Equal(name))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item *Property
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Property{Name: "def"}
		})

		It("Should return error for object with no items", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"invalid schema: property %s does not have a type",
				addName(prefix, item.Name),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for property with invalid type", func() {
			object = map[interface{}]interface{}{
				"type": 1,
			}
			expected := fmt.Errorf(
				"invalid schema: property %s - unsupported type: %T",
				addName(prefix, item.Name),
				object["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for invalid item", func() {
			object = map[interface{}]interface{}{
				"type": "array",
			}
			expected := fmt.Errorf(
				"invalid schema: array %s does not have items",
				addName(prefix, item.Name),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid property", func() {
			object = map[interface{}]interface{}{
				"type": "array",
				"items": map[interface{}]interface{}{
					"type": "string",
				},
			}
			itemType := object["items"].(map[interface{}]interface{})["type"]
			expected := "[]" + itemType.(string)
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.Item.Type("")).To(Equal(expected))
		})
	})

	Describe("collect tests", func() {
		It("Should collect object", func() {
			object := &Object{"abc", nil}
			item := &Property{"", object}
			expected := []*Object{object}
			result := item.Collect(1)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("generate property tests", func() {
		var (
			prefix = "abc"
			suffix = "xyz"
			annotation = "123"
			item *Property
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Property{Name:"def_id"}
		})

		It("Should generate correct property for plain item", func() {
			object = map[interface{}]interface{}{
				"type": "boolean",
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			expected := fmt.Sprintf(
				"\tDefID bool `%s:\"%s\"`\n",
				annotation,
				item.Name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for array", func() {
			object = map[interface{}]interface{}{
				"type": "array",
				"items": map[interface{}]interface{}{
					"type": "integer",
				},
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			expected := fmt.Sprintf(
				"\tDefID []int64 `%s:\"%s\"`\n",
				annotation,
				item.Name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for object", func() {
			object = map[interface{}]interface{}{
				"type": "object",
				"properties": map[interface{}]interface{}{
					"test": map[interface{}]interface{}{
						"type": "string",
					},
				},
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			expected := fmt.Sprintf(
				"\tDefID AbcDefIDXyz `%s:\"%s\"`\n",
				annotation,
				item.Name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})
	})
})
