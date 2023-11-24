package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OpsWindowRequest struct {

	// - 开始时间， UTC时间
	StartTime string `json:"start_time"`

	// - 结束时间，UTC时间
	EndTime string `json:"end_time"`
}

func (o OpsWindowRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpsWindowRequest struct{}"
	}

	return strings.Join([]string{"OpsWindowRequest", string(data)}, " ")
}
