package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RecycleInstsanceV3 struct {

	// 实例id
	Id *string `json:"id,omitempty"`

	// 实例名
	Name *string `json:"name,omitempty"`

	// 实例主备模式，取值：Ha（主备），不区分大小写。
	HaMode *string `json:"ha_mode,omitempty"`

	// 引擎名
	EngineName *string `json:"engine_name,omitempty"`

	// 数据库引擎版本
	EngineVersion *string `json:"engine_version,omitempty"`

	// 计费方式
	PayModel *string `json:"pay_model,omitempty"`

	// 创建时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	CreatedAt *string `json:"created_at,omitempty"`

	// 删除时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	DeletedAt *string `json:"deleted_at,omitempty"`

	// 磁盘类型。 取值范围如下，区分大小写： - COMMON，表示SATA。 - HIGH，表示SAS。 - ULTRAHIGH，表示SSD。 - ULTRAHIGHPRO，表示SSD尊享版，仅支持超高性能型尊享版（需申请权限）。 - CLOUDSSD，表示SSD云盘，仅支持通用型和独享型规格实例。 - LOCALSSD，表示本地SSD。
	VolumeType *string `json:"volume_type,omitempty"`

	// 磁盘大小，单位为GB。 取值范围：40GB~4000GB，必须为10的整数倍。  部分用户支持40GB~6000GB，如果您想创建存储空间最大为6000GB的数据库实例，或提高扩容上限到10000GB，请联系客服开通。  说明：对于只读实例，该参数无效，磁盘大小，默认和主实例相同。
	VolumeSize *int32 `json:"volume_size,omitempty"`

	// 内网地址
	DataVip *string `json:"data_vip,omitempty"`

	// ipv6内网地址
	DataVipV6 *string `json:"data_vip_v6,omitempty"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 保留时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	RetainedUntil *string `json:"retained_until,omitempty"`

	// 备份id
	RecycleBackupId *string `json:"recycle_backup_id,omitempty"`

	// 备份状态 取值范围如下，区分大小写: - BUILDING 备份中，不能进行重建 - COMPLETED，标识备份完成，可以重建
	RecycleStatus *string `json:"recycle_status,omitempty"`
}

func (o RecycleInstsanceV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecycleInstsanceV3 struct{}"
	}

	return strings.Join([]string{"RecycleInstsanceV3", string(data)}, " ")
}
