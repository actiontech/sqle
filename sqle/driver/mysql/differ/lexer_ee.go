//go:build enterprise
// +build enterprise

package differ

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MalformedSQLError represents a fatal problem parsing or lexing SQL: the input
// contains an unterminated quote or unterminated multi-line comment.
type MalformedSQLError struct {
	str        string
	filePath   string
	lineNumber int
	colNumber  int
}

// Error satisfies the builtin error interface.
func (mse *MalformedSQLError) Error() string {
	var parts []string
	if mse.filePath != "" {
		parts = append(parts, "File "+mse.filePath+": ")
	}
	if mse.str != "" {
		parts = append(parts, mse.str)
	} else {
		parts = append(parts, "Malformed SQL")
	}
	if mse.lineNumber > 0 {
		parts = append(parts, fmt.Sprintf(" at line %d", mse.lineNumber))
		if mse.colNumber > 1 {
			parts = append(parts, fmt.Sprintf(", column %d", mse.colNumber))
		}
	}
	return strings.Join(parts, "")
}

// TokenType represents the category of a lexical token.
type TokenType uint32

// Constants enumerating TokenType values
const (
	TokenNone       TokenType = iota // zero value for TokenType
	TokenWord                        // bare word, either a keyword or an unquoted identifier
	TokenIdent                       // backtick-wrapped identifier
	TokenString                      // string wrapped in either single quotes or double quotes
	TokenNumeric                     // int or float (unsigned; leading - will be treated as symbol by lexer)
	TokenSymbol                      // single operator or other symbol (always just one rune)
	TokenExtComment                  // C-style comment with contents beginning with !, +, or M! TODO not really handled well yet
	TokenDelimiter                   // token equal to current delimiter (whatever that happened to be)
	TokenFiller                      // mix of whitespaces and/or comments (other than tokenExtComment)
)

// Lexer is a simple partial SQL lexical tokenizer. It intentionally does not
// aim to be a complete solution for use in fully parsing SQL. See comment at
// the top of parser.go for the purpose of this implementation.
type Lexer struct {
	reader          *bufio.Reader
	bufferSize      int          // size of reader's buffer
	largeData       bytes.Buffer // temporary buffer for large token values
	err             error
	delimiter       string
	delimBytes      []byte
	delimTricky     bool // true if delimBytes[0] could appear at the start of an unquoted identifier or numeric literal
	commandMode     bool // true if \n (or \r\n) should emit TokenDelimiter instead of TokenFiller
	prevTokenFiller bool // true if previous Scan token was filler (whitespaces/comments)
}

// NewLexer returns a lexer which reads from r. The supplied delimiter string
// will be used initially, but may be changed by calling ChangeDelimiter.
func NewLexer(r io.Reader, delimiter string, bufferSize int) *Lexer {
	lex := &Lexer{
		reader:     bufio.NewReaderSize(r, bufferSize),
		bufferSize: bufferSize,
	}
	lex.ChangeDelimiter(delimiter)
	return lex
}

// Delimiter returns the lexer's current delimiter string.
func (lex *Lexer) Delimiter() string {
	return lex.delimiter
}

// ChangeDelimiter changes the lexer's delimiter string for all subsequent
// calls to Scan.
func (lex *Lexer) ChangeDelimiter(newDelimiter string) {
	lex.delimiter = newDelimiter
	lex.delimBytes = []byte(newDelimiter)
	lex.delimTricky = len(lex.delimBytes) > 0 && (lex.delimBytes[0] >= utf8.RuneSelf || isWord(lex.delimBytes[0]) || lex.delimBytes[0] == '-' || lex.delimBytes[0] == '.')
}

