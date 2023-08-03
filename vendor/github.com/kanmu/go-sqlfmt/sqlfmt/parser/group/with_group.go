package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// With clause
type With struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (w *With) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(w.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, w.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (w *With) IncrementIndentLevel(lev int) {
	w.IndentLevel += lev
}
