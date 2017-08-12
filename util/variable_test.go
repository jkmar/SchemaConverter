package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("string util tests", func() {
	Describe("ToLowerFirst tests", func() {
		It("Should return an empty string for an empty input", func() {
			Expect(ToLowerFirst("")).To(Equal(""))
		})

		It("Should change the first letter", func() {
			Expect(ToLowerFirst("Abc")).To(Equal("abc"))
		})

		It("Should not change the first letter", func() {
			Expect(ToLowerFirst("123")).To(Equal("123"))
		})
	})

	Describe("VariableName tests", func() {
		It("Should return unchanged variable name", func() {
			Expect(VariableName("Abc")).To(Equal("abc"))
		})

		It("Should return variable name with prefix to avoid keywords", func() {
			Expect(VariableName("Range")).To(Equal("rangeObject"))
		})
	})

	Describe("IndexVariable tests", func() {
		It("Should return correct index variable name - i", func() {
			Expect(IndexVariable(1)).To(Equal('i'))
		})

		It("Should return correct index variable name - j", func() {
			Expect(IndexVariable(2)).To(Equal('j'))
		})
	})
})
