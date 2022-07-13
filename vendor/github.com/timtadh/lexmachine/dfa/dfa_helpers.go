package dfa

import (
	"fmt"

	"github.com/timtadh/lexmachine/frontend"
)

// LabeledAST is a post-order labeled version of the AST. The root the node will be Order[len(Order)-1].
type LabeledAST struct {
	Root      frontend.AST   // root of the AST
	Order     []frontend.AST // post-order labeling of the all the nodes
	Kids      [][]int        // a lookup table of location for each of the nodes children
	Positions []int          // a labeling of all the position nodes (Character/Range) in the AST.
	posmap    map[int]int    // maps an order index to a pos index.
	Matches   []int          // the index (in Positions) of all the EOS nodes
	nullable  []bool
	first     [][]int
	last      [][]int
	follow    []map[int]bool
}

// Label a tree and its "positions" (Character/Range nodes) in post-order
// notation.
func Label(ast frontend.AST) *LabeledAST {
	type entry struct {
		n    frontend.AST
		i    int
		kids []int
	}
	order := make([]frontend.AST, 0, 10)
	children := make([][]int, 0, 10)
	positions := make([]int, 0, 10)
	matches := make([]int, 0, 10)
	posmap := make(map[int]int)
	stack := make([]entry, 0, 10)
	stack = append(stack, entry{ast, 0, []int{}})
	for len(stack) > 0 {
		var c entry
		stack, c = stack[:len(stack)-1], stack[len(stack)-1]
		kids := c.n.Children()
		for c.i < len(kids) {
			kid := kids[c.i]
			stack = append(stack, entry{c.n, c.i + 1, c.kids})
			c = entry{kid, 0, []int{}}
			kids = c.n.Children()
		}
		oid := len(order)
		if len(stack) > 0 {
			stack[len(stack)-1].kids = append(stack[len(stack)-1].kids, oid)
		}
		order = append(order, c.n)
		children = append(children, c.kids)
		switch c.n.(type) {
		case *frontend.Character, *frontend.Range:
			posmap[oid] = len(positions)
			positions = append(positions, oid)
		case *frontend.EOS:
			pid := len(positions)
			posmap[oid] = pid
			positions = append(positions, oid)
			matches = append(matches, pid)
		}
	}
	return &LabeledAST{
		Root:      ast,
		Order:     order,
		Kids:      children,
		Positions: positions,
		posmap:    posmap,
		Matches:   matches,
	}
}

func (a *LabeledAST) pos(oid int) int {
	if pid, has := a.posmap[oid]; !has {
		panic("Passed a bad order id into Position (likely used a non-position node's id)")
	} else {
		return pid
	}
}

// Follow computes look up tables for each Position (leaf node) in the tree
// which indicates what other Position could follow the current position in a
// matching string. It also computes what positions appear first in the tree.
func (a *LabeledAST) Follow() (firstOfTree []int, follow []map[int]bool) {
	positions := a.Positions
	nullable := a.MatchesEmptyString()
	first := a.First()
	last := a.Last()

	// get the first of the whole ast by retrieving the first for the root (len(order)-1).
	firstOfTree = make([]int, 0, len(first[len(a.Order)-1]))
	for _, p := range first[len(a.Order)-1] {
		firstOfTree = append(firstOfTree, p)
	}

	if a.follow != nil {
		return firstOfTree, a.follow
	}

	follow = make([]map[int]bool, len(positions))
	for i := range follow {
		follow[i] = make(map[int]bool)
	}
	for i, node := range a.Order {
		switch n := node.(type) {
		case *frontend.Concat:
			for x := 0; x < len(n.Items)-1; x++ {
				j := a.Kids[i][x]
				kFirst := make([]int, 0, 10)
				for y := x + 1; y < len(n.Items); y++ {
					k := a.Kids[i][y]
					for _, p := range first[k] {
						kFirst = append(kFirst, p)
					}
					if !nullable[k] {
						break
					}
				}
				for _, p := range last[j] {
					for _, q := range kFirst {
						follow[p][q] = true
					}
				}
			}
		case *frontend.Star, *frontend.Plus:
			nFirst := make([]int, 0, 10)
			for _, p := range first[i] {
				nFirst = append(nFirst, p)
			}
			for _, p := range last[i] {
				for _, q := range nFirst {
					follow[p][q] = true
				}
			}
		}
	}

	a.follow = follow
	return firstOfTree, follow
}

