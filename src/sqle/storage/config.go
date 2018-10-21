package storage

import "github.com/jinzhu/gorm"

type InspectConfig struct {
	Model
	Code       int
	ConfigType int
	StmtType   int
	Variable   string
	Desc       string
	Level      int
	Disable    bool
}

type RuleTemplate struct {
	Model
	Name  string
	Desc  string
	Rules []Rule `gorm:"foreignkey:TemplateId"`
}

type Rule struct {
	Model      `json:"-"`
	TemplateId uint   `json:"-"`
	Code       string `json:"code"`
	Value      string `json:"value"`
	Level      string `json:"level"`
}

// inspector rule code
const (
	CHECK_NAME_LENGTH = "check_name_length"
)

const (
	SELECT_STMT_TABLE_MUST_EXIST = iota
)

// inspector rule level
const (
	RULE_LEVEL_SUGGEST = iota
	RULE_LEVEL_WARN
	RULE_LEVEL_ERROR
)

var DfConfigMap = []*InspectConfig{
	&InspectConfig{
		Code:       SELECT_STMT_TABLE_MUST_EXIST,
		ConfigType: 0,
		Variable:   "",
		StmtType:   0,
		Level:      RULE_LEVEL_WARN,
		Disable:    false,
	},
}

func (s *Storage) GetTemplateByName(name string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Preload("Rules").Where("name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, err
}

func (s *Storage) UpdateRules(tpl *RuleTemplate, rules ...Rule) error {
	return s.db.Model(tpl).Association("Rules").Append(rules).Error
}

func (s *Storage) GetAllTpl() ([]RuleTemplate, error) {
	ts := []RuleTemplate{}
	err := s.db.Preload("Rules").Find(&ts).Error
	return ts, err
}
