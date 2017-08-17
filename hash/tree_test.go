package hash

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type IHashableImplementation struct {
	value string
}

func (item *IHashableImplementation) ToString() string {
	return item.value
}

var _ = Describe("tree tests", func() {
	Describe("create tree tests", func() {
		It("Should create tree", func() {
			item := &IHashableImplementation{"abc"}

			tree := CreateTreeNode(item, nil)

			Expect(tree.item).To(Equal(item))
			Expect(tree.nodes).To(BeNil())
		})
	})

	Describe("setup tests", func() {
		It("Should use the same hash", func() {
			tree := TreeNode{
				nodes: []*TreeNode{
					{
						nodes: []*TreeNode{
							{},
							{},
						},
					},
					{},
				},
			}

			tree.Setup()
			result := tree.hash

			Expect(tree.nodes[0].nodes[0].hash).To(Equal(result))
			Expect(tree.nodes[0].nodes[1].hash).To(Equal(result))
			Expect(tree.nodes[1].hash).To(Equal(result))
		})
	})

	Describe("calc all hashes tests", func() {
		It("Should calc hashes for one node tree", func() {
			item := &IHashableImplementation{"abc"}
			tree := CreateTreeNode(item, nil)
			tree.Setup()

			tree.CalcAllHashes()

			hash := &Hash{}
			Expect(tree.value).To(Equal(hash.Calc("abc()")))
		})

		It("Should calc hashes for tree", func() {
			item := &IHashableImplementation{"1"}
			left := &IHashableImplementation{"2"}
			right := &IHashableImplementation{"3"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(left, nil),
					CreateTreeNode(right, nil),
				},
			)
			tree.Setup()

			tree.CalcAllHashes()

			hash := &Hash{}
			Expect(tree.value).To(Equal(hash.Calc("1(2(),3(),)")))
			Expect(tree.nodes[0].value).To(Equal(hash.Calc("2()")))
			Expect(tree.nodes[1].value).To(Equal(hash.Calc("3()")))
		})
	})

	Describe("compression tests", func() {
		It("Should not compress different trees", func() {
			item := &IHashableImplementation{"1"}
			left := &IHashableImplementation{"2"}
			right := &IHashableImplementation{"3"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(left, nil),
					CreateTreeNode(right, nil),
				},
			)
			tree.Setup()
			tree.CalcAllHashes()

			tree.Compress()

			Expect(tree.nodes[0].item).ToNot(Equal(tree.nodes[1].item))
		})

		It("Should not compress tree with different parents", func() {
			item1 := &IHashableImplementation{"X"}
			item2 := &IHashableImplementation{"X"}
			item := &IHashableImplementation{"1"}
			left := &IHashableImplementation{"2"}
			right := &IHashableImplementation{"3"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						left,
						[]*TreeNode{CreateTreeNode(item1, nil)},
					),
					CreateTreeNode(
						right,
						[]*TreeNode{CreateTreeNode(item2, nil)},
					),
				},
			)
			tree.Setup()
			tree.CalcAllHashes()

			tree.Compress()

			Expect(tree.nodes[0].nodes[0].item == tree.nodes[1].nodes[0].item).To(BeFalse())
		})

		It("Should compress same trees", func() {
			item := &IHashableImplementation{"1"}
			left := &IHashableImplementation{"2"}
			right := &IHashableImplementation{"2"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(left, nil),
					CreateTreeNode(right, nil),
				},
			)
			tree.Setup()
			tree.CalcAllHashes()
			Expect(tree.nodes[0].item == tree.nodes[1].item).To(BeFalse())

			tree.Compress()

			Expect(tree.nodes[0].item == tree.nodes[1].item).To(BeTrue())
		})
	})

	Describe("power tests", func() {
		It("Should get correct powers", func() {
			Expect(powers(0)).To(Equal([]int{}))
		})

		It("Should get correct powers", func() {
			Expect(powers(4)).To(Equal([]int{2}))
		})

		It("Should get correct powers", func() {
			Expect(powers(7)).To(Equal([]int{0,1,2}))
		})
	})
})