// Scan returns the next token data from the reader, returning the raw data and
// token type. If an error occurs mid-scan (*after* at least one rune was
// successfully read), it will be suppressed until the next call to Scan. In
// other words, if the returned data is non-nil, err will always be nil; and
// likewise if err is non-nil then data will be nil and typ will be TokenNone.
// Data returned by Scan() is only valid until the subsequent call to Scan().
func (lex *Lexer) Scan() (data []byte, typ TokenType, err error) {
	// Peek ahead in the buffered reader, in a way which only performs an
	// underlying read if the buffer is less than half full. Peek returns data
	// directly from the bufio.Reader's underlying buffer, so this is ideal for
	// minimizing the number of allocations and copies.
	peekSize := lex.reader.Buffered()
	if peekSize == 0 && lex.err != nil {
		return nil, TokenNone, lex.err
	} else if peekSize < (lex.bufferSize/2) && lex.err == nil {
		peekSize = lex.bufferSize / 2 // ensure we're doing an underlying read
	}
	var p []byte
	p, err = lex.reader.Peek(peekSize)
	if len(p) == 0 {
		if err == nil {
			err = io.EOF
		}
		lex.err = err
		return nil, TokenNone, lex.err
	} else if err != nil {
		lex.err = err
		err = nil
	}

	// Bookkeeping before examining next token
	lex.largeData.Reset()

	// Client commands such as USE may optionally be terminated by only a newline
	// instead of the usual delimiter. Meanwhile the DELIMITER command *must* be
	// newline-terminated. Newlines will emit TokenDelimiter in these situations.
	if lex.commandMode {
		var newlineDelimFoundLen int
		if p[0] == '\n' {
			newlineDelimFoundLen = 1
		} else if p[0] == '\r' && len(p) > 1 && p[1] == '\n' {
			newlineDelimFoundLen = 2
		}
		if newlineDelimFoundLen > 0 {
			lex.commandMode = false
			lex.prevTokenFiller = false
			lex.reader.Discard(newlineDelimFoundLen)
			return p[0:newlineDelimFoundLen], TokenDelimiter, nil
		}
	}

	if !lex.prevTokenFiller {
		if ok, _, _ := isFiller(p); ok {
			return lex.scanFiller(p)
		}
	}
	lex.prevTokenFiller = false // scanFiller ensures there are never 2 filler tokens in a row

	// Check for delimiter token. This must be done before any other non-filler
	// comparisons, since delimiter string can consist of any arbitrary rune(s).
	// A single trailing newline (LF or CRLF), if present, will be included in the
	// delimiter token.
	if bytes.HasPrefix(p, lex.delimBytes) {
		n := len(lex.delimBytes)
		if len(p) > n+1 && p[n] == '\r' && p[n+1] == '\n' {
			n += 2
		} else if len(p) > n && p[n] == '\n' {
			n++
		}
		lex.commandMode = false
		lex.reader.Discard(n)
		return p[0:n], TokenDelimiter, nil
	}

	var sawDecimalPoint, sawE bool // used in building numerics
	var size int
	if p[0] < utf8.RuneSelf { // single-byte rune
		if p[0] == '\'' || p[0] == '"' || p[0] == '`' {
			return lex.scanString(p)
		} else if p[0] >= '0' && p[0] <= '9' {
			typ = TokenNumeric
		} else if p[0] == '.' { // might be a symbol or might be a float/decimal numeric below 0 without the leading 0
			if len(p) == 1 || p[1] < '0' || p[1] > '9' {
				typ = TokenSymbol
			} else {
				typ = TokenNumeric
				sawDecimalPoint = true
			}
		} else if !isWord(p[0]) {
			// note: negative numbers are intentionally emitted as '-' symbol token
			// followed by a separate positive numeric token. Parser can re-combine if
			// ever needed.
			typ = TokenSymbol
		} else {
			typ = TokenWord
		}
		size = 1
	} else if _, size = utf8.DecodeRune(p); size > 3 {
		typ = TokenSymbol // 4-byte rune handled as symbol, since cannot be in unquoted identifier and won't be a keyword
	} else {
		typ = TokenWord
	}

	// Symbols (operators, periods, parens, etc) are always lexed as a single rune.
	// Parser could combine multi-rune operators if ever needed.
	if typ == TokenSymbol {
		lex.reader.Discard(size)
		return p[0:size], TokenSymbol, nil
	}

	// Remaining situations are handled as TokenWord or TokenNumeric. Care must be
	// taken to differentiate between the two, since unquoted identifiers can
	// include digits in any position, as long as the entire identifier isn't
	// digits.

	// It's unlikely we will need more data than p in legitimate SQL cases, but
	// handle very long tokens anyway. We'll try to read more data when we don't
	// have enough for one rune or the delimiter string, whichever is longer.
	minBytes := 4
	if len(lex.delimBytes) > minBytes {
		minBytes = len(lex.delimBytes)
	}

	var r rune
	n := size // skip past the first rune, which we already examined
	for n < len(p) {
		if p[n] < utf8.RuneSelf { // single-byte rune
			size = 1
			if isSpace(p[n]) {
				break
			} else if p[n] == '.' { // decimal point only part of the token if permitted in this position of a numeric token
				if typ != TokenNumeric || sawDecimalPoint || sawE || len(p) == n+1 || p[n+1] < '0' || p[n+1] > '9' {
					break
				}
				sawDecimalPoint = true
			} else if !isWord(p[n]) {
				break
			} else if typ == TokenNumeric && (p[n] < '0' || p[n] > '9') {
				if (p[n] == 'e' || p[n] == 'E') && !sawE && len(p) > n+1 && (p[n+1] == '-' || (p[n+1] >= '0' && p[n+1] <= '9')) {
					sawE = true // allow a single e in a specific allowed position of numerics
					size++      // skip past that next digit or minus sign as well; simplifies handling for numbers in form 123e-4
				} else if sawDecimalPoint {
					break // demical/float numeric token ends upon encountering non-digit, aside from specific 'e'/'E' cases above
				} else {
					typ = TokenWord // otherwise it was actually a word (unquoted identifier which happened to contain digits)
				}
			}
		} else if r, size = utf8.DecodeRune(p[n:]); size > 3 || unicode.IsSpace(r) {
			break // token definitely ends at 4-byte rune or multibyte/extended space
		} else if typ == TokenNumeric {
			if sawDecimalPoint {
				break // decimal/float numeric token ends upon encountering 2-3 byte rune
			}
			typ = TokenWord // otherwise consider it a word (unquoted identifier which happened to contain digits)
		}
		if lex.delimTricky && bytes.HasPrefix(p[n:], lex.delimBytes) {
			break
		}

		n += size
		if lex.err == nil && len(p)-n < minBytes {
			p, n = lex.bufferAndPeek(p[0:n]), 0
		}
	}
	return lex.buildReturn(p[0:n], typ)
}

