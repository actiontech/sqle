package sqlFile

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/common"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/sirupsen/logrus"
)

type SQLFile struct {
	l *logrus.Entry
	c *scanner.Client

	sqlDir           string
	skipErrorSqlFile bool
	dbType           string
	instName         string
	schemaName       string
}

type Params struct {
	SQLDir           string
	SkipErrorQuery   bool
	SkipErrorSqlFile bool
	DbType           string
	InstName         string
	SchemaName       string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SQLFile, error) {
	return &SQLFile{
		sqlDir:           params.SQLDir,
		skipErrorSqlFile: params.SkipErrorSqlFile,
		dbType:           params.DbType,
		instName:         params.InstName,
		schemaName:       params.SchemaName,
		l:                l,
		c:                c,
	}, nil
}

func (sf *SQLFile) Run(ctx context.Context) error {
	sqls, err := common.GetSQLFromPath(sf.sqlDir, false, sf.skipErrorSqlFile, utils.SQLFileSuffix)
	if err != nil {
		return fmt.Errorf("failed to get sql from path: %v", err)
	}

	return common.DirectAudit(ctx, sf.c, sqls, sf.dbType, sf.instName, sf.schemaName)
}

func (sf *SQLFile) SQLs() <-chan scanners.SQL {
	return nil
}

func (sf *SQLFile) Upload(ctx context.Context, sqls []scanners.SQL, errorMessage string) error {
	return nil
}
