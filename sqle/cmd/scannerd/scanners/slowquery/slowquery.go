package slowquery

import (
	"context"
	"fmt"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/scanners"
	"actiontech.cloud/sqle/sqle/sqle/pkg/scanner"
	pkg "actiontech.cloud/sqle/sqle/sqle/pkg/scanner"

	"github.com/go-sql-driver/mysql"
	"github.com/percona/pmm-agent/agents/mysql/perfschema"
	"github.com/sirupsen/logrus"
)

type SlowQuery struct {
	l  *logrus.Entry
	c  *pkg.Client
	ps *perfschema.PerfSchema

	apName          string
	slowQuerySecond int
}

type Params struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string

	APName          string
	SlowQuerySecond int

	// todo: support PG
	// DBType string
}

func New(params *Params, l *logrus.Entry, c *pkg.Client) (*SlowQuery, error) {
	dsn := mysql.NewConfig()
	dsn.Addr = fmt.Sprintf("%v:%v", params.DBHost, params.DBPort)
	dsn.Net = "tcp"
	dsn.User = params.DBUser
	dsn.Passwd = params.DBPass
	pp := &perfschema.Params{
		DSN: dsn.FormatDSN(),
	}
	ps, err := perfschema.New(pp, l)
	if err != nil {
		return nil, err
	}

	return &SlowQuery{
		l:               l,
		c:               c,
		ps:              ps,
		apName:          params.APName,
		slowQuerySecond: params.SlowQuerySecond}, nil
}

func (sq *SlowQuery) Run(ctx context.Context) error {
	sq.ps.Run(ctx)
	return nil
}

func (sq *SlowQuery) SQLs() <-chan scanners.SQL {
	sqlCh := make(chan scanners.SQL, 10)
	go func() {
		for change := range sq.ps.Changes() {
			sq.l.Infoln("change status", change.Status)
			for _, mb := range change.MetricsBucket {
				if int(mb.Common.MQueryTimeSum) < sq.slowQuerySecond {
					continue
				}
				sq.l.Infoln("%+v", mb)
				sql := scanners.SQL{
					Counter:     int(mb.Common.MQueryTimeCnt),
					Fingerprint: mb.Common.Fingerprint,
					RawText:     mb.Common.Example,
				}
				sqlCh <- sql
			}
		}
		close(sqlCh)
	}()
	return sqlCh
}

func (sq *SlowQuery) Upload(ctx context.Context, sqls []scanners.SQL) error {
	var sqlsReq []scanner.AuditPlanSQLReq
	now := time.Now()
	for _, sql := range sqls {
		sqlsReq = append(sqlsReq, scanner.AuditPlanSQLReq{
			Fingerprint:          sql.Fingerprint,
			LastReceiveText:      sql.RawText,
			LastReceiveTimestamp: now.Format(time.RFC3339),
			Counter:              "1",
		})
	}
	return sq.c.UploadReq(scanner.PartialUpload, sq.apName, sqlsReq)
}
