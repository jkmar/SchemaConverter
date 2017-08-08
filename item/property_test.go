package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
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
			typeOfItem := "object"
			property := CreatePropertyWithType(name, typeOfItem)

			Expect(property.name).To(Equal(name))
			Expect(property.IsObject()).To(BeTrue())
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
			property := &Property{"", object}

			err := property.AddProperties(properties, true)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties).To(Equal(properties))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix   = "abc"
			property *Property
			data     map[interface{}]interface{}
		)

		BeforeEach(func() {
			property = &Property{name: "def"}
		})

		It("Should return an error for an object with no items", func() {
			data = map[interface{}]interface{}{}

			err := property.Parse(prefix, data)

			expected := fmt.Errorf(
				"property %s does not have a type",
				util.AddName(prefix, property.name),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for a property with an invalid type", func() {
			data = map[interface{}]interface{}{
				"type": 1,
			}

			err := property.Parse(prefix, data)

			expected := fmt.Errorf(
				"property %s: unsupported type: %T",
				util.AddName(prefix, property.name),
				data["type"],
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an invalid item", func() {
			data = map[interface{}]interface{}{
				"type": "array",
			}

			err := property.Parse(prefix, data)

			expected := fmt.Errorf(
				"array %s does not have items",
				util.AddName(prefix, property.name),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid property", func() {
			data = map[interface{}]interface{}{
				"type": "array",
				"items": map[interface{}]interface{}{
					"type": "string",
				},
			}

			err := property.Parse(prefix, data)

			typeOfItem := data["items"].(map[interface{}]interface{})["type"]
			expected := "[]" + typeOfItem.(string)
			Expect(err).ToNot(HaveOccurred())
			Expect(property.item.Type("")).To(Equal(expected))
		})
	})

	Describe("collect objects tests", func() {
		It("Should collect object", func() {
			object := &Object{"abc", nil}
			property := &Property{"", object}
			objects := set.New()
			objects.Insert(object)

			result, err := property.CollectObjects(1, 0)

			expected := objects
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Describe("generate property tests", func() {
		var (
			prefix     = "abc"
			suffix     = "xyz"
			annotation = "123"
			property   *Property
			data       map[interface{}]interface{}
		)

		BeforeEach(func() {
			property = &Property{name: "def_id"}
		})

		It("Should generate correct property for a plain item", func() {
			data = map[interface{}]interface{}{
				"type": "boolean",
			}

			err := property.Parse(prefix, data)
			Expect(err).ToNot(HaveOccurred())

			result := property.GenerateProperty(suffix, annotation)

			expected := fmt.Sprintf(
				"\tDefID bool `%s:\"%s\"`\n",
				annotation,
				property.name,
			)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for an array", func() {
			data = map[interface{}]interface{}{
				"type": "array",
				"items": map[interface{}]interface{}{
					"type": "integer",
				},
			}

			err := property.Parse(prefix, data)
			Expect(err).ToNot(HaveOccurred())

			result := property.GenerateProperty(suffix, annotation)

			expected := fmt.Sprintf(
				"\tDefID []int64 `%s:\"%s\"`\n",
				annotation,
				property.name,
			)
			Expect(result).To(Equal(expected))
		})

		It("Should generate correct property for an object", func() {
			data = map[interface{}]interface{}{
				"type": "object",
				"properties": map[interface{}]interface{}{
					"test": map[interface{}]interface{}{
						"type": "string",
					},
				},
			}

			err := property.Parse(prefix, data)
			Expect(err).ToNot(HaveOccurred())

			result := property.GenerateProperty(suffix, annotation)

			expected := fmt.Sprintf(
				"\tDefID AbcDefIDXyz `%s:\"%s\"`\n",
				annotation,
				property.name,
			)
			Expect(result).To(Equal(expected))
		})
	})
})
