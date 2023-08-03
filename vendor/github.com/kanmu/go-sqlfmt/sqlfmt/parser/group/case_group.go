package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Case Clause
type Case struct {
	Element        []Reindenter
	IndentLevel    int
	hasCommaBefore bool
}

// Reindent reindents its elements
func (c *Case) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(c.Element)
	if err != nil {
		return err
	}
	for _, v := range elements {
		if token, ok := v.(lexer.Token); ok {
			writeCase(buf, token, c.IndentLevel, c.hasCommaBefore)
		} else {
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified increment level
func (c *Case) IncrementIndentLevel(lev int) {
	c.IndentLevel += lev
}
