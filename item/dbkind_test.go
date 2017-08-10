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

			newItem.Parse(
				"",
				0,
				true,
				map[interface{}]interface{}{"type": list},
			)

			result := dbKind.Type("", newItem)

			expected := "sql." + util.ToGoName("null", typeOfItem)
			Expect(result).To(Equal(expected))
		})

		It("Should return a correct type for a not null item", func() {
			typeOfItem := "string"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			newItem.Parse(
				"",
				0,
				true,
				map[interface{}]interface{}{"type": typeOfItem},
			)

			result := dbKind.Type("", newItem)

			expected := typeOfItem
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

			newItem.Parse(
				"",
				0,
				true,
				map[interface{}]interface{}{"type": list},
			)

			result := dbKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`db:\"%s\"`",
				name,
			)
			Expect(result).To(Equal(expected))
		})

		It("Should return a correct annotation for a not null item", func() {
			name := "name"
			typeOfItem := "string"
			newItem, err := CreateItem(typeOfItem)
			Expect(err).ToNot(HaveOccurred())

			newItem.Parse(
				"",
				0,
				true,
				map[interface{}]interface{}{"type": typeOfItem},
			)

			result := dbKind.Annotation(name, newItem)

			expected := fmt.Sprintf(
				"`db:\"%s\"`",
				name,
			)
			Expect(result).To(Equal(expected))
		})
	})
})
