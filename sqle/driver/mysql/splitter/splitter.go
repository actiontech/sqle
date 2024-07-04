package splitter

import (
	"bytes"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"strings"
)

type splitter struct {
	parser    *parser.Parser
	delimiter *Delimiter
	scanner   *parser.Scanner
}

func NewSplitter() *splitter {
	return &splitter{
		parser:    parser.New(),
		delimiter: NewDelimiter(),
		scanner:   parser.NewScanner(""),
	}
}

func (s *splitter) ParseSqlText(sqlText string) ([]ast.StmtNode, error) {
	s.delimiter.reset()
	results, err := s.splitSqlText(sqlText)
	if err != nil {
		return nil, err
	}
	return s.processToExecutableNodes(results), nil
}

func (s *splitter) processToExecutableNodes(results []*sqlWithLineNumber) []ast.StmtNode {
	s.delimiter.reset()

	var executableNodes []ast.StmtNode
	for _, result := range results {
		if matched, _ := s.matchAndSetCustomDelimiter(result.sql); matched {
			continue
		}
		if strings.HasSuffix(result.sql, s.delimiter.DelimiterStr) {
			trimmedSQL := strings.TrimSuffix(result.sql, s.delimiter.DelimiterStr)
			if trimmedSQL == "" {
				continue
			}
			result.sql = trimmedSQL + ";"
		}
		s.parser.Parse(result.sql, "", "")
		if len(s.parser.Result()) == 1 {
			// 若结果集长度为1，则为单条且可解析的SQL
			stmt := s.parser.Result()[0]
			stmt.SetStartLine(result.line)
			executableNodes = append(executableNodes, stmt)
		} else {
			// 若结果集长度大于1，则为多条合并的SQL
			// 若结果集长度为0，则不可解析的SQL
			unParsedStmt := &ast.UnparsedStmt{}
			unParsedStmt.SetStartLine(result.line)
			unParsedStmt.SetText(result.sql)
			executableNodes = append(executableNodes, unParsedStmt)
		}
	}
	return executableNodes
}

type sqlWithLineNumber struct {
	sql  string
	line int
}

func (s *splitter) splitSqlText(sqlText string) (results []*sqlWithLineNumber, err error) {
	result, err := s.getNextSql(sqlText)
	if err != nil {
		return nil, err
	}
	if result != nil {
		results = append(results, result)
	}
	// 递归切分剩余SQL
	if s.scanner.Offset() < len(sqlText) {
		subResults, _ := s.splitSqlText(sqlText[s.scanner.Offset():])
		results = append(results, subResults...)
	}
	return results, nil
}

func (s *splitter) getNextSql(sqlText string) (*sqlWithLineNumber, error) {
	matcheDelimiterCommand, err := s.matchAndSetCustomDelimiter(sqlText)
	if err != nil {
		return nil, err
	}
	// 若匹配到自定义分隔符语法，则输出结果，否则匹配分隔符，输出结果
	if matcheDelimiterCommand || s.matcheSql(sqlText) {
		buff := bytes.Buffer{}
		buff.WriteString(sqlText[:s.scanner.Offset()])
		lineBeforeStart := strings.Count(sqlText[:s.delimiter.startPos], "\n")
		result := &sqlWithLineNumber{
			sql:  strings.TrimSpace(buff.String()),
			line: s.delimiter.line + lineBeforeStart + 1,
		}
		s.delimiter.line += s.scanner.ScannedLines() // pos().Line-1表示的是该SQL中有多少换行
		return result, nil
	}
	restOfSql := strings.TrimSpace(sqlText)
	if restOfSql == "" {
		return nil, nil
	}
	return &sqlWithLineNumber{
		sql:  restOfSql,
		line: s.delimiter.line + strings.Count(sqlText[:s.delimiter.startPos], "\n") + 1,
	}, nil
}

func (s *splitter) matcheSql(sql string) bool {
	s.scanner.Reset(sql)
	token := &parser.Token{}
	var isFirstToken bool = true

	for s.scanner.Offset() < len(sql) {
		token = s.scanner.NextToken()
		if isFirstToken {
			s.delimiter.startPos = s.scanner.Offset()
			isFirstToken = false
		}
		token = s.skipBeginEndBlock(token)
		if s.isTokenMatchDelimiter(token) {
			return true
		}
	}
	return false
}

