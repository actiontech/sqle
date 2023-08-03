package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Update clause
type Update struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindents its elements
func (u *Update) Reindent(buf *bytes.Buffer) error {
	columnCount = 0

	src, err := processPunctuation(u.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if err := writeWithComma(buf, v, u.IndentLevel); err != nil {
				return err
			}
		case Reindenter:
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (u *Update) IncrementIndentLevel(lev int) {
	u.IndentLevel += lev
}
