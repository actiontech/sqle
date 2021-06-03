package parser

import (
	"bytes"
	"github.com/pingcap/parser/ast"
)

// PerfectParse parses a query string to raw ast.StmtNode. support parses query string
// who contains unparsed SQL, the unparsed SQL will be parses to ast.UnparsedStmt.
func (parser *Parser) PerfectParse(sql, charset, collation string) (stmt []ast.StmtNode, warns []error, err error) {
	_, warns, err = parser.Parse(sql, charset, collation)
	if err == nil {
		return parser.result, warns, nil
	}
	// if err is not nil, the query string must be contains unparsed sql.

	if len(parser.result) > 0 {
		for _, stmt := range parser.result {
			ast.SetFlag(stmt)
		}
		stmt = append(stmt, parser.result...)
	}

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
	for {
		result := l.Lex(&v)
		if result == 0 {
			endOffset = l.lastScanOffset - 1
			break
		}
		if result == ';' {
			endOffset = l.lastScanOffset
			break
		}
	}
	unparsedStmtBuf := bytes.Buffer{}
	unparsedStmtBuf.WriteString(sql[start:cur])
	unparsedStmtBuf.WriteString(remainingSql[:endOffset+1])

	unparsedSql := unparsedStmtBuf.String()
	if len(unparsedSql) > 0 {
		un := &ast.UnparsedStmt{}
		un.SetText(unparsedSql)
		stmt = append(stmt, un)
	}

	if len(remainingSql) <= endOffset {
		return stmt, warns, nil
	}

	cStmt, cWarn, cErr := parser.PerfectParse(remainingSql[endOffset+1:], charset, collation)
	warns = append(warns, cWarn...)
	if len(cStmt) > 0 {
		stmt = append(stmt, cStmt...)
	}
	if cErr == nil {
		return stmt, warns, cErr
	}
	return stmt, warns, nil
}
