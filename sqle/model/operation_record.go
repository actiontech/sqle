package model

import (
	"time"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"golang.org/x/text/language"
)

type OperationRecord struct {
	Model
	OperationTime        time.Time       `gorm:"column:operation_time;type:datetime;" json:"operation_time"`
	OperationUserName    string          `gorm:"column:operation_user_name;type:varchar(255);not null" json:"operation_user_name"`
	OperationReqIP       string          `gorm:"column:operation_req_ip; type:varchar(255)" json:"operation_req_ip"`
	OperationTypeName    string          `gorm:"column:operation_type_name; type:varchar(255)" json:"operation_type_name"`
	OperationAction      string          `gorm:"column:operation_action; type:varchar(255)" json:"operation_action"`
	OperationContent     string          `gorm:"column:operation_content; type:varchar(255)" json:"operation_content"` // Deprecated: use OperationI18nContent instead
	OperationProjectName string          `gorm:"column:operation_project_name; type:varchar(255)" json:"operation_project_name"`
	OperationStatus      string          `gorm:"column:operation_status; type:varchar(255)" json:"operation_status"`
	OperationI18nContent i18nPkg.I18nStr `gorm:"column:operation_i18n_content; type:json" json:"operation_i18n_content"`
}

func (o *OperationRecord) GetOperationContentByLangTag(lang language.Tag) string {
	if o.OperationContent != "" {
		// 兼容老sqle的数据
		o.OperationI18nContent.SetStrInLang(i18nPkg.DefaultLang, o.OperationContent)
	}
	return o.OperationI18nContent.GetStrInLang(lang)
}

func (s *Storage) GetOperationRecordProjectNameList() ([]string, error) {
	var projectNameList []string
	err := s.db.Model(&OperationRecord{}).Group("operation_project_name").Pluck("operation_project_name", &projectNameList).Error
	if err != nil {
		return nil, err
	}
	return projectNameList, err
}

func (s *Storage) GetExpiredOperationRecordIDListByStartTime(start time.Time) ([]string, error) {
	var idList []string
	err := s.db.Model(&OperationRecord{}).Where("operation_time < ?", start).Pluck("id", &idList).Error
	if err != nil {
		return nil, err
	}
	return idList, err
}

func (s *Storage) DeleteExpiredOperationRecordByIDList(idList []string) error {
	return s.db.Exec("DELETE FROM operation_records WHERE id IN (?)", idList).Error
}
