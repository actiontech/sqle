package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Returning clause
type Returning struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (r *Returning) Reindent(buf *bytes.Buffer) error {
	columnCount = 0

	src, err := processPunctuation(r.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if err := writeWithComma(buf, v, r.IndentLevel); err != nil {
				return err
			}
		case Reindenter:
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (r *Returning) IncrementIndentLevel(lev int) {
	r.IndentLevel += lev
}
