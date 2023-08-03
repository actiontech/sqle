package lexer

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

// Tokenizer tokenizes SQL statements
type Tokenizer struct {
	r      *bufio.Reader
	w      *bytes.Buffer // w  writes token value. It resets its value when the end of token appears
	result []Token
}

// rune that can't be contained in SQL statement
// TODO: I have to make better solution of making rune of eof in stead of using '∂'
var eof = '∂'

// value of literal
const (
	Comma            = ","
	StartParenthesis = "("
	EndParenthesis   = ")"
	StartBracket     = "["
	EndBracket       = "]"
	StartBrace       = "{"
	EndBrace         = "}"
	SingleQuote      = "'"
	NewLine          = "\n"
)

// NewTokenizer creates Tokenizer
func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		r: bufio.NewReader(strings.NewReader(src)),
		w: &bytes.Buffer{},
	}
}

// GetTokens returns tokens for parsing
func (t *Tokenizer) GetTokens() ([]Token, error) {
	var result []Token

	tokens, err := t.Tokenize()
	if err != nil {
		return nil, errors.Wrap(err, "Tokenize failed")
	}
	// replace all tokens without whitespaces and new lines
	// if "AND" or "OR" appears after new line, token value will be ANDGROUP, ORGROUP
	for i, tok := range tokens {
		if tok.Type == AND && tokens[i-1].Type == NEWLINE {
			andGroupToken := Token{Type: ANDGROUP, Value: tok.Value}
			result = append(result, andGroupToken)
			continue
		}
		if tok.Type == OR && tokens[i-1].Type == NEWLINE {
			orGroupToken := Token{Type: ORGROUP, Value: tok.Value}
			result = append(result, orGroupToken)
			continue
		}
		if tok.Type == WS || tok.Type == NEWLINE {
			continue
		}
		result = append(result, tok)
	}
	return result, nil
}

// Tokenize analyses every rune in SQL statement
// every token is identified when whitespace appears
func (t *Tokenizer) Tokenize() ([]Token, error) {
	for {
		isEOF, err := t.scan()

		if isEOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return t.result, nil
}

// unread undoes t.r.readRune method to get last character
func (t *Tokenizer) unread() { t.r.UnreadRune() }

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '　'
}

func isComma(ch rune) bool {
	return ch == ','
}

func isStartParenthesis(ch rune) bool {
	return ch == '('
}

func isEndParenthesis(ch rune) bool {
	return ch == ')'
}

func isSingleQuote(ch rune) bool {
	return ch == '\''
}

func isStartBracket(ch rune) bool {
	return ch == '['
}

func isEndBracket(ch rune) bool {
	return ch == ']'
}

func isStartBrace(ch rune) bool {
	return ch == '{'
}

func isEndBrace(ch rune) bool {
	return ch == '}'
}

// scan scans each character and appends to result until "eof" appears
// when it finishes scanning all characters, it returns true
func (t *Tokenizer) scan() (bool, error) {
	ch, _, err := t.r.ReadRune()
	if err != nil {
		if err.Error() == "EOF" {
			ch = eof
		} else {
			return false, errors.Wrap(err, "read rune failed")
		}
	}

	switch {
	case ch == eof:
		tok := Token{Type: EOF, Value: "EOF"}
		t.result = append(t.result, tok)
		return true, nil
	case isWhiteSpace(ch):
		if err := t.scanWhiteSpace(); err != nil {
			return false, err
		}
		return false, nil
	// extract string
	case isSingleQuote(ch):
		if err := t.scanString(); err != nil {
			return false, err
		}
		return false, nil
	case isComma(ch):
		token := Token{Type: COMMA, Value: Comma}
		t.result = append(t.result, token)
		return false, nil
	case isStartParenthesis(ch):
		token := Token{Type: STARTPARENTHESIS, Value: StartParenthesis}
		t.result = append(t.result, token)
		return false, nil
	case isEndParenthesis(ch):
		token := Token{Type: ENDPARENTHESIS, Value: EndParenthesis}
		t.result = append(t.result, token)
		return false, nil
	case isStartBracket(ch):
		token := Token{Type: STARTBRACKET, Value: StartBracket}
		t.result = append(t.result, token)
		return false, nil
	case isEndBracket(ch):
		token := Token{Type: ENDBRACKET, Value: EndBracket}
		t.result = append(t.result, token)
		return false, nil
	case isStartBrace(ch):
		token := Token{Type: STARTBRACE, Value: StartBrace}
		t.result = append(t.result, token)
		return false, nil
	case isEndBrace(ch):
		token := Token{Type: ENDBRACE, Value: EndBrace}
		t.result = append(t.result, token)
		return false, nil
	default:
		if err := t.scanIdent(); err != nil {
			return false, err
		}
		return false, nil
	}
}

