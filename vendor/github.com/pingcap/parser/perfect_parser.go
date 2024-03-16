package parser

import (
	"bytes"

	"github.com/pingcap/parser/ast"
)

// PerfectParse parses a query string to raw ast.StmtNode. support parses query string
// who contains unparsed SQL, the unparsed SQL will be parses to ast.UnparsedStmt.
func (parser *Parser) PerfectParse(sql, charset, collation string) (stmt []ast.StmtNode, warns []error, err error) {
	_, warns, err = parser.Parse(sql, charset, collation)
	stmts := parser.result
	parser.updateStartLineWithOffset(stmts)
	if err == nil {
		return stmts, warns, nil
	}
	// if err is not nil, the query string must be contains unparsed sql.

	if len(stmts) > 0 {
		for _, stmt := range stmts {
			ast.SetFlag(stmt)
		}
		stmt = append(stmt, stmts...)
	}
	parser.startLineOffset = parser.lexer.r.pos().Line - 1
	// The origin SQL text(input args `sql`) consists of many SQL segments,
	// each SQL segments is a complete SQL and be parsed into `ast.StmtNode`.
	//
	//     good SQL segment       bad SQL segment
	// |---------------------|---------------------|---------------------|---------------------|    origin SQL text
	//			     		 ^				^
	//		            stmtStartPos   lastScanOffset
	//										|------|---------------------|---------------------|    remaining SQL text
	//
	//                       |<   unparsed stmt   >|<          continue to parse it           >|

	start := parser.lexer.stmtStartPos
	cur := parser.lexer.lastScanOffset

	remainingSql := sql[cur:]
	l := NewScanner(remainingSql)
	var v yySymType
	var endOffset int
	var scanEnd = 0
	var defaultDelimiter int = ';'
	delimiter := defaultDelimiter
ScanLoop:
	for {
		result := l.Lex(&v)
		switch result {
		case scanEnd:
			endOffset = l.lastScanOffset - 1
			break ScanLoop
		case delimiter:
			endOffset = l.lastScanOffset
			break ScanLoop
		case begin:
			// ref: https://dev.mysql.com/doc/refman/8.0/en/begin-end.html
			// ref: https://dev.mysql.com/doc/refman/8.0/en/stored-programs-defining.html
			// Support match:
			// BEGIN
			// ...
			// END;
			//
			delimiter = scanEnd
		case end:
			// match `end;`
			var ny yySymType
			next := l.Lex(&ny)
			if next == defaultDelimiter {
				delimiter = defaultDelimiter
				endOffset = l.lastScanOffset
				break ScanLoop
			}
		case invalid:
			// `Lex`内`scan`在进行token遍历时，当有特殊字符时返回invalid，此时未调用`inc`进行滑动，导致每次遍历同一个pos点位触发死循环。有多种情况会返回invalid。
			// 对于解析器本身没影响，因为 token 提取失败就退出了，但是我们需要继续遍历。
			if l.lastScanOffset == l.r.p.Offset {
				l.r.inc()
			}
		}
	}
	unparsedStmtBuf := bytes.Buffer{}
	unparsedStmtBuf.WriteString(sql[start:cur])
	unparsedStmtBuf.WriteString(remainingSql[:endOffset+1])

	unparsedSql := unparsedStmtBuf.String()
	if len(unparsedSql) > 0 {
		un := &ast.UnparsedStmt{}
		un.SetStartLine(parser.startLineOffset + 1)
		un.SetText(unparsedSql)
		stmt = append(stmt, un)
	}

	if len(remainingSql) > endOffset {
		cStmt, cWarn, cErr := parser.PerfectParse(remainingSql[endOffset+1:], charset, collation)
		warns = append(warns, cWarn...)
		if len(cStmt) > 0 {
			stmt = append(stmt, cStmt...)
		}
		if cErr == nil {
			return stmt, warns, cErr
		}
	}
	return stmt, warns, nil
}

func (parser *Parser) updateStartLineWithOffset(stmts []ast.StmtNode) {
	for i := range stmts {
		stmts[i].SetStartLine(stmts[i].StartLine() + parser.startLineOffset)
	}
}
