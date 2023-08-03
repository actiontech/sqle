package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Subquery group
type Subquery struct {
	Element      []Reindenter
	IndentLevel  int
	InColumnArea bool
	ColumnCount  int
}

// Reindent reindents its elements
func (s *Subquery) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(s.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			writeSubquery(buf, token, s.IndentLevel, s.ColumnCount, s.InColumnArea)
		} else {
			if s.InColumnArea {
				el.IncrementIndentLevel(1)
				el.Reindent(buf)
			} else {
				el.Reindent(buf)
			}
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (s *Subquery) IncrementIndentLevel(lev int) {
	s.IndentLevel += lev
}
