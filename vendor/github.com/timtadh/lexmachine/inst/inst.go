package inst

import (
	"fmt"
	"strings"
)

const (
	CHAR  = iota // CHAR instruction op code: match a byte in the range [X, Y] (inclusive)
	SPLIT        // SPLIT instruction op code: split jump to both X and Y
	JMP          // JMP instruction op code: jmp to X
	MATCH        // MATCH instruction op code: match the string
)

// Inst represents an NFA byte code instruction
type Inst struct {
	Op uint8
	X  uint32
	Y  uint32
}

// Slice is a list of NFA instructions
type Slice []*Inst

// New creates a new instruction
func New(op uint8, x, y uint32) *Inst {
	return &Inst{
		Op: op,
		X:  x,
		Y:  y,
	}
}

// String humanizes the byte code
func (i Inst) String() (s string) {
	switch i.Op {
	case CHAR:
		if i.X == i.Y {
			s = fmt.Sprintf("CHAR   %d (%q)", i.X, string([]byte{byte(i.X)}))
		} else {
			s = fmt.Sprintf("CHAR   %d (%q), %d (%q)", i.X, string([]byte{byte(i.X)}), i.Y, string([]byte{byte(i.Y)}))
		}
	case SPLIT:
		s = fmt.Sprintf("SPLIT  %v, %v", i.X, i.Y)
	case JMP:
		s = fmt.Sprintf("JMP    %v", i.X)
	case MATCH:
		s = "MATCH"
	}
	return
}

// Serialize outputs machine readable assembly
func (i Inst) Serialize() (s string) {
	switch i.Op {
	case CHAR:
		s = fmt.Sprintf("CHAR   %d, %d", i.X, i.Y)
	case SPLIT:
		s = fmt.Sprintf("SPLIT  %v, %v", i.X, i.Y)
	case JMP:
		s = fmt.Sprintf("JMP    %v", i.X)
	case MATCH:
		s = "MATCH"
	}
	return
}

// String humanizes the byte code
func (is Slice) String() (s string) {
	s = "{\n"
	for i, inst := range is {
		if inst == nil {
			continue
		}
		if i < 10 {
			s += fmt.Sprintf("    0%v %v\n", i, inst)
		} else {
			s += fmt.Sprintf("    %v %v\n", i, inst)
		}
	}
	s += "}"
	return
}

// Serialize outputs machine readable assembly
func (is Slice) Serialize() (s string) {
	lines := make([]string, 0, len(is))
	for i, inst := range is {
		lines = append(lines, fmt.Sprintf("%3d %s", i, inst.Serialize()))
	}
	return strings.Join(lines, "\n")
}
