package avl

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/tree"
	"github.com/timtadh/data-structures/types"
)

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type AvlTree struct {
	root *AvlNode
}

func NewAvlTree() *AvlTree {
	return &AvlTree{}
}

func (self *AvlTree) Root() types.TreeNode {
	return self.root
}

func (self *AvlTree) Size() int {
	return self.root.Size()
}

func (self *AvlTree) Has(key types.Hashable) bool {
	return self.root.Has(key)
}

func (self *AvlTree) Put(key types.Hashable, value interface{}) (err error) {
	self.root, _ = self.root.Put(key, value)
	return nil
}

func (self *AvlTree) Get(key types.Hashable) (value interface{}, err error) {
	return self.root.Get(key)
}

func (self *AvlTree) Remove(key types.Hashable) (value interface{}, err error) {
	new_root, value, err := self.root.Remove(key)
	if err != nil {
		return nil, err
	}
	self.root = new_root
	return value, nil
}

func (self *AvlTree) Iterate() types.KVIterator {
	return self.root.Iterate()
}

func (self *AvlTree) Items() (vi types.KIterator) {
	return types.MakeItemsIterator(self)
}

func (self *AvlTree) Values() types.Iterator {
	return self.root.Values()
}

func (self *AvlTree) Keys() types.KIterator {
	return self.root.Keys()
}

type AvlNode struct {
	key    types.Hashable
	value  interface{}
	height int
	left   *AvlNode
	right  *AvlNode
}

func (self *AvlNode) Has(key types.Hashable) (has bool) {
	if self == nil {
		return false
	}
	if self.key.Equals(key) {
		return true
	} else if key.Less(self.key) {
		return self.left.Has(key)
	} else {
		return self.right.Has(key)
	}
}

func (self *AvlNode) Get(key types.Hashable) (value interface{}, err error) {
	if self == nil {
		return nil, errors.NotFound(key)
	}
	if self.key.Equals(key) {
		return self.value, nil
	} else if key.Less(self.key) {
		return self.left.Get(key)
	} else {
		return self.right.Get(key)
	}
}

func (self *AvlNode) pop_node(node *AvlNode) *AvlNode {
	if node == nil {
		panic("node can't be nil")
	} else if node.left != nil && node.right != nil {
		panic("node must not have both left and right")
	}

	if self == nil {
		return nil
	} else if self == node {
		var n *AvlNode
		if node.left != nil {
			n = node.left
		} else if node.right != nil {
			n = node.right
		} else {
			n = nil
		}
		node.left = nil
		node.right = nil
		return n
	}

	if node.key.Less(self.key) {
		self.left = self.left.pop_node(node)
	} else {
		self.right = self.right.pop_node(node)
	}

	self.height = max(self.left.Height(), self.right.Height()) + 1
	return self
}

func (self *AvlNode) push_node(node *AvlNode) *AvlNode {
	if node == nil {
		panic("node can't be nil")
	} else if node.left != nil || node.right != nil {
		panic("node now be a leaf")
	}

	if self == nil {
		node.height = 1
		return node
	} else if node.key.Less(self.key) {
		self.left = self.left.push_node(node)
	} else {
		self.right = self.right.push_node(node)
	}
	self.height = max(self.left.Height(), self.right.Height()) + 1
	return self
}

func (self *AvlNode) rotate_right() *AvlNode {
	if self == nil {
		return self
	}
	if self.left == nil {
		return self
	}
	new_root := self.left.rmd()
	self = self.pop_node(new_root)
	new_root.left = self.left
	new_root.right = self.right
	self.left = nil
	self.right = nil
	return new_root.push_node(self)
}

func (self *AvlNode) rotate_left() *AvlNode {
	if self == nil {
		return self
	}
	if self.right == nil {
		return self
	}
	new_root := self.right.lmd()
	self = self.pop_node(new_root)
	new_root.left = self.left
	new_root.right = self.right
	self.left = nil
	self.right = nil
	return new_root.push_node(self)
}

func (self *AvlNode) balance() *AvlNode {
	if self == nil {
		return self
	}
	for abs(self.left.Height()-self.right.Height()) > 2 {
		if self.left.Height() > self.right.Height() {
			self = self.rotate_right()
		} else {
			self = self.rotate_left()
		}
	}
	return self
}

func (self *AvlNode) Put(key types.Hashable, value interface{}) (_ *AvlNode, updated bool) {
	if self == nil {
		return &AvlNode{key: key, value: value, height: 1}, false
	}

	if self.key.Equals(key) {
		self.value = value
		return self, true
	}

	if key.Less(self.key) {
		self.left, updated = self.left.Put(key, value)
	} else {
		self.right, updated = self.right.Put(key, value)
	}
	if !updated {
		self.height += 1
		return self.balance(), updated
	}
	return self, updated
}

func (self *AvlNode) Remove(key types.Hashable) (_ *AvlNode, value interface{}, err error) {
	if self == nil {
		return nil, nil, errors.NotFound(key)
	}

	if self.key.Equals(key) {
		if self.left != nil && self.right != nil {
			if self.left.Size() < self.right.Size() {
				lmd := self.right.lmd()
				lmd.left = self.left
				return self.right, self.value, nil
			} else {
				rmd := self.left.rmd()
				rmd.right = self.right
				return self.left, self.value, nil
			}
		} else if self.left == nil {
			return self.right, self.value, nil
		} else if self.right == nil {
			return self.left, self.value, nil
		} else {
			return nil, self.value, nil
		}
	}
	if key.Less(self.key) {
		self.left, value, err = self.left.Remove(key)
	} else {
		self.right, value, err = self.right.Remove(key)
	}
	if err != nil {
		return self.balance(), value, err
	}
	return self, value, err
}

func (self *AvlNode) Height() int {
	if self == nil {
		return 0
	}
	return self.height
}

func (self *AvlNode) Size() int {
	if self == nil {
		return 0
	}
	return 1 + self.left.Size() + self.right.Size()
}

func (self *AvlNode) Key() types.Hashable {
	return self.key
}

func (self *AvlNode) Value() interface{} {
	return self.value
}

func (self *AvlNode) Left() types.BinaryTreeNode {
	if self.left == nil {
		return nil
	}
	return self.left
}

func (self *AvlNode) Right() types.BinaryTreeNode {
	if self.right == nil {
		return nil
	}
	return self.right
}

func (self *AvlNode) GetChild(i int) types.TreeNode {
	return types.DoGetChild(self, i)
}

func (self *AvlNode) ChildCount() int {
	return types.DoChildCount(self)
}

func (self *AvlNode) Children() types.TreeNodeIterator {
	return types.MakeChildrenIterator(self)
}

func (self *AvlNode) Iterate() types.KVIterator {
	tni := tree.TraverseBinaryTreeInOrder(self)
	return types.MakeKVIteratorFromTreeNodeIterator(tni)
}

func (self *AvlNode) Keys() types.KIterator {
	return types.MakeKeysIterator(self)
}

func (self *AvlNode) Values() types.Iterator {
	return types.MakeValuesIterator(self)
}

func (self *AvlNode) _md(side func(*AvlNode) *AvlNode) *AvlNode {
	if self == nil {
		return nil
	} else if side(self) != nil {
		return side(self)._md(side)
	} else {
		return self
	}
}

func (self *AvlNode) lmd() *AvlNode {
	return self._md(func(node *AvlNode) *AvlNode { return node.left })
}

func (self *AvlNode) rmd() *AvlNode {
	return self._md(func(node *AvlNode) *AvlNode { return node.right })
}
