package avl

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/tree"
	"github.com/timtadh/data-structures/types"
)

type ImmutableAvlTree struct {
	root *ImmutableAvlNode
}

func NewImmutableAvlTree() *ImmutableAvlTree {
	return &ImmutableAvlTree{}
}

func (self *ImmutableAvlTree) Root() types.TreeNode {
	return self.root.Copy()
}

func (self *ImmutableAvlTree) Size() int {
	return self.root.Size()
}

func (self *ImmutableAvlTree) Has(key types.Hashable) bool {
	return self.root.Has(key)
}

func (self *ImmutableAvlTree) Put(key types.Hashable, value interface{}) (err error) {
	self.root, _ = self.root.Put(key, value)
	return nil
}

func (self *ImmutableAvlTree) Get(key types.Hashable) (value interface{}, err error) {
	return self.root.Get(key)
}

func (self *ImmutableAvlTree) Remove(key types.Hashable) (value interface{}, err error) {
	new_root, value, err := self.root.Remove(key)
	if err != nil {
		return nil, err
	}
	self.root = new_root
	return value, nil
}

func (self *ImmutableAvlTree) Iterate() types.KVIterator {
	return self.root.Iterate()
}

func (self *ImmutableAvlTree) Items() (vi types.KIterator) {
	return types.MakeItemsIterator(self)
}

func (self *ImmutableAvlTree) Values() types.Iterator {
	return self.root.Values()
}

func (self *ImmutableAvlTree) Keys() types.KIterator {
	return self.root.Keys()
}

type ImmutableAvlNode struct {
	key    types.Hashable
	value  interface{}
	height int
	left   *ImmutableAvlNode
	right  *ImmutableAvlNode
}

func (self *ImmutableAvlNode) Copy() *ImmutableAvlNode {
	if self == nil {
		return nil
	}
	return &ImmutableAvlNode{
		key:    self.key,
		value:  self.value,
		height: self.height,
		left:   self.left,
		right:  self.right,
	}
}

