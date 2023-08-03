package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Where clause
type Where struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (w *Where) Reindent(buf *bytes.Buffer) error {
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
func (w *Where) IncrementIndentLevel(lev int) {
	w.IndentLevel += lev
}
