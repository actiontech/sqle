package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ListInstancesSupportFastRestoreRequestBody struct {

	// 要恢复的时间点，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	RestoreTime string `json:"restore_time"`

	// 实例id列表。
	InstanceIds []string `json:"instance_ids"`
}

func (o ListInstancesSupportFastRestoreRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesSupportFastRestoreRequestBody struct{}"
	}

	return strings.Join([]string{"ListInstancesSupportFastRestoreRequestBody", string(data)}, " ")
}
