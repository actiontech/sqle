package parser

// 这个Reset方法在原有reset()方法上增加了将最后扫描的偏移量置零
func (s *Scanner) Reset(sql string) {
	s.reset(sql)
	s.lastScanOffset = 0
}

// 该方法返回Scanner最后扫描到的位置
func (s *Scanner) Offset() int {
	return s.lastScanOffset
}

// 该方法修改Scanner最后扫描到的位置
func (s *Scanner) SetCursor(offset int) {
	s.lastScanOffset = offset
}

const (
	Identifier int = identifier
	YyEOFCode  int = yyEOFCode
	YyDefault  int = yyDefault
	IfKwd      int = ifKwd
	CaseKwd    int = caseKwd
	Repeat     int = repeat
	Begin      int = begin
	End        int = end
	StringLit  int = stringLit
	Invalid    int = invalid
)

type TokenValue yySymType

type Token struct {
	tokenType  int
	tokenValue *yySymType
}

func (t Token) Ident() string {
	return t.tokenValue.ident
}

func (t Token) TokenType() int {
	return t.tokenType
}

func (s *Scanner) NextToken() *Token {
	tokenValue := &yySymType{}
	tokenType := s.Lex(tokenValue)
	return &Token{
		tokenType:  tokenType,
		tokenValue: tokenValue,
	}
}

func (s *Scanner) ScannedLines() int {
	return s.r.pos().Line - 1
}

func (s *Scanner) Text() string {
	return s.r.s
}

func (s *Scanner) HandleInvalid() {
	if s.lastScanOffset == s.r.p.Offset {
		s.r.inc()
	}
}
