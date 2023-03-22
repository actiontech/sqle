//go:build enterprise
// +build enterprise

package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

func (s *Storage) GetPersonaliseConfig() (*PersonaliseConfig, bool, error) {
	pc := new(PersonaliseConfig)
	err := s.db.Last(&pc).Error
	if err == gorm.ErrRecordNotFound {
		return pc, false, nil
	}
	return pc, true, errors.New(errors.ConnectStorageError, err)
}
