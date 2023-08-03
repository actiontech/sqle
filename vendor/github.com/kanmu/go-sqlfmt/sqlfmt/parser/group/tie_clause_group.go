package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// TieClause such as UNION, EXCEPT, INTERSECT
type TieClause struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (tie *TieClause) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(tie.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, tie.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (tie *TieClause) IncrementIndentLevel(lev int) {
	tie.IndentLevel += lev
}