func (t *Tokenizer) scanWhiteSpace() error {
	t.unread()

	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return err
			}
		}
		if !isWhiteSpace(ch) {
			t.unread()
			break
		} else {
			t.w.WriteRune(ch)
		}
	}

	if strings.Contains(t.w.String(), "\n") {
		tok := Token{Type: NEWLINE, Value: "\n"}
		t.result = append(t.result, tok)
	} else {
		tok := Token{Type: WS, Value: t.w.String()}
		t.result = append(t.result, tok)
	}
	t.w.Reset()
	return nil
}

// scan string token including single quotes
func (t *Tokenizer) scanString() error {
	var counter int
	t.unread()

	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return err
			}
		}
		// ignore the first single quote
		if counter != 0 && isSingleQuote(ch) {
			t.w.WriteRune(ch)
			break
		} else {
			t.w.WriteRune(ch)
		}
		counter++
	}
	tok := Token{Type: STRING, Value: t.w.String()}
	t.result = append(t.result, tok)
	t.w.Reset()
	return nil
}

// append all ch to result until ch is a white space
// if ident is keyword, Type will be the keyword and value will be the uppercase keyword
func (t *Tokenizer) scanIdent() error {
	t.unread()

	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return err
			}
		}
		if isWhiteSpace(ch) {
			t.unread()
			break
		} else if isComma(ch) {
			t.unread()
			break
		} else if isStartParenthesis(ch) {
			t.unread()
			break
		} else if isEndParenthesis(ch) {
			t.unread()
			break
		} else if isSingleQuote(ch) {
			t.unread()
			break
		} else if isStartBracket(ch) {
			t.unread()
			break
		} else if isEndBracket(ch) {
			t.unread()
			break
		} else if isStartBrace(ch) {
			t.unread()
			break
		} else if isEndBrace(ch) {
			t.unread()
			break
		} else {
			t.w.WriteRune(ch)
		}
	}
	t.append(t.w.String())
	return nil
}

func (t *Tokenizer) append(v string) {
	upperValue := strings.ToUpper(v)

	if ttype, ok := t.isSQLKeyWord(upperValue); ok {
		t.result = append(t.result, Token{
			Type:  ttype,
			Value: upperValue,
		})
	} else {
		t.result = append(t.result, Token{
			Type:  ttype,
			Value: v,
		})
	}
	t.w.Reset()
}

func (t *Tokenizer) isSQLKeyWord(v string) (TokenType, bool) {
	if ttype, ok := sqlKeywordMap[v]; ok {
		return ttype, ok
	} else if ttype, ok := typeWithParenMap[v]; ok {
		if r, _, err := t.r.ReadRune(); err == nil && string(r) == StartParenthesis {
			t.unread()
			return ttype, ok
		}
		t.unread()
		return IDENT, ok
	}
	return IDENT, false
}