// ScanBOM peeks at the next 3 bytes of the reader to see if they contain a UTF8
// byte-order marker. This method should generally only be called on a new lexer
// prior to any other Scan().
func (lex *Lexer) ScanBOM() bool {
	var r rune
	r, _, lex.err = lex.reader.ReadRune()
	if r == '\uFEFF' {
		return true
	}
	lex.reader.UnreadRune() // harmless even if lex.err is non-nil (UnreadRune will just fail)
	return false
}

// bufferAndPeek is a helper method which copies p into lex.largeData, advances
// the reader to p's length, and then peeks for more data. If an error occurs,
// it is stored into lex.err for use in a subsequent Scan. The caller should
// always ensure lex.err is nil before calling bufferAndPeek.
func (lex *Lexer) bufferAndPeek(p []byte) []byte {
	lex.largeData.Write(p)
	lex.reader.Discard(len(p))
	p, lex.err = lex.reader.Peek(lex.bufferSize - lex.reader.Buffered())
	return p
}

// buildReturn is a helper method for returning values from Scan. It advances
// the reader to be right after p, and returns data that properly accounts for
// lex.largeData if non-empty.
func (lex *Lexer) buildReturn(p []byte, typ TokenType) ([]byte, TokenType, error) {
	lex.reader.Discard(len(p))
	if lex.largeData.Len() > 0 { // read a large enough token to need external buffer
		lex.largeData.Write(p)
		return lex.largeData.Bytes(), typ, nil
	}
	return p, typ, nil
}

var (
	needleNewline      = []byte{'\n'}
	needleCloseComment = []byte{'*', '/'}
)

// isFiller is a helper method for identifying the start of a span of spaces or
// comments. The supplied p must have len(1) or more, otherwise isFiller will
// panic. If p begins with whitespace or marks the beginning of a comment, ok
// will be true, and prefixLen will indicate how many bytes of the start of p
// are definitely part of the filler token. In the case of a comment,
// closerNeedle will be the byte(s) that indicate the end of the comment.
func isFiller(p []byte) (ok bool, prefixLen int, closerNeedle []byte) {
	if p[0] >= utf8.RuneSelf { // multi-byte rune at the start of p, or non-utf8 e.g. extended ascii nbsp
		if r, size := utf8.DecodeRune(p); unicode.IsSpace(r) {
			return true, size, nil
		}
		return false, 0, nil
	}
	if isSpace(p[0]) { // efficiently check for single-byte spaces
		return true, 1, nil
	} else if p[0] == '#' {
		return true, 1, needleNewline
	} else if len(p) == 1 {
		return false, 0, nil // already know the one byte in p isn't a space, and can't start /* or -- with one byte
	} else if p[0] == '-' && p[1] == '-' {
		if r, _ := utf8.DecodeRune(p[2:]); unicode.IsSpace(r) || len(p) == 2 {
			return true, 2, needleNewline // don't include r's length in prefixLen, since r may already be the terminating newline!
		}
		return false, 0, nil // "--" followed by non-space: not a comment, and n marks the beginning of next token
	} else if p[0] == '/' && p[1] == '*' {
		return true, 2, needleCloseComment
	}
	return false, 0, nil
}