func (s *splitter) skipBeginEndBlock(token *parser.Token) *parser.Token {
	var blockStack []Block
	if token.TokenType() == parser.Begin {
		blockStack = append(blockStack, BeginEndBlock{})
	}
	for len(blockStack) > 0 {
		token = s.scanner.NextToken()
		for _, block := range allBlocks {
			if block.MatchBegin(token) {
				blockStack = append(blockStack, block)
				break
			}
		}
		// 如果匹配到END，则需要判断END后的token是否匹配当前的Block
		if token.TokenType() == parser.End {
			currentBlock := blockStack[len(blockStack)-1]
			token = s.scanner.NextToken()
			if currentBlock.MatchEnd(token) {
				blockStack = blockStack[:len(blockStack)-1]
			}
		}
		if len(blockStack) == 0 {
			break
		}
	}
	return token
}

// ref:https://dev.mysql.com/doc/refman/8.4/en/flow-control-statements.html
func (s *splitter) isTokenMatchDelimiter(token *parser.Token) bool {
	switch token.TokenType() {
	case s.delimiter.FirstTokenTypeOfDelimiter:
		/*
			在mysql client的语法中需要跳过注释以及分隔符处于引号中的情况，由于scanner.Lex会自动跳过注释，因此，仅需要判断分隔符处于引号中的情况。对于该方法，以分隔符的第一个token作为特征仅需匹配，可能会匹配到由引号括起的情况，存在stringLit和identifier两种token需要进一步判断：
				1. 当匹配到identifier时，identifier有可能由反引号括起:
					1. 若identifier没有反引号括起，则不需要判断是否跳过
					2. 若identifier被反引号括起，匹配的字符串会带上反引号，能在匹配字符串时能够检查出是否需要跳过
				2. 当匹配到stringLit时，stringLit一定是由单引号或双引号括起:
					1. 当分隔符第一个token值与stringLit的token值不等，那么一定不是分隔符，则跳过
					2. 当分隔符第一个token值与stringLit的token值相等， 如："'abc'd" '"abc"d'会因为字符串不匹配而跳过
		*/
		// 1. 当分隔符第一个token值与stringLit的token值不等，那么一定不是分隔符，则跳过
		if token.TokenType() == parser.StringLit && token.Ident() != s.delimiter.FirstTokenValueOfDelimiter {
			return false
		}
		// 2. 定位特征的第一个字符所处的位置
		indexIntoken := strings.Index(token.Ident(), s.delimiter.FirstTokenValueOfDelimiter)
		if indexIntoken == -1 {
			return false
		}
		// 3. 字符串匹配
		begin := s.scanner.Offset() + indexIntoken
		end := begin + len(s.delimiter.DelimiterStr)
		if begin < 0 || end > len(s.scanner.Text()) {
			return false
		}
		expected := s.scanner.Text()[begin:end]
		if expected != s.delimiter.DelimiterStr {
			return false
		}
		s.scanner.SetCursor(end)
		return true

	case parser.Invalid:
		s.scanner.HandleInvalid()
	}
	return false
}

/*
该方法检测sql文本开头是否是自定义分隔符语法，若是匹配并更新分隔符:

 1. 分隔符语法满足：delimiter str 或者 \d str
 2. 参考链接：https://dev.mysql.com/doc/refman/5.7/en/mysql-commands.html
*/
func (s *splitter) matchAndSetCustomDelimiter(sql string) (bool, error) {
	// 重置扫描器
	s.scanner.Reset(sql)
	var sqlAfterDelimiter string
	token := s.scanner.NextToken()
	switch token.TokenType() {
	case BackSlash:
		if s.delimiter.isSortDelimiterCommand(sql, s.scanner.Offset()) {
			sqlAfterDelimiter = sql[s.scanner.Offset()+2:] // \d的长度是2字节
			s.delimiter.startPos = s.scanner.Offset()
			s.scanner.SetCursor(s.scanner.Offset() + 2)
		}
	case parser.Identifier:
		if s.delimiter.isDelimiterCommand(token.Ident()) {
			sqlAfterDelimiter = sql[s.scanner.Offset()+9:] //DELIMITER的长度是9字节
			s.delimiter.startPos = s.scanner.Offset()
			s.scanner.SetCursor(s.scanner.Offset() + 9)
		}
	default:
		return false, nil
	}
	// 处理自定义分隔符
	if sqlAfterDelimiter != "" {
		restOfThisLine := strings.Index(sqlAfterDelimiter, "\n")
		if restOfThisLine == -1 {
			restOfThisLine = len(sqlAfterDelimiter)
		}
		newDelimiter := getDelimiter(sqlAfterDelimiter[:restOfThisLine])
		if err := s.delimiter.setDelimiter(newDelimiter); err != nil {
			return false, err
		}
		// 若识别到分隔符，则这一整行都为定义分隔符的sql，
		// 例如 delimiter ;; xx 其中;;为分隔符，而xx不产生任何影响，但属于这条语句
		s.scanner.SetCursor(s.scanner.Offset() + restOfThisLine)
		return true, nil
	}
	return false, nil
}
