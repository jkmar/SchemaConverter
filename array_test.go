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

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item *Array
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Array{}
		})

		It("Should return error for object with no items", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"invalid schema: array %s does not have items",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for item with no type", func() {
			object = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"a": 1,
				},
			}
			expected := fmt.Errorf(
				"invalid schema: items of array %s do not have a type",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for invalid item", func() {
			object = map[interface{}]interface{}{
				"items": map[interface{}]interface{}{
					"type": 1,
				},
			}
			expected := fmt.Errorf(
				"invalid schema: array %s - unsupported type: %T",
				prefix,
				object["items"].(map[interface{}]interface{})["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid array", func() {
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

		Describe("collect tests", func() {
			It("Should return nil for array of plain items", func() {
				item := &Array{&PlainItem{}}
				Expect(item.Collect(1)).To(BeNil())
			})

			It("Should return object for array of objects", func() {
				name := "Test"
				item := &Array{&Object{name, nil}}
				result := item.Collect(1)
				Expect(len(result)).To(Equal(1))
				Expect(result[0].Type("")).To(Equal(name))
			})
		})
	})
})