// scanFiller combines contiguous whitespace and/or comments into a single
// token.
func (lex *Lexer) scanFiller(p []byte) (data []byte, typ TokenType, err error) {
	var n int
	var needle []byte
	for n < len(p) {
		if lex.commandMode && (p[n] == '\n' || (p[n] == '\r' && len(p) > n+1 && p[n+1] == '\n')) {
			// in command mode, we exclude the newline from the filler token, since it
			// will become a separate TokenDelimiter. Note that we know p doesn't *begin*
			// with a newline by virtue of Scan() handling that case before anything else.
			break
		}
		if needle == nil { // not currently in a comment of any type
			if ok, prefixLen, closerNeedle := isFiller(p[n:]); ok {
				n += prefixLen
				needle = closerNeedle
			} else {
				break // not starting a new comment, and not a space, so n is the end of the filler token
			}
		} else if i := bytes.Index(p[n:], needle); i >= 0 { // in a comment, and the closing needle is found in p
			n += i
			if !lex.commandMode || p[n] != '\n' {
				// include the needle in the filler token, UNLESS it's a newline in command
				// mode, in which case it gets emitted as a separate TokenDelimiter
				n += len(needle)
			}
			needle = nil // go back to "not in a comment" mode
		} else if lex.err == nil { // didn't find needle in p, but we can read and buffer more data
			n = len(p) + 1 - len(needle) // leave some n if multi-byte needle is split between p and next chunk
		} else {
			n = len(p) // there's no more data, so no need to worry about a split needle
		}

		// Unless we're at EOF (or other i/o error), ensure p has at least 6 bytes,
		// which is enough to hold "--" followed by a 4-byte rune
		if lex.err == nil && len(p)-n < 6 {
			p, n = lex.bufferAndPeek(p[0:n]), 0 // move data into external buffer and refill p
		}
	}
	if needle != nil && bytes.Equal(needle, needleCloseComment) {
		lex.err = &MalformedSQLError{str: "Comment starting with /* is never closed"}
	}
	lex.prevTokenFiller = true
	return lex.buildReturn(p[0:n], TokenFiller)
}

// scanString scans a quote-wrapped string wrapped in single-quotes, double-
// quotes, or backticks. The token will be considered TokenIdent in the case of
// backticks, otherwise TokenString.
func (lex *Lexer) scanString(p []byte) (data []byte, typ TokenType, err error) {
	c := p[0]
	if c == '`' {
		typ = TokenIdent
	} else {
		typ = TokenString
	}

	var done, skipNext, keepLast bool
	n := 1 // start right after the opening quote
	for !done && n < len(p) {
		if skipNext { // previous byte escaped this one
			skipNext = false
		} else if p[n] == '\\' && c != '`' { // backslash-escape only possible for c=='\'' or c=='"'
			skipNext = true
		} else if p[n] == c {
			if n >= len(p)-1 { // not enough data to see if c is escaped by doubling
				if lex.err == nil { // but there's more data to read
					keepLast = true // so keep this last byte at the head of the next refill of p
				} else { // if there isn't more data to read,
					done = true // then we're done since the last byte was the closing quote symbol
				}
			} else if p[n+1] == c { // c is definitely escaped by doubling
				skipNext = true
			} else { // c is not escaped
				done = true
			}
		}
		n++
		if n == len(p) && lex.err == nil && !done {
			// keep current last byte of p to become start of new p, to later see if
			// there's escape-by-doubling or not
			if keepLast {
				n--
				keepLast = false
			}
			p, n = lex.bufferAndPeek(p[0:n]), 0
		}
	}
	if !done && lex.err == io.EOF { // never found closing quote
		var noun string
		if c == '`' {
			noun = "Identifier"
		} else {
			noun = "String"
		}
		lex.err = &MalformedSQLError{str: noun + " is missing closing quote"}
	}
	return lex.buildReturn(p[0:n], typ)
}

type asciiSet [4]uint32 // only supports single-byte runes!

var wordSet = buildASCIISet("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$_")
var spaceSet = buildASCIISet(" \t\n\v\f\r")

func buildASCIISet(chars string) (as asciiSet) {
	for n := 0; n < len(chars); n++ {
		c := chars[n]
		if c >= utf8.RuneSelf {
			// TODO 错误处理
			// fs: assertion failed: ascii set cannot contain multibyte rune
			return as
		}
		as[c/32] |= 1 << (c % 32)
	}
	return as
}

// Returns true if b is within range 0-9a-zA-Z$_
func isWord(b byte) bool {
	return b < 128 && (wordSet[b/32]&(1<<(b%32))) != 0
}

// Returns true if b is an ascii whitespace character. This does handle ascii
// NBSP or NEL as whitespace despite them not being valid utf8.
func isSpace(b byte) bool {
	if b < 128 {
		return (spaceSet[b/32] & (1 << (b % 32))) != 0
	}
	return b == '\x85' || b == '\xA0'
}
