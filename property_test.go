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
			Expect(item.Item.Type()).To(Equal(expected))
		})
	})
})
