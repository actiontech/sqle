package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// OrGroup clause
type OrGroup struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (o *OrGroup) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(o.Element)
	if err != nil {
		return err
	}

	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, o.IndentLevel)
		} else {
			el.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified increment level
func (o *OrGroup) IncrementIndentLevel(lev int) {
	o.IndentLevel += lev
}
