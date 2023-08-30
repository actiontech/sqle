//go:build enterprise
// +build enterprise

package tidb_audit_log

import (
	"bytes"
	"fmt"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

const (
	TokenValue = "TokenValue"
	TokenStart = "TokenStart"
	TokenStop  = "TokenStop"
	TokenQuote = "TokenQuote"
	TokenKV    = "TokenKV"
)

var (
	tokens = []string{
		TokenValue,
		TokenStart,
		TokenStop,
		TokenQuote,
		TokenKV,
	}
	tokenMap = map[string] /*token name*/ int /*lexmachine.Token.Type*/ {}
)

func GetTokenName(tokenType int) string {
	return tokens[tokenType]
}

func GetTokenType(tokenName string) int {
	return tokenMap[tokenName]
}

func token(tokenName string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(GetTokenType(tokenName), string(m.Bytes), m), nil
	}
}

type byteScanner struct {
	text []byte
	pos  int
}

func (bs *byteScanner) Next() (byte, bool) {
	if bs.pos >= len(bs.text) {
		return 0, false
	}
	b := bs.text[bs.pos]
	bs.pos += 1
	return b, true
}

func (bs *byteScanner) LastN(count int) []byte {
	if bs.pos <= count {
		return bs.text[:bs.pos]
	}
	return bs.text[bs.pos-count : bs.pos]
}

var lexer = lexmachine.NewLexer()

func init() {
	for i, token := range tokens {
		tokenMap[token] = i
	}
	lexer.Add([]byte("\\["), token(TokenStart))
	lexer.Add([]byte("\\]"), token(TokenStop))
	lexer.Add([]byte("="), token(TokenKV))
	AddTokenBetween([]byte("\""), []byte("\""), false, token(TokenQuote))
}

func AddTokenBetween(left []byte, right []byte, matchEnd bool, action lexmachine.Action) {
	lexer.Add(left, func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		match.EndLine = match.StartLine
		match.EndColumn = match.StartColumn
		matchRight := false
		bs := byteScanner{text: scan.Text[scan.TC:]}
		for {
			b, has := bs.Next()
			if !has {
				break
			}
			match.EndColumn += 1
			if b == '\n' {
				match.EndLine += 1
			}
			if bytes.Equal(right, bs.LastN(len(right))) {
				matchRight = true
				break
			}
		}
		if matchRight {
			match.Bytes = scan.Text[scan.TC : scan.TC+bs.pos-len(right)]
			scan.TC = scan.TC + bs.pos
			match.TC = scan.TC
			return action(scan, match)
		}
		if matchEnd {
			match.Bytes = scan.Text[scan.TC:]
			scan.TC = scan.TC + bs.pos
			match.TC = scan.TC
			return action(scan, match)
		}

		return nil, fmt.Errorf("unclosed %s with %s, staring at %d, (%d, %d)",
			string(left), string(right), match.TC, match.StartLine, match.StartColumn)
	})
}
