//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) createSqlManageRawSQLs(sqls []*SQLManageRawSQL) error {
	return nil
}
