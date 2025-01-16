//go:build enterprise
// +build enterprise

package slowquery

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/percona/go-mysql/log"
	"github.com/percona/pmm-agent/agents/mysql/slowlog/parser"

	"github.com/sirupsen/logrus"
)

type SlowQuery struct {
	l *logrus.Entry
	c *scanner.Client

	logFilePath string

	sqlCh chan scanners.SQL

	auditPlanID string

	includeUsers   map[string]struct{} // 用户白名单: 仅采集这些用户的 SQL;
	excludeUsers   map[string]struct{} // 用户黑名单: 不采集这些用户的 SQL;
	includeSchemas map[string]struct{} // Schema 白名单: 仅采集这些库的 SQL;
	excludeSchemas map[string]struct{} // Schema 黑名单: 不采集这些库的 SQL;
}

type Params struct {
	LogFilePath    string
	AuditPlanID    string
	IncludeUsers   string
	ExcludeUsers   string
	IncludeSchemas string
	ExcludeSchemas string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SlowQuery, error) {
	slowQuery := &SlowQuery{
		l:              l,
		c:              c,
		logFilePath:    params.LogFilePath,
		auditPlanID:    params.AuditPlanID,
		sqlCh:          make(chan scanners.SQL, 10),
		includeUsers:   map[string]struct{}{},
		excludeUsers:   map[string]struct{}{},
		includeSchemas: map[string]struct{}{},
		excludeSchemas: map[string]struct{}{},
	}
	if params.IncludeUsers != "" {
		for _, user := range strings.Split(params.IncludeUsers, ",") {
			slowQuery.includeUsers[strings.TrimSpace(user)] = struct{}{}
		}
	}
	if params.ExcludeUsers != "" {
		for _, user := range strings.Split(params.ExcludeUsers, ",") {
			slowQuery.excludeUsers[strings.TrimSpace(user)] = struct{}{}
		}
	}
	if params.IncludeSchemas != "" {
		for _, schema := range strings.Split(params.IncludeSchemas, ",") {
			slowQuery.includeSchemas[strings.TrimSpace(schema)] = struct{}{}
		}
	}
	if params.ExcludeSchemas != "" {
		for _, schema := range strings.Split(params.ExcludeSchemas, ",") {
			slowQuery.excludeSchemas[strings.TrimSpace(schema)] = struct{}{}
		}
	}
	return slowQuery, nil
}

// parser从sql string中解析出慢日志的时间度量值
const (
	QueryTime   string = "Query_time"    //执行时间
	LockTime    string = "Lock_time"     //加锁时间
	RowExamined string = "Rows_examined" //扫描行数
)

func (s *SlowQuery) Run(ctx context.Context) error {
	reader, err := scanners.NewContinuousFileReader(s.logFilePath, logrus.WithField("filename", s.logFilePath), true)
	if err != nil {
		return err
	}

	p := parser.NewSlowLogParser(reader, log.Options{FilterAdminCommand: map[string]bool{
		"Binlog Dump":      true,
		"Binlog Dump GTID": true,
	}})
	go p.Run()

	events := make(chan *log.Event, 1000)
	go func() {
		for {
			event := p.Parse()
			if event != nil {
				events <- event
				continue
			}

			if err := p.Err(); err != io.EOF {
				s.l.Warnf("Parser error: %v.", err)
			}
			close(events)
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			err = reader.Close()
			s.l.Infof("context done: %s. reader closed: %v.", ctx.Err(), err)
			close(s.sqlCh)
			return nil

		case e, ok := <-events:
			if !ok {
				return nil
			}
			s.l.Debugf("parsed slowlog event: %+v.", e)

			if strings.TrimSpace(e.Query) == "" {
				continue
			}

			if len(s.includeUsers) > 0 {
				if _, ok := s.includeUsers[e.User]; !ok {
					continue
				}
			}
			if len(s.excludeUsers) > 0 {
				if _, ok := s.excludeUsers[e.User]; ok {
					continue
				}
			}
			if len(s.includeSchemas) > 0 {
				if _, ok := s.includeSchemas[e.Db]; !ok {
					continue
				}
			}
			if len(s.excludeSchemas) > 0 {
				if _, ok := s.excludeSchemas[e.Db]; ok {
					continue
				}
			}
			s.sqlCh <- scanners.SQL{
				RawText:     e.Query,
				Counter:     1,
				QueryTime:   e.TimeMetrics[QueryTime],
				RowExamined: float64(e.NumberMetrics[RowExamined]),
				QueryAt:     e.Ts,
				DBUser:      e.User,
				Schema:      e.Db,
				Endpoint:    e.Host,
			}
		}
	}
}

