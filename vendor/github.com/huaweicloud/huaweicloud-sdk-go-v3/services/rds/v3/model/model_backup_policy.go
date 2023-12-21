package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BackupPolicy 备份策略信息。
type BackupPolicy struct {

	// 指定已生成的备份文件可以保存的天数。  取值范围：0～732。取0值，表示关闭自动备份策略。如果需要延长保留时间请联系客服人员申请，自动备份最长可以保留2562天。  注意： 关闭备份策略后，备份任务将立即停止，所有增量备份任务将立即删除，使用增量备份的相关操作可能失败，相关操作不限于下载、复制、恢复、重建等，请谨慎操作。
	KeepDays int32 `json:"keep_days"`

	// 备份时间段。自动备份将在该时间段内触发。除关闭自动备份策略外，必选。  取值范围：格式必须为hh:mm-HH:MM且有效，当前时间指UTC时间。  - HH取值必须比hh大1。 - mm和MM取值必须相同，且取值必须为00、15、30或45。
	StartTime *string `json:"start_time,omitempty"`

	// 备份周期配置。自动备份将在每星期指定的天进行。除关闭自动备份策略外，必选。  取值范围：格式为逗号隔开的数字，数字代表星期。
	Period *string `json:"period,omitempty"`
}

func (o BackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupPolicy struct{}"
	}

	return strings.Join([]string{"BackupPolicy", string(data)}, " ")
}
