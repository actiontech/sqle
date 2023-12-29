package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GetRestoreTimeResponseRestoreTime struct {

	// 可恢复时间段的起始时间点，UNIX时间戳格式，单位是毫秒，时区是UTC。
	StartTime int64 `json:"start_time"`

	// 可恢复时间段的结束时间点，UNIX时间戳格式，单位是毫秒，时区是UTC。
	EndTime int64 `json:"end_time"`
}

func (o GetRestoreTimeResponseRestoreTime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetRestoreTimeResponseRestoreTime struct{}"
	}

	return strings.Join([]string{"GetRestoreTimeResponseRestoreTime", string(data)}, " ")
}
