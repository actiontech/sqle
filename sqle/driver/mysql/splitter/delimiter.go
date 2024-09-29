package splitter

import (
	"errors"
	"github.com/pingcap/parser"
	"strings"
)

const (
	BackSlash              int    = '\\'
	BackSlashString        string = "\\"
	BlankSpace             string = " "
	DefaultDelimiterString string = ";"
	DelimiterCommand       string = "DELIMITER"
	DelimiterCommandSort   string = `\d`
)

type Delimiter struct {
	FirstTokenTypeOfDelimiter  int
	FirstTokenValueOfDelimiter string
	DelimiterStr               string
	line                       int
	startPos                   int
}

func NewDelimiter() *Delimiter {
	return &Delimiter{}
}

/*
	根据传入的SQL和位置，判断当前位置开始是否是一个缩写分隔符的语法

判断依据：从当前位置开始是否紧跟着一个\\d
不使用Lex的原因：
 1. \\d会被识别为三个token，即： \ \ d
 2. Lex可能会跳过空格和注释，因此这里使用字符串匹配
*/
func (d *Delimiter) isSortDelimiterCommand(sql string, index int) bool {
	return index+2 < len(sql) && sql[index:index+2] == "\\d"
}

// DELIMITER会被识别为identifier，因此这里仅需识别其值是否相等
func (d *Delimiter) isDelimiterCommand(token string) bool {
	return strings.ToUpper(token) == DelimiterCommand
}

// 该函数翻译自MySQL Client获取delimiter值的代码，参考：https://github.com/mysql/mysql-server/blob/824e2b4064053f7daf17d7f3f84b7a3ed92e5fb4/client/mysql.cc#L4866
func getDelimiter(line string) string {
	ptr := 0
	start := 0
	quoted := false
	qtype := byte(0)

	// 跳过开头的空格
	for ptr < len(line) && isSpace(line[ptr]) {
		ptr++
	}

	if ptr == len(line) {
		return ""
	}

	// 检查是否为引号字符串
	if line[ptr] == '\'' || line[ptr] == '"' || line[ptr] == '`' {
		qtype = line[ptr]
		quoted = true
		ptr++
	}

	start = ptr

	// 找到字符串结尾
	for ptr < len(line) {
		if !quoted && line[ptr] == '\\' && ptr+1 < len(line) { // 跳过转义字符
			ptr += 2
		} else if (!quoted && isSpace(line[ptr])) || (quoted && line[ptr] == qtype) {
			break
		} else {
			ptr++
		}
	}

	return line[start:ptr]
}

// 辅助函数,判断字符是否为空格
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

var ErrDelimiterCanNotExtractToken = errors.New("sorry, we cannot extract any token form the delimiter you provide, please change a delimiter")
var ErrDelimiterContainsBackslash = errors.New("DELIMITER cannot contain a backslash character")
var ErrDelimiterContainsBlankSpace = errors.New("DELIMITER should not contain blank space")
var ErrDelimiterMissing = errors.New("DELIMITER must be followed by a 'delimiter' character or string")
var ErrDelimiterReservedKeyword = errors.New("delimiter should not use a reserved keyword")

/*
该方法设置分隔符，对分隔符的内容有一定的限制：

 1. 不允许分隔符内部包含反斜杠
 2. 不允许分隔符为空字符串
 3. 不允许分隔符为mysql的保留字，因为这样会被scanner扫描为其他类型的token，从而绕过判断分隔符的逻辑

注：其中1和2与MySQL客户端对分隔符内容一致，错误内容参考MySQL客户端源码中的com_delimiter函数
https://github.com/mysql/mysql-server/blob/824e2b4064053f7daf17d7f3f84b7a3ed92e5fb4/client/mysql.cc#L4621
*/
func (d *Delimiter) setDelimiter(delimiter string) (err error) {
	if delimiter == "" {
		return ErrDelimiterMissing
	}
	if strings.Contains(delimiter, BackSlashString) {
		return ErrDelimiterContainsBackslash
	}
	if strings.Contains(delimiter, BlankSpace) {
		return ErrDelimiterContainsBlankSpace
	}
	if isReservedKeyWord(delimiter) {
		return ErrDelimiterReservedKeyword
	}
	token := parser.NewScanner(delimiter).NextToken()
	d.FirstTokenTypeOfDelimiter = token.TokenType()
	if d.FirstTokenTypeOfDelimiter == 0 {
		return ErrDelimiterCanNotExtractToken
	}
	d.FirstTokenValueOfDelimiter = token.Ident()
	d.DelimiterStr = delimiter
	return nil
}

func isReservedKeyWord(input string) bool {
	token := parser.NewScanner(input).NextToken()
	tokenType := token.TokenType()
	if len(token.Ident()) < len(input) {
		// 如果分隔符无法识别为一个token，则一定不是关键字
		return false
	}
	// 如果分隔符识别为一个关键字，但不知道是哪个关键字，则为identifier，此时就非保留字
	return tokenType != parser.Identifier && tokenType > parser.YyEOFCode && tokenType < parser.YyDefault
}

func (d *Delimiter) reset() error {
	d.line = 0
	d.startPos = 0
	return d.setDelimiter(DefaultDelimiterString)
}
