package frontend

import (
	"fmt"
	"strings"
)

// AST is an abstract syntax tree for a regular expression.
type AST interface {
	String() string
	Children() []AST
	Equals(AST) bool
}

// AltMatch either match A or B and then finalize matching the string
type AltMatch struct {
	A AST
	B AST
}

// Children returns a list of the child nodes
func (a *AltMatch) Children() []AST {
	return []AST{a.A, a.B}
}

// String humanizes the subtree
func (a *AltMatch) String() string {
	return fmt.Sprintf("(AltMatch %v, %v)", a.A, a.B)
}

// EOS end of string
type EOS struct{}

// Children returns a list of the child nodes
func (e *EOS) Children() []AST {
	return []AST{}
}

// String humanizes the subtree
func (e *EOS) String() string {
	return "(EOS)"
}

// Match the tree AST finalizes the matching
type Match struct {
	AST
}

// Children returns a list of the child nodes
func (m *Match) Children() []AST {
	return []AST{m.AST}
}

// String humanizes the subtree
func (m *Match) String() string {
	return fmt.Sprintf("(Match %v)", m.AST)
}

// Alternation matches A or B
type Alternation struct {
	A AST
	B AST
}

// Children returns a list of the child nodes
func (a *Alternation) Children() []AST {
	return []AST{a.A, a.B}
}

// String humanizes the subtree
func (a *Alternation) String() string {
	return fmt.Sprintf("(Alternation %v, %v)", a.A, a.B)
}

// Star is a kleene star (that is a repetition operator). Matches 0 or more times.
type Star struct {
	AST
}

// Children returns a list of the child nodes
func (s *Star) Children() []AST {
	return []AST{s.AST}
}

// String humanizes the subtree
func (s *Star) String() string {
	return fmt.Sprintf("(* %v)", s.AST)
}

// Plus matches 1 or more times
type Plus struct {
	AST
}

// Children returns a list of the child nodes
func (p *Plus) Children() []AST {
	return []AST{p.AST}
}

// String humanizes the subtree
func (p *Plus) String() string {
	return fmt.Sprintf("(+ %v)", p.AST)
}

// Maybe matches 0 or 1 times
type Maybe struct {
	AST
}

// Children returns a list of the child nodes
func (m *Maybe) Children() []AST {
	return []AST{m.AST}
}

// String humanizes the subtree
func (m *Maybe) String() string {
	return fmt.Sprintf("(? %v)", m.AST)
}

// Concat matches each item in sequence
type Concat struct {
	Items []AST
}

// Children returns a list of the child nodes
func (c *Concat) Children() []AST {
	return c.Items
}

// String humanizes the subtree
func (c *Concat) String() string {
	s := "(Concat "
	items := make([]string, 0, len(c.Items))
	for _, i := range c.Items {
		items = append(items, i.String())
	}
	s += strings.Join(items, ", ") + ")"
	return s
}

// Range matches byte ranges From-To inclusive
type Range struct {
	From byte
	To   byte
}

// Children returns a list of the child nodes
func (r *Range) Children() []AST {
	return []AST{}
}

// String humanizes the subtree
func (r *Range) String() string {
	return fmt.Sprintf(
		"(Range %d %d)",
		r.From,
		r.To,
	)
}

// Character matches a single byte
type Character struct {
	Char byte
}

// Children returns a list of the child nodes
func (c *Character) Children() []AST {
	return []AST{}
}

// String humanizes the subtree
func (c *Character) String() string {
	return fmt.Sprintf(
		"(Character %s)",
		string([]byte{c.Char}),
	)
}

// NewAltMatch creates an AltMatch
func NewAltMatch(a, b AST) AST {
	if a == nil || b == nil {
		panic("Alt match does not except nils")
	}
	return &AltMatch{a, b}
}

// NewMatch create a Match
func NewMatch(ast AST) AST {
	return &Match{NewConcat(ast, NewEOS())}
}

// NewEOS creates a EOS
func NewEOS() AST {
	return &EOS{}
}

// NewAlternation creates an Alternation
func NewAlternation(choice, alternation AST) AST {
	if alternation == nil {
		return choice
	}
	return &Alternation{choice, alternation}
}

// NewApplyOp applies a given op (Star, Plus, Maybe) to a tree
func NewApplyOp(op, atomic AST) AST {
	switch o := op.(type) {
	case *Star:
		o.AST = atomic
	case *Plus:
		o.AST = atomic
	case *Maybe:
		o.AST = atomic
	default:
		panic("unexpected op")
	}
	return op
}

// NewOp constructs a Star, Plus, or Maybe from *, +, ? respectively
func NewOp(op string) AST {
	switch op {
	case "*":
		return &Star{}
	case "+":
		return &Plus{}
	case "?":
		return &Maybe{}
	default:
		panic("unexpected op")
	}
}

// NewConcat concatenates two tree together
func NewConcat(char, concat AST) AST {
	if concat == nil {
		return char
	}
	if cc, ok := concat.(*Concat); ok {
		items := make([]AST, len(cc.Items)+1)
		items[0] = char
		for i, item := range cc.Items {
			items[i+1] = item
		}
		return &Concat{items}
	}
	return &Concat{[]AST{char, concat}}
}

// NewCharacter constructs a character
func NewCharacter(b byte) *Character {
	return &Character{b}
}

// NewAny constructs the . operator
func NewAny() *Range {
	return &Range{From: 0, To: 255}
}

// NewRange constructs a range operator
func NewRange(from, to byte) *Range {
	if from <= to {
		return &Range{From: from, To: to}
	}
	return &Range{From: to, To: from}
}