var sqlKeywordMap = map[string]TokenType{
	"SELECT":      SELECT,
	"FROM":        FROM,
	"WHERE":       WHERE,
	"CASE":        CASE,
	"ORDER":       ORDER,
	"BY":          BY,
	"AS":          AS,
	"JOIN":        JOIN,
	"LEFT":        LEFT,
	"RIGHT":       RIGHT,
	"INNER":       INNER,
	"OUTER":       OUTER,
	"ON":          ON,
	"WHEN":        WHEN,
	"END":         END,
	"GROUP":       GROUP,
	"DESC":        DESC,
	"ASC":         ASC,
	"LIMIT":       LIMIT,
	"AND":         AND,
	"OR":          OR,
	"IN":          IN,
	"IS":          IS,
	"NOT":         NOT,
	"NULL":        NULL,
	"DISTINCT":    DISTINCT,
	"LIKE":        LIKE,
	"BETWEEN":     BETWEEN,
	"UNION":       UNION,
	"ALL":         ALL,
	"HAVING":      HAVING,
	"EXISTS":      EXISTS,
	"UPDATE":      UPDATE,
	"SET":         SET,
	"RETURNING":   RETURNING,
	"DELETE":      DELETE,
	"INSERT":      INSERT,
	"INTO":        INTO,
	"DO":          DO,
	"VALUES":      VALUES,
	"FOR":         FOR,
	"THEN":        THEN,
	"ELSE":        ELSE,
	"DISTINCTROW": DISTINCTROW,
	"FILTER":      FILTER,
	"WITHIN":      WITHIN,
	"COLLATE":     COLLATE,
	"INTERSECT":   INTERSECT,
	"EXCEPT":      EXCEPT,
	"OFFSET":      OFFSET,
	"FETCH":       FETCH,
	"FIRST":       FIRST,
	"ROWS":        ROWS,
	"USING":       USING,
	"OVERLAPS":    OVERLAPS,
	"NATURAL":     NATURAL,
	"CROSS":       CROSS,
	"ZONE":        ZONE,
	"NULLS":       NULLS,
	"LAST":        LAST,
	"AT":          AT,
	"LOCK":        LOCK,
	"WITH":        WITH,
}

var typeWithParenMap = map[string]TokenType{
	"SUM":             FUNCTION,
	"AVG":             FUNCTION,
	"MAX":             FUNCTION,
	"MIN":             FUNCTION,
	"COUNT":           FUNCTION,
	"COALESCE":        FUNCTION,
	"EXTRACT":         FUNCTION,
	"OVERLAY":         FUNCTION,
	"POSITION":        FUNCTION,
	"CAST":            FUNCTION,
	"SUBSTRING":       FUNCTION,
	"TRIM":            FUNCTION,
	"XMLELEMENT":      FUNCTION,
	"XMLFOREST":       FUNCTION,
	"XMLCONCAT":       FUNCTION,
	"RANDOM":          FUNCTION,
	"DATE_PART":       FUNCTION,
	"DATE_TRUNC":      FUNCTION,
	"ARRAY_AGG":       FUNCTION,
	"PERCENTILE_DISC": FUNCTION,
	"GREATEST":        FUNCTION,
	"LEAST":           FUNCTION,
	"OVER":            FUNCTION,
	"ROW_NUMBER":      FUNCTION,
	"BIG":             TYPE,
	"BIGSERIAL":       TYPE,
	"BOOLEAN":         TYPE,
	"CHAR":            TYPE,
	"BIT":             TYPE,
	"TEXT":            TYPE,
	"INTEGER":         TYPE,
	"NUMERIC":         TYPE,
	"DECIMAL":         TYPE,
	"DEC":             TYPE,
	"FLOAT":           TYPE,
	"CUSTOMTYPE":      TYPE,
	"VARCHAR":         TYPE,
	"VARBIT":          TYPE,
	"TIMESTAMP":       TYPE,
	"TIME":            TYPE,
	"SECOND":          TYPE,
	"INTERVAL":        TYPE,
}
