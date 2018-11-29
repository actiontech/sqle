package inspector

import (
	"sqle/model"
	"sqle/sqlserverClient"
)

type SqlserverInspect struct {
	*Inspect
}

func NeSqlserverInspect(task *model.Task)  Inspector {
	return &SqlserverInspect{
		Inspect:NewInspect(task),
	}
}

func (i *SqlserverInspect) SplitSql(sql string) ([]string, error) {
	return sqlserverClient.GetClient().SplitSql(sql)
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
	return sqlserverClient.GetClient().Advise(i.Task.CommitSqls, rules)
}

func (i *SqlserverInspect) GenerateAllRollbackSql() ([]*model.RollbackSql, error) {
	return []*model.RollbackSql{}, nil
}

