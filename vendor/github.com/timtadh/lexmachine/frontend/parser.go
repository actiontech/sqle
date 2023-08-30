// Package frontend parses regular expressions and compiles them into NFA
// bytecode
package frontend

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
)

// Turn on debug prints
var DEBUG = false

// ParseError gives structured errors for parsing problems.
type ParseError struct {
	Reason     string
	Production string
	TC         int
	text       []byte
	chain      []*ParseError
}

// Errorf constructs a parse error with format for a particular location.
func Errorf(text []byte, tc int, format string, args ...interface{}) *ParseError {
	pc, _, _, ok := runtime.Caller(1)
	return errorf(pc, ok, text, tc, format, args...)
}

func matchErrorf(text []byte, tc int, format string, args ...interface{}) *ParseError {
	pc, _, _, ok := runtime.Caller(2)
	return errorf(pc, ok, text, tc, format, args...)
}

func errorf(pc uintptr, ok bool, text []byte, tc int, format string, args ...interface{}) *ParseError {
	var fn = "unknown"
	if ok {
		fn = runtime.FuncForPC(pc).Name()
		split := strings.Split(fn, ".")
		fn = split[len(split)-1]
	}
	msg := fmt.Sprintf(format, args...)
	return &ParseError{
		Reason:     msg,
		Production: fn,
		TC:         tc,
		text:       text,
	}
}

// Error implements the error interface
func (p *ParseError) Error() string {
	errs := make([]string, 0, len(p.chain)+1)
	for i := len(p.chain) - 1; i >= 0; i-- {
		errs = append(errs, p.chain[i].Error())
	}
	errs = append(errs, p.error())
	return strings.Join(errs, "\n")
}

func (p *ParseError) error() string {
	line, col := LineCol(p.text, p.TC)
	return fmt.Sprintf("Regex parse error in production '%v' : at index %v line %v column %v '%s' : %v",
		p.Production, p.TC, line, col, p.text[p.TC:], p.Reason)
}

// String formats the error for humans
func (p *ParseError) String() string {
	return p.Error()
}

// Chain joins multiple ParseErrors together
func (p *ParseError) Chain(e *ParseError) *ParseError {
	p.chain = append(p.chain, e)
	return p
}

// LineCol computes the line and column of a particular index inside of a byte
// slice.
func LineCol(text []byte, tc int) (line int, col int) {
	for i := 0; i <= tc && i < len(text); i++ {
		if text[i] == '\n' {
			col = 0
			line++
		} else {
			col++
		}
	}
	if tc == 0 && tc < len(text) {
		if text[tc] == '\n' {
			line++
			col--
		}
	}
	return line, col
}

// Parse a regular expression into an Abstract Syntax Tree (AST)
func Parse(text []byte) (AST, error) {
	a, err := (&parser{
		text:      text,
		lastError: Errorf(text, 0, "unconsumed input"),
	}).regex()
	if err != nil {
		return nil, err
	}
	return a, nil
}

type parser struct {
	text      []byte
	lastError *ParseError
}

func (p *parser) regex() (AST, *ParseError) {
	i, ast, err := p.alternation(0)
	if err != nil {
		return nil, err
	} else if i != len(p.text) {
		return nil, p.lastError
	}
	return NewMatch(ast), nil
}

func (p *parser) alternation(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter alternation %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit alternation %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, A, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	i, B, err := p.alternation_(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewAlternation(A, B), nil
}

func (p *parser) alternation_(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter alternation_ %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit alternation_ %v '%v'", i, string(p.text[i:]))
		}()
	}
	if i >= len(p.text) {
		return i, nil, nil
	}
	i, err := p.match(i, '|')
	if err != nil {
		return i, nil, nil
	}
	i, A, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	i, B, err := p.alternation_(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewAlternation(A, B), nil
}

func (p *parser) atomicOps(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomicOps %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit atomicOps %v '%v'", i, string(p.text[i:]))
		}()
	}
	if i >= len(p.text) {
		return i, nil, nil
	}
	i, A, err := p.atomicOp(i)
	if err != nil {
		p.lastError.Chain(err)
		return i, nil, nil
	}
	i, B, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewConcat(A, B), nil
}

func (p *parser) atomicOp(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomicOp %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit atomicOp %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, A, err := p.atomic(i)
	if DEBUG {
		log.Printf("atomic %v", err)
	}
	if err != nil {
		return i, nil, err
	}
	i, OPS, err := p.ops(i)
	if err != nil && err.Reason == "No Operator" {
		return i, A, nil
	} else if err != nil {
		return i, A, err
	}
	var N = A
	for _, OP := range OPS {
		N = NewApplyOp(OP, N)
	}
	return i, N, err
}

