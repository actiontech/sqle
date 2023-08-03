package parser

import (
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// TODO: calling each Retrieve function is not smart, so should be refactored

// Parser parses Token Source
type parser struct {
	offset int
	result []group.Reindenter
	err    error
}

// ParseTokens parses Tokens, creating slice of Reindenter
// each Reindenter is group of SQL Clause such as SelectGroup, FromGroup ...etc
func ParseTokens(tokens []lexer.Token) ([]group.Reindenter, error) {
	if !isSQL(tokens[0].Type) {
		return nil, errors.New("can not parse no sql statement")
	}

	var (
		offset int
		result []group.Reindenter
	)

	for {
		if tokens[offset].Type == lexer.EOF {
			break
		}

		r := NewRetriever(tokens[offset:])
		element, endIdx, err := r.Retrieve()
		if err != nil {
			return nil, errors.Wrap(err, "ParseTokens failed")
		}

		group := createGroup(element)
		result = append(result, group)

		offset += endIdx
	}
	return result, nil
}

func isSQL(ttype lexer.TokenType) bool {
	return ttype == lexer.SELECT || ttype == lexer.UPDATE || ttype == lexer.DELETE || ttype == lexer.INSERT || ttype == lexer.LOCK || ttype == lexer.WITH
}
