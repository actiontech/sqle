package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Set clause
type Set struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (s *Set) Reindent(buf *bytes.Buffer) error {
	columnCount = 0

	src, err := processPunctuation(s.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if err := writeWithComma(buf, v, s.IndentLevel); err != nil {
				return err
			}
		case Reindenter:
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (s *Set) IncrementIndentLevel(lev int) {
	s.IndentLevel += lev
}
