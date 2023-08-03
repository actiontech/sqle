package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// Reindenter interface
// specific values of Reindenter would be clause group or token
type Reindenter interface {
	Reindent(buf *bytes.Buffer) error
	IncrementIndentLevel(lev int)
}

// count of ident appearing in column area
var columnCount int

// to reindent
const (
	NewLine          = "\n"
	WhiteSpace       = " "
	DoubleWhiteSpace = "  "
)

func write(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch {
	case token.IsNeedNewLineBefore():
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
	case token.Type == lexer.COMMA:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case token.Type == lexer.DO:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, token.Value, WhiteSpace))
	case strings.HasPrefix(token.Value, "::"):
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case token.Type == lexer.WITH:
		buf.WriteString(fmt.Sprintf("%s%s", NewLine, token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}

func writeWithComma(buf *bytes.Buffer, v interface{}, indent int) error {
	if token, ok := v.(lexer.Token); ok {
		switch {
		case token.IsNeedNewLineBefore():
			buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
		case token.Type == lexer.BY:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
		case token.Type == lexer.COMMA:
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.Value))
		default:
			return fmt.Errorf("can not reindent %#v", token.Value)
		}
	} else if str, ok := v.(string); ok {
		str = strings.TrimRight(str, " ")
		if columnCount == 0 {
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, str))
		} else if strings.HasPrefix(token.Value, "::") {
			buf.WriteString(fmt.Sprintf("%s", str))
		} else {
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, str))
		}
		columnCount++
	}
	return nil
}

func writeSelect(buf *bytes.Buffer, el interface{}, indent int) error {
	if token, ok := el.(lexer.Token); ok {
		switch token.Type {
		case lexer.SELECT, lexer.INTO:
			buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
		case lexer.AS, lexer.DISTINCT, lexer.DISTINCTROW, lexer.GROUP, lexer.ON:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
		case lexer.EXISTS:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
			columnCount++
		case lexer.COMMA:
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.Value))
		default:
			return fmt.Errorf("can not reindent %#v", token.Value)
		}
	} else if str, ok := el.(string); ok {
		str = strings.Trim(str, WhiteSpace)
		if columnCount == 0 {
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, str))
		} else {
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, str))
		}
		columnCount++
	}
	return nil
}

func writeCase(buf *bytes.Buffer, token lexer.Token, indent int, hasCommaBefore bool) {
	if hasCommaBefore {
		switch token.Type {
		case lexer.CASE:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
		case lexer.WHEN, lexer.ELSE:
			buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, WhiteSpace, DoubleWhiteSpace, token.Value))
		case lexer.END:
			buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, WhiteSpace, token.Value))
		case lexer.COMMA:
			buf.WriteString(fmt.Sprintf("%s", token.Value))
		default:
			if strings.HasPrefix(token.Value, "::") {
				buf.WriteString(fmt.Sprintf("%s", token.Value))
			} else {
				buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
			}
		}
	} else {
		switch token.Type {
		case lexer.CASE, lexer.END:
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.Value))
		case lexer.WHEN, lexer.ELSE:
			buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, DoubleWhiteSpace, token.Value))
		case lexer.COMMA:
			buf.WriteString(fmt.Sprintf("%s", token.Value))
		default:
			if strings.HasPrefix(token.Value, "::") {
				buf.WriteString(fmt.Sprintf("%s", token.Value))
			} else {
				buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
			}
		}
	}
}

func writeJoin(buf *bytes.Buffer, token lexer.Token, indent int, isFirst bool) {
	switch {
	case isFirst && token.IsJoinStart():
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
	case token.Type == lexer.ON || token.Type == lexer.USING:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
	case strings.HasPrefix(token.Value, "::"):
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}

func writeFunction(buf *bytes.Buffer, token, prev lexer.Token, indent, columnCount int, inColumnArea bool) {
	switch {
	case prev.Type == lexer.STARTPARENTHESIS || token.Type == lexer.STARTPARENTHESIS || token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case token.Type == lexer.FUNCTION && columnCount == 0 && inColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.Value))
	case token.Type == lexer.FUNCTION:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	case token.Type == lexer.COMMA:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case strings.HasPrefix(token.Value, "::"):
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}

func writeParenthesis(buf *bytes.Buffer, token lexer.Token, indent, columnCount int, inColumnArea, hasStartBefore bool) {
	switch {
	case token.Type == lexer.STARTPARENTHESIS && columnCount == 0 && inColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.Value))
	case token.Type == lexer.STARTPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	case token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case token.Type == lexer.COMMA:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case hasStartBefore:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	case strings.HasPrefix(token.Value, "::"):
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}

func writeSubquery(buf *bytes.Buffer, token lexer.Token, indent, columnCount int, inColumnArea bool) {
	switch {
	case token.Type == lexer.STARTPARENTHESIS && columnCount == 0 && inColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
	case token.Type == lexer.STARTPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	case token.Type == lexer.ENDPARENTHESIS && columnCount > 0:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.Value))
	case token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent-1), token.Value))
	case strings.HasPrefix(token.Value, "::"):
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}

func writeTypeCast(buf *bytes.Buffer, token lexer.Token) {
	switch token.Type {
	case lexer.TYPE:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	case lexer.COMMA:
		buf.WriteString(fmt.Sprintf("%s%s", token.Value, WhiteSpace))
	default:
		buf.WriteString(fmt.Sprintf("%s", token.Value))
	}
}

func writeLock(buf *bytes.Buffer, token lexer.Token) {
	switch token.Type {
	case lexer.LOCK:
		buf.WriteString(fmt.Sprintf("%s%s", NewLine, token.Value))
	case lexer.IN:
		buf.WriteString(fmt.Sprintf("%s%s", NewLine, token.Value))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.Value))
	}
}
