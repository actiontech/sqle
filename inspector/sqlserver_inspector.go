package inspector

import (
	"github.com/sirupsen/logrus"
	"sqle/model"
	"sqle/sqlserverClient"
)

type SqlserverInspect struct {
	*Inspect
}

func NeSqlserverInspect(entry *logrus.Entry, task *model.Task) Inspector {
	return &SqlserverInspect{
		Inspect: NewInspect(entry, task),
	}
}

func (i *SqlserverInspect) SplitSql(sql string) ([]string, error) {
	sqls, err := sqlserverClient.GetClient().SplitSql(sql)
	if err != nil {
		i.Logger().Errorf("parser t-sql from ms grpc server failed, error: %v", err)
	}
	return sqls, err
}

func (i *SqlserverInspect) Add(sql *model.Sql, action func(sql *model.Sql) error) error {
	i.SqlArray = append(i.SqlArray, sql)
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *SqlserverInspect) Do() error {
	for n, sql := range i.SqlArray {
		err := i.SqlAction[n](sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *SqlserverInspect) Advise(rules []model.Rule) error {
	i.Logger().Info("start advise sql")
	var meta = sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	err := sqlserverClient.GetClient().Advise(i.Task.CommitSqls, rules, meta)
	if err != nil {
		i.Logger().Errorf("advise t-sql from ms grpc server failed, error: %v", err)
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}

func (i *SqlserverInspect) GenerateAllRollbackSql() ([]*model.RollbackSql, error) {
	i.Logger().Info("start generate rollback sql")
	i.Logger().Info("generate rollback sql finish")
	return []*model.RollbackSql{}, nil
}
