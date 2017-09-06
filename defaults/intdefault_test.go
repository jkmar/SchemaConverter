package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("int defaults tests", func() {
	Describe("write tests", func() {
		It("Should write a correct value", func() {
			defaults := &IntDefault{value: 1}
			Expect(defaults.Write()).To(Equal("1"))
		})
	})
})
