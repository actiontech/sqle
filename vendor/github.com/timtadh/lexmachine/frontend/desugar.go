package frontend

import "fmt"

// DesugarRanges transform all Range nodes into Alternatives with individual characters
func DesugarRanges(ast AST) AST {
	switch n := ast.(type) {
	case *AltMatch:
		return &AltMatch{A: DesugarRanges(n.A), B: DesugarRanges(n.B)}
	case *Match:
		return &Match{AST: DesugarRanges(n.AST)}
	case *Alternation:
		return &Alternation{A: DesugarRanges(n.A), B: DesugarRanges(n.B)}
	case *Star:
		return &Star{AST: DesugarRanges(n.AST)}
	case *Plus:
		return &Plus{AST: DesugarRanges(n.AST)}
	case *Maybe:
		return &Maybe{AST: DesugarRanges(n.AST)}
	case *Concat:
		items := make([]AST, 0, len(n.Items))
		for _, i := range n.Items {
			items = append(items, DesugarRanges(i))
		}
		return &Concat{Items: items}
	case *Character:
		return n
	case *EOS:
		return n
	case *Range:
		chars := make([]*Character, 0, n.To-n.From+1)
		for i := int(n.From); i <= int(n.To); i++ {
			chars = append(chars, NewCharacter(byte(i)))
		}
		if len(chars) <= 0 {
			panic(fmt.Errorf("Empty, unmatchable range: %v", n))
		}
		if len(chars) == 1 {
			return chars[0]
		}
		alt := NewAlternation(chars[0], chars[1])
		for i := 2; i < len(chars); i++ {
			alt = NewAlternation(alt, chars[i])
		}
		return alt
	default:
		panic(fmt.Errorf("Unexpected node type %T", n))
	}
}
