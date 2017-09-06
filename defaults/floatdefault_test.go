package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("float defaults tests", func() {
	Describe("write tests", func() {
		It("Should write a correct value", func() {
			defaults := &FloatDefault{value: 1.5}
			Expect(defaults.Write()).To(Equal("1.5E+00"))
		})
	})
})
