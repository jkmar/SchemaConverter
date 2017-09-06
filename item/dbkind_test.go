package item

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

var _ = Describe("json kind tests", func() {
	var dbKind = &DBKind{}

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

			result := dbKind.Type("", newItem)

			expected := "goext." + util.ToGoName("null", typeOfItem)
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

			result := dbKind.Type("", newItem)

			expected := typeOfItem
			Expect(result).To(Equal(expected))
		})
	})

	Describe("interface type tests", func() {
		It("Should return a correct interface type for a null item", func() {
			newItem, err := CreateItem("int64")
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Required: false,
				Data:     map[interface{}]interface{}{"type": "int64"},
			})
			Expect(err).ToNot(HaveOccurred())

			result := dbKind.InterfaceType("", newItem)

			expected := "goext.NullInt"
			Expect(result).To(Equal(expected))
		})

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

			result := dbKind.InterfaceType("", newItem)

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

			result := dbKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`db:\"%s\" json:\"%s,omitempty\"`",
				name, name,
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

			result := dbKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`db:\"%s\" json:\"%s\"`",
				name, name,
			)
			Expect(result).To(Equal(expected))
		})
	})

	Describe("default value tests", func() {
		It("Should return a correct default value for a null item", func() {
			value := 1
			typeOfItem := "int"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			err = newItem.Parse(ParseContext{
				Defaults: value,
				Data:     map[interface{}]interface{}{"type": typeOfItem},
			})
			Expect(err).ToNot(HaveOccurred())

			result := dbKind.Default("", newItem)

			expected := fmt.Sprintf(
				"goext.MakeNullInt(%d)",
				value,
			)
			Expect(result).To(Equal(expected))
		})

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

			result := dbKind.Default("", newItem)

			expected := "1"
			Expect(result).To(Equal(expected))
		})
	})
})
