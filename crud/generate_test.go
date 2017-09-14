package crud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate tests", func() {
	Describe("Fetch tests", func() {
		It("Should generate a correct fetch function", func() {
			expected := `func FetchA(` +
				`schema goext.ISchema, ` +
				`id string, ` +
				`context goext.Context` +
				`) (esi.IA, error) {
	result, err := schema.Fetch(id, context)
	if err != nil {
		return nil, err
	}
	return result.(esi.IA), nil
}
`
			result := GenerateFetch(
				"goext",
				"A",
				"esi.IA",
				"",
				"",
			)

			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct lock fetch raw function", func() {
			expected := `func LockFetchRawB(` +
				`schema goext.ISchema, ` +
				`id string, ` +
				`context goext.Context` +
				`) (*resources.B, error) {
	result, err := schema.LockFetchRaw(id, context)
	if err != nil {
		return nil, err
	}
	return result.(*resources.B), nil
}
`
			result := GenerateFetch(
				"goext",
				"B",
				"*resources.B",
				"Lock",
				"Raw",
			)

			Expect(result).To(Equal(expected))
		})
	})

	Describe("List tests", func() {
		It("Should generate a correct list function", func() {
			expected := `func ListA(` +
				`schema goext.ISchema, ` +
				`filter goext.Filter, ` +
				`paginator *goext.Paginator, ` +
				`context goext.Context` +
				`) ([]esi.IA, error) {
	list, err := schema.List(filter, paginator, context)
	if err != nil {
		return nil, err
	}
	result := make([]esi.IA, len(list))
	for i, object := range list {
		result[i] = object.(esi.IA)
	}
	return result, nil
}
`
			result := GenerateList(
				"goext",
				"A",
				"esi.IA",
				"",
				"",
			)
			Expect(result).To(Equal(expected))
		})

		It("Should generate a correct lock list raw function", func() {
			expected := `func LockListRawB(` +
				`schema goext.ISchema, ` +
				`filter goext.Filter, ` +
				`paginator *goext.Paginator, ` +
				`context goext.Context` +
				`) ([]*resources.B, error) {
	list, err := schema.LockListRaw(filter, paginator, context)
	if err != nil {
		return nil, err
	}
	result := make([]*resources.B, len(list))
	for i, object := range list {
		result[i] = object.(*resources.B)
	}
	return result, nil
}
`
			result := GenerateList(
				"goext",
				"B",
				"*resources.B",
				"Lock",
				"Raw",
			)
			Expect(result).To(Equal(expected))
		})
	})
})
