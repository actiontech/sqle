package server

import (
	"context"

	sqlop "github.com/actiontech/dms/pkg/dms-common/sql_op"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/sirupsen/logrus"
)

func GetSQLOperations(l *logrus.Entry, sql string, dbType string) ([]*sqlop.SQLObjectOps, error) {
	plugin, err := common.NewDriverManagerWithoutCfg(l, dbType)
	if err != nil {
		return nil, err
	}

	defer plugin.Close(context.TODO())

	return plugin.GetSQLOp(context.TODO(), sql)
}
