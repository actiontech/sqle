package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Insert clause
type Insert struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (insert *Insert) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(insert.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, insert.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (insert *Insert) IncrementIndentLevel(lev int) {
	insert.IndentLevel += lev
}
