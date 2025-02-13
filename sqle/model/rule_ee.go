//go:build enterprise
// +build enterprise

package model

import "github.com/actiontech/sqle/sqle/errors"

// 获取所有的自定义规则以及Knowledge
func (s *Storage) GetAllCustomRules() ([]*CustomRule, error) {
	rules := []*CustomRule{}
	err := s.db.Preload("Knowledge").Preload("Categories").Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}
