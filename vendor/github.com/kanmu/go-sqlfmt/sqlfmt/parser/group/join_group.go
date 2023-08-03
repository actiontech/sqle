package group

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Join clause
type Join struct {
	Element     []Reindenter
	IndentLevel int
}

// Reindent reindent its elements
func (j *Join) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(j.Element)
	if err != nil {
		return err
	}
	for i, v := range elements {
		if token, ok := v.(lexer.Token); ok {
			writeJoin(buf, token, j.IndentLevel, i == 0)
		} else {
			v.Reindent(buf)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified increment level
func (j *Join) IncrementIndentLevel(lev int) {
	j.IndentLevel += lev
}
