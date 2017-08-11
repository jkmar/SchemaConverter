package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("item tests", func() {
	Describe("item creation tests", func() {
		Describe("createItemFromString tests", func() {
			It("Should create an array", func() {
				typeName := "array"

				result := createItemFromString(typeName)

				expected := &Array{}
				Expect(result).To(Equal(expected))
			})

			It("Should create an object", func() {
				typeName := "object"

				result := createItemFromString(typeName)

				expected := &Object{}
				Expect(result).To(Equal(expected))
			})

			It("Should create a plain item", func() {
				typeName := "string"

				result := createItemFromString(typeName)

				expected := &PlainItem{}
				Expect(result).To(Equal(expected))
			})
		})

		Describe("CreateItem", func() {
			var typeOfItem interface{}

			It("Should return error for an invalid type", func() {
				typeOfItem = 1

				_, err := CreateItem(typeOfItem)

				expected := fmt.Errorf("unsupported type: %T", typeOfItem)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})

			It("Should create item with a correct type", func() {
				typeOfItem = []interface{}{"null", 1, "object", "array"}

				result, err := CreateItem(typeOfItem)

				expected := &Object{}
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expected))
			})
		})
	})
})
