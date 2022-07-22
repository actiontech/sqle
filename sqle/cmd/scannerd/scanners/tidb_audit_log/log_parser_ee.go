//go:build enterprise
// +build enterprise

package tidb_audit_log

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type TiDBAuditLogParser interface {
	Parse(line string) (*TiDBAuditLog, error)
}

/*
日志格式大致如下:
[2022/04/20 12:32:56.450 +08:00] [INFO] [logger.go:76] [ID=1650429176916] [TIMESTAMP=2022/04/20 12:32:56.450 +08:00] [EVENT_CLASS=GENERAL] [EVENT_SUBCLASS=] [STATUS_CODE=0] [COST_TIME=0] [HOST=127.0.0.1] [CLIENT_IP=127.0.0.1] [USER=root] [DATABASES="[]"] [TABLES="[]"] [SQL_TEXT="select sum ( `k` ) from `sbtest1` where `id` between ? and ?"] [ROWS=0] [CONNECTION_ID=573] [CLIENT_PORT=39046] [PID=29277] [COMMAND=Execute] [SQL_STATEMENTS=Select]
*/
type TiDBAuditLog struct {
	Timestamp     time.Time
	Level         TiDBAuditLogLevel
	ID            string
	EventClass    TiDBAuditEventClass
	EventSubclass TiDBAuditEventSubclass
	StatusCode    int
	CostTime      int64
	Host          string
	ClientIp      string
	User          string
	Databases     []string
	Tables        []string
	SQLText       string
	Rows          int64
	ConnectionID  int
	ClientPort    int
	PID           int
	Command       string
	SQLStatements string
}

type TiDBAuditLogLevel string
type TiDBAuditEventClass string
type TiDBAuditEventSubclass string

type lexerParser struct {
	*lexmachine.Lexer
}

var stdLogStructor = &lexerParser{
	Lexer: lexer,
}

func GetLexerParser() *lexerParser {
	return stdLogStructor
}

func (l *lexerParser) Parse(line string) (*TiDBAuditLog, error) {
	log := &TiDBAuditLog{}
	scanner, err := l.Scanner([]byte(line))
	if err != nil {
		return nil, err
	}
	err = l.traverse(scanner, line, log)
	return log, err

}

func (l *lexerParser) traverse(scanner *lexmachine.Scanner, line string, result *TiDBAuditLog) error {
	resultSet, err := l.scanner(scanner)
	if err != nil {
		return err
	}
	return l.parse(resultSet, line, result)
}

// 遍历Token, 解析出各项属性的值
func (l *lexerParser) parse(tokens []*lexmachine.Token, line string, log *TiDBAuditLog) error {
	startIndex := 0
	var err error
	for i := 0; i < len(tokens); i++ {
		switch GetTokenName(tokens[i].Type) {
		case TokenStart:
			startIndex = i

		case TokenStop: // 开始token后下一个token就是结束token意味着两个这个属性中没有等号,比如[INFO]
			start := tokens[startIndex]
			stop := tokens[i]
			err = l.assembleNotKV(line[start.EndColumn:stop.StartColumn-1], log)
			if err != nil {
				return err
			}

		case TokenKV: // 开始token后下一个是等号意味着这是一个键值对, 比如 [ID=1650429176916]
			stopIndex := 0
			for j := i; j < len(tokens); j++ {
				if TokenStop == GetTokenName(tokens[j].Type) {
					stopIndex = j
					break
				}
			}
			if stopIndex == 0 {
				return errors.New("could not match end symbol")
			}
			err = l.assembleKV(tokens[startIndex:stopIndex+1], line, log)
			if err != nil {
				return err
			}
			i = stopIndex
		}
	}
	return nil
}

// 日志中有三个不是kv格式的信息, 其中时间戳可以从后面的kv中拿到,不在这里拿, 代码位置不做记录, 只有错误等级全部为字母,此方法主要就是为了获取错误等级
// 实际案例： [2022/04/20 12:32:56.450 +08:00] [INFO] [logger.go:76] , 只需要INFO
func (l *lexerParser) assembleNotKV(str string, log *TiDBAuditLog) error {
	for _, s := range str {
		if !unicode.IsLetter(s) {
			return nil
		}
	}
	log.Level = TiDBAuditLogLevel(str)
	return nil
}

type tidbValueKV struct {
	Key   string
	Value string
}

