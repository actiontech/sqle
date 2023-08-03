package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// LimitClause such as LIMIT, OFFSET, FETCH FIRST
type LimitClause struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (l *LimitClause) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(l.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, l.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (l *LimitClause) IncrementIndentLevel(lev int) {
	l.IndentLevel += lev
}
