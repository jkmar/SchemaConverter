package hash

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hash tests", func() {
	Describe("mod add tests", func() {
		It("Should add numbers with sum less than mod", func() {
			Expect(AddMod(1, 2, 4)).To(Equal(uint32(3)))
		})

		It("Should add numbers with sum greater than mod", func() {
			Expect(AddMod(2, 3, 4)).To(Equal(uint32(1)))
		})

		It("Should add numbers with overflow", func() {
			Expect(AddMod(0xFFFFFFFE, 2, 0xFFFFFFFF)).To(Equal(uint32(1)))
		})
	})

	Describe("mod mul tests", func() {
		It("Should multiply numbers with product less than mod", func() {
			Expect(MulMod(2, 3, 7)).To(Equal(uint32(6)))
		})

		It("Should multiply numbers with product greater than mod", func() {
			Expect(MulMod(10, 20, 3)).To(Equal(uint32(2)))
		})

		It("Should multiply numbers with overflow", func() {
			Expect(MulMod(0xFFFFFFFE, 5, 0xFFFFFFFF)).To(Equal(uint32(0xFFFFFFFA)))
		})
	})

	Describe("calc hash tests", func() {
		It("Should calc correct empty hash", func() {
			hash := Hash{}

			result := hash.Calc("")

			Expect(result.length).To(Equal(0))
			Expect(result.value).To(Equal(uint32(0)))
		})

		It("Should calc correct hash", func() {
			hash := Hash{}

			result := hash.Calc("a")

			Expect(result.length).To(Equal(1))
			Expect(result.value).To(Equal(uint32('a')))
		})
	})

	Describe("join tests", func() {
		It("Should join hash with zero value", func() {
			hash := Hash{}
			first := Node{1, 5}
			second := Node{0, 4}

			result := hash.Join(first, second)

			Expect(result.length).To(Equal(9))
			Expect(result.value).To(Equal(uint32(1)))
		})

		It("Should join hashes", func() {
			hash := Hash{}
			first := Node{1, 2}
			second := Node{3, 4}

			result := hash.Join(first, second)

			Expect(result.length).To(Equal(6))
			Expect(result.value).To(Equal(uint32(198148)))
		})
	})

	Describe("general tests", func() {
		It("Should compare equal strings", func() {
			hash := Hash{}
			strings := []string{"a", "bc", "ab", "c"}
			node := make([]Node, len(strings))

			for i, string := range strings {
				node[i] = hash.Calc(string)
			}

			first := hash.Join(node[0], node[1])
			second := hash.Join(node[2], node[3])

			Expect(first.value).To(Equal(second.value))
		})

		It("Should compare equal strings", func() {
			hash := Hash{}
			strings := []string{"abc", "def", "abc", "a", "bcdefabc"}
			node := make([]Node, len(strings))

			for i, string := range strings {
				node[i] = hash.Calc(string)
			}

			first := hash.Join(hash.Join(node[0], node[1]), node[2])
			second := hash.Join(node[3], node[4])

			Expect(first.value).To(Equal(second.value))
		})

		It("Should compare non equal strings", func() {
			hash := Hash{}
			strings := []string{"aaaaaaaaaaa", "aaaaaaaaaa"}

			first := hash.Calc(strings[0])
			second := hash.Calc(strings[1])

			Expect(first.value).ToNot(Equal(second.value))
		})

		It("Should compare non equal strings", func() {
			hash := Hash{}
			strings := []string{"abc", "def", "abc", "aa", "bcdefabc"}
			node := make([]Node, len(strings))

			for i, string := range strings {
				node[i] = hash.Calc(string)
			}

			first := hash.Join(hash.Join(node[0], node[1]), node[2])
			second := hash.Join(node[3], node[4])

			Expect(first.value).ToNot(Equal(second.value))
		})
	})
})
