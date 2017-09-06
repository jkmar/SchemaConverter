package defaults

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("string defaults tests", func() {
	Describe("write tests", func() {
		It("Should write a correct value", func() {
			defaults := &StringDefault{value: "test"}
			Expect(defaults.Write()).To(Equal(`"test"`))
		})
	})
})