// MatchesEmptyString computes a look up table for each node in the tree (in
// post order, matching the Order member) on whether or not the subtree rooted
// in each node matches the empty string.
func (a *LabeledAST) MatchesEmptyString() (nullable []bool) {
	if a.nullable != nil {
		return a.nullable
	}
	nullable = make([]bool, 0, len(a.Order))
	for i, node := range a.Order {
		switch n := node.(type) {
		case *frontend.EOS:
			nullable = append(nullable, true)
		case *frontend.Character:
			nullable = append(nullable, false)
		case *frontend.Range:
			nullable = append(nullable, false)
		case *frontend.Maybe:
			nullable = append(nullable, true)
		case *frontend.Plus:
			nullable = append(nullable, nullable[a.Kids[i][0]])
		case *frontend.Star:
			nullable = append(nullable, true)
		case *frontend.Match:
			nullable = append(nullable, nullable[a.Kids[i][0]])
		case *frontend.Alternation:
			nullable = append(nullable, nullable[a.Kids[i][0]] || nullable[a.Kids[i][1]])
		case *frontend.AltMatch:
			nullable = append(nullable, nullable[a.Kids[i][0]] || nullable[a.Kids[i][1]])
		case *frontend.Concat:
			epsilon := true
			for _, j := range a.Kids[i] {
				epsilon = epsilon && nullable[j]
			}
			nullable = append(nullable, epsilon)
		default:
			panic(fmt.Errorf("Unexpected type %T", n))
		}
	}
	a.nullable = nullable
	return nullable
}

// First computes a look up table for each node in the tree (in post order,
// matching the Order member) indicating for the subtree rooted at each node
// which positions (leaf nodes indexes in the Positions slice) will appear at
// the beginning of a string matching that subtree.
func (a *LabeledAST) First() (first [][]int) {
	if a.first != nil {
		return a.first
	}
	nullable := a.MatchesEmptyString()
	first = make([][]int, 0, len(a.Order))
	for i, node := range a.Order {
		switch n := node.(type) {
		case *frontend.EOS:
			first = append(first, []int{a.pos(i)})
		case *frontend.Character:
			first = append(first, []int{a.pos(i)})
		case *frontend.Range:
			first = append(first, []int{a.pos(i)})
		case *frontend.Maybe:
			first = append(first, first[a.Kids[i][0]])
		case *frontend.Plus:
			first = append(first, first[a.Kids[i][0]])
		case *frontend.Star:
			first = append(first, first[a.Kids[i][0]])
		case *frontend.Match:
			first = append(first, first[a.Kids[i][0]])
		case *frontend.Alternation:
			first = append(first, append(first[a.Kids[i][0]], first[a.Kids[i][1]]...))
		case *frontend.AltMatch:
			first = append(first, append(first[a.Kids[i][0]], first[a.Kids[i][1]]...))
		case *frontend.Concat:
			f := make([]int, 0, len(n.Items))
			for _, j := range a.Kids[i] {
				f = append(f, first[j]...)
				if !nullable[j] {
					break
				}
			}
			first = append(first, f)
		default:
			panic(fmt.Errorf("Unexpected type %T", n))
		}
	}
	a.first = first
	return first
}

// Last computes a look up table for each node in the tree (in post order,
// matching the Order member) indicating for the subtree rooted at each node
// which positions (leaf nodes indexes in the Positions slice) will appear at
// the end of a string matching that subtree.
func (a *LabeledAST) Last() (last [][]int) {
	if a.last != nil {
		return a.last
	}
	nullable := a.MatchesEmptyString()
	last = make([][]int, 0, len(a.Order))
	for i, node := range a.Order {
		switch n := node.(type) {
		case *frontend.EOS:
			last = append(last, []int{a.pos(i)})
		case *frontend.Character:
			last = append(last, []int{a.pos(i)})
		case *frontend.Range:
			last = append(last, []int{a.pos(i)})
		case *frontend.Maybe:
			last = append(last, last[a.Kids[i][0]])
		case *frontend.Plus:
			last = append(last, last[a.Kids[i][0]])
		case *frontend.Star:
			last = append(last, last[a.Kids[i][0]])
		case *frontend.Match:
			last = append(last, last[a.Kids[i][0]])
		case *frontend.Alternation:
			last = append(last, append(last[a.Kids[i][0]], last[a.Kids[i][1]]...))
		case *frontend.AltMatch:
			last = append(last, append(last[a.Kids[i][0]], last[a.Kids[i][1]]...))
		case *frontend.Concat:
			l := make([]int, 0, len(n.Items))
			for x := len(n.Items) - 1; x >= 0; x-- {
				j := a.Kids[i][x]
				l = append(l, last[j]...)
				if !nullable[j] {
					break
				}
			}
			last = append(last, l)
		default:
			panic(fmt.Errorf("Unexpected type %T", n))
		}
	}
	a.last = last
	return last
}
