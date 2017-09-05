package hash

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hash util tests", func() {
	Describe("power tests", func() {
		It("Should get correct powers", func() {
			log, powers := powers(0)

			Expect(log).To(Equal(-1))
			Expect(powers).To(BeNil())
		})

		It("Should get correct powers", func() {
			log, powers := powers(1)

			Expect(log).To(Equal(0))
			Expect(powers).To(Equal([]int{0}))
		})

		It("Should get correct powers", func() {
			log, powers := powers(4)

			Expect(log).To(Equal(2))
			Expect(powers).To(Equal([]int{2}))
		})

		It("Should get correct powers", func() {
			log, powers := powers(7)

			Expect(log).To(Equal(2))
			Expect(powers).To(Equal([]int{0, 1, 2}))
		})

		It("Should get correct powers", func() {
			log, powers := powers(10)

			Expect(log).To(Equal(3))
			Expect(powers).To(Equal([]int{1, 3}))
		})
	})
})
