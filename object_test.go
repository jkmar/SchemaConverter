package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("object tests", func() {
	Describe("type tests", func() {
		It("Should return correct object type", func() {
			itemType := "ab"
			object := Object{objectType: itemType}
			expected := toGoName(itemType, "")
			result := object.Type("")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("is object tests", func() {
		It("Should return true for is object query", func() {
			item := &Object{}
			Expect(item.IsObject()).To(BeTrue())
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			item   *Object
			object map[interface{}]interface{}
		)

		BeforeEach(func() {
			item = &Object{}
		})

		It("Should return error for object with no properties", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"invalid schema: object %s does not have properties",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for object with invalid property definition", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": 1,
				},
			}
			expected := fmt.Errorf(
				"invalid schema: object %s has invalid property a",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for object with invalid property", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"a": 1,
					},
				},
			}
			expected := fmt.Errorf(
				"invalid schema: property %s does not have a type",
				addName(prefix, "a"),
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid object", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"type": "string",
					},
					"b": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "string",
						},
					},
					"c": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"x": map[interface{}]interface{}{
								"type": "boolean",
							},
							"y": map[interface{}]interface{}{
								"type": "string",
							},
						},
					},
				},
			}
			err := item.Parse(prefix, object)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.objectType).To(Equal(prefix))
			Expect(len(item.Properties)).To(Equal(len(object["properties"].(map[interface{}]interface{}))))
			names := make([]string, len(item.Properties))
			for i, property := range item.Properties {
				names[i] = property.Name
				switch names[i] {
				case "a":
					Expect(property.Item.Type("")).To(Equal("string"))
				case "b":
					Expect(property.Item.Type("")).To(Equal("[]string"))
				case "c":
					Expect(property.Item.Type("")).To(Equal("AbcC"))
				}
			}
			sort.Strings(names)
			Expect(names).To(Equal([]string{"a", "b", "c"}))
		})
	})

	Describe("collect tests", func() {
		var (
			names = []string{"A", "B", "C"}
			nested = &Object{
				names[0],
				[]*Property{
					&Property{
						"1",
						&Object{
							names[1],
							[]*Property{
								&Property{
									"2",
									&Object{
										names[2],
										nil,
									},
								},
							},
						},
					},
				},
			}
		)

		It("Should return nil for non positive depth", func() {
			item := &Object{"", nil}
			Expect(item.Collect(0)).To(BeNil())
		})

		It("Should return correct objects", func() {
			item := nested
			result := item.Collect(2)
			Expect(len(result)).To(Equal(2))
			Expect(result[0].Type("")).To(Equal(names[0]))
			Expect(result[1].Type("")).To(Equal(names[1]))
		})

		It("Should return all objects for negative depth", func() {
			item := nested
			result := item.Collect(-1)
			Expect(len(result)).To(Equal(3))
			Expect(result[0].Type("")).To(Equal(names[0]))
			Expect(result[1].Type("")).To(Equal(names[1]))
			Expect(result[2].Type("")).To(Equal(names[2]))
		})
	})

	Describe("generate struct tests", func() {
		It("Should generate correct struct", func() {
			annotate := func(name string) string {
				return fmt.Sprintf("`annotation:\"%s\"`", name)
			}

			item := Object{
				"abc_def",
				[]*Property{
					&Property{"id", &PlainItem{"string"}},
					&Property{"ip", &Array{&PlainItem{"int64"}}},
					&Property{"abc", &Object{"xyz", []*Property{}}}},
			}
			expected := `type AbcDefSuffix struct {
	ID string ` + annotate("id") + `
	IP []int64 ` + annotate("ip") + `
	Abc XyzSuffix ` + annotate("abc") + `
}
`
			result := item.GenerateStruct("suffix", "annotation")
			Expect(result).To(Equal(expected))
		})
	})
})
