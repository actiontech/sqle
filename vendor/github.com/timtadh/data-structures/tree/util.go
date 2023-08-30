package tree

import (
	"github.com/timtadh/data-structures/types"
)

func pop(stack []types.TreeNode) ([]types.TreeNode, types.TreeNode) {
	if len(stack) <= 0 {
		return stack, nil
	} else {
		return stack[0 : len(stack)-1], stack[len(stack)-1]
	}
}

func btn_expose_nil(node types.BinaryTreeNode) types.BinaryTreeNode {
	if types.IsNil(node) {
		return nil
	}
	return node
}

func tn_expose_nil(node types.TreeNode) types.TreeNode {
	if types.IsNil(node) {
		return nil
	}
	return node
}

func TraverseBinaryTreeInOrder(node types.BinaryTreeNode) types.TreeNodeIterator {
	stack := make([]types.TreeNode, 0, 10)
	var cur types.TreeNode = btn_expose_nil(node)
	var tn_iterator types.TreeNodeIterator
	tn_iterator = func() (tn types.TreeNode, next types.TreeNodeIterator) {
		if len(stack) > 0 || cur != nil {
			for cur != nil {
				stack = append(stack, cur)
				cur = cur.(types.BinaryTreeNode).Left()
			}
			stack, cur = pop(stack)
			tn = cur
			cur = cur.(types.BinaryTreeNode).Right()
			return tn, tn_iterator
		} else {
			return nil, nil
		}
	}
	return tn_iterator
}

func TraverseTreePreOrder(node types.TreeNode) types.TreeNodeIterator {
	stack := append(make([]types.TreeNode, 0, 10), tn_expose_nil(node))
	var tn_iterator types.TreeNodeIterator
	tn_iterator = func() (tn types.TreeNode, next types.TreeNodeIterator) {
		if len(stack) <= 0 {
			return nil, nil
		}
		stack, tn = pop(stack)
		kid_count := 1
		if tn.ChildCount() >= 0 {
			kid_count = tn.ChildCount()
		}
		kids := make([]types.TreeNode, 0, kid_count)
		for child, next := tn.Children()(); next != nil; child, next = next() {
			kids = append(kids, child)
		}
		for i := len(kids) - 1; i >= 0; i-- {
			stack = append(stack, kids[i])
		}
		return tn, tn_iterator
	}
	return tn_iterator
}

func TraverseTreePostOrder(node types.TreeNode) types.TreeNodeIterator {
	type entry struct {
		tn types.TreeNode
		i  int
	}

	pop := func(stack []entry) ([]entry, types.TreeNode, int) {
		if len(stack) <= 0 {
			return stack, nil, 0
		} else {
			e := stack[len(stack)-1]
			return stack[0 : len(stack)-1], e.tn, e.i
		}
	}

	stack := append(make([]entry, 0, 10), entry{tn_expose_nil(node), 0})

	var tn_iterator types.TreeNodeIterator
	tn_iterator = func() (tn types.TreeNode, next types.TreeNodeIterator) {
		var i int

		if len(stack) <= 0 {
			return nil, nil
		}

		stack, tn, i = pop(stack)
		for i < tn.ChildCount() {
			kid := tn.GetChild(i)
			stack = append(stack, entry{tn, i + 1})
			tn = kid
			i = 0
		}
		return tn, tn_iterator
	}
	return tn_iterator
}
