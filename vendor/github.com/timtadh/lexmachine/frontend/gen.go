package frontend

import (
	"fmt"
)

import (
	"github.com/timtadh/lexmachine/inst"
)

type generator struct {
	program inst.Slice
}

// Generate an NFA program from the AST for a regular expression.
func Generate(ast AST) (inst.Slice, error) {
	g := &generator{
		program: make([]*inst.Inst, 0, 100),
	}
	fill := g.gen(ast)
	if len(fill) != 0 {
		return nil, fmt.Errorf("unconnected instructions")
	}
	return g.program, nil
}

func (g *generator) gen(ast AST) (fill []*uint32) {
	switch n := ast.(type) {
	case *AltMatch:
		fill = g.altMatch(n)
	case *Match:
		fill = g.match(n)
	case *Alternation:
		fill = g.alt(n)
	case *Star:
		fill = g.star(n)
	case *Plus:
		fill = g.plus(n)
	case *Maybe:
		fill = g.maybe(n)
	case *Concat:
		fill = g.concat(n)
	case *Character:
		fill = g.character(n)
	case *Range:
		fill = g.rangeGen(n)
	}
	return fill
}

func (g *generator) dofill(fill []*uint32) {
	for _, jmp := range fill {
		*jmp = uint32(len(g.program))
	}
}

func (g *generator) altMatch(a *AltMatch) []*uint32 {
	split := inst.New(inst.SPLIT, 0, 0)
	g.program = append(g.program, split)
	split.X = uint32(len(g.program))
	g.gen(a.A)
	split.Y = uint32(len(g.program))
	g.gen(a.B)
	return nil
}

func (g *generator) match(m *Match) []*uint32 {
	g.dofill(g.gen(m.AST))
	g.program = append(
		g.program, inst.New(inst.MATCH, 0, 0))
	return nil
}

func (g *generator) alt(a *Alternation) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	g.program = append(g.program, split)
	split.X = uint32(len(g.program))
	g.dofill(g.gen(a.A))
	jmp := inst.New(inst.JMP, 0, 0)
	g.program = append(g.program, jmp)
	split.Y = uint32(len(g.program))
	fill = g.gen(a.B)
	fill = append(fill, &jmp.X)
	return fill
}

func (g *generator) repeat(ast AST) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	splitPos := uint32(len(g.program))
	g.program = append(g.program, split)
	split.X = uint32(len(g.program))
	g.dofill(g.gen(ast))
	jmp := inst.New(inst.JMP, splitPos, 0)
	g.program = append(g.program, jmp)
	return []*uint32{&split.Y}
}

func (g *generator) star(s *Star) (fill []*uint32) {
	return g.repeat(s.AST)
}

func (g *generator) plus(p *Plus) (fill []*uint32) {
	g.dofill(g.gen(p.AST))
	return g.repeat(p.AST)
}

func (g *generator) maybe(m *Maybe) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	g.program = append(g.program, split)
	split.X = uint32(len(g.program))
	fill = g.gen(m.AST)
	fill = append(fill, &split.Y)
	return fill
}

func (g *generator) concat(c *Concat) (fill []*uint32) {
	for _, ast := range c.Items {
		g.dofill(fill)
		fill = g.gen(ast)
	}
	return fill
}

func (g *generator) character(ch *Character) []*uint32 {
	g.program = append(
		g.program,
		inst.New(inst.CHAR, uint32(ch.Char), uint32(ch.Char)))
	return nil
}

func (g *generator) rangeGen(r *Range) []*uint32 {
	g.program = append(
		g.program,
		inst.New(inst.CHAR, uint32(r.From), uint32(r.To)))
	return nil
}
