package splitter

import (
	"strings"

	"github.com/pingcap/parser"
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
	// 如果下一个 token 是 CASE，说明是 END CASE 形式（控制流语句）
	if token.TokenType() == parser.CaseKwd {
		return true
	}

	// 如果下一个 token 不是其他控制流关键字，则认为这个 END 就是结束当前 CASE 块的
	// 这处理了表达式形式的 CASE 语句，如 CASE WHEN ... END
	switch token.TokenType() {
	case parser.IfKwd, parser.Repeat:
		// 如果是 IF 或 REPEAT，说明这个 END 是开始新的控制流块，不是结束 CASE
		return false
	case parser.Identifier:
		// 检查标识符是否是其他控制流关键字
		upperIdent := strings.ToUpper(token.Ident())
		if upperIdent == "WHILE" || upperIdent == "LOOP" {
			// 如果是 WHILE 或 LOOP，说明这个 END 是开始新的控制流块，不是结束 CASE
			return false
		}
	}

	// 其他情况下，认为这个 END 是结束当前 CASE 块的（表达式形式）
	return true
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