func (sq *SlowQuery) SQLs() <-chan scanners.SQL {
	return sq.sqlCh
}

// 使用fingerprint和schema区分相同fingerprint但不同schema的sql
type sqlIdentity struct {
	Fingerprint string
	Schema      string
}

func (s sqlIdentity) jsonString() string {
	keyByte, _ := json.Marshal(s)
	return string(keyByte)
}

func (sq *SlowQuery) Upload(ctx context.Context, sqls []scanners.SQL, errorMessage string) error {
	var sqlListReq []*scanner.AuditPlanSQLReq
	now := time.Now()
	auditPlanSqlMap := make(map[string]*scanner.AuditPlanSQLReq, 0)

	var nodes []driverV2.Node
	var sqlIdentity sqlIdentity
	var processedRawText, processedFingerPrint string
	var err error
	for _, sql := range sqls {
		nodes, err = common.Parse(ctx, sql.RawText)
		if err != nil {
			// 解析出错 节点为空 保留SQL原语句
			processedRawText = sql.RawText
			processedFingerPrint = sql.RawText
		} else if len(nodes) > 0 {
			// 解析正常 节点非空 保留解析出的SQL
			// 若节点数量大于1 则第一个节点对应解析出的SQL 后续节点为;后的字符串 不保留
			processedFingerPrint = nodes[0].Fingerprint
			processedRawText = strings.Trim(nodes[0].Text, ";")
		} else {
			// 解析正常 节点为空 不保留该语句
			continue
		}

		// 使用指纹和Schema的JSON String作为SQL的唯一标识符
		sqlIdentity.Fingerprint = processedFingerPrint
		sqlIdentity.Schema = sql.Schema
		key := sqlIdentity.jsonString()

		if sqlReq, ok := auditPlanSqlMap[key]; ok {
			atoi, err := strconv.Atoi(sqlReq.Counter)
			if err != nil {
				return err
			}

			tMax := utils.MaxFloat64(*sqlReq.QueryTimeMax, sql.QueryTime)
			tAvg := utils.IncrementalAverageFloat64(*sqlReq.QueryTimeAvg, sql.QueryTime, atoi, 1)
			RowExaminedAvg := utils.IncrementalAverageFloat64(*sqlReq.RowExaminedAvg, sql.RowExamined, atoi, 1)
			sqlReq.QueryTimeMax = &tMax
			sqlReq.QueryTimeAvg = &tAvg
			sqlReq.RowExaminedAvg = &RowExaminedAvg
			sqlReq.Counter = strconv.Itoa(atoi + 1)
			sqlReq.LastReceiveText = processedRawText
			sqlReq.LastReceiveTimestamp = now.Format(time.RFC3339)
			if sql.Endpoint != "" {
				sqlReq.Endpoints = append(sqlReq.Endpoints, sql.Endpoint)
			}
		} else {
			sqlReq := &scanner.AuditPlanSQLReq{
				Fingerprint:          processedFingerPrint,
				LastReceiveText:      processedRawText,
				LastReceiveTimestamp: now.Format(time.RFC3339),
				Counter:              "1",
				QueryTimeAvg:         &sql.QueryTime,
				QueryTimeMax:         &sql.QueryTime,
				RowExaminedAvg:       &sql.RowExamined,
				FirstQueryAt:         sql.QueryAt,
				DBUser:               sql.DBUser,
				Schema:               sql.Schema,
			}
			if sql.Endpoint != "" {
				sqlReq.Endpoints = append(sqlReq.Endpoints, sql.Endpoint)
			}

			auditPlanSqlMap[key] = sqlReq
			sqlListReq = append(sqlListReq, sqlReq)
		}
	}

	return sq.c.UploadReq(scanner.UploadSQL, sq.auditPlanID, errorMessage, sqlListReq)
}
