package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
)

var _ = Describe("object tests", func() {
	Describe("name tests", func() {
		It("Should return a correct object name", func() {
			name := "abc_abc"
			object := Object{objectType: name}
			Expect(object.Name()).To(Equal(name))
		})
	})

	Describe("type tests", func() {
		It("Should return a correct object type", func() {
			itemType := "ab"
			object := Object{objectType: itemType}
			expected := toGoName(itemType, "")
			result := object.Type("")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("is object tests", func() {
		It("Should return true for an object query", func() {
			item := &Object{}
			Expect(item.IsObject()).To(BeTrue())
		})
	})

	Describe("add properties tests", func() {
		var (
			item       *Object
			properties set.Set
			names      = []string{"a", "b"}
		)

		BeforeEach(func() {
			properties = set.New()
			for _, name := range names {
				properties.Insert(CreateProperty(name))
			}
			item = &Object{"", properties}
		})

		It("Should ignore an empty set", func() {
			err := item.AddProperties(nil, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.properties).To(Equal(properties))
		})

		It("Should add new property in safe mode", func() {
			newProperty := CreateProperty("c")
			newProperties := set.New()
			newProperties.Insert(newProperty)
			err := item.AddProperties(newProperties, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.properties.Size()).To(Equal(3))
			Expect(item.properties.Contains(newProperty)).To(BeTrue())
		})

		It("Should return an error for duplicate property in safe mode", func() {
			newProperty := CreateProperty("b")
			newProperties := set.New()
			newProperties.Insert(newProperty)
			expected := fmt.Errorf(
				"object %s: multiple properties have the same name",
				"",
			)
			err := item.AddProperties(newProperties, true)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
			Expect(item.properties.Size()).To(Equal(2))
		})

		It("Should not override a property in unsafe mode", func() {
			newType := "string"
			newProperty := CreatePropertyWithType("b", newType)
			newProperties := set.New()
			newProperties.Insert(newProperty)
			err := item.AddProperties(newProperties, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.properties.Size()).To(Equal(2))
			Expect(item.properties.Contains(newProperty)).To(BeTrue())
			array := item.properties.ToArray()
			Expect(array[1].(*Property).item).To(BeNil())
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

		It("Should return an error for an object with no properties", func() {
			object = map[interface{}]interface{}{}
			expected := fmt.Errorf(
				"object %s does not have properties",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an object with a non string property name", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					1: "string",
				},
			}
			expected := fmt.Errorf(
				"object %s has property which name is not a string",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an object with invalid property definition", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": 1,
				},
			}
			expected := fmt.Errorf(
				"object %s has invalid property a",
				prefix,
			)
			err := item.Parse(prefix, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an object with an invalid property", func() {
			object = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"a": 1,
					},
				},
			}
			expected := fmt.Errorf(
				"property %s does not have a type",
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
			Expect(item.properties.Size()).To(Equal(len(object["properties"].(map[interface{}]interface{}))))
			names := item.properties.ToArray()
			Expect(names[0].(*Property).item.Type("")).To(Equal("string"))
			Expect(names[1].(*Property).item.Type("")).To(Equal("[]string"))
			Expect(names[2].(*Property).item.Type("")).To(Equal("AbcC"))
			Expect(names[0].(*Property).Name()).To(Equal("a"))
			Expect(names[1].(*Property).Name()).To(Equal("b"))
			Expect(names[2].(*Property).Name()).To(Equal("c"))
		})
	})

	Describe("collect objects tests", func() {
		var (
			names      []string
			objects    []*Object
			properties []*Property
			nested     *Object
		)

		BeforeEach(func() {
			names = []string{"a", "b", "c", "d", "b"}
			objects = make([]*Object, len(names))
			properties = make([]*Property, len(names)-1)
			objects[0] = &Object{}
			objects[0].objectType = names[0]
			for i := 0; i < len(properties); i++ {
				properties[i] = &Property{}
				objects[i+1] = &Object{}
				properties[i].item = objects[i+1]
				objects[i+1].objectType = names[i+1]
				set := set.New()
				set.Insert(properties[i])
				objects[i].properties = set
			}
			objects[len(objects)-1].properties = set.New()
			nested = objects[0]
		})

		It("Should return an empty set for offset greater or equal limit", func() {
			result, err := nested.CollectObjects(1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Empty()).To(BeTrue())
		})

		It("Should return correct objects", func() {
			result, err := nested.CollectObjects(3, 1)
			Expect(err).ToNot(HaveOccurred())
			expected := []set.Element{objects[1], objects[2]}
			array := result.ToArray()
			Expect(array).To(Equal(expected))
		})

		It("Should return all objects for negative depth", func() {
			result, err := nested.CollectObjects(-1, 2)
			Expect(err).ToNot(HaveOccurred())
			expected := []set.Element{objects[4], objects[2], objects[3]}
			array := result.ToArray()
			Expect(array).To(Equal(expected))
		})

		It("Should return an error for multiple objects with the same name", func() {
			_, err := nested.CollectObjects(-1, 0)
			expected := fmt.Errorf(
				"multiple objects with the same type at object %s",
				objects[len(objects)-1].Name(),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("collect tests", func() {
		var (
			names      []string
			objects    []*Object
			properties []*Property
			nested     *Object
		)

		BeforeEach(func() {
			names = []string{"a", "b", "c"}
			objects = make([]*Object, 4)
			properties = make([]*Property, 5)
			for i := 0; i < 4; i++ {
				objects[i] = &Object{objectType: string(i)}
				properties[i] = &Property{}
			}
			properties[4] = &Property{names[2], &PlainItem{"string"}}
			properties[3].item = &PlainItem{"string"}
			properties[3].name = names[2]
			objectSet := set.New()
			objectSet.Insert(properties[4])
			objects[3].properties = objectSet
			objectSet = set.New()
			objectSet.Insert(properties[3])
			objects[2].properties = objectSet
			properties[2].name = names[1]
			properties[2].item = objects[3]
			properties[1].name = names[0]
			properties[1].item = objects[2]
			objectSet = set.New()
			objectSet.Insert(properties[1])
			objectSet.Insert(properties[2])
			objects[1].properties = objectSet
			properties[0].name = names[0]
			properties[0].item = objects[1]
			objectSet = set.New()
			objectSet.Insert(properties[0])
			objects[0].properties = objectSet
			nested = objects[0]
		})

		It("Should return an empty set for offset greater or equal limit", func() {
			result, err := nested.CollectProperties(1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Empty()).To(BeTrue())
		})

		It("Should return correct objects", func() {
			result, err := nested.CollectProperties(2, 1)
			Expect(err).ToNot(HaveOccurred())
			expected := []set.Element{properties[1], properties[2]}
			array := result.ToArray()
			Expect(array).To(Equal(expected))
		})

		It("Should return an error for multiple properties with the same name at property", func() {
			_, err := nested.CollectProperties(2, 0)
			expected := fmt.Errorf(
				"multiple properties with the same name: %s",
				names[0],
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for multiple properties with the same name at object", func() {
			_, err := nested.CollectProperties(-1, 0)
			expected := fmt.Errorf(
				"multiple properties with the same name at object %s",
				objects[1].Name(),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("generate struct tests", func() {
		It("Should generate correct struct", func() {
			annotate := func(name string) string {
				return fmt.Sprintf("`annotation:\"%s\"`", name)
			}

			properties := set.New()
			properties.Insert(&Property{"id", &PlainItem{"string"}})
			properties.Insert(&Property{"ip", &Array{&PlainItem{"int64"}}})
			properties.Insert(&Property{"abc", &Object{"xyz", nil}})

			item := Object{
				"abc_def",
				properties,
			}
			expected := `type AbcDefSuffix struct {
	Abc XyzSuffix ` + annotate("abc") + `
	ID string ` + annotate("id") + `
	IP []int64 ` + annotate("ip") + `
}
`
			result := item.GenerateStruct("suffix", "annotation")
			Expect(result).To(Equal(expected))
		})
	})
})