func (p *parser) ops(i int) (int, []AST, *ParseError) {
	if DEBUG {
		log.Printf("enter ops %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit ops %v '%v'", i, string(p.text[i:]))
		}()
	}
	ops := make([]AST, 0, 2)
	var err *ParseError
	var O AST
	for {
		i, O, err = p.op(i)
		if err != nil {
			if len(ops) <= 0 {
				return i, nil, err
			}
			return i, ops, nil
		}
		ops = append(ops, O)
	}
}

func (p *parser) op(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter op %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit op %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, err := p.match(i, '+')
	if err == nil {
		return i, NewOp("+"), nil
	}
	i, err = p.match(i, '*')
	if err == nil {
		return i, NewOp("*"), nil
	}
	i, err = p.match(i, '?')
	if err == nil {
		return i, NewOp("?"), nil
	}
	return i, nil, Errorf(p.text, i, "No Operator")
}

func (p *parser) atomic(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomic %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit atomic %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, ast, errChar := p.char(i)
	if errChar == nil {
		return i, ast, nil
	}
	if DEBUG {
		log.Printf("char %v", errChar)
	}
	i, ast, errGroup := p.group(i)
	if errGroup == nil {
		return i, ast, nil
	}
	if DEBUG {
		log.Printf("group %v", errGroup)
	}
	return i, nil, Errorf(p.text, i, "Expected group or char").Chain(errChar).Chain(errGroup)
}

func (p *parser) group(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter group %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit group %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, err := p.match(i, '(')
	if err != nil {
		return i, nil, err
	}
	i, A, err := p.alternation(i)
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, ')')
	if err != nil {
		return i, nil, err
	}
	return i, A, nil
}

func (p *parser) char(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter char %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit char %v '%v'", i, string(p.text[i:]))
		}()
	}
	i, C, errCHAR := p.CHAR(i)
	if errCHAR == nil {
		return i, C, nil
	}
	i, R, errRange := p.charClass(i)
	if errRange == nil {
		return i, R, nil
	}
	return i, nil, Errorf(p.text, i,
		"Expected a CHAR or charRange at %d, %v", i, string(p.text)).Chain(errCHAR).Chain(errRange)
}

// The CHAR token
func (p *parser) CHAR(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter CHAR %v '%v'", i, string(p.text[i:]))
		defer func() {
			log.Printf("exit CHAR %v '%v'", i, string(p.text[i:]))
		}()
	}
	if i >= len(p.text) {
		return i, nil, Errorf(p.text, i, "out of input %v, %v", i, string(p.text))
	}
	if p.text[i] == '\\' {
		i, cls, err := p.builtInClass(i)
		if err == nil {
			return i, cls, nil
		}
		i, b, err := p.getByte(i)
		if err != nil {
			return i, nil, err
		}
		return i, NewCharacter(b), nil
	}
	switch p.text[i] {
	case '|', '+', '*', '?', '(', ')', '[', ']', '^':
		return i, nil, Errorf(p.text, i,
			"unexpected operator, %s", string([]byte{p.text[i]}))
	case '.':
		return i + 1, NewAny(), nil
	default:
		return i + 1, NewCharacter(p.text[i]), nil
	}
}

var (
	builtInd = canonizeRanges([]*Range{NewRange(48, 57)})
	builtInD = invertRanges(builtInd)
	builtIns = canonizeRanges([]*Range{
		NewRange(9, 9),   // \t
		NewRange(10, 10), // \n
		NewRange(12, 12), // \f
		NewRange(13, 13), // \r
		NewRange(32, 32), // ' ' (a space)
	})
	builtInS = invertRanges(builtIns)
	builtInw = canonizeRanges([]*Range{
		NewRange(48, 57),  // 0-9
		NewRange(65, 90),  // A-Z
		NewRange(97, 122), // a-z
		NewRange(95, 95),  // _
	})
	builtInW = invertRanges(builtInw)
)

func (p *parser) builtInClass(i int) (int, AST, *ParseError) {
	if p.text[i] != '\\' {
		return i, nil, Errorf(p.text, i, "Not the start of built-in character class %q", string([]byte{p.text[i]}))
	}
	if i+1 < len(p.text) {
		if p.text[i+1] == 'd' {
			return i + 2, rangesToAST(builtInd), nil
		} else if p.text[i+1] == 'D' {
			return i + 2, rangesToAST(builtInD), nil
		} else if p.text[i+1] == 's' {
			return i + 2, rangesToAST(builtIns), nil
		} else if p.text[i+1] == 'S' {
			return i + 2, rangesToAST(builtInS), nil
		} else if p.text[i+1] == 'w' {
			return i + 2, rangesToAST(builtInw), nil
		} else if p.text[i+1] == 'W' {
			return i + 2, rangesToAST(builtInW), nil
		}
		return i, nil, Errorf(p.text, i, "Unknown class %q", string([]byte{p.text[i+1]}))
	}
	return i, nil, Errorf(p.text, i, "Unexpected EOS")
}

