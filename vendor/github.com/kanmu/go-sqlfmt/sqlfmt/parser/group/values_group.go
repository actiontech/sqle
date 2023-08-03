package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Values clause
type Values struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (val *Values) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(val.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, val.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (val *Values) IncrementIndentLevel(lev int) {
	val.IndentLevel += lev
}
