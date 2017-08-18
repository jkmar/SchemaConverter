package hash

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
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

			tree.setup()
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
			tree.setup()

			tree.calcAllHashes()

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
			tree.setup()

			tree.calcAllHashes()

			hash := &Hash{}
			Expect(tree.value).To(Equal(hash.Calc("1(2(),3(),)")))
			Expect(tree.nodes[0].value).To(Equal(hash.Calc("2()")))
			Expect(tree.nodes[1].value).To(Equal(hash.Calc("3()")))
		})
	})

	Describe("compression tests", func() {
		Describe("level 1", func() {
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
				tree.setup()
				tree.calcAllHashes()

				tree.compress(1)

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
				tree.setup()
				tree.calcAllHashes()

				tree.compress(1)

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
				tree.setup()
				tree.calcAllHashes()
				Expect(tree.nodes[0].item == tree.nodes[1].item).To(BeFalse())

				tree.compress(1)

				Expect(tree.nodes[0].item == tree.nodes[1].item).To(BeTrue())
			})
		})

		Describe("higher level tests", func() {
			It("Should compress same trees with level 2", func() {
				tree := CreateTreeNode(
					&IHashableImplementation{"A"},
					[]*TreeNode{
						CreateTreeNode(
							&IHashableImplementation{"A"},
							[]*TreeNode{
								CreateTreeNode(
									&IHashableImplementation{"A"},
									[]*TreeNode{
										CreateTreeNode(
											&IHashableImplementation{"A"},
											[]*TreeNode{
												CreateTreeNode(
													&IHashableImplementation{"A"},
													nil,
												),
												CreateTreeNode(
													&IHashableImplementation{"B"},
													nil,
												),
											},
										),
									},
								),
								CreateTreeNode(
									&IHashableImplementation{"B"},
									[]*TreeNode{
										CreateTreeNode(
											&IHashableImplementation{"A"},
											[]*TreeNode{
												CreateTreeNode(
													&IHashableImplementation{"A"},
													nil,
												),
												CreateTreeNode(
													&IHashableImplementation{"B"},
													nil,
												),
											},
										),
									},
								),
							},
						),
						CreateTreeNode(
							&IHashableImplementation{"A"},
							[]*TreeNode{
								CreateTreeNode(
									&IHashableImplementation{"A"},
									[]*TreeNode{
										CreateTreeNode(
											&IHashableImplementation{"A"},
											nil,
										),
										CreateTreeNode(
											&IHashableImplementation{"B"},
											nil,
										),
									},
								),
							},
						),
					},
				)

				tree.setup()
				tree.calcAllHashes()
				tree.compress(2)

				type pair struct {
					first,
					second int
				}
				ok := map[pair]bool{
					{3, 7}: true,
				}
				nodes := tree.getAllNodes()
				for i, node := range nodes {
					for j := i + 1; j < len(nodes); j++ {
						Expect(node.item == nodes[j].item).To(Equal(ok[pair{i, j}]))
					}
				}
			})
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
			Expect(powers(7)).To(Equal([]int{0, 1, 2}))
		})
	})

	Describe("get all nodes tests", func() {
		It("Should get all nodes visited in pre order", func() {
			item := &IHashableImplementation{"test"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(item, nil),
							CreateTreeNode(item, nil),
						},
					),
					CreateTreeNode(item, nil),
				},
			)

			result := tree.getAllNodes()

			Expect(len(result)).To(Equal(5))
			Expect(result[0]).To(Equal(tree))
			Expect(result[1]).To(Equal(tree.nodes[0]))
			Expect(result[2]).To(Equal(tree.nodes[0].nodes[0]))
			Expect(result[3]).To(Equal(tree.nodes[0].nodes[1]))
			Expect(result[4]).To(Equal(tree.nodes[1]))
		})
	})

	Describe("get ancestors tests", func() {
		It("Should get self for level 0", func() {
			item := &IHashableImplementation{"test"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(item, nil),
							CreateTreeNode(item, nil),
						},
					),
					CreateTreeNode(item, nil),
				},
			)
			nodes := tree.getAllNodes()

			tree.getAncestors(0)

			for _, node := range nodes {
				Expect(node.parent).To(Equal(node))
			}
		})

		It("Should get self for level 1", func() {
			item := &IHashableImplementation{"test"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(item, nil),
							CreateTreeNode(item, nil),
						},
					),
					CreateTreeNode(item, nil),
				},
			)
			nodes := tree.getAllNodes()
			parents := make([]*TreeNode, len(nodes))
			for i, node := range nodes {
				parents[i] = node.parent
			}

			tree.getAncestors(1)

			for i, node := range nodes {
				Expect(node.parent).To(Equal(parents[i]))
			}
		})

		It("Should get self for level 3", func() {
			item := &IHashableImplementation{"test"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(item, nil),
							CreateTreeNode(
								item,
								[]*TreeNode{
									CreateTreeNode(item, nil),
									CreateTreeNode(item, nil),
								},
							),
						},
					),
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(
								item,
								[]*TreeNode{
									CreateTreeNode(item, nil),
									CreateTreeNode(item, nil),
								},
							),
						},
					),
				},
			)
			nodes := tree.getAllNodes()

			tree.getAncestors(3)

			withRoot := []int{4, 5, 8, 9}
			withNil := []int{0, 1, 2, 3, 6, 7}
			for _, val := range withRoot {
				Expect(nodes[val].parent).To(Equal(tree))
			}
			for _, val := range withNil {
				Expect(nodes[val].parent).To(BeNil())
			}
		})

		It("Should get self for level 9", func() {
			item := &IHashableImplementation{"test"}
			tree := CreateTreeNode(
				item,
				[]*TreeNode{
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(item, nil),
							CreateTreeNode(
								item,
								[]*TreeNode{
									CreateTreeNode(item, nil),
									CreateTreeNode(item, nil),
								},
							),
						},
					),
					CreateTreeNode(
						item,
						[]*TreeNode{
							CreateTreeNode(
								item,
								[]*TreeNode{
									CreateTreeNode(item, nil),
									CreateTreeNode(item, nil),
								},
							),
						},
					),
				},
			)
			nodes := tree.getAllNodes()

			tree.getAncestors(9)

			for _, node := range nodes {
				Expect(node.parent).To(BeNil())
			}
		})
	})

	Describe("run tests", func() {
		It("Should compress", func() {
			items := []*IHashableImplementation{
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
				&IHashableImplementation{"A"},
			}
			tree := CreateTreeNode(
				items[0],
				[]*TreeNode{
					CreateTreeNode(
						items[1],
						[]*TreeNode{
							CreateTreeNode(
								items[2],
								[]*TreeNode{
									CreateTreeNode(
										items[3],
										nil,
									),
								},
							),
						},
					),
					CreateTreeNode(
						items[4],
						[]*TreeNode{
							CreateTreeNode(
								items[5],
								[]*TreeNode{
									CreateTreeNode(
										items[6],
										nil,
									),
								},
							),
						},
					),
					CreateTreeNode(
						items[7],
						[]*TreeNode{
							CreateTreeNode(
								items[8],
								[]*TreeNode{
									CreateTreeNode(
										items[9],
										nil,
									),
								},
							),
						},
					),
				},
			)

			tree.Run(3)

			type pair struct {
				first,
				second int
			}
			ok := map[pair]bool{
				{3, 6}: true,
				{3, 9}: true,
				{6, 9}: true,
			}
			//nodes := tree.getAllNodes()
			//for i, node := range nodes {
			//	for j := i + 1; j < len(nodes); j++ {
			//		Expect(node.item == nodes[j].item).To(Equal(ok[pair{i, j}]))
			//	}
			//}
			for _, item := range items {
				fmt.Printf("%p %v\n", item, item)
			}
			for i, item := range items {
				for j := i + 1; j < len(items); j++ {
					Expect(item == items[j]).To(Equal(ok[pair{i, j}]))
				}
			}
		})
	})
})
