package storage

import "github.com/jinzhu/gorm"

type InspectConfig struct {
	gorm.Model
	Code       int
	ConfigType int
	StmtType   int
	Variable   string
	Desc       string
	Level      int
	Disable    bool
}

type ConfigTemplate struct {
	gorm.Model

}

type ConfigRule struct {
	gorm.Model
	
}

type ConfigMeta struct {
	gorm.Model
	Code string
	Name string

}

// inspector rule code
const (
	SELECT_STMT_TABLE_MUST_EXIST = iota
)

// inspector rule level
const (
	RULE_LEVEL_ERROR = iota
	RULE_LEVEL_WARN
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
