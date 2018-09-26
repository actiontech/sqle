package util

import (
	"strings"
)

type KvTree struct {
	val   string
	nodes map[string]*KvTree
}

func LoadKvTreeFromMap(m map[string]string) (*KvTree) {
	tree := &KvTree{
		nodes: map[string]*KvTree{},
	}

	for k, v := range m {
		t := tree
		for _, seg := range strings.Split(k, "/") {
			if "" == seg {
				continue
			}
			if nil == t.nodes[seg] {
				t.nodes[seg] = &KvTree{
					nodes: map[string]*KvTree{},
				}
			}
			t = t.nodes[seg]
		}
		t.val = v
	}

	return tree
}

func (tree *KvTree) GetNextLevelKeys(root string) []string {
	t := tree
	for _, seg := range strings.Split(root, "/") {
		if "" == seg {
			continue
		}
		if nil == t.nodes[seg] {
			return []string{}
		}
		t = t.nodes[seg]
	}
	ret := []string{}
	for k := range t.nodes {
		ret = append(ret, k)
	}
	return ret
}

func (tree *KvTree) GetVal(key string) string {
	t := tree
	for _, seg := range strings.Split(key, "/") {
		if "" == seg {
			continue
		}
		if nil == t.nodes[seg] {
			return ""
		}
		t = t.nodes[seg]
	}
	return t.val
}