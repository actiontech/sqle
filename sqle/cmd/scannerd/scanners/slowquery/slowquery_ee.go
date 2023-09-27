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
	"github.com/percona/go-mysql/query"
	"github.com/percona/pmm-agent/agents/mysql/slowlog/parser"
	"github.com/sirupsen/logrus"
)

type SlowQuery struct {
	l *logrus.Entry
	c *scanner.Client

	logFilePath string

	sqlCh chan scanners.SQL

	apName string

	includeUsers   map[string]struct{} // 用户白名单: 仅采集这些用户的 SQL;
	excludeUsers   map[string]struct{} // 用户黑名单: 不采集这些用户的 SQL;
	includeSchemas map[string]struct{} // Schema 白名单: 仅采集这些库的 SQL;
	excludeSchemas map[string]struct{} // Schema 黑名单: 不采集这些库的 SQL;
}

type Params struct {
	LogFilePath    string
	APName         string
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
		apName:         params.APName,
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
	QueryTime string = "Query_time" //执行时间
	LockTime  string = "Lock_time"  //加锁时间
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
				Fingerprint: query.Fingerprint(e.Query),
				RawText:     e.Query,
				Counter:     1,
				QueryTime:   e.TimeMetrics[QueryTime],
				QueryAt:     e.Ts,
				DBUser:      e.User,
				Schema:      e.Db,
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

func (sq *SlowQuery) Upload(ctx context.Context, sqls []scanners.SQL) error {
	var sqlListReq []*scanner.AuditPlanSQLReq
	now := time.Now()
	auditPlanSqlMap := make(map[string]*scanner.AuditPlanSQLReq, 0)

	var nodes []driverV2.Node
	var sqlIdentity sqlIdentity
	var processedRawText, processedFingerPrint string

	for _, sql := range sqls {
		// 在聚合之前对sql.Fingerprint和sql.RawText进行预处理，去除SQL中的脏数据
		processedFingerPrint = sql.Fingerprint
		processedRawText = sql.RawText

		// 当解析结果中有多个node的时候，说明该SQL带有冗余的脏数据
		// 这里把第一个node作为fingerprint，并且使用node.Text作为最后匹配到该指纹的SQL
		// 为了和其他SQL一致，这里还需要去除node.Text字符串结尾的';'
		// 若解析失败则nodes为nil len(nil)=0 此时Fingerprint和RawText保持解析前的状态
		nodes, _ = common.Parse(ctx, sql.RawText)
		if len(nodes) > 1 {
			processedFingerPrint = nodes[0].Fingerprint
			processedRawText = strings.Trim(nodes[0].Text, ";")
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
			sqlReq.QueryTimeMax = &tMax
			sqlReq.QueryTimeAvg = &tAvg
			sqlReq.Counter = strconv.Itoa(atoi + 1)
			sqlReq.LastReceiveText = sql.RawText
			sqlReq.LastReceiveTimestamp = now.Format(time.RFC3339)
		} else {
			sqlReq := &scanner.AuditPlanSQLReq{
				Fingerprint:          processedFingerPrint,
				LastReceiveText:      processedRawText,
				LastReceiveTimestamp: now.Format(time.RFC3339),
				Counter:              "1",
				QueryTimeAvg:         &sql.QueryTime,
				QueryTimeMax:         &sql.QueryTime,
				FirstQueryAt:         sql.QueryAt,
				DBUser:               sql.DBUser,
				Schema:               sql.Schema,
			}
			auditPlanSqlMap[key] = sqlReq
			sqlListReq = append(sqlListReq, sqlReq)
		}
	}

	return sq.c.UploadReq(scanner.PartialUpload, sq.apName, sqlListReq)
}
