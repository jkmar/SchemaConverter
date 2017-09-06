package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("defaults tests", func() {
	Describe("defaults creation tests", func() {
		It("Should create a nil defaults for a complex type", func() {
			Expect(CreatePlainDefaults([]string{})).To(BeNil())
		})

		It("Should create a string defaults", func() {
			//Expect(CreatePlainDefaults("test")).To(Type)
		})
	})
})
