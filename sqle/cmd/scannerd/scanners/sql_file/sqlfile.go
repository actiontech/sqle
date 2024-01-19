package sqlFile

import (
	"context"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/sirupsen/logrus"
)

type SQLFile struct {
	l *logrus.Entry
	c *scanner.Client

	sqls []scanners.SQL

	allSQL []driverV2.Node
	getAll chan struct{}

	apName           string
	sqlDir           string
	skipErrorSqlFile bool
	skipAudit        bool
}

type Params struct {
	SQLDir           string
	APName           string
	SkipErrorQuery   bool
	SkipErrorSqlFile bool
	SkipAudit        bool
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SQLFile, error) {
	return &SQLFile{
		sqlDir:           params.SQLDir,
		apName:           params.APName,
		skipErrorSqlFile: params.SkipErrorSqlFile,
		skipAudit:        params.SkipAudit,
		l:                l,
		c:                c,
		getAll:           make(chan struct{}),
	}, nil
}

func (sf *SQLFile) Run(ctx context.Context) error {
	sqls, err := common.GetSQLFromPath(sf.sqlDir, false, sf.skipErrorSqlFile, utils.SQLFileSuffix)
	if err != nil {
		return err
	}

	sf.allSQL = sqls
	close(sf.getAll)

	<-ctx.Done()
	return nil
}

func (sf *SQLFile) SQLs() <-chan scanners.SQL {
	// todo: channel size configurable
	sqlCh := make(chan scanners.SQL, 10240)

	go func() {
		<-sf.getAll
		for _, sql := range sf.allSQL {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
			}
		}
		close(sqlCh)
	}()
	return sqlCh
}

func (sf *SQLFile) Upload(ctx context.Context, sqls []scanners.SQL) error {
	sf.sqls = append(sf.sqls, sqls...)
	err := common.Upload(ctx, sf.sqls, sf.c, sf.apName)
	if err != nil {
		return err
	}

	if sf.skipAudit {
		return nil
	}

	return common.Audit(sf.c, sf.apName)
}
