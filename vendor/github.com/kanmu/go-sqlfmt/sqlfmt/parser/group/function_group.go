package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Function clause
type Function struct {
	Element      []Reindenter
	IndentLevel  int
	InColumnArea bool
	ColumnCount  int
}

// Reindent reindents its elements
func (f *Function) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(f.Element)
	if err != nil {
		return err
	}

	for i, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			var prev lexer.Token

			if i > 0 {
				if preToken, ok := elements[i-1].(lexer.Token); ok {
					prev = preToken
				}
			}
			writeFunction(buf, token, prev, f.IndentLevel, f.ColumnCount, f.InColumnArea)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (f *Function) IncrementIndentLevel(lev int) {
	f.IndentLevel += lev
}
