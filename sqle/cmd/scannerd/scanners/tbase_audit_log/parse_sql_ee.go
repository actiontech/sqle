//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"fmt"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/pganalyze/collector/logs"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
	"github.com/sirupsen/logrus"
)

const (
	PgLogTotalColumns       = 26
	TxStartTimeIndex        = 0
	UserIndex               = 1
	SchemaIndex             = 2
	ClientHostWithPortIndex = 5
	SqlTypeIndex            = 10
	ConnectionTimeIndex     = 11
	SqlStatementIndex       = 16
	ParamsIndex             = 17
)

type TbaseLog struct {
	TxStartTime        time.Time
	SqlType            string
	User               string
	Schema             string
	ClientHostWithPort string
	ConnectionTime     time.Time
	Duration           float64
	SQLText            string
}

type TbaseParser struct {
	l *logrus.Entry
}

func NewTbaseParser() *TbaseParser {
	return &TbaseParser{
		l: logrus.WithField("scanner", "tbase-audit-log"),
	}
}

func (t *TbaseParser) parseRecord(record []string) (*TbaseLog, error) {
	record = t.removeCoord(record)
	if len(record) < PgLogTotalColumns {
		return nil, fmt.Errorf("record is not a standard log output, record: %v", record)
	}

	layout := "2006-01-02 15:04:05.999 CST"
	txStartTime, err := time.Parse(layout, record[TxStartTimeIndex])
	if err != nil {
		return nil, err
	}
	connectionTime, err := time.Parse(layout, record[ConnectionTimeIndex])
	if err != nil {
		return nil, err
	}

	duration, sql, err := t.parseSQLStatement(record[SqlStatementIndex], record[ParamsIndex])
	if err != nil {
		return nil, err
	}
	if duration == 0 || sql == "" {
		return nil, nil
	}

	return &TbaseLog{
		TxStartTime:        txStartTime,
		SqlType:            record[SqlTypeIndex],
		User:               record[UserIndex],
		Schema:             record[SchemaIndex],
		ClientHostWithPort: record[ClientHostWithPortIndex],
		ConnectionTime:     connectionTime,
		Duration:           duration,
		SQLText:            sql,
	}, nil
}

// 客户的pg_log日志文件中，存在这种格式的值 coord(xxx, xxx)
// 但是该值没有被双引号包裹起来，csv解析器会将该值解析成 "coord(xxx" 和 "xxx)"这两个值
// 为了不影响解析，使用空格替代coord(xxx, xxx)
func (t *TbaseParser) removeCoord(record []string) []string {
	newRecord := []string{}
	index := 0
	for index < len(record) {
		value := record[index]
		if strings.HasPrefix(value, "coord(") {
			index += 2
			newRecord = append(newRecord, "")
			continue
		} else {
			newRecord = append(newRecord, value)
			index += 1
		}
	}
	return newRecord
}

func (t *TbaseParser) parseSQLStatement(statement, params string) (duration float64, sql string, err error) {
	if statement == "" {
		return
	}
	logLines := []state.LogLine{
		{
			Content: statement,
		},
	}
	if params != "" {
		logLines = append(logLines, state.LogLine{
			Content:  params,
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		})
	}

	_, results := logs.AnalyzeLogLines(logLines)
	if len(results) == 0 {
		return
	}
	result := results[0]
	duration = result.RuntimeMs
	sql = replaceParams(result.Query, result.Parameters)
	return
}

// 临时方案：
// 使用字符串替换的方式将参数拼接到sql语句中
// 参数以char类型进行拼接
func replaceParams(sql string, params []null.String) string {
	for i, param := range params {
		sql = strings.Replace(sql, fmt.Sprintf("$%v", i+1), fmt.Sprintf("'%s'", param.String), 1)
	}
	return sql
}
