package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostgresqlHbaHistory struct {

	// 修改结果，    success：已生效     failed：未生效     setting：设置中\",
	Status *string `json:"status,omitempty"`

	// 修改时间
	Time *sdktime.SdkTime `json:"time,omitempty"`

	// 修改失败原因
	FailReason *string `json:"fail_reason,omitempty"`

	// 修改之前的值
	BeforeConfs *[]PostgresqlHbaConf `json:"before_confs,omitempty"`

	// 修改之后的值
	AfterConfs *[]PostgresqlHbaConf `json:"after_confs,omitempty"`
}

func (o PostgresqlHbaHistory) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlHbaHistory struct{}"
	}

	return strings.Join([]string{"PostgresqlHbaHistory", string(data)}, " ")
}
