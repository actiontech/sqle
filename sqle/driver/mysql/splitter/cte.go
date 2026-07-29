package splitter

import (
	"strings"
	"unicode"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

// tryParseMySQLCTE 在 pingcap parser 无法识别 CTE 时，保守剥离
// WITH [RECURSIVE] name AS (...) [, ...] 后解析外层语句。
// 仅当确认 CTE 形态且外层解析成功时返回 (stmt, true)；否则 (nil, false)。
// 禁止匹配裸首词 WITH（如 WITH GRANT OPTION）。
func tryParseMySQLCTE(p *parser.Parser, originSQL string) (ast.StmtNode, bool) {
	outer, ok := extractOuterSQLAfterCTE(originSQL)
	if !ok {
		return nil, false
	}
	stmt, err := p.ParseOneStmt(outer, "", "")
	if err != nil {
		return nil, false
	}
	return stmt, true
}

// extractOuterSQLAfterCTE 识别 MySQL CTE 前缀并返回外层 SQL。
// 形态：WITH [RECURSIVE] cte_name [(cols)] AS (subquery) [, ...] <outer>
func extractOuterSQLAfterCTE(sql string) (string, bool) {
	s := strings.TrimSpace(sql)
	if s == "" {
		return "", false
	}
	i := skipSQLTrivia(s, 0)
	if i >= len(s) {
		return "", false
	}
	if !hasKeywordAt(s, i, "WITH") {
		return "", false
	}
	i += len("WITH")
	i = skipSQLTrivia(s, i)
	if hasKeywordAt(s, i, "RECURSIVE") {
		i += len("RECURSIVE")
		i = skipSQLTrivia(s, i)
	}

	for {
		nameEnd, ok := readSQLIdent(s, i)
		if !ok {
			return "", false
		}
		i = skipSQLTrivia(s, nameEnd)

		// optional column list
		if i < len(s) && s[i] == '(' {
			end, ok := skipBalancedParens(s, i)
			if !ok {
				return "", false
			}
			i = skipSQLTrivia(s, end)
		}

		if !hasKeywordAt(s, i, "AS") {
			return "", false
		}
		i += len("AS")
		i = skipSQLTrivia(s, i)
		if i >= len(s) || s[i] != '(' {
			return "", false
		}
		end, ok := skipBalancedParens(s, i)
		if !ok {
			return "", false
		}
		i = skipSQLTrivia(s, end)

		if i < len(s) && s[i] == ',' {
			i++
			i = skipSQLTrivia(s, i)
			continue
		}
		break
	}

	if i >= len(s) {
		return "", false
	}
	outer := strings.TrimSpace(s[i:])
	if outer == "" {
		return "", false
	}
	return outer, true
}

func hasKeywordAt(s string, i int, keyword string) bool {
	if i+len(keyword) > len(s) {
		return false
	}
	if !strings.EqualFold(s[i:i+len(keyword)], keyword) {
		return false
	}
	end := i + len(keyword)
	if end < len(s) {
		r := rune(s[end])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return false
		}
	}
	return true
}

func skipSQLTrivia(s string, i int) int {
	for i < len(s) {
		switch {
		case unicode.IsSpace(rune(s[i])):
			i++
		case i+1 < len(s) && s[i] == '-' && s[i+1] == '-':
			i += 2
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case s[i] == '#':
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case i+1 < len(s) && s[i] == '/' && s[i+1] == '*':
			i += 2
			for i+1 < len(s) && !(s[i] == '*' && s[i+1] == '/') {
				i++
			}
			if i+1 < len(s) {
				i += 2
			}
		default:
			return i
		}
	}
	return i
}

func readSQLIdent(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	switch s[i] {
	case '`':
		j := i + 1
		for j < len(s) {
			if s[j] == '`' {
				if j+1 < len(s) && s[j+1] == '`' {
					j += 2
					continue
				}
				return j + 1, j > i+1
			}
			j++
		}
		return i, false
	case '"', '\'':
		// MySQL ansi_quotes / rare; treat quoted as ident if double-quoted
		quote := s[i]
		if quote == '\'' {
			return i, false
		}
		j := i + 1
		for j < len(s) {
			if s[j] == '\\' && j+1 < len(s) {
				j += 2
				continue
			}
			if s[j] == quote {
				return j + 1, j > i+1
			}
			j++
		}
		return i, false
	default:
		if !isIdentStart(rune(s[i])) {
			return i, false
		}
		j := i + 1
		for j < len(s) && isIdentPart(rune(s[j])) {
			j++
		}
		return j, true
	}
}

func isIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == '$'
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}

// skipBalancedParens starts at '(' and returns index after matching ')'.
func skipBalancedParens(s string, i int) (int, bool) {
	if i >= len(s) || s[i] != '(' {
		return i, false
	}
	depth := 0
	for i < len(s) {
		c := s[i]
		switch c {
		case '(':
			depth++
			i++
		case ')':
			depth--
			i++
			if depth == 0 {
				return i, true
			}
		case '\'', '"', '`':
			end, ok := skipQuoted(s, i)
			if !ok {
				return i, false
			}
			i = end
		case '-':
			if i+1 < len(s) && s[i+1] == '-' {
				i += 2
				for i < len(s) && s[i] != '\n' {
					i++
				}
			} else {
				i++
			}
		case '#':
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case '/':
			if i+1 < len(s) && s[i+1] == '*' {
				i += 2
				for i+1 < len(s) && !(s[i] == '*' && s[i+1] == '/') {
					i++
				}
				if i+1 < len(s) {
					i += 2
				}
			} else {
				i++
			}
		default:
			i++
		}
	}
	return i, false
}

func skipQuoted(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	quote := s[i]
	j := i + 1
	for j < len(s) {
		if quote == '`' {
			if s[j] == '`' {
				if j+1 < len(s) && s[j+1] == '`' {
					j += 2
					continue
				}
				return j + 1, true
			}
			j++
			continue
		}
		if s[j] == '\\' && j+1 < len(s) {
			j += 2
			continue
		}
		if s[j] == quote {
			// MySQL '' escape inside strings
			if quote == '\'' && j+1 < len(s) && s[j+1] == '\'' {
				j += 2
				continue
			}
			return j + 1, true
		}
		j++
	}
	return i, false
}
