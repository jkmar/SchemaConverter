package set

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("element tests", func() {
	Describe("len tests", func() {
		It("Should return correct length", func() {
			var array byName = []Element{test("a"), test("b")}
			Expect(array.Len()).To(Equal(len(array)))
		})
	})

	Describe("swap tests", func() {
		It("Should swap elements for array", func() {
			var array byName = []Element{test("a"), test("b")}
			array.Swap(0, 1)
			Expect(array[0]).To(Equal(test("b")))
			Expect(array[1]).To(Equal(test("a")))
		})
	})

	Describe("less tests", func() {
		It("Should compare elements", func() {
			var array byName = []Element{test("a"), test("b")}
			Expect(array.Less(0, 1)).To(BeTrue())
		})
	})
})
