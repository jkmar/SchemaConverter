package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("defaults tests", func() {
	Describe("defaults creation tests", func() {
		It("Should create nil defaults for a complex type", func() {
			Expect(CreatePlainDefaults([]string{})).To(BeNil())
		})

		It("Should create string defaults", func() {
			Expect(CreatePlainDefaults("test")).To(Equal(&StringDefault{value: "test"}))
		})

		It("Should create int defaults", func() {
			Expect(CreatePlainDefaults(1)).To(Equal(&IntDefault{value: 1}))
		})

		It("Should create float defaults", func() {
			Expect(CreatePlainDefaults(1.5)).To(Equal(&FloatDefault{value: 1.5}))
		})

		It("Should create bool defaults", func() {
			Expect(CreatePlainDefaults(true)).To(Equal(&BoolDefault{value: true}))
		})
	})
})
