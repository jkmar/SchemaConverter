package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

var _ = Describe("property tests", func() {
	Describe("creation tests", func() {
		It("Should create a property with a correct name", func() {
			name := "name"
			property := CreateProperty(name)
			Expect(property.name).To(Equal(name))
		})

		It("Should create a property with a correct name and type", func() {
			name := "name"
			itemType := "object"
			property := CreatePropertyWithType(name, itemType)
			Expect(property.name).To(Equal(name))
			Expect(property.item.IsObject()).To(BeTrue())
		})
	})

	Describe("is object tests", func() {
		It("Should return false for not an object", func() {
			property := CreatePropertyWithType("", "string")
			Expect(property.IsObject()).To(BeFalse())
		})

		It("Should return false for an object", func() {
			property := CreatePropertyWithType("", "object")
			Expect(property.IsObject()).To(BeTrue())
		})
	})

	Describe("add properties tests", func() {
		It("Should add properties for an item", func() {
			properties := set.New()
			properties.Insert(CreateProperty("a"))
			object := &Object{}
			item := &Property{"", object}
			err := item.AddProperties(properties, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties).To(Equal(properties))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item *Property
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Property{name: "def"}
		})

		It("Should return an error for an object with no items", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"property %s does not have a type",
				addName(prefix, item.name),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for a property with an invalid type", func() {
			object = map[interface{}]interface{}{
				"type": 1,
			}
			expected := fmt.Errorf(
				"property %s: unsupported type: %T",
				addName(prefix, item.name),
				object["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an invalid item", func() {
			object = map[interface{}]interface{}{
				"type": "array",
			}
			expected := fmt.Errorf(
				"array %s does not have items",
				addName(prefix, item.name),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid property", func() {
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
			Expect(item.item.Type("")).To(Equal(expected))
		})
	})

	Describe("collect objects tests", func() {
		It("Should collect object", func() {
			object := &Object{"abc", nil}
			item := &Property{"", object}
			objects := set.New()
			objects.Insert(object)
			expected := objects
			result, err := item.CollectObjects(1, 0)
			Expect(err).ToNot(HaveOccurred())
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
			item = &Property{name: "def_id"}
		})

		It("Should generate correct property for a plain item", func() {
			object = map[interface{}]interface{}{
				"type": "boolean",
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			expected := fmt.Sprintf(
				"\tDefID bool `%s:\"%s\"`\n",
				annotation,
				item.name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for an array", func() {
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
				item.name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for an object", func() {
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
				item.name,
			)
			result := item.GenerateProperty(suffix, annotation)
			Expect(result).To(Equal(expected))
		})
	})
})
