package types

import (
	"reflect"
)

func IsNil(object interface{}) bool {
	return object == nil || reflect.ValueOf(object).IsNil()
}

func MakeKVIteratorFromTreeNodeIterator(tni TreeNodeIterator) KVIterator {
	var kv_iterator KVIterator
	kv_iterator = func() (key Hashable, value interface{}, next KVIterator) {
		var tn TreeNode
		tn, tni = tni()
		if tni == nil {
			return nil, nil, nil
		}
		return tn.Key(), tn.Value(), kv_iterator
	}
	return kv_iterator
}

func ChainTreeNodeIterators(tnis ...TreeNodeIterator) TreeNodeIterator {
	var make_tni func(int) TreeNodeIterator
	make_tni = func(i int) (tni_iterator TreeNodeIterator) {
		if i >= len(tnis) {
			return nil
		}
		var next TreeNodeIterator = tnis[i]
		tni_iterator = func() (TreeNode, TreeNodeIterator) {
			var tn TreeNode
			tn, next = next()
			if next == nil {
				tni := make_tni(i + 1)
				if tni == nil {
					return nil, nil
				} else {
					return tni()
				}
			}
			return tn, tni_iterator
		}
		return tni_iterator
	}
	return make_tni(0)
}

func MakeKeysIterator(obj KVIterable) KIterator {
	kv_iterator := obj.Iterate()
	var k_iterator KIterator
	k_iterator = func() (key Hashable, next KIterator) {
		key, _, kv_iterator = kv_iterator()
		if kv_iterator == nil {
			return nil, nil
		}
		return key, k_iterator
	}
	return k_iterator
}

func MakeValuesIterator(obj KVIterable) Iterator {
	kv_iterator := obj.Iterate()
	var v_iterator Iterator
	v_iterator = func() (value interface{}, next Iterator) {
		_, value, kv_iterator = kv_iterator()
		if kv_iterator == nil {
			return nil, nil
		}
		return value, v_iterator
	}
	return v_iterator
}

func MakeItemsIterator(obj KVIterable) (kit KIterator) {
	kv_iterator := obj.Iterate()
	kit = func() (item Hashable, next KIterator) {
		var key Hashable
		var value interface{}
		key, value, kv_iterator = kv_iterator()
		if kv_iterator == nil {
			return nil, nil
		}
		return &MapEntry{key, value}, kit
	}
	return kit
}

func make_child_slice(node BinaryTreeNode) []BinaryTreeNode {
	nodes := make([]BinaryTreeNode, 0, 2)
	if !IsNil(node) {
		if !IsNil(node.Left()) {
			nodes = append(nodes, node.Left())
		}
		if !IsNil(node.Right()) {
			nodes = append(nodes, node.Right())
		}
	}
	return nodes
}

func DoGetChild(node BinaryTreeNode, i int) TreeNode {
	return make_child_slice(node)[i]
}

func DoChildCount(node BinaryTreeNode) int {
	return len(make_child_slice(node))
}

func MakeChildrenIterator(node BinaryTreeNode) TreeNodeIterator {
	nodes := make_child_slice(node)
	var make_tn_iterator func(int) TreeNodeIterator
	make_tn_iterator = func(i int) TreeNodeIterator {
		return func() (kid TreeNode, next TreeNodeIterator) {
			if i < len(nodes) {
				return nodes[i], make_tn_iterator(i + 1)
			}
			return nil, nil
		}
	}
	return make_tn_iterator(0)
}

func MakeMarshals(empty func() MHashable) (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Marshaler)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := empty()
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func Int8Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Int8)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := Int8(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func UInt8Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(UInt8)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := UInt8(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func Int16Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Int16)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := Int16(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func UInt16Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(UInt16)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := UInt16(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func Int32Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Int32)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := Int32(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func UInt32Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(UInt32)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := UInt32(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func Int64Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Int64)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := Int64(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func UInt64Marshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(UInt64)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := UInt64(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func IntMarshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(Int)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := Int(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func UIntMarshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(UInt)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := UInt(0)
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func StringMarshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		i := item.(String)
		return i.MarshalBinary()
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		i := String("")
		err := i.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return marshal, unmarshal
}

func ByteSliceMarshals() (ItemMarshal, ItemUnmarshal) {
	marshal := func(item Hashable) ([]byte, error) {
		return []byte(item.(ByteSlice)), nil
	}
	unmarshal := func(bytes []byte) (Hashable, error) {
		return ByteSlice(bytes), nil
	}
	return marshal, unmarshal
}
