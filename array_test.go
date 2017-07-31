package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("array tests", func() {
	Describe("type tests", func() {
		It("Should return correct array type", func() {
			itemType := "int64"
			array := Array{&PlainItem{itemType}}
			expected := "[]" + itemType
			result := array.Type("")
			Expect(result).To(Equal(expected))
		})

		It("Should return correct array type for nested arrays", func() {
			itemType := "string"
			nested := Array{&PlainItem{itemType}}
			array := Array{&nested}
			expected := "[][]" + itemType
			result := array.Type("")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("is object tests", func() {
		It("Should return false for is object query", func() {
			item := &Array{}
			Expect(item.IsObject()).To(BeFalse())
		})
	})

	Describe("add properties tests", func() {
		It("Should return an error", func() {
			item := &Array{}
			expected := fmt.Errorf("cannot add properties to an array")
			err := item.AddProperties(nil, false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item *Array
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Array{}
		})

		It("Should return an error for an object with no items", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"array %s does not have items",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an item with no type", func() {
			object = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"a": 1,
				},
			}
			expected := fmt.Errorf(
				"items of array %s do not have a type",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an invalid item", func() {
			object = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"type": 1,
				},
			}
			expected := fmt.Errorf(
				"array %s: unsupported type: %T",
				prefix,
				object["items"].(map[interface{}]interface{})["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid array", func() {
			object = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"type": "string",
				},
			}
			itemType := object["items"].(map[interface{}]interface{})["type"]
			expected := "[]" + itemType.(string)
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.Type("")).To(Equal(expected))
		})

		Describe("collect object tests", func() {
			It("Should return nil for an array of plain items", func() {
				item := &Array{&PlainItem{}}
				Expect(item.CollectObjects(1, 0)).To(BeNil())
			})

			It("Should return object for an array of objects", func() {
				name := "Test"
				item := &Array{&Object{name, nil}}
				result, err := item.CollectObjects(1, 0)
				Expect(err).ToNot(HaveOccurred())
				array := result.ToArray()
				Expect(len(array)).To(Equal(1))
				Expect(array[0].(*Object).Type("")).To(Equal(name))
			})
		})

		Describe("collect properties tests", func() {
			It("Should return nil for an array of plain items", func() {
				item := &Array{&PlainItem{}}
				Expect(item.CollectProperties(1, 0)).To(BeNil())
			})
		})
	})
})