func (self *ImmutableAvlNode) Has(key types.Hashable) (has bool) {
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

func (self *ImmutableAvlNode) Get(key types.Hashable) (value interface{}, err error) {
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

func (self *ImmutableAvlNode) pop_node(node *ImmutableAvlNode) (new_self, new_node *ImmutableAvlNode) {
	if node == nil {
		panic("node can't be nil")
	} else if node.left != nil && node.right != nil {
		panic("node must not have both left and right")
	}

	if self == nil {
		return nil, node.Copy()
	} else if self == node {
		var n *ImmutableAvlNode
		if node.left != nil {
			n = node.left
		} else if node.right != nil {
			n = node.right
		} else {
			n = nil
		}
		node = node.Copy()
		node.left = nil
		node.right = nil
		return n, node
	}

	self = self.Copy()

	if node.key.Less(self.key) {
		self.left, node = self.left.pop_node(node)
	} else {
		self.right, node = self.right.pop_node(node)
	}

	self.height = max(self.left.Height(), self.right.Height()) + 1
	return self, node
}

func (self *ImmutableAvlNode) push_node(node *ImmutableAvlNode) *ImmutableAvlNode {
	if node == nil {
		panic("node can't be nil")
	} else if node.left != nil || node.right != nil {
		panic("node must now be a leaf")
	}

	self = self.Copy()

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

func (self *ImmutableAvlNode) rotate_right() *ImmutableAvlNode {
	if self == nil {
		return self
	}
	if self.left == nil {
		return self
	}
	return self.rotate(self.left.rmd)
}

func (self *ImmutableAvlNode) rotate_left() *ImmutableAvlNode {
	if self == nil {
		return self
	}
	if self.right == nil {
		return self
	}
	return self.rotate(self.right.lmd)
}

func (self *ImmutableAvlNode) rotate(get_new_root func() *ImmutableAvlNode) *ImmutableAvlNode {
	self, new_root := self.pop_node(get_new_root())
	new_root.left = self.left
	new_root.right = self.right
	self.left = nil
	self.right = nil
	return new_root.push_node(self)
}

func (self *ImmutableAvlNode) balance() *ImmutableAvlNode {
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

func (self *ImmutableAvlNode) Put(key types.Hashable, value interface{}) (_ *ImmutableAvlNode, updated bool) {
	if self == nil {
		return &ImmutableAvlNode{key: key, value: value, height: 1}, false
	}

	self = self.Copy()

	if self.key.Equals(key) {
		self.value = value
		return self, true
	}

	if key.Less(self.key) {
		self.left, updated = self.left.Put(key, value)
	} else {
		self.right, updated = self.right.Put(key, value)
	}
	self.height = max(self.left.Height(), self.right.Height()) + 1

	if !updated {
		self.height += 1
		return self.balance(), updated
	}
	return self, updated
}

func (self *ImmutableAvlNode) Remove(key types.Hashable) (_ *ImmutableAvlNode, value interface{}, err error) {
	if self == nil {
		return nil, nil, errors.NotFound(key)
	}

	if self.key.Equals(key) {
		if self.left != nil && self.right != nil {
			var new_root *ImmutableAvlNode
			if self.left.Size() < self.right.Size() {
				self, new_root = self.pop_node(self.right.lmd())
			} else {
				self, new_root = self.pop_node(self.left.rmd())
			}
			new_root.left = self.left
			new_root.right = self.right
			return new_root, self.value, nil
		} else if self.left == nil {
			return self.right, self.value, nil
		} else if self.right == nil {
			return self.left, self.value, nil
		} else {
			return nil, self.value, nil
		}
	}

	self = self.Copy()

	if key.Less(self.key) {
		self.left, value, err = self.left.Remove(key)
	} else {
		self.right, value, err = self.right.Remove(key)
	}
	self.height = max(self.left.Height(), self.right.Height()) + 1
	if err != nil {
		return self.balance(), value, err
	}
	return self, value, err
}

func (self *ImmutableAvlNode) Height() int {
	if self == nil {
		return 0
	}
	return self.height
}

func (self *ImmutableAvlNode) Size() int {
	if self == nil {
		return 0
	}
	return 1 + self.left.Size() + self.right.Size()
}

func (self *ImmutableAvlNode) Key() types.Hashable {
	return self.key
}

func (self *ImmutableAvlNode) Value() interface{} {
	return self.value
}

func (self *ImmutableAvlNode) Left() types.BinaryTreeNode {
	if self.left == nil {
		return nil
	}
	return self.left
}

func (self *ImmutableAvlNode) Right() types.BinaryTreeNode {
	if self.right == nil {
		return nil
	}
	return self.right
}

func (self *ImmutableAvlNode) GetChild(i int) types.TreeNode {
	return types.DoGetChild(self, i)
}

func (self *ImmutableAvlNode) ChildCount() int {
	return types.DoChildCount(self)
}

func (self *ImmutableAvlNode) Children() types.TreeNodeIterator {
	return types.MakeChildrenIterator(self)
}

func (self *ImmutableAvlNode) Iterate() types.KVIterator {
	tni := tree.TraverseBinaryTreeInOrder(self)
	return types.MakeKVIteratorFromTreeNodeIterator(tni)
}

func (self *ImmutableAvlNode) Keys() types.KIterator {
	return types.MakeKeysIterator(self)
}

func (self *ImmutableAvlNode) Values() types.Iterator {
	return types.MakeValuesIterator(self)
}

func (self *ImmutableAvlNode) _md(side func(*ImmutableAvlNode) *ImmutableAvlNode) *ImmutableAvlNode {
	if self == nil {
		return nil
	} else if side(self) != nil {
		return side(self)._md(side)
	} else {
		return self
	}
}

func (self *ImmutableAvlNode) lmd() *ImmutableAvlNode {
	return self._md(func(node *ImmutableAvlNode) *ImmutableAvlNode { return node.left })
}

func (self *ImmutableAvlNode) rmd() *ImmutableAvlNode {
	return self._md(func(node *ImmutableAvlNode) *ImmutableAvlNode { return node.right })
}
