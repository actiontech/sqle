package model

import (
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
)

type RuleTemplateRule struct {
	gorm.JoinTableHandler
	RuleTemplateId uint   `json:"rule_template_id" gorm:"column:rule_template_id;primary_key;auto_increment:false"`
	RuleName       string `json:"name" gorm:"column:rule_name;primary_key;"`
	RuleLevel      string `json:"level" gorm:"column:level;"`
	RuleValue      string `json:"value" gorm:"column:value;" `
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}

// todo 增改时填充 来自rules表的level、value的默认值 如不需要可移除
func (rtr *RuleTemplateRule) Add(_ gorm.JoinTableHandlerInterface, db *gorm.DB, foreignValue interface{}, associationValue interface{}) error {
	association := db.NewScope(associationValue)
	associationPrimaryKey := fmt.Sprint(association.PrimaryKeyValue())
	foreignPrimaryKey, _ := strconv.Atoi(fmt.Sprint(db.NewScope(foreignValue).PrimaryKeyValue()))

	attrMap := make(map[string]interface{}, 0)
	for _, target := range []string{"value", "level"} {
		if f, ok := association.FieldByName(target); ok {
			attrMap[target] = f.Field.Interface()
		}
	}

	ruleTR := &RuleTemplateRule{
		RuleTemplateId: uint(foreignPrimaryKey),
		RuleName:       associationPrimaryKey,
	}
	return db.Where(*ruleTR).Assign(attrMap).FirstOrCreate(ruleTR).Error
}

// 完成连接表的level value字段更新
func (s *Storage) AfterUpdateRuleTemplateRules(tpl *RuleTemplate, rules ...RuleTemplateRule) error {
	err := s.db.Model(tpl).Association("RTR").Replace(rules).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
