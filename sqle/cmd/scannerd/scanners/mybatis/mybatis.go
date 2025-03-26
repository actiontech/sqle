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

	xmlDir          string
	skipErrorQuery  bool
	skipErrorXml    bool
	dbType          string
	instName        string
	schemaName      string
	showFileContent bool
}

type Params struct {
	XMLDir          string
	SkipErrorQuery  bool
	SkipErrorXml    bool
	DbType          string
	InstName        string
	SchemaName      string
	ShowFileContent bool
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*MyBatis, error) {
	return &MyBatis{
		xmlDir:          params.XMLDir,
		skipErrorQuery:  params.SkipErrorQuery,
		skipErrorXml:    params.SkipErrorXml,
		dbType:          params.DbType,
		instName:        params.InstName,
		schemaName:      params.SchemaName,
		showFileContent: params.ShowFileContent,
		l:               l,
		c:               c,
	}, nil
}

func (mb *MyBatis) Run(ctx context.Context) error {
	sqls, err := common.GetSQLFromPath(mb.xmlDir, mb.skipErrorQuery, mb.skipErrorXml, utils.MybatisFileSuffix, mb.showFileContent)
	if err != nil {
		return err
	}

	return common.DirectAudit(ctx, mb.c, sqls, mb.dbType, mb.instName, mb.schemaName)
}

func (mb *MyBatis) SQLs() <-chan scanners.SQL {
	return nil
}

func (mb *MyBatis) Upload(ctx context.Context, sqls []scanners.SQL, errorMessage string) error {
	return nil
}
