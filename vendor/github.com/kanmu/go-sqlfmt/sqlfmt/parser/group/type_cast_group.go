package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// TypeCast group
type TypeCast struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (t *TypeCast) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(t.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			writeTypeCast(buf, token)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (t *TypeCast) IncrementIndentLevel(lev int) {
	t.IndentLevel += lev
}
