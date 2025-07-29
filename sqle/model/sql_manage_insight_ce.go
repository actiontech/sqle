//go:build !enterprise
// +build !enterprise

package model

import "time"

func (s *Storage) createSqlManageRawSQLs(sqls []*SQLManageRawSQL) error {
	return nil
}

func (s *Storage) RemoveExpiredSqlInsightRecord(expiredTime time.Time) (int64, error) {
	return 0, nil
}
