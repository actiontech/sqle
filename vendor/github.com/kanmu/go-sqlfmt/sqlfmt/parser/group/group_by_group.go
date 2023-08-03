package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// GroupBy clause
type GroupBy struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (g *GroupBy) Reindent(buf *bytes.Buffer) error {
	columnCount = 0

	elements, err := processPunctuation(g.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(elements) {
		switch v := el.(type) {
		case lexer.Token, string:
			if err := writeWithComma(buf, v, g.IndentLevel); err != nil {
				return err
			}
		case Reindenter:
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (g *GroupBy) IncrementIndentLevel(lev int) {
	g.IndentLevel += lev
}
