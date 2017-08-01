package util

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type parsing tests", func() {
	Describe("AddName tests", func() {
		It("Should not add _ to an empty prefix", func() {
			name := "abc"
			result := AddName("", name)
			Expect(result).To(Equal(name))
		})

		It("Should add _ to a nonempty prefix", func() {
			prefix := "a"
			sufix := "b"
			expected := prefix + "_" + sufix
			result := AddName(prefix, sufix)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("ToGoName tests", func() {
		It("Should return correct go name", func() {
			name := "aa-bb_cc-dd-eE"
			expected := "AaBbCcDdEE"
			result := ToGoName(name, "")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("mapType tests", func() {
		It("Should return the mapped type", func() {
			typeName := "number"
			expected := typeMapping[typeName]
			result := mapType(typeName)
			Expect(result).To(Equal(expected))
		})

		It("Should return given type for a type with no match", func() {
			typeName := "string"
			expected := typeName
			result := mapType(typeName)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("ParseType tests", func() {
		var itemType interface{}

		It("Should return error for an unsupported argument type", func() {
			itemType = 1
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := ParseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for an array with no string", func() {
			itemType = []interface{}{1, false}
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := ParseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for an array with nulls", func() {
			itemType = []interface{}{"null", "null"}
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := ParseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return the mapped type for a string input", func() {
			itemType = "number"
			expected := typeMapping[itemType.(string)]
			result, err := ParseType(itemType)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})

		It("Should return the mapped type for an array input", func() {
			itemType = []interface{}{"null", "number", "boolean"}
			expected := typeMapping[itemType.([]interface{})[1].(string)]
			result, err := ParseType(itemType)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})