// kv格式的信息有两种情况,第一种value用引号包裹, 此时tokens长度应当为4, 另一种情况value没有引号包裹, 此时token长度应当为3
// 实际案例:
// [DATABASES="[]"] --> token1: 最左边的[ | token2: = | token3: "[]" | token4: 最右边的]
// [EVENT_CLASS=GENERAL] --> token1: 最左边的[ | token2: = | token3: 最右边的]
func (l *lexerParser) assembleKV(tokens []*lexmachine.Token, str string, log *TiDBAuditLog) error {
	var kv tidbValueKV
	var err error
	switch len(tokens) {
	case 3:
		kv = l.assembleNotQuote(tokens, str)
	case 4:
		kv, err = l.assembleQuote(tokens, str)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown log content")
	}

	return l.parserOneKVToLog(kv, log)
}

// 处理值被引号包裹的键值对
func (l *lexerParser) assembleQuote(tokens []*lexmachine.Token, str string) (tidbValueKV, error) {
	if tokens[2].Type != GetTokenType(TokenQuote) {
		return tidbValueKV{}, fmt.Errorf("wrong log format")
	}
	kv := tidbValueKV{
		Key:   str[tokens[0].EndColumn : tokens[1].StartColumn-1],
		Value: string(tokens[2].Lexeme),
	}
	return kv, nil
}

// 处理值没被引号包裹的键值对
func (l *lexerParser) assembleNotQuote(tokens []*lexmachine.Token, str string) tidbValueKV {
	kv := tidbValueKV{
		Key:   str[tokens[0].EndColumn : tokens[1].StartColumn-1],
		Value: str[tokens[1].EndColumn : tokens[2].EndColumn-1],
	}
	return kv
}

// 获取所有Token
func (l *lexerParser) scanner(scanner *lexmachine.Scanner) ([]*lexmachine.Token, error) {
	result := []*lexmachine.Token{}
	for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); is {
			// skip the unconsumed input
			scanner.TC = ui.FailTC
			continue
		} else if err != nil {
			return nil, err
		}

		result = append(result, tok.(*lexmachine.Token))
	}
	return result, nil
}

// 日志文件中存在的Key
const (
	KeyID            = "ID"
	KeyTimestamp     = "TIMESTAMP"
	KeyEventClass    = "EVENT_CLASS"
	KeyEventSubclass = "EVENT_SUBCLASS"
	KeyStatusCode    = "STATUS_CODE"
	KeyCostTime      = "COST_TIME"
	KeyHost          = "HOST"
	KeyClientIP      = "CLIENT_IP"
	KeyUser          = "USER"
	KeyDatabases     = "DATABASES"
	KeyTables        = "TABLES"
	KeySQLText       = "SQL_TEXT"
	KeyRows          = "ROWS"
	KeyConnectionID  = "CONNECTION_ID"
	KeyClientPort    = "CLIENT_PORT"
	KeyPID           = "PID"
	KeyCommand       = "COMMAND"
	KeySQLStatements = "SQL_STATEMENTS"
)

// 将解析出的值填充到结构体中
func (l *lexerParser) parserOneKVToLog(kv tidbValueKV, log *TiDBAuditLog) error {
	mp := map[string]interface{}{}
	switch kv.Key {
	case KeyStatusCode, KeyConnectionID, KeyClientPort, KeyPID:
		v, err := strconv.Atoi(kv.Value)
		if err != nil {
			return fmt.Errorf("the value of %v is malformed", kv.Key)
		}
		mp[kv.Key] = v
	case KeyID, KeyEventClass, KeyEventSubclass, KeyHost, KeyClientIP, KeyUser, KeySQLText, KeyCommand, KeySQLStatements:
		mp[kv.Key] = kv.Value
	case KeyTimestamp:
		//TODO 时间目前用不上, 所以忽略错误, 需要用时间时再判断
		mp[kv.Key], _ = time.Parse("2006/01/02 15:04:05", kv.Value)
	case KeyCostTime, KeyRows:
		//TODO 这两个目前用不上, 暂时忽略错误
		mp[kv.Key], _ = strconv.ParseInt(kv.Value, 10, 64)
	case KeyDatabases, KeyTables:
		v := strings.TrimPrefix(strings.TrimSuffix(kv.Value, "]"), "[")
		s := []string{}
		for _, value := range strings.Split(v, ",") {
			if value != "" {
				s = append(s, value)
			}
		}
		mp[kv.Key] = s
	default:
		// 读不出来的Key程序也用不上
		return nil
	}
	return gconv.Scan(mp, log)
}
