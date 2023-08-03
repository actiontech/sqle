package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// AndGroup is AND clause not AND operator
// AndGroup is made after new line
//// select xxx and xxx  <= this is not AndGroup
//// select xxx from xxx where xxx
//// and xxx      <= this is AndGroup
type AndGroup struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (a *AndGroup) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(a.Element)
	if err != nil {
		return err
	}

	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, a.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (a *AndGroup) IncrementIndentLevel(lev int) {
	a.IndentLevel += lev
}
