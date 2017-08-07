package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plain item tests", func() {
	Describe("type tests", func() {
		It("Should return correct item type", func() {
			itemType := "int64"
			item := PlainItem{itemType}
			expected := itemType
			result := item.Type("")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("add properties tests", func() {
		It("Should return an error", func() {
			item := &PlainItem{}
			expected := fmt.Errorf("cannot add properties to a plain item")
			err := item.AddProperties(nil, false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item   *PlainItem
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &PlainItem{}
		})

		It("Should return an error for an object with no type", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"item %s does not have a type",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an unsupported type", func() {
			object = map[interface{}]interface{}{"type": 1}
			expected := fmt.Errorf(
				"item %s: unsupported type: %T",
				prefix,
				object["type"],
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid object", func() {
			object = map[interface{}]interface{}{"type": "number"}
			expected := "int64"
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.Type("")).To(Equal(expected))
		})
	})

	Describe("collect objects tests", func() {
		It("Should return nil for a plain item", func() {
			item := &PlainItem{}
			Expect(item.CollectObjects(1, 0)).To(BeNil())
		})
	})

	Describe("collect objects tests", func() {
		It("Should return nil for a plain item", func() {
			item := &PlainItem{}
			Expect(item.CollectProperties(1, 0)).To(BeNil())
		})
	})
})
