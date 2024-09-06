package model

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

func (s *Storage) GetCompanyNotice() (*CompanyNotice, error) {
	notice := new(CompanyNotice)
	err := s.db.First(&notice).Error
	if e.Is(err, gorm.ErrRecordNotFound) {
		return notice, nil
	}
	return notice, errors.New(errors.ConnectStorageError, err)
}
