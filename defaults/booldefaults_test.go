package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("bool defaults tests", func() {
	Describe("write tests", func() {
		It("Should write a correct value", func() {
			defaults := &BoolDefault{value: true}
			Expect(defaults.Write()).To(Equal("true"))
		})
	})
})
