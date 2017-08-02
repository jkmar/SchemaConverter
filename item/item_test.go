package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("item tests", func() {
	Describe("item creation tests", func() {
		Describe("CreateItemFromString tests", func() {
			It("Should create an array", func() {
				typeName := "array"
				expected := &Array{}
				result := CreateItemFromString(typeName)
				Expect(result).To(Equal(expected))
			})

			It("Should create an object", func() {
				typeName := "object"
				expected := &Object{}
				result := CreateItemFromString(typeName)
				Expect(result).To(Equal(expected))
			})

			It("Should create a plain item", func() {
				typeName := "string"
				expected := &PlainItem{}
				result := CreateItemFromString(typeName)
				Expect(result).To(Equal(expected))
			})
		})

		Describe("CreateItem", func() {
			var itemType interface{}

			It("Should return error for an invalid type", func() {
				itemType = 1
				expected := fmt.Errorf("unsupported type: %T", itemType)
				_, err := CreateItem(itemType)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})

			It("Should create item with a correct type", func() {
				itemType = []interface{}{"null", 1, "object", "array"}
				expected := &Object{}
				result, err := CreateItem(itemType)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expected))
			})
		})
	})
})
