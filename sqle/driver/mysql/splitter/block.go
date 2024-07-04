package splitter

import (
	"github.com/pingcap/parser"
	"strings"
)

type Block interface {
	MatchBegin(token *parser.Token) bool
	MatchEnd(token *parser.Token) bool
}

var allBlocks []Block = []Block{
	BeginEndBlock{},
	IfEndIfBlock{},
	CaseEndCaseBlock{},
	RepeatEndRepeatBlock{},
	WhileEndWhileBlock{},
	LoopEndLoopBlock{},
}

type LoopEndLoopBlock struct{}

func (b BeginEndBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.Begin
}

func (b BeginEndBlock) MatchEnd(token *parser.Token) bool {
	return true
}

type IfEndIfBlock struct{}

func (b IfEndIfBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.IfKwd
}

func (b IfEndIfBlock) MatchEnd(token *parser.Token) bool {
	return token.TokenType() == parser.IfKwd
}

type CaseEndCaseBlock struct{}

func (b CaseEndCaseBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.CaseKwd
}

func (b CaseEndCaseBlock) MatchEnd(token *parser.Token) bool {
	return token.TokenType() == parser.CaseKwd
}

type RepeatEndRepeatBlock struct{}

func (b RepeatEndRepeatBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.Repeat
}

func (b RepeatEndRepeatBlock) MatchEnd(token *parser.Token) bool {
	return token.TokenType() == parser.Repeat
}

type WhileEndWhileBlock struct{}

func (b WhileEndWhileBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.Identifier && strings.ToUpper(token.Ident()) == "WHILE"
}

func (b WhileEndWhileBlock) MatchEnd(token *parser.Token) bool {
	return token.TokenType() == parser.Identifier && strings.ToUpper(token.Ident()) == "WHILE"
}

type BeginEndBlock struct{}

func (b LoopEndLoopBlock) MatchBegin(token *parser.Token) bool {
	return token.TokenType() == parser.Identifier && strings.ToUpper(token.Ident()) == "LOOP"
}

func (b LoopEndLoopBlock) MatchEnd(token *parser.Token) bool {
	return token.TokenType() == parser.Identifier && strings.ToUpper(token.Ident()) == "LOOP"
}
