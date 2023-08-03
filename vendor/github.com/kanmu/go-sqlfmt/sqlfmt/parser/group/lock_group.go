package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Lock clause
type Lock struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindent its elements
func (l *Lock) Reindent(buf *bytes.Buffer) error {
	for _, v := range l.Element {
		if token, ok := v.(lexer.Token); ok {
			writeLock(buf, token)
		} else {
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified increment level
func (l *Lock) IncrementIndentLevel(lev int) {
	l.IndentLevel += lev
}
