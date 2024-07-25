package mybatis

import (
	"context"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/sirupsen/logrus"
)

type MyBatis struct {
	l *logrus.Entry
	c *scanner.Client

	xmlDir         string
	skipErrorQuery bool
	skipErrorXml   bool
	skipAudit      bool
	dbType         string
	instName       string
	schemaName     string
}

type Params struct {
	XMLDir         string
	SkipErrorQuery bool
	SkipErrorXml   bool
	SkipAudit      bool
	DbType         string
	InstName       string
	SchemaName     string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*MyBatis, error) {
	return &MyBatis{
		xmlDir:         params.XMLDir,
		skipErrorQuery: params.SkipErrorQuery,
		skipErrorXml:   params.SkipErrorXml,
		skipAudit:      params.SkipAudit,
		dbType:         params.DbType,
		instName:       params.InstName,
		schemaName:     params.SchemaName,
		l:              l,
		c:              c,
	}, nil
}

func (mb *MyBatis) Run(ctx context.Context) error {
	sqls, err := common.GetSQLFromPath(mb.xmlDir, mb.skipErrorQuery, mb.skipErrorXml, utils.MybatisFileSuffix)
	if err != nil {
		return err
	}

	if mb.skipAudit {
		return nil
	}

	return common.DirectAudit(ctx, mb.c, sqls, mb.dbType, mb.instName, mb.schemaName)
}

func (mb *MyBatis) SQLs() <-chan scanners.SQL {
	return nil
}

func (mb *MyBatis) Upload(ctx context.Context, sqls []scanners.SQL) error {
	return nil
}
