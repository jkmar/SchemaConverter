package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("json kind tests", func() {
	var jsonKind = &JSONKind{}

	Describe("type tests", func() {
		It("Should return a correct type for a null item", func() {
			typeOfItem := "string"
			list := []interface{}{typeOfItem, "null"}
			newItem, err := CreateItem(list)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Required: true,
				Data:     map[interface{}]interface{}{"type": list},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.Type("", newItem)

			expected := typeOfItem
			Expect(result).To(Equal(expected))
		})

		It("Should return a correct type for a not null item", func() {
			typeOfItem := "string"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Required: true,
				Data:     map[interface{}]interface{}{"type": typeOfItem},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.Type("", newItem)

			expected := typeOfItem
			Expect(result).To(Equal(expected))
		})
	})

	Describe("interface type tests", func() {
		It("Should return a correct interface type for an object", func() {
			newItem, err := CreateItem("object")
			Expect(err).ToNot(HaveOccurred())

			name := "Test"
			err = newItem.Parse(ParseContext{
				Prefix:   name,
				Required: false,
				Data: map[interface{}]interface{}{
					"type":       "object",
					"properties": map[interface{}]interface{}{},
				},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.InterfaceType("", newItem)

			expected := "I" + name
			Expect(result).To(Equal(expected))
		})
	})

	Describe("annotation tests", func() {
		It("Should return a correct annotation for a null item", func() {
			name := "name"
			typeOfItem := "string"
			list := []interface{}{typeOfItem, "null"}
			newItem, err := CreateItem(list)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Required: true,
				Data:     map[interface{}]interface{}{"type": list},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`json:\"%s,omitempty\"`",
				name,
			)
			Expect(result).To(Equal(expected))
		})

		It("Should return a correct annotation for a not null item", func() {
			name := "name"
			typeOfItem := "string"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Required: true,
				Data:     map[interface{}]interface{}{"type": typeOfItem},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`json:\"%s\"`",
				name,
			)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("default value tests", func() {
		It("Should return a correct default value for a not null item", func() {
			typeOfItem := "int"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Defaults: 1,
				Required: true,
				Data:     map[interface{}]interface{}{"type": typeOfItem},
			})
			Expect(err).ToNot(HaveOccurred())

			result := jsonKind.Default("", newItem)

			expected := "1"
			Expect(result).To(Equal(expected))
		})
	})
})
