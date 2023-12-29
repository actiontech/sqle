package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SimplifiedInstanceEntry struct {

	// 实例id
	Id string `json:"id"`

	// 创建的实例名称
	Name string `json:"name"`

	// 引擎名字
	EngineName string `json:"engine_name"`

	// 引擎版本
	EngineVersion string `json:"engine_version"`

	// 实例状态。 normal,表示正常 abnormal,表示异常 creating,表示创建中 createfail,表示创建失败 data_disk_full,表示磁盘满 deleted,表示删除 shutdown,表示关机
	InstanceStatus string `json:"instance_status"`

	// 是否冻结
	Frozen bool `json:"frozen"`

	// 按照实例类型查询。取值Single、Ha、Replica、Enterprise，分别对应于单实例、主备实例和只读实例、分布式实例（企业版）。
	Type string `json:"type"`

	// 按需还是包周期
	PayModel string `json:"pay_model"`

	// 规格码
	SpecCode string `json:"spec_code"`

	// 可用区集合
	AvailabilityZoneIds []string `json:"availability_zone_ids"`

	// 只读实例id集合
	ReadOnlyInstances []string `json:"read_only_instances"`

	// 当前实例操作动作集合
	CurrentActions []string `json:"current_actions"`

	// 磁盘类型。
	VolumeType string `json:"volume_type"`

	// 磁盘大小(单位:G)。
	VolumeSize int64 `json:"volume_size"`

	// 企业项目标签ID。
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

func (o SimplifiedInstanceEntry) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SimplifiedInstanceEntry struct{}"
	}

	return strings.Join([]string{"SimplifiedInstanceEntry", string(data)}, " ")
}
