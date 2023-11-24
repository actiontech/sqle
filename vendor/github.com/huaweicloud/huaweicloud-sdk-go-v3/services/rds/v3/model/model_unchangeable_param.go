package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UnchangeableParam struct {

	// 表名大小写是否敏感，默认值是“1”，当前仅MySQL 8.0支持。 取值范围： - 0：表名被存储成固定且表名称大小写敏感。 - 1：表名将被存储成小写且表名称大小写不敏感。
	LowerCaseTableNames *string `json:"lower_case_table_names,omitempty"`
}

func (o UnchangeableParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnchangeableParam struct{}"
	}

	return strings.Join([]string{"UnchangeableParam", string(data)}, " ")
}
