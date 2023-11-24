package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTdeStatusRequestBody sqlserverTDE开关信息
type UpdateTdeStatusRequestBody struct {

	// 轮转天数
	RotateDay *int32 `json:"rotate_day,omitempty"`

	// 密钥ID
	SecretId *string `json:"secret_id,omitempty"`

	// 密钥名称
	SecretName *string `json:"secret_name,omitempty"`

	// 密钥版本
	SecretVersion *string `json:"secret_version,omitempty"`
}

func (o UpdateTdeStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTdeStatusRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTdeStatusRequestBody", string(data)}, " ")
}
