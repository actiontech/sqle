//go:build enterprise
// +build enterprise

package slowquery

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
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
}

type Params struct {
	LogFilePath string
	APName      string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SlowQuery, error) {
	return &SlowQuery{
		l:           l,
		c:           c,
		logFilePath: params.LogFilePath,
		apName:      params.APName,
		sqlCh:       make(chan scanners.SQL, 10),
	}, nil
}

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
			s.sqlCh <- scanners.SQL{
				Fingerprint: query.Fingerprint(e.Query),
				RawText:     e.Query,
				Counter:     1}

		}
	}

}

func (sq *SlowQuery) SQLs() <-chan scanners.SQL {
	return sq.sqlCh
}

func (sq *SlowQuery) Upload(ctx context.Context, sqls []scanners.SQL) error {
	var sqlListReq []*scanner.AuditPlanSQLReq
	now := time.Now()
	auditPlanSqlMap := make(map[string]*scanner.AuditPlanSQLReq, 0)
	for _, sql := range sqls {
		if sqlReq, ok := auditPlanSqlMap[sql.Fingerprint]; ok {
			atoi, err := strconv.Atoi(sqlReq.Counter)
			if err != nil {
				return err
			}

			sqlReq.Counter = strconv.Itoa(atoi + 1)
			sqlReq.LastReceiveText = sql.RawText
			sqlReq.LastReceiveTimestamp = now.Format(time.RFC3339)
		} else {
			sqlReq := &scanner.AuditPlanSQLReq{
				Fingerprint:          sql.Fingerprint,
				LastReceiveText:      sql.RawText,
				LastReceiveTimestamp: now.Format(time.RFC3339),
				Counter:              "1",
			}
			auditPlanSqlMap[sql.Fingerprint] = sqlReq
			sqlListReq = append(sqlListReq, sqlReq)
		}
	}

	return sq.c.UploadReq(scanner.PartialUpload, sq.apName, sqlListReq)
}
