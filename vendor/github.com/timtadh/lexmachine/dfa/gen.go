package dfa

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/linked"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
	"github.com/timtadh/lexmachine/frontend"
	"github.com/timtadh/lexmachine/machines"
)

// DFA is a Deterministic Finite-state Automaton
type DFA struct {
	minimal   bool
	Start     int                   // the starting state
	Error     int                   // the error state (should be 0)
	Accepting machines.DFAAccepting // state-idx to match-id
	Trans     machines.DFATrans     // the transition matrix
	Matches   [][]int               // match-id to list of accepting states
}

// Generate a DFA from a regular expressions AST. The generated DFA is
// minimized during the generation process.
func Generate(root frontend.AST) *DFA {
	ast := Label(root)
	positions := ast.Positions
	first, follow := ast.Follow()
	trans := hashtable.NewLinearHash()
	states := set.NewSortedSet(len(positions))
	accepting := set.NewSortedSet(len(positions))
	matchSet := set.NewSortedSet(len(ast.Matches))
	unmarked := linked.New()
	start := makeDState(first)
	trans.Put(start, make(map[byte]*set.SortedSet))
	states.Add(start)
	unmarked.Push(start)

	for _, m := range ast.Matches {
		matchSet.Add(types.Int(m))
	}

	for unmarked.Size() > 0 {
		x, err := unmarked.Pop()
		if err != nil {
			panic(err)
		}
		s := x.(*set.SortedSet)
		posBySymbol := make(map[int][]int)
		for pos, next := s.Items()(); next != nil; pos, next = next() {
			p := int(pos.(types.Int))
			switch n := ast.Order[positions[p]].(type) {
			case *frontend.EOS:
				posBySymbol[-1] = append(posBySymbol[-1], p)
			case *frontend.Character:
				sym := int(n.Char)
				posBySymbol[sym] = append(posBySymbol[sym], p)
			case *frontend.Range:
				for i := int(n.From); i <= int(n.To); i++ {
					posBySymbol[i] = append(posBySymbol[i], p)
				}
			}
		}
		for symbol, positions := range posBySymbol {
			if symbol == -1 {
				accepting.Add(s)
			} else if 0 <= symbol && symbol < 256 {
				// pFollow will be a new DState
				pFollow := set.NewSortedSet(len(positions) * 2)
				for _, p := range positions {
					for next := range follow[p] {
						pFollow.Add(types.Int(next))
					}
				}
				if !states.Has(pFollow) {
					trans.Put(pFollow, make(map[byte]*set.SortedSet))
					states.Add(pFollow)
					unmarked.Push(pFollow)
				}
				x, err := trans.Get(s)
				if err != nil {
					panic(err)
				}
				t := x.(map[byte]*set.SortedSet)
				t[byte(symbol)] = pFollow
			} else {
				panic("symbol outside of range")
			}
		}
	}

	idx := func(state *set.SortedSet) int {
		i, has, err := states.Find(state)
		if err != nil {
			panic(err)
		}
		if !has {
			panic(fmt.Errorf("missing state %v", state))
		}
		return i
	}

	dfa := &DFA{
		Start:     idx(start) + 1,
		Matches:   make([][]int, len(ast.Matches)),
		Accepting: make(machines.DFAAccepting),
		Trans:     make(machines.DFATrans, trans.Size()+1),
	}
	for k, v, next := trans.Iterate()(); next != nil; k, v, next = next() {
		from := k.(*set.SortedSet)
		toMap := v.(map[byte]*set.SortedSet)
		fromIdx := idx(from) + 1
		for symbol, to := range toMap {
			dfa.Trans[fromIdx][symbol] = idx(to) + 1
		}
		if accepting.Has(from) {
			idx := 0
			for ; idx < len(ast.Matches); idx++ {
				if from.Has(types.Int(ast.Matches[idx])) {
					break
				}
			}
			if idx >= len(ast.Matches) {
				panic(fmt.Errorf("Could not find any of %v in %v", ast.Matches, from))
			}
			dfa.Matches[idx] = append(dfa.Matches[idx], fromIdx)
			dfa.Accepting[fromIdx] = idx
		}
	}

	return dfa.minimize()
}

