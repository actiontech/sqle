//go:build enterprise
// +build enterprise

package tidb_audit_log

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

type AuditLog struct {
	l *logrus.Entry
	c *scanner.Client

	logFilePath string

	sqlCh chan scanners.SQL

	auditPlanID string
	apType      string

	parser TiDBAuditLogParser
}

type Params struct {
	AuditLogPath string

	AuditPlanID string
	APType      string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*AuditLog, error) {
	return &AuditLog{
		l:           l,
		c:           c,
		logFilePath: params.AuditLogPath,
		auditPlanID: params.AuditPlanID,
		sqlCh:       make(chan scanners.SQL, 10),
		parser:      GetLexerParser(),
	}, nil
}

func (s *AuditLog) parseLine(line string) (sql *TiDBAuditLog, schema string, err error) {
	defer func() {
		if e := recover(); e != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("parse log panic, error: %v | stack: %v", e, string(buf[:n]))
			err = errors.New("panic")
		}
	}()

	sql, err = s.parser.Parse(line)
	if err != nil {
		s.l.Errorf("Malformed log, skip parsing: %v", line)
		return sql, "", err
	}
	schema, err = GetMissingSchema(sql.SQLText, sql.Databases)
	if err != nil {
		s.l.Warnf("sql execution location not found: %v", line)
		return sql, schema, err
	}
	return sql, schema, nil

}

func (s *AuditLog) Run(ctx context.Context) error {
	reader, err := scanners.NewContinuousFileReader(s.logFilePath, logrus.WithField("filename", s.logFilePath), true)
	if err != nil {
		return err
	}

	sqls := make(chan *scanners.SQL, 1000)
	go func() {
		for {
			line, err := reader.NextLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				s.l.Errorf("read log file error: %v", err)
				return
			}

			sql, schema, err := s.parseLine(line)
			if err != nil {
				continue
			}

			sqls <- &scanners.SQL{
				RawText: sql.SQLText,
				Schema:  schema,
			}

		}
	}()

	for {
		select {
		case <-ctx.Done():
			err = reader.Close()
			s.l.Infof("context done: %s. reader closed: %v.", ctx.Err(), err)
			close(s.sqlCh)
			return nil

		case e, ok := <-sqls:
			if !ok {
				return nil
			}
			s.l.Debugf("parsed tidb audit log event: %+v.", e)
			fingerprint, err := util.Fingerprint(e.RawText, true)
			if err != nil {
				return err
			}
			if e.Schema != "" {
				fingerprint = fmt.Sprintf("%v -- current schema: %v", fingerprint, e.Schema)
			}
			s.sqlCh <- scanners.SQL{
				Fingerprint: fingerprint,
				RawText:     e.RawText,
				Counter:     1,
				Schema:      e.Schema,
			}

		}
	}

}

func (sq *AuditLog) SQLs() <-chan scanners.SQL {
	return sq.sqlCh
}

func (sq *AuditLog) Upload(ctx context.Context, sqls []scanners.SQL, errorMessage string) error {
	sqlsReq := []*scanner.AuditPlanSQLReq{}
	now := time.Now()

	counter := map[string]int /*slice subscript*/ {}

	for _, sql := range sqls {
		if subscript, ok := counter[sql.Fingerprint]; ok {
			count, _ := strconv.Atoi(sqlsReq[subscript].Counter)
			sqlsReq[subscript].Counter = strconv.Itoa(count + 1)
		} else {
			sqlsReq = append(sqlsReq, &scanner.AuditPlanSQLReq{
				Fingerprint:          sql.Fingerprint,
				LastReceiveText:      sql.RawText,
				LastReceiveTimestamp: now.Format(time.RFC3339),
				Counter:              "1",
				Schema:               sql.Schema,
			})
			counter[sql.Fingerprint] = len(sqlsReq) - 1
		}
	}
	return sq.c.UploadReq(scanner.UploadSQL, sq.auditPlanID, errorMessage, sqlsReq)

}
