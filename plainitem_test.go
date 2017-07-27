package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("plain item tests", func() {
	Describe("type tests", func() {
		It("Should return correct item type", func() {
			itemType := "int64"
			item := PlainItem{itemType}
			expected := itemType
			result := item.Type()
			Expect(result).To(Equal(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item *PlainItem
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &PlainItem{}
		})

		It("Should return error for object with no type", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"invalid schema: item %s does not have a type",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for unsupported type", func() {
			object = map[interface{}]interface{}{"type": 1}
			expected := fmt.Errorf(
				"invalid schema: item %s - unsupported type: %T",
				prefix,
				object["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid object", func() {
			object = map[interface{}]interface{}{"type": "number"}
			expected := typeMapping[object["type"].(string)]
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.Type()).To(Equal(expected))
		})
	})
})
