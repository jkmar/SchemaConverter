package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("package util tests", func() {
	Describe("collect data tests", func() {
		It("Should generate an empty string for an empty data", func() {
			name := "name"
			data := []string{}

			result := CollectData(name, data)

			Expect(result).To(BeEmpty())
		})

		It("Should collect data", func() {
			name := "name"
			data := []string{"a", "", "b"}

			result := CollectData(name, data)

			expected := `package name

a
b`
			Expect(result).To(Equal(expected))
		})
	})
})
