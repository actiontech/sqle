package model

type CompanyNotice struct {
	Model
	NoticeStr string `gorm:"type:mediumtext;comment:'企业公告'" json:"notice_str"`
}

func (s *Storage) TableName() string {
	return "company_notices"
}
