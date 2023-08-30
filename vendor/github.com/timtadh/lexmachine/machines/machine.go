// Package machines implements the lexing algorithms.
package machines

import (
	"bytes"
	"fmt"
)

import (
	"github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/queue"
)

// EmptyMatchError is returned when a pattern would have matched the empty
// string
type EmptyMatchError struct {
	TC      int
	Line    int
	Column  int
	MatchID int
}

func (e *EmptyMatchError) Error() string {
	return fmt.Sprintf("Lexer error: matched the empty string at %d:%d (tc=%d) for match id %d.",
		e.Line, e.Column, e.TC, e.MatchID,
	)
}

// UnconsumedInput error type
type UnconsumedInput struct {
	StartTC     int
	FailTC      int
	StartLine   int
	StartColumn int
	FailLine    int
	FailColumn  int
	Text        []byte
}

// Error implements the error interface
func (u *UnconsumedInput) Error() string {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	stc := min(u.StartTC, len(u.Text)-1)
	etc := min(max(u.StartTC+1, u.FailTC), len(u.Text))
	return fmt.Sprintf("Lexer error: could not match text starting at %v:%v failing at %v:%v.\n\tunmatched text: %q",
		u.StartLine, u.StartColumn,
		u.FailLine, u.FailColumn,
		string(u.Text[stc:etc]),
	)
}

// A Match represents the positional and textual information from a match.
type Match struct {
	PC          int
	TC          int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Bytes       []byte // the actual bytes matched during scanning.
}

func computeLineCol(text []byte, prevTC, tc, line, col int) (int, int) {
	if tc < 0 {
		return line, col
	}
	if tc < prevTC {
		for i := prevTC; i > tc && i > 0; i-- {
			if text[i] == '\n' {
				line--
			}
		}
		col = 0
		for i := tc; i >= 0; i-- {
			if text[i] == '\n' {
				break
			}
			col++
		}
		return line, col
	}
	for i := prevTC + 1; i <= tc && i < len(text); i++ {
		if text[i] == '\n' {
			col = 0
			line++
		} else {
			col++
		}
	}
	if prevTC == tc && tc == 0 && tc < len(text) {
		if text[tc] == '\n' {
			line++
			col--
		}
	}
	return line, col
}

// Equals checks two matches for equality
func (m *Match) Equals(other *Match) bool {
	if m == nil && other == nil {
		return true
	} else if m == nil {
		return false
	} else if other == nil {
		return false
	}
	return m.PC == other.PC &&
		m.StartLine == other.StartLine &&
		m.StartColumn == other.StartColumn &&
		m.EndLine == other.EndLine &&
		m.EndColumn == other.EndColumn &&
		bytes.Equal(m.Bytes, other.Bytes)
}

// String formats the match for humans
func (m Match) String() string {
	return fmt.Sprintf("<Match %d %d (%d, %d)-(%d, %d) '%v'>", m.PC, m.TC, m.StartLine, m.StartColumn, m.EndLine, m.EndColumn, string(m.Bytes))
}

// Scanner is a functional iterator returned by the LexerEngine. See
// http://hackthology.com/functional-iteration-in-go.html
type Scanner func(int) (int, *Match, error, Scanner)

// LexerEngine does the actual tokenization of the byte slice text using the
// NFA bytecode in program. If the lexing process fails the Scanner will return
// an UnconsumedInput error.
func LexerEngine(program inst.Slice, text []byte) Scanner {
	done := false
	matchPC := -1
	matchTC := -1

	prevTC := 0
	line := 1
	col := 1

	var scan Scanner
	var cqueue, nqueue *queue.Queue = queue.New(len(program)), queue.New(len(program))
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done && tc == len(text) {
			return tc, nil, nil, nil
		}
		startTC := tc
		if tc < matchTC {
			// we back-tracked so reset the last matchTC
			matchTC = -1
		} else if tc == matchTC {
			// the caller did not reset the tc, we are where we left
		} else if matchTC != -1 && tc > matchTC {
			// we skipped text
			matchTC = tc
		}
		cqueue.Clear()
		nqueue.Clear()
		cqueue.Push(0)
		for ; tc <= len(text); tc++ {
			if cqueue.Empty() {
				break
			}
			for !cqueue.Empty() {
				pc := cqueue.Pop()
				i := program[pc]
				switch i.Op {
				case inst.CHAR:
					x := byte(i.X)
					y := byte(i.Y)
					if tc < len(text) && x <= text[tc] && text[tc] <= y {
						nqueue.Push(pc + 1)
					}
				case inst.MATCH:
					if matchTC < tc {
						matchPC = int(pc)
						matchTC = tc
					} else if matchPC > int(pc) {
						matchPC = int(pc)
						matchTC = tc
					}
				case inst.JMP:
					cqueue.Push(i.X)
				case inst.SPLIT:
					cqueue.Push(i.X)
					cqueue.Push(i.Y)
				default:
					panic(fmt.Errorf("unexpected instruction %v", i))
				}
			}
			cqueue, nqueue = nqueue, cqueue
			if cqueue.Empty() && matchPC > -1 {
				line, col = computeLineCol(text, prevTC, startTC, line, col)
				eLine, eCol := computeLineCol(text, startTC, matchTC-1, line, col)
				match := &Match{
					PC:          matchPC,
					TC:          startTC,
					StartLine:   line,
					StartColumn: col,
					EndLine:     eLine,
					EndColumn:   eCol,
					Bytes:       text[startTC:matchTC],
				}
				if matchTC == startTC {
					err := &EmptyMatchError{
						MatchID: matchPC,
						TC:      tc,
						Line:    line,
						Column:  col,
					}
					return startTC, nil, err, scan
				}
				prevTC = startTC
				matchPC = -1
				return matchTC, match, nil, scan
			}
		}
		if matchTC != len(text) && startTC >= len(text) {
			// the user has moved us farther than the text. Assume that was
			// the intent and return EOF.
			return tc, nil, nil, nil
		} else if matchTC != len(text) {
			done = true
			if matchTC == -1 {
				matchTC = 0
			}
			sline, scol := computeLineCol(text, 0, startTC, 1, 1)
			fline, fcol := computeLineCol(text, 0, tc, 1, 1)
			err := &UnconsumedInput{
				StartTC:     startTC,
				FailTC:      tc,
				StartLine:   sline,
				StartColumn: scol,
				FailLine:    fline,
				FailColumn:  fcol,
				Text:        text,
			}
			return tc, nil, err, scan
		} else {
			return tc, nil, nil, nil
		}
	}
	return scan
}
