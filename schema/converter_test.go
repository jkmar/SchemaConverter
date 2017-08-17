package schema

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("converter tess", func() {
	Describe("error tests", func() {
		var (
			validSchema = []map[interface{}]interface{}{
				{
					"id":     "my_id",
					"schema": map[interface{}]interface{}{},
				},
			}
			invalidSchema = []map[interface{}]interface{}{
				{
					"invalid schema": "test",
				},
			}
		)

		Describe("parse all errors", func() {
			var expected = fmt.Errorf("schema does not have an id")

			It("Should return error for invalid other schema", func() {
				_, _, _, err := Convert(validSchema, invalidSchema, "")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})

			It("Should return error for invalid other schema", func() {
				_, _, _, err := Convert(invalidSchema, validSchema, "")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})
		})

		Describe("collect errors", func() {
			It("Should return error for multiple schemas with the same name", func() {
				_, _, _, err := Convert(validSchema, validSchema, "")

				expected := fmt.Errorf("multiple schemas with the same name")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})

			It("Should return error for multiple objects with the same name", func() {
				name := "a"
				schemas := []map[interface{}]interface{}{
					{
						"id": name,
						"schema": map[interface{}]interface{}{
							"type": "object",
							"properties": map[interface{}]interface{}{
								"__": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"_": map[interface{}]interface{}{
											"type":       "object",
											"properties": map[interface{}]interface{}{},
										},
									},
								},
								"_": map[interface{}]interface{}{
									"type": "object",
									"properties": map[interface{}]interface{}{
										"__": map[interface{}]interface{}{
											"type":       "object",
											"properties": map[interface{}]interface{}{},
										},
									},
								},
							},
						},
					},
				}

				_, _, _, err := Convert(nil, schemas, "")

				expected := fmt.Errorf(
					"invalid schema %s: multiple objects with the same type at object %s",
					name,
					name,
				)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expected))
			})
		})
	})
	Describe("valid data tests", func() {
		It("Should convert valid schemas", func() {
			other := []map[interface{}]interface{}{
				{
					"id": "base",
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"id": map[interface{}]interface{}{
								"type": "string",
							},
							"ip": map[interface{}]interface{}{
								"type": "number",
							},
							"object": map[interface{}]interface{}{
								"type": "object",
								"properties": map[interface{}]interface{}{
									"x": map[interface{}]interface{}{
										"type":    "string",
										"default": "abc",
									},
									"y": map[interface{}]interface{}{
										"type": "string",
									},
								},
								"required": []interface{}{
									"y",
								},
							},
						},
					},
				},
				{
					"id":     "middle",
					"parent": "parent",
					"extends": []interface{}{
						"base",
					},
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"null": map[interface{}]interface{}{
								"type": []interface{}{
									"boolean",
									"null",
								},
							},
							"array": map[interface{}]interface{}{
								"type": "array",
								"items": map[interface{}]interface{}{
									"type": "array",
									"items": map[interface{}]interface{}{
										"type": "number",
									},
								},
							},
							"nested": map[interface{}]interface{}{
								"type": "object",
								"properties": map[interface{}]interface{}{
									"first": map[interface{}]interface{}{
										"type": "object",
										"properties": map[interface{}]interface{}{
											"second": map[interface{}]interface{}{
												"type":       "object",
												"properties": map[interface{}]interface{}{},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			toConvert := []map[interface{}]interface{}{
				{
					"id": "general",
					"extends": []interface{}{
						"middle",
						"base",
					},
					"schema": map[interface{}]interface{}{
						"type": "object",
						"properties": map[interface{}]interface{}{
							"complex": map[interface{}]interface{}{
								"type": "array",
								"items": map[interface{}]interface{}{
									"type": "array",
									"items": map[interface{}]interface{}{
										"type": "object",
										"properties": map[interface{}]interface{}{
											"for": map[interface{}]interface{}{
												"type": "number",
											},
											"int": map[interface{}]interface{}{
												"type": "boolean",
											},
										},
									},
								},
							},
							"tree": map[interface{}]interface{}{
								"type": "object",
								"properties": map[interface{}]interface{}{
									"left": map[interface{}]interface{}{
										"type": "object",
										"properties": map[interface{}]interface{}{
											"leaf_first": map[interface{}]interface{}{
												"type": "string",
											},
											"leaf_second": map[interface{}]interface{}{
												"type": "object",
												"properties": map[interface{}]interface{}{
													"value": map[interface{}]interface{}{
														"type": "number",
													},
												},
											},
										},
									},
									"right": map[interface{}]interface{}{
										"type": "object",
										"properties": map[interface{}]interface{}{
											"leaf_third": map[interface{}]interface{}{
												"type": "array",
												"items": map[interface{}]interface{}{
													"type": "boolean",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					"id": "only_derive",
					"extends": []interface{}{
						"base",
					},
					"schema": map[interface{}]interface{}{},
				},
				{
					"id":     "empty",
					"schema": map[interface{}]interface{}{},
				},
			}

			interfaces, structs, implementations, err := Convert(other, toConvert, "")
			Expect(err).ToNot(HaveOccurred())

			generalInterface := `type IGeneral interface {
	GetArray() [][]int64
	SetArray([][]int64)
	GetComplex() [][]IGeneralComplex
	SetComplex([][]IGeneralComplex)
	GetID() string
	SetID(string)
	GetIP() goext.NullInt
	SetIP(goext.NullInt)
	GetNested() IMiddleNested
	SetNested(IMiddleNested)
	GetNull() goext.NullBool
	SetNull(goext.NullBool)
	GetObject() IBaseObject
	SetObject(IBaseObject)
	GetParentID() string
	SetParentID(string)
	GetTree() IGeneralTree
	SetTree(IGeneralTree)
}
`
			onlyDeriveInterface := `type IOnlyDerive interface {
	GetID() string
	SetID(string)
	GetIP() goext.NullInt
	SetIP(goext.NullInt)
	GetObject() IBaseObject
	SetObject(IBaseObject)
}
`
			emptyInterface := `type IEmpty interface {
}
`
			generalTreeLeftLeafSecondInterface := `type IGeneralTreeLeftLeafSecond interface {
	GetValue() int64
	SetValue(int64)
}
`
			middleNestedFirstInterface := `type IMiddleNestedFirst interface {
	GetSecond() IMiddleNestedFirstSecond
	SetSecond(IMiddleNestedFirstSecond)
}
`
			middleNestedInterface := `type IMiddleNested interface {
	GetFirst() IMiddleNestedFirst
	SetFirst(IMiddleNestedFirst)
}
`
			generalComplexInterface := `type IGeneralComplex interface {
	GetFor() int64
	SetFor(int64)
	GetInt() bool
	SetInt(bool)
}
`
			generalTreeLeftInterface := `type IGeneralTreeLeft interface {
	GetLeafFirst() string
	SetLeafFirst(string)
	GetLeafSecond() IGeneralTreeLeftLeafSecond
	SetLeafSecond(IGeneralTreeLeftLeafSecond)
}
`
			generalTreeInterface := `type IGeneralTree interface {
	GetLeft() IGeneralTreeLeft
	SetLeft(IGeneralTreeLeft)
	GetRight() IGeneralTreeRight
	SetRight(IGeneralTreeRight)
}
`
			baseObjectInterface := `type IBaseObject interface {
	GetX() string
	SetX(string)
	GetY() string
	SetY(string)
}
`
			generalTreeRightInterface := `type IGeneralTreeRight interface {
	GetLeafThird() []bool
	SetLeafThird([]bool)
}
`
			middleNestedFirstSecondInterface := `type IMiddleNestedFirstSecond interface {
}
`
			generalStruct := `type General struct {
	Array [][]int64 ` + "`" + `db:"array"` + "`" + `
	Complex [][]*GeneralComplex ` + "`" + `db:"complex"` + "`" + `
	ID string ` + "`" + `db:"id"` + "`" + `
	IP goext.NullInt ` + "`" + `db:"ip"` + "`" + `
	Nested *MiddleNested ` + "`" + `db:"nested"` + "`" + `
	Null goext.NullBool ` + "`" + `db:"null"` + "`" + `
	Object *BaseObject ` + "`" + `db:"object"` + "`" + `
	ParentID string ` + "`" + `db:"parent_id"` + "`" + `
	Tree *GeneralTree ` + "`" + `db:"tree"` + "`" + `
}
`
			onlyDeriveStruct := `type OnlyDerive struct {
	ID string ` + "`" + `db:"id"` + "`" + `
	IP goext.NullInt ` + "`" + `db:"ip"` + "`" + `
	Object *BaseObject ` + "`" + `db:"object"` + "`" + `
}
`
			emptyStruct := `type Empty struct {
}
`
			generalTreeLeftLeafSecondStruct := `type GeneralTreeLeftLeafSecond struct {
	Value int64 ` + "`" + `json:"value,omitempty"` + "`" + `
}
`
			middleNestedFirstStruct := `type MiddleNestedFirst struct {
	Second *MiddleNestedFirstSecond ` + "`" + `json:"second"` + "`" + `
}
`
			middleNestedStruct := `type MiddleNested struct {
	First *MiddleNestedFirst ` + "`" + `json:"first"` + "`" + `
}
`
			generalComplexStruct := `type GeneralComplex struct {
	For int64 ` + "`" + `json:"for,omitempty"` + "`" + `
	Int bool ` + "`" + `json:"int,omitempty"` + "`" + `
}
`
			generalTreeLeftStruct := `type GeneralTreeLeft struct {
	LeafFirst string ` + "`" + `json:"leaf_first,omitempty"` + "`" + `
	LeafSecond *GeneralTreeLeftLeafSecond ` + "`" + `json:"leaf_second"` + "`" + `
}
`
			generalTreeStruct := `type GeneralTree struct {
	Left *GeneralTreeLeft ` + "`" + `json:"left"` + "`" + `
	Right *GeneralTreeRight ` + "`" + `json:"right"` + "`" + `
}
`
			baseObjectStruct := `type BaseObject struct {
	X string ` + "`" + `json:"x"` + "`" + `
	Y string ` + "`" + `json:"y"` + "`" + `
}
`
			generalTreeRightStruct := `type GeneralTreeRight struct {
	LeafThird []bool ` + "`" + `json:"leaf_third"` + "`" + `
}
`
			middleNestedFirstSecondStruct := `type MiddleNestedFirstSecond struct {
}
`
			generalImplementation := `func (general *General) GetArray() [][]int64 {
	return general.Array
}

func (general *General) SetArray(array [][]int64) {
	general.Array = make([][]int64, len(array))
	for i := range array {
		general.Array[i] = array[i]
	}
}

func (general *General) GetComplex() [][]IGeneralComplex {
	result := make([][]IGeneralComplex, len(general.Complex))
	for i := range general.Complex {
		result[i] = make([]IGeneralComplex, len(general.Complex[i]))
		for j := range general.Complex[i] {
			result[i][j] = general.Complex[i][j]
		}
	}
	return result
}

func (general *General) SetComplex(complex [][]IGeneralComplex) {
	general.Complex = make([][]*GeneralComplex, len(complex))
	for i := range complex {
		general.Complex[i] = make([]*GeneralComplex, len(complex[i]))
		for j := range complex[i] {
			general.Complex[i][j] = complex[i][j].(*GeneralComplex)
		}
	}
}

func (general *General) GetID() string {
	return general.ID
}

func (general *General) SetID(id string) {
	general.ID = id
}

func (general *General) GetIP() goext.NullInt {
	return general.IP
}

func (general *General) SetIP(ip goext.NullInt) {
	general.IP = ip
}

func (general *General) GetNested() IMiddleNested {
	return general.Nested
}

func (general *General) SetNested(nested IMiddleNested) {
	general.Nested = nested.(*MiddleNested)
}

func (general *General) GetNull() goext.NullBool {
	return general.Null
}

func (general *General) SetNull(null goext.NullBool) {
	general.Null = null
}

func (general *General) GetObject() IBaseObject {
	return general.Object
}

func (general *General) SetObject(object IBaseObject) {
	general.Object = object.(*BaseObject)
}

func (general *General) GetParentID() string {
	return general.ParentID
}

func (general *General) SetParentID(parentID string) {
	general.ParentID = parentID
}

func (general *General) GetTree() IGeneralTree {
	return general.Tree
}

func (general *General) SetTree(tree IGeneralTree) {
	general.Tree = tree.(*GeneralTree)
}
`
			onlyDeriveImplementation := `func (onlyDerive *OnlyDerive) GetID() string {
	return onlyDerive.ID
}

func (onlyDerive *OnlyDerive) SetID(id string) {
	onlyDerive.ID = id
}

func (onlyDerive *OnlyDerive) GetIP() goext.NullInt {
	return onlyDerive.IP
}

func (onlyDerive *OnlyDerive) SetIP(ip goext.NullInt) {
	onlyDerive.IP = ip
}

func (onlyDerive *OnlyDerive) GetObject() IBaseObject {
	return onlyDerive.Object
}

func (onlyDerive *OnlyDerive) SetObject(object IBaseObject) {
	onlyDerive.Object = object.(*BaseObject)
}
`
			emptyImplementation := ``
			generalTreeLeftLeafSecondImplementation := `func (generalTreeLeftLeafSecond *GeneralTreeLeftLeafSecond) GetValue() int64 {
	return generalTreeLeftLeafSecond.Value
}

func (generalTreeLeftLeafSecond *GeneralTreeLeftLeafSecond) SetValue(value int64) {
	generalTreeLeftLeafSecond.Value = value
}
`
			middleNestedFirstImplementation := `func (middleNestedFirst *MiddleNestedFirst) GetSecond() IMiddleNestedFirstSecond {
	return middleNestedFirst.Second
}

func (middleNestedFirst *MiddleNestedFirst) SetSecond(second IMiddleNestedFirstSecond) {
	middleNestedFirst.Second = second.(*MiddleNestedFirstSecond)
}
`
			middleNestedImplementation := `func (middleNested *MiddleNested) GetFirst() IMiddleNestedFirst {
	return middleNested.First
}

func (middleNested *MiddleNested) SetFirst(first IMiddleNestedFirst) {
	middleNested.First = first.(*MiddleNestedFirst)
}
`
			generalComplexImplementation := `func (generalComplex *GeneralComplex) GetFor() int64 {
	return generalComplex.For
}

func (generalComplex *GeneralComplex) SetFor(forObject int64) {
	generalComplex.For = forObject
}

func (generalComplex *GeneralComplex) GetInt() bool {
	return generalComplex.Int
}

func (generalComplex *GeneralComplex) SetInt(int bool) {
	generalComplex.Int = int
}
`
			generalTreeLeftImplementation := `func (generalTreeLeft *GeneralTreeLeft) GetLeafFirst() string {
	return generalTreeLeft.LeafFirst
}

func (generalTreeLeft *GeneralTreeLeft) SetLeafFirst(leafFirst string) {
	generalTreeLeft.LeafFirst = leafFirst
}

func (generalTreeLeft *GeneralTreeLeft) GetLeafSecond() IGeneralTreeLeftLeafSecond {
	return generalTreeLeft.LeafSecond
}

func (generalTreeLeft *GeneralTreeLeft) SetLeafSecond(leafSecond IGeneralTreeLeftLeafSecond) {
	generalTreeLeft.LeafSecond = leafSecond.(*GeneralTreeLeftLeafSecond)
}
`
			generalTreeImplementation := `func (generalTree *GeneralTree) GetLeft() IGeneralTreeLeft {
	return generalTree.Left
}

func (generalTree *GeneralTree) SetLeft(left IGeneralTreeLeft) {
	generalTree.Left = left.(*GeneralTreeLeft)
}

func (generalTree *GeneralTree) GetRight() IGeneralTreeRight {
	return generalTree.Right
}

func (generalTree *GeneralTree) SetRight(right IGeneralTreeRight) {
	generalTree.Right = right.(*GeneralTreeRight)
}
`
			baseObjectImplementation := `func (baseObject *BaseObject) GetX() string {
	return baseObject.X
}

func (baseObject *BaseObject) SetX(x string) {
	baseObject.X = x
}

func (baseObject *BaseObject) GetY() string {
	return baseObject.Y
}

func (baseObject *BaseObject) SetY(y string) {
	baseObject.Y = y
}
`
			generalTreeRightImplementation := `func (generalTreeRight *GeneralTreeRight) GetLeafThird() []bool {
	return generalTreeRight.LeafThird
}

func (generalTreeRight *GeneralTreeRight) SetLeafThird(leafThird []bool) {
	generalTreeRight.LeafThird = leafThird
}
`
			middleNestedFirstSecondImplementation := ``

			sort.Strings(interfaces)
			sort.Strings(structs)
			sort.Strings(implementations)

			expectedInterfaces := []string{
				baseObjectInterface,
				emptyInterface,
				generalInterface,
				generalComplexInterface,
				generalTreeInterface,
				generalTreeLeftInterface,
				generalTreeLeftLeafSecondInterface,
				generalTreeRightInterface,
				middleNestedInterface,
				middleNestedFirstInterface,
				middleNestedFirstSecondInterface,
				onlyDeriveInterface,
			}
			expectedStructs := []string{
				baseObjectStruct,
				emptyStruct,
				generalStruct,
				generalComplexStruct,
				generalTreeStruct,
				generalTreeLeftStruct,
				generalTreeLeftLeafSecondStruct,
				generalTreeRightStruct,
				middleNestedStruct,
				middleNestedFirstStruct,
				middleNestedFirstSecondStruct,
				onlyDeriveStruct,
			}
			expectedImplementations := []string{
				emptyImplementation,
				middleNestedFirstSecondImplementation,
				baseObjectImplementation,
				generalImplementation,
				generalComplexImplementation,
				generalTreeImplementation,
				generalTreeLeftImplementation,
				generalTreeLeftLeafSecondImplementation,
				generalTreeRightImplementation,
				middleNestedImplementation,
				middleNestedFirstImplementation,
				onlyDeriveImplementation,
			}
			Expect(interfaces).To(Equal(expectedInterfaces))
			Expect(structs).To(Equal(expectedStructs))
			Expect(implementations).To(Equal(expectedImplementations))
		})
	})
})
