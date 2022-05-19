package model

import "fmt"

func (s *Storage) GetSqlQueryExecSqlByQueryId(id int) (*SqlQueryExecutionSql, error) {
	sql := &SqlQueryExecutionSql{}
	if err := s.db.Where("id = ?", id).Find(sql).Error; err != nil {
		return nil, err
	}
	return sql, nil
}

func (s *Storage) GetSqlQueryHistoryById(id uint) (*SqlQueryHistory, error) {
	history := &SqlQueryHistory{}
	if err := s.db.Where("id = ?", id).Find(history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (s *Storage) GetSqlQueryRawSqlByUserId(userID uint, pageIndex, pageSize uint32, fuzzyKey string) ([]SqlQueryHistory, error) {
	var res []SqlQueryHistory
	query := s.db.Select("raw_sql").Where("create_user_id = ?", userID)
	if fuzzyKey != "" {
		query = query.Where("raw_sql LIKE ?", fmt.Sprintf("%%%s%%", fuzzyKey))
	}
	if err := query.Group("raw_sql").Order("max(created_at) desc").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
