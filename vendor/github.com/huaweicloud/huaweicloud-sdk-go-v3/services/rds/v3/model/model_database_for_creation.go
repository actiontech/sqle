package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DatabaseForCreation 数据库信息。
type DatabaseForCreation struct {

	// 数据库名称。 数据库名称长度可在1～64个字符之间，由字母、数字、中划线、下划线或$组成，$累计总长度小于等于10个字符，（MySQL 8.0不可包含$）。
	Name string `json:"name"`

	// 数据库使用的字符集，例如utf8、gbk、ascii等MySQL支持的字符集。
	CharacterSet string `json:"character_set"`

	// 数据库备注，最大长度512
	Comment *string `json:"comment,omitempty"`
}

func (o DatabaseForCreation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DatabaseForCreation struct{}"
	}

	return strings.Join([]string{"DatabaseForCreation", string(data)}, " ")
}