func makeDState(positions []int) *set.SortedSet {
	s := set.NewSortedSet(len(positions))
	for _, p := range positions {
		s.Add(types.Int(p))
	}
	return s
}

// String humanizes the DFA
func (dfa *DFA) String() string {
	lines := make([]string, 0, len(dfa.Trans))
	lines = append(lines, fmt.Sprintf("start: %d", dfa.Start))
	lines = append(lines, "accepting:")
	for i, matches := range dfa.Matches {
		t := make([]string, 0, len(matches))
		for _, m := range matches {
			t = append(t, fmt.Sprintf("%d", m))
		}
		lines = append(lines, fmt.Sprintf("    %d {%v}", i, strings.Join(t, ", ")))
	}
	for i, row := range dfa.Trans {
		t := make([]string, 0, 10)
		for sym, to := range row {
			if to == 0 {
				continue
			}
			t = append(t, fmt.Sprintf("%q->%v", string([]byte{byte(sym)}), to))
		}

		lines = append(lines, fmt.Sprintf("%d %v", i, strings.Join(t, ", ")))
	}
	return strings.Join(lines, "\n")
}

// Dotty translates the DFA into Graphviz source code suitable for visualization with
// the `dot` command.
func (dfa *DFA) Dotty() string {
	lines := make([]string, 0, len(dfa.Trans))
	lines = append(lines, "digraph DFA {")
	lines = append(lines, "rankdir=LR")
	lines = append(lines, "node [shape=circle]")
	lines = append(lines, `start [shape="none", label=""]`)
	lines = append(lines, fmt.Sprintf("start -> %d [label=start]", dfa.Start))
	for i, matches := range dfa.Matches {
		for _, m := range matches {
			lines = append(lines, fmt.Sprintf(`%d [style=bold, peripheries=2, xlabel="matches %d"]`, m, i))
		}
	}
	quote := func(s string) string {
		q := strconv.Quote(s)
		return q[1 : len(q)-1]
	}
	for i, row := range dfa.Trans {
		ranges := make(map[int][]struct{ beg, end int })
		target := -1
		beg := -1
		for sym, to := range row {
			if to != dfa.Error && target < 0 {
				beg = sym
				target = to
			}
			if to != target && target != -1 {
				end := sym - 1
				ranges[target] = append(ranges[target], struct{ beg, end int }{beg, end})
				if to == 0 {
					target = -1
					beg = -1
				} else {
					target = to
					beg = sym
				}
			}
		}
		if target != -1 {
			end := 255
			ranges[target] = append(ranges[target], struct{ beg, end int }{beg, end})
		}
		for to, trans := range ranges {
			sort.Slice(trans, func(i, j int) bool {
				x := byte(trans[i].beg)
				y := byte(trans[j].beg)
				if 'a' <= x && x <= 'z' && (('A' <= y && y <= 'Z') || ('0' <= y && y <= '9' || y == '_')) {
					return true
				}
				if 'A' <= x && x <= 'Z' && '0' <= y && y <= '9' {
					return true
				}
				if 'a' <= y && y <= 'z' && (('A' <= x && x <= 'Z') || ('0' <= x && x <= '9' || x == '_')) {
					return false
				}
				if 'A' <= y && y <= 'Z' && '0' <= x && x <= '9' {
					return false
				}
				return x < y
			})
			parts := make([]string, 0, len(trans))
			for _, t := range trans {
				if t.beg == t.end {
					parts = append(parts, quote(string([]byte{byte(t.beg)})))
				} else {
					parts = append(parts,
						quote(string([]byte{byte(t.beg)}))+"-"+
							quote(string([]byte{byte(t.end)})))
				}
			}
			lines = append(lines,
				fmt.Sprintf(
					`%d -> %d [label=%q]`,
					i, to, strings.Join(parts, "")))
		}
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// match tests whether a text matches a single expression in the DFA. The
// match-id is returned on success, otherwise a -1 is the result.
func (dfa *DFA) match(text string) int {
	s := dfa.Start
	for tc := 0; tc < len(text); tc++ {
		s = dfa.Trans[s][text[tc]]
	}
	if mid, has := dfa.Accepting[s]; has {
		return mid
	}
	return -1
}

func (dfa *DFA) minimize() *DFA {
	if dfa.minimal {
		return dfa
	}

	accepting := set.NewSortedSet(10)
	partition := set.NewSortedSet(10)
	for _, states := range dfa.Matches {
		group := set.NewSortedSet(10)
		for _, state := range states {
			group.Add(types.Int(state))
			accepting.Add(types.Int(state))
		}
		if group.Size() > 0 {
			partition.Add(group)
		}
	}
	nonAccepting := set.NewSortedSet(10)
	for state := range dfa.Trans {
		if state == dfa.Error {
			errGroup := set.NewSortedSet(1)
			errGroup.Add(types.Int(state))
			partition.Add(errGroup)
		} else {
			if !accepting.Has(types.Int(state)) {
				nonAccepting.Add(types.Int(state))
			}
		}
	}
	if nonAccepting.Size() > 0 {
		partition.Add(nonAccepting)
	}

	replace := func(i int, replacement *set.SortedSet) int {
		err := partition.Remove(i)
		if err != nil {
			panic(err)
		}
		err = partition.Extend(replacement.Items())
		if err != nil {
			panic(err)
		}
		first, err := replacement.Get(0)
		if err != nil {
			panic(err)
		}
		i, has, err := partition.Find(first)
		if err != nil {
			panic(err)
		} else if !has {
			panic(fmt.Errorf("Could not find %v in %v", first, partition))
		}
		return i
	}

	findGroup := func(s int) int {
		i := 0
		for v, next := partition.Items()(); next != nil; v, next = next() {
			g := v.(*set.SortedSet)
			if g.Has(types.Int(s)) {
				return i
			}
			i++
		}
		panic(fmt.Errorf("Could not find a group for %v in %v", s, partition))
	}

	equivalent := func(s int, ec *set.SortedSet) bool {
		x, err := ec.Get(0)
		if err != nil {
			panic(err)
		}
		t := int(x.(types.Int))
		for sym := 0; sym < 256; sym++ {
			a := findGroup(dfa.Trans[s][sym])
			b := findGroup(dfa.Trans[t][sym])
			if a != b || a < 0 || b < 0 {
				return false
			}
		}
		return true
	}

	for i := 0; i < partition.Size(); i++ {
		g, err := partition.Get(i)
		if err != nil {
			panic(err)
		}
		group := g.(*set.SortedSet)
		subgroups := set.NewSortedSet(10)
		for s, next := group.Items()(); next != nil; s, next = next() {
			state := int(s.(types.Int))
			found := false
			for ec, next := subgroups.Items()(); next != nil; ec, next = next() {
				eqClass := ec.(*set.SortedSet)
				if equivalent(state, eqClass) {
					eqClass.Add(types.Int(state))
					found = true
					break
				}
			}
			if !found {
				ec := set.NewSortedSet(10)
				ec.Add(s)
				subgroups.Add(ec)
			}
		}
		if subgroups.Size() > 1 {
			i = replace(i, subgroups) - 1
		}
	}

	// if the dfa is already minimal return it
	if partition.Size() == len(dfa.Trans) {
		dfa.minimal = true
		return dfa
	}

	newdfa := &DFA{
		minimal:   true,
		Error:     findGroup(dfa.Error),
		Start:     findGroup(dfa.Start),
		Matches:   make([][]int, len(dfa.Matches)),
		Accepting: make(machines.DFAAccepting),
		Trans:     make(machines.DFATrans, partition.Size()),
	}
	for gid := 0; gid < partition.Size(); gid++ {
		g, err := partition.Get(gid)
		if err != nil {
			panic(err)
		}
		group := g.(*set.SortedSet)
		r, err := group.Get(0)
		if err != nil {
			panic(err)
		}
		rep := int(r.(types.Int))
		for sym := 0; sym < 256; sym++ {
			newdfa.Trans[gid][sym] = findGroup(dfa.Trans[rep][sym])
		}
		if matchID, has := dfa.Accepting[rep]; has {
			newdfa.Matches[matchID] = append(newdfa.Matches[matchID], gid)
			newdfa.Accepting[gid] = matchID
		}
	}
	return newdfa
}
