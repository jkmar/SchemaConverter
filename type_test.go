package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("type parsing tests", func() {
	Describe("addName tests", func() {
		It("Should not add _ to empty prefix", func() {
			name := "abc"
			result := addName("", name)
			Expect(result).To(Equal(name))
		})

		It("Should add _ to nonempty prefix", func() {
			prefix := "a"
			sufix := "b"
			expected := prefix + "_" + sufix
			result := addName(prefix, sufix)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("toGoName tests", func() {
		It("Should return correct go name", func() {
			name := "aa-bb_cc-dd-eE"
			expected := "AaBbCcDdEE"
			result := toGoName(name, "")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("mapType tests", func() {
		It("Should return mapped type", func() {
			typeName := "number"
			expected := typeMapping[typeName]
			result := mapType(typeName)
			Expect(result).To(Equal(expected))
		})

		It("Should return given type for type with no match", func() {
			typeName := "string"
			expected := typeName
			result := mapType(typeName)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("parseType tests", func() {
		var itemType interface{}

		It("Should return error for unsupported argument type", func() {
			itemType = 1
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := parseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for array with no string", func() {
			itemType = []interface{}{1, false}
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := parseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for array with nulls", func() {
			itemType = []interface{}{"null", "null"}
			expected := fmt.Errorf("unsupported type: %T", itemType)
			_, err := parseType(itemType)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return mapped type for string input", func() {
			itemType = "number"
			expected := typeMapping[itemType.(string)]
			result, err := parseType(itemType)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})

		It("Should return mapped type for array input", func() {
			itemType = []interface{}{"null", "number", "boolean"}
			expected := typeMapping[itemType.([]interface{})[1].(string)]
			result, err := parseType(itemType)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})