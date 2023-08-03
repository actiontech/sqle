package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Delete clause
type Delete struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (d *Delete) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(d.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, d.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (d *Delete) IncrementIndentLevel(lev int) {
	d.IndentLevel += lev
}
