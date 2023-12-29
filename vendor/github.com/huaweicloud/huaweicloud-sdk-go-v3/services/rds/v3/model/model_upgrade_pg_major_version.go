package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradePgMajorVersion struct {

	// 目标版本。 高于实例当前的大版本，如当前为12，目标版本需要是13或14。
	TargetVersion string `json:"target_version"`

	// 是否将实例内网IP切换到大版本实例  true：升级后切换当前实例的内网IP到大版本实例 false：升级后当前实例的内网IP不变，大版本实例使用新的内网IP
	IsChangePrivateIp bool `json:"is_change_private_ip"`

	// 统计信息收集方式。is_change_private_ip为true时必选  before_change_private_ip：将实例内网IP切换到大版本实例前收集  after_change_private_ip：将实例内网IP切换到大版本实例后收集
	StatisticsCollectionMode string `json:"statistics_collection_mode"`
}

func (o UpgradePgMajorVersion) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradePgMajorVersion struct{}"
	}

	return strings.Join([]string{"UpgradePgMajorVersion", string(data)}, " ")
}
