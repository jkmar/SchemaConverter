package name

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mark tests", func() {
	Describe("create mark tests", func() {
		It("Should craete a correct mark", func() {
			prefix := "test"

			result := CreateMark(prefix)

			expected := &CommonMark{
				used: false,
				begin: len(prefix),
				end: 0,
			}
			Expect(result).To(Equal(expected))
		})
	})
})
