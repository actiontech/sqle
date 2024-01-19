package mybatis

import (
	"context"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/sirupsen/logrus"
)

type MyBatis struct {
	l *logrus.Entry
	c *scanner.Client

	sqls []scanners.SQL

	allSQL []driverV2.Node
	getAll chan struct{}

	apName         string
	xmlDir         string
	skipErrorQuery bool
	skipErrorXml   bool
	skipAudit      bool
}

type Params struct {
	XMLDir         string
	APName         string
	SkipErrorQuery bool
	SkipErrorXml   bool
	SkipAudit      bool
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*MyBatis, error) {
	return &MyBatis{
		xmlDir:         params.XMLDir,
		apName:         params.APName,
		skipErrorQuery: params.SkipErrorQuery,
		skipErrorXml:   params.SkipErrorXml,
		skipAudit:      params.SkipAudit,
		l:              l,
		c:              c,
		getAll:         make(chan struct{}),
	}, nil
}

func (mb *MyBatis) Run(ctx context.Context) error {
	sqls, err := common.GetSQLFromPath(mb.xmlDir, mb.skipErrorQuery, mb.skipErrorXml, utils.MybatisFileSuffix)
	if err != nil {
		return err
	}

	mb.allSQL = sqls
	close(mb.getAll)

	<-ctx.Done()
	return nil
}

func (mb *MyBatis) SQLs() <-chan scanners.SQL {
	// todo: channel size configurable
	sqlCh := make(chan scanners.SQL, 10240)

	go func() {
		<-mb.getAll
		for _, sql := range mb.allSQL {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
			}
		}
		close(sqlCh)
	}()
	return sqlCh
}

func (mb *MyBatis) Upload(ctx context.Context, sqls []scanners.SQL) error {
	mb.sqls = append(mb.sqls, sqls...)
	err := common.Upload(ctx, mb.sqls, mb.c, mb.apName)
	if err != nil {
		return err
	}
	if mb.skipAudit {
		return nil
	}
	return common.Audit(mb.c, mb.apName)
}
