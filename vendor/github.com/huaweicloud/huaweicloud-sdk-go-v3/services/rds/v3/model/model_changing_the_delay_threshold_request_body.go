package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ChangingTheDelayThresholdRequestBody struct {

	// 延时阈值（单位：KB），取值范围为0~10485760。
	DelayThresholdInKilobytes int32 `json:"delay_threshold_in_kilobytes"`
}

func (o ChangingTheDelayThresholdRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangingTheDelayThresholdRequestBody struct{}"
	}

	return strings.Join([]string{"ChangingTheDelayThresholdRequestBody", string(data)}, " ")
}
