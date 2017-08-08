package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
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
			typeOfItem := "ab"
			object := Object{objectType: typeOfItem}

			result := object.Type("")

			expected := util.ToGoName(typeOfItem, "")
			Expect(result).To(Equal(expected))
		})
	})

	Describe("add properties tests", func() {
		var (
			object     *Object
			properties set.Set
			names      = []string{"a", "b"}
		)

		BeforeEach(func() {
			properties = set.New()
			for _, name := range names {
				properties.Insert(CreateProperty(name))
			}
			object = &Object{"", properties}
		})

		It("Should ignore an empty set", func() {
			err := object.AddProperties(nil, false)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties).To(Equal(properties))
		})

		It("Should add new property in safe mode", func() {
			newProperty := CreateProperty("c")
			newProperties := set.New()
			newProperties.Insert(newProperty)

			err := object.AddProperties(newProperties, true)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties.Size()).To(Equal(3))
			Expect(object.properties.Contains(newProperty)).To(BeTrue())
		})

		It("Should return an error for duplicate property in safe mode", func() {
			newProperty := CreateProperty("b")
			newProperties := set.New()
			newProperties.Insert(newProperty)

			err := object.AddProperties(newProperties, true)

			expected := fmt.Errorf(
				"object %s: multiple properties have the same name",
				"",
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
			Expect(object.properties.Size()).To(Equal(2))
		})

		It("Should not override a property in unsafe mode", func() {
			newType := "string"
			newProperty := CreatePropertyWithType("b", newType)
			newProperties := set.New()
			newProperties.Insert(newProperty)

			err := object.AddProperties(newProperties, false)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties.Size()).To(Equal(2))
			Expect(object.properties.Contains(newProperty)).To(BeTrue())

			array := object.properties.ToArray()
			Expect(array[1].(*Property).item).To(BeNil())
		})
	})

	Describe("parse tests", func() {
		var (
			prefix = "abc"
			object *Object
			data   map[interface{}]interface{}
		)

		BeforeEach(func() {
			object = &Object{}
		})

		It("Should return an error for an object with invalid properties type", func() {
			data = map[interface{}]interface{}{
				"properties": "string",
			}

			err := object.Parse(prefix, data)

			expected := fmt.Errorf(
				"object %s has invalid properties",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an empty object for an object with no properties", func() {
			data = map[interface{}]interface{}{}

			err := object.Parse(prefix, data)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties.Empty()).To(BeTrue())
			Expect(object.Name()).To(Equal(prefix))
		})

		It("Should return an error for an object with a non string property name", func() {
			data = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					1: "string",
				},
			}

			err := object.Parse(prefix, data)

			expected := fmt.Errorf(
				"object %s has property which name is not a string",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an object with invalid property definition", func() {
			data = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": 1,
				},
			}

			err := object.Parse(prefix, data)

			expected := fmt.Errorf(
				"object %s has invalid property a",
				prefix,
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an object with an invalid property", func() {
			data = map[interface{}]interface{}{
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"a": 1,
					},
				},
			}

			err := object.Parse(prefix, data)

			expected := fmt.Errorf(
				"property %s does not have a type",
				util.AddName(prefix, "a"),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse valid object", func() {
			data = map[interface{}]interface{}{
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

			err := object.Parse(prefix, data)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.objectType).To(Equal(prefix))
			Expect(object.properties.Size()).To(Equal(len(data["properties"].(map[interface{}]interface{}))))

			names := object.properties.ToArray()
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

			array := result.ToArray()

			expected := []set.Element{objects[1], objects[2]}
			Expect(array).To(Equal(expected))
		})

		It("Should return all objects for negative depth", func() {
			result, err := nested.CollectObjects(-1, 2)

			Expect(err).ToNot(HaveOccurred())

			array := result.ToArray()

			expected := []set.Element{objects[4], objects[2], objects[3]}
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

			array := result.ToArray()

			expected := []set.Element{properties[1], properties[2]}
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

			result := item.GenerateStruct("suffix", "annotation")

			expected := `type AbcDefSuffix struct {
	Abc XyzSuffix ` + annotate("abc") + `
	ID string ` + annotate("id") + `
	IP []int64 ` + annotate("ip") + `
}
`
			Expect(result).To(Equal(expected))
		})
	})
})
