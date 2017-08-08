package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plain item tests", func() {
	Describe("type tests", func() {
		It("Should return correct item type", func() {
			typeOfItem := "int64"
			plainItem := PlainItem{typeOfItem}

			result := plainItem.Type("")

			expected := typeOfItem
			Expect(result).To(Equal(expected))
		})
	})

	Describe("add properties tests", func() {
		It("Should return an error", func() {
			plainItem := &PlainItem{}

			err := plainItem.AddProperties(nil, false)

			expected := fmt.Errorf("cannot add properties to a plain item")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix    = "abc"
			plainItem *PlainItem
			data      map[interface{}]interface{}
		)

		BeforeEach(func() {
			plainItem = &PlainItem{}
		})

		It("Should return an error for an object with no type", func() {
			data = map[interface{}]interface{}{}

			err := plainItem.Parse(prefix, data)

			expected := fmt.Errorf(
				"item %s does not have a type",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an unsupported type", func() {
			data = map[interface{}]interface{}{"type": 1}

			err := plainItem.Parse(prefix, data)

			expected := fmt.Errorf(
				"item %s: unsupported type: %T",
				prefix,
				data["type"],
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid object", func() {
			data = map[interface{}]interface{}{"type": "number"}

			err := plainItem.Parse(prefix, data)

			expected := "int64"
			Expect(err).ToNot(HaveOccurred())
			Expect(plainItem.Type("")).To(Equal(expected))
		})
	})

	Describe("collect objects tests", func() {
		It("Should return nil for a plain item", func() {
			plainItem := &PlainItem{}

			Expect(plainItem.CollectObjects(1, 0)).To(BeNil())
		})
	})

	Describe("collect objects tests", func() {
		It("Should return nil for a plain item", func() {
			plainItem := &PlainItem{}

			Expect(plainItem.CollectProperties(1, 0)).To(BeNil())
		})
	})
})
