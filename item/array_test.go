package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("array tests", func() {
	Describe("type tests", func() {
		It("Should return correct array type", func() {
			typeOfItem := "int64"
			array := Array{&PlainItem{itemType: typeOfItem}}

			result := array.Type("")

			expected := "[]" + typeOfItem
			Expect(result).To(Equal(expected))
		})

		It("Should return correct array type for nested arrays", func() {
			typeOfItem := "string"
			nested := Array{&PlainItem{itemType: typeOfItem}}
			array := Array{&nested}

			result := array.Type("")

			expected := "[][]" + typeOfItem
			Expect(result).To(Equal(expected))
		})
	})

	Describe("add properties tests", func() {
		It("Should return an error", func() {
			array := &Array{}

			err := array.AddProperties(nil, false)

			expected := fmt.Errorf("cannot add properties to an array")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			array  *Array
			data   map[interface{}]interface{}
		)

		BeforeEach(func() {
			array = &Array{}
		})

		It("Should return an error for an object with no items", func() {
			data = map[interface{}]interface{}{}

			err := array.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"array %s does not have items",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an item with no type", func() {
			data = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"a": 1,
				},
			}

			err := array.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"items of array %s do not have a type",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an invalid item", func() {
			data = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"type": 1,
				},
			}

			err := array.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"array %s: unsupported type: %T",
				prefix,
				data["items"].(map[interface{}]interface{})["type"],
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid array", func() {
			data = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"type": "string",
				},
			}

			err := array.Parse(prefix, 0, true, data)

			typeOfItem := data["items"].(map[interface{}]interface{})["type"]
			expected := "[]" + typeOfItem.(string)
			Expect(err).ToNot(HaveOccurred())
			Expect(array.Type("")).To(Equal(expected))
		})

		Describe("collect object tests", func() {
			It("Should return nil for an array of plain items", func() {
				array := &Array{&PlainItem{}}

				Expect(array.CollectObjects(1, 0)).To(BeNil())
			})

			It("Should return object for an array of objects", func() {
				name := "Test"
				array := &Array{&Object{name, nil}}

				result, err := array.CollectObjects(1, 0)

				Expect(err).ToNot(HaveOccurred())
				objects := result.ToArray()
				Expect(len(objects)).To(Equal(1))
				Expect(objects[0].(*Object).Type("")).To(Equal(name))
			})
		})

		Describe("collect properties tests", func() {
			It("Should return nil for an array of plain items", func() {
				array := &Array{&PlainItem{}}

				Expect(array.CollectProperties(1, 0)).To(BeNil())
			})
		})
	})
})
