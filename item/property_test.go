package item

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zimnx/YamlSchemaToGoStruct/set"
	"github.com/zimnx/YamlSchemaToGoStruct/util"
)

var _ = Describe("property tests", func() {
	Describe("hash tests", func() {
		Describe("to string tests", func() {
			It("Should return a correct string representation of a property", func() {
				name := "name"
				property := CreateProperty(name)

				result := property.ToString()

				Expect(result).To(Equal(name))
			})
		})

		Describe("compress tests", func() {
			It("Should compress if destination is owned by the array", func() {
				source := &PlainItem{itemType: "1"}
				destination := &PlainItem{itemType: "2"}
				property := Property{item: destination}

				property.Compress(source, destination)

				Expect(source).ToNot(BeIdenticalTo(destination))
				Expect(property.item).To(BeIdenticalTo(source))
			})

			It("Should not compress if destination is not owned by the array", func() {
				source := &PlainItem{itemType: "1"}
				destination := &PlainItem{itemType: "2"}
				property := Property{item: destination}

				property.Compress(destination, source)

				Expect(source).ToNot(BeIdenticalTo(destination))
				Expect(property.item).To(BeIdenticalTo(destination))
			})
		})

		Describe("get children tests", func() {
			It("Should return a correct children set", func() {
				plainItem := &PlainItem{itemType: "test"}
				property := Property{item: plainItem}

				result := property.GetChildren()

				Expect(len(result)).To(Equal(1))
				Expect(result[0]).To(BeIdenticalTo(plainItem))
			})
		})
	})

	Describe("creation tests", func() {
		It("Should create a property with a correct name", func() {
			name := "name"
			property := CreateProperty(name)

			Expect(property.name).To(Equal(name))
		})
	})

	Describe("is object tests", func() {
		It("Should return false for not an object", func() {
			property := CreateProperty("")
			property.item, _ = CreateItem("string")

			Expect(property.IsObject()).To(BeFalse())
		})

		It("Should return false for an object", func() {
			property := CreateProperty("")
			property.item, _ = CreateItem("object")

			Expect(property.IsObject()).To(BeTrue())
		})
	})

	Describe("add properties tests", func() {
		It("Should add properties for an item", func() {
			properties := set.New()
			properties.Insert(CreateProperty("a"))
			object := &Object{}
			property := &Property{
				name: "",
				item: object,
			}

			err := property.AddProperties(properties, true)

			Expect(err).ToNot(HaveOccurred())
			Expect(object.properties).To(Equal(properties))
		})
	})

	Describe("parse tests", func() {
		var (
			prefix   = "abc"
			property *Property
			data     map[interface{}]interface{}
		)

		BeforeEach(func() {
			property = &Property{name: "def"}
		})

		It("Should return an error for an object with no items", func() {
			data = map[interface{}]interface{}{}

			err := property.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"property %s does not have a type",
				util.AddName(prefix, property.name),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for a property with an invalid type", func() {
			data = map[interface{}]interface{}{
				"type": 1,
			}

			err := property.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"property %s: unsupported type: %T",
				util.AddName(prefix, property.name),
				data["type"],
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return an error for an invalid item", func() {
			data = map[interface{}]interface{}{
				"type": "array",
			}

			err := property.Parse(prefix, 0, true, data)

			expected := fmt.Errorf(
				"array %s does not have items",
				util.AddName(prefix, property.name),
			)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should parse a valid property", func() {
			data = map[interface{}]interface{}{
				"type": "array",
				"items": map[interface{}]interface{}{
					"type": "string",
				},
			}

			err := property.Parse(prefix, 0, true, data)

			typeOfItem := data["items"].(map[interface{}]interface{})["type"]
			expected := "[]" + typeOfItem.(string)
			Expect(err).ToNot(HaveOccurred())
			Expect(property.item.Type("")).To(Equal(expected))
		})
	})

	Describe("collect objects tests", func() {
		It("Should collect object", func() {
			object := &Object{objectType: "abc"}
			property := &Property{
				name: "",
				item: object,
			}
			objects := set.New()
			objects.Insert(object)

			result, err := property.CollectObjects(1, 0)

			expected := objects
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Describe("generate property tests", func() {
		var (
			prefix   = "abc"
			suffix   = "xyz"
			property *Property
			data     map[interface{}]interface{}
		)

		Describe("db property tests", func() {
			const annotation = "db"

			BeforeEach(func() {
				property = &Property{
					name: "def_id",
				}
			})

			It("Should generate a correct property for a null item", func() {
				data = map[interface{}]interface{}{
					"type": []interface{}{
						"string",
						"null",
					},
				}

				err := property.Parse(prefix, 0, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID goext.NullString `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for a plain item", func() {
				data = map[interface{}]interface{}{
					"type": "boolean",
				}

				err := property.Parse(prefix, 0, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID bool `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for an array", func() {
				data = map[interface{}]interface{}{
					"type": "array",
					"items": map[interface{}]interface{}{
						"type": "integer",
					},
				}

				err := property.Parse(prefix, 0, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID []int64 `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for an object", func() {
				data = map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"test": map[interface{}]interface{}{
							"type": "string",
						},
					},
				}

				err := property.Parse(prefix, 0, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID *AbcDefIDXyz `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})
		})

		Describe("json property tests", func() {
			const annotation = "json"

			BeforeEach(func() {
				property = &Property{
					name: "def_id",
				}
			})

			It("Should generate a correct property for a null item", func() {
				data = map[interface{}]interface{}{
					"type": []interface{}{
						"string",
						"null",
					},
				}

				err := property.Parse(prefix, 2, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID string `%s:\"%s,omitempty\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for a plain item", func() {
				data = map[interface{}]interface{}{
					"type": "boolean",
				}

				err := property.Parse(prefix, 2, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID bool `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for an array", func() {
				data = map[interface{}]interface{}{
					"type": "array",
					"items": map[interface{}]interface{}{
						"type": "integer",
					},
				}

				err := property.Parse(prefix, 2, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID []int64 `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})

			It("Should generate a correct property for an object", func() {
				data = map[interface{}]interface{}{
					"type": "object",
					"properties": map[interface{}]interface{}{
						"test": map[interface{}]interface{}{
							"type": "string",
						},
					},
				}

				err := property.Parse(prefix, 2, true, data)
				Expect(err).ToNot(HaveOccurred())

				result := property.GenerateProperty(suffix)

				expected := fmt.Sprintf(
					"\tDefID *AbcDefIDXyz `%s:\"%s\"`\n",
					annotation,
					property.name,
				)
				Expect(result).To(Equal(expected))
			})
		})
	})

	Describe("getter header tests", func() {
		It("Should generate a correct getter header for a plain item", func() {
			property := &Property{
				name: "name",
				item: &PlainItem{itemType: "string", null: true},
				kind: &DBKind{},
			}

			result := property.GetterHeader("suffix")

			expected := "GetName() goext.NullString"
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct getter header for an object", func() {
			property := &Property{
				name: "name",
				item: &Object{objectType: "test"},
				kind: &DBKind{},
			}

			result := property.GetterHeader("suffix")

			expected := "GetName() ITestSuffix"
			Expect(result).To(Equal(expected))
		})
	})

	Describe("setter header tests", func() {
		It("Should generate a correct setter header for a plain item", func() {
			property := &Property{
				name: "name",
				item: &PlainItem{itemType: "string", null: true},
				kind: &DBKind{},
			}

			result := property.SetterHeader("suffix", true)

			expected := "SetName(name goext.NullString)"
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct setter header for an object", func() {
			property := &Property{
				name: "name",
				item: &Object{objectType: "test"},
				kind: &DBKind{},
			}

			result := property.SetterHeader("suffix", false)

			expected := "SetName(ITestSuffix)"
			Expect(result).To(Equal(expected))
		})
	})

	Describe("generate getter tests", func() {
		It("Should generate a correct getter for a plain item", func() {
			property := &Property{
				name: "def",
				item: &PlainItem{itemType: "int64", null: true},
				kind: &DBKind{},
			}

			result := property.GenerateGetter("var", "")

			expected := `GetDef() goext.NullInt {
	return var.Def
}`
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct getter for an object", func() {
			property := &Property{
				name: "abc",
				item: &Object{objectType: "xyz"},
				kind: &DBKind{},
			}

			result := property.GenerateGetter("var", "")

			expected := `GetAbc() IXyz {
	return var.Abc
}`
			Expect(result).To(Equal(expected))
		})
	})

	Describe("generate setter tests", func() {
		It("Should generate a correct setter for a plain item", func() {
			property := &Property{
				name: "def",
				item: &PlainItem{itemType: "int64", null: true},
				kind: &DBKind{},
			}

			result := property.GenerateSetter("var", "", "")

			expected := `SetDef(def goext.NullInt) {
	var.Def = def
}`
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct setter for an object", func() {
			property := &Property{
				name: "range",
				item: &Object{objectType: "xyz"},
				kind: &DBKind{},
			}

			result := property.GenerateSetter("var", "", "")

			expected := `SetRange(rangeObject IXyz) {
	var.Range = rangeObject.(*Xyz)
}`
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct setter for an array", func() {
			property := &Property{
				name: "a",
				item: &Array{&Array{&Object{objectType: "object"}}},
				kind: &DBKind{},
			}

			result := property.GenerateSetter("var", "", "")

			expected := `SetA(a [][]IObject) {
	var.A = make([][]*Object, len(a))
	for i := range a {
		var.A[i] = make([]*Object, len(a[i]))
		for j := range a[i] {
			var.A[i][j] = a[i][j].(*Object)
		}
	}
}`
			Expect(result).To(Equal(expected))
		})
	})

	Describe("compression tests", func() {
		It("Should compress exactly identical objects", func() {
			data := map[interface{}]interface{}{
				"type": "object",
				"properties": map[interface{}]interface{}{
					"a": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"x": map[interface{}]interface{}{
									"type": "string",
								},
								"y": map[interface{}]interface{}{
									"type": "number",
								},
							},
						},
					},
					"b": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"x": map[interface{}]interface{}{
									"type": "string",
								},
								"y": map[interface{}]interface{}{
									"type": "number",
								},
							},
						},
					},
					"c": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"x": map[interface{}]interface{}{
									"type": "string",
								},
								"y": map[interface{}]interface{}{
									"type": "number",
								},
							},
						},
					},
					"d": map[interface{}]interface{}{
						"type": "array",
						"items": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"x": map[interface{}]interface{}{
									"type": "string",
								},
								"y": map[interface{}]interface{}{
									"type": "string",
								},
							},
						},
					},
				},
			}
			property := &Property{name: "test"}
			err := property.Parse("test", 0, true, data)
			Expect(err).ToNot(HaveOccurred())

			property.CompressObjects()

			properties := property.item.GetChildren()
			Expect(len(properties)).To(Equal(4))
			Expect(properties[0].(*Property).item).To(
				BeIdenticalTo(properties[1].(*Property).item),
			)
			Expect(properties[1].(*Property).item).To(
				BeIdenticalTo(properties[2].(*Property).item),
			)
			Expect(properties[2].(*Property).item).ToNot(
				BeIdenticalTo(properties[3].(*Property).item),
			)

			objects, err := property.CollectObjects(-1, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(objects.Size()).To(Equal(3))
		})
	})
})
