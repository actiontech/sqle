package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// From clause
type From struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (f *From) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(f.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, f.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel indents by its specified indent level
func (f *From) IncrementIndentLevel(lev int) {
	f.IndentLevel += lev
}