func (p *parser) getByte(i int) (int, byte, *ParseError) {
	i, err := p.match(i, '\\')
	if err == nil {
		if i >= len(p.text) {
			return len(p.text), p.text[len(p.text)-1], nil
		} else if i < len(p.text) && p.text[i] == 'n' {
			return i + 1, '\n', nil
		} else if i < len(p.text) && p.text[i] == 'r' {
			return i + 1, '\r', nil
		} else if i < len(p.text) && p.text[i] == 't' {
			return i + 1, '\t', nil
		}
		return i + 1, p.text[i], nil
	}
	if i >= len(p.text) {
		return i, 0, Errorf(p.text, i, "ran out of p.text at %d", i)
	}
	return i + 1, p.text[i], nil
}

func (p *parser) charClass(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter charRange %v '%v'", i, string(p.text[i:]))
	}
	exclude := false
	i, err := p.match(i, '[')
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, '^')
	if err == nil {
		exclude = true
	}
	ranges := make([]*Range, 0, 10)
	for {
		var r *Range
		i, r, err = p.charClassItem(i)
		if err != nil {
			return i, nil, err
		}
		ranges = append(ranges, r)
		if i < len(p.text) && p.text[i] == ']' {
			break
		} else if i >= len(p.text) {
			break
		}
	}
	i, err = p.match(i, ']')
	if err != nil {
		return i, nil, err
	}
	ranges = canonizeRanges(ranges)
	if exclude {
		ranges = invertRanges(ranges)
	}
	ast := rangesToAST(ranges)
	return i, ast, err
}

func (p *parser) charClassItem(i int) (int, *Range, *ParseError) {
	i, S, err := p.getByte(i)
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, '-')
	if err != nil {
		return i, NewRange(S, S), nil
	}
	i, T, err := p.getByte(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewRange(S, T), nil
}

func canonizeRanges(from []*Range) []*Range {
	sort.SliceStable(from, func(i, j int) bool {
		return from[i].From < from[j].From || (from[i].From == from[j].From && from[i].To < from[j].To)
	})
	to := make([]*Range, 0, len(from))
	var prev *Range
	for _, r := range from {
		if prev != nil {
			if prev.To+1 >= r.From {
				if prev.From >= r.From {
					// no change drop r it is a subset or prev
				} else {
					// extend prev to the end of r
					prev.To = r.To
				}
			} else {
				// r start after prev ends add a new range
				to = append(to, r)
				prev = r
			}
		} else {
			to = append(to, r)
			prev = r
		}
	}
	return to
}

// Expects p.combineOverlaps() to have been run on orig s.t. there are no
// overlapping ranges, the ranges are sorted, and all ranges are separated by at
// least one character.
func invertRanges(orig []*Range) []*Range {
	if len(orig) <= 0 {
		return []*Range{NewAny()}
	}
	invrt := make([]*Range, 0, len(orig)+1)
	if orig[0].From > 0 {
		invrt = append(invrt, NewRange(0, orig[0].From-1))
	}
	for i := 0; i+1 < len(orig); i++ {
		if orig[i].To < 255 {
			invrt = append(invrt, NewRange(orig[i].To+1, orig[i+1].From-1))
		}
	}
	if orig[len(orig)-1].To < 255 {
		invrt = append(invrt, NewRange(orig[len(orig)-1].To+1, 255))
	}
	return invrt
}

// Expects p.combineOverlaps() to have been run on orig s.t. there are no
// overlapping ranges, the ranges are sorted, and all ranges are separated by at
// least one character.
func rangesToAST(ranges []*Range) AST {
	if len(ranges) == 0 {
		panic("no ranges")
	} else if len(ranges) == 1 {
		return ranges[0]
	}
	ast := NewAlternation(
		ranges[len(ranges)-2],
		ranges[len(ranges)-1],
	)
	for j := len(ranges) - 3; j >= 0; j-- {
		ast = NewAlternation(
			ranges[j],
			ast,
		)
	}
	return ast
}

func (p *parser) matchAny(i int) (int, *Character, *ParseError) {
	if i >= len(p.text) {
		return i, nil, Errorf(p.text, i, "out of p.text, %d", i)
	}
	return i + 1, NewCharacter(p.text[i]), nil
}

func (p *parser) match(i int, c byte) (int, *ParseError) {
	if i >= len(p.text) {
		return i, matchErrorf(p.text, i, "out of p.text, %d", i)
	} else if p.text[i] == c {
		i++
		return i, nil
	}
	return i, matchErrorf(p.text, i,
		"expected '%v' at %v got '%v' of '%v'",
		string([]byte{c}),
		i,
		string(p.text[i:i+1]),
		string(p.text[i:]),
	)
}
