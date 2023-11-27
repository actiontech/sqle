package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeReports struct {

	// 升级报告ID。
	Id string `json:"id"`

	// 升级开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	StartTime string `json:"start_time"`

	// 升级结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	EndTime string `json:"end_time"`

	// 原实例ID。
	SrcInstanceId string `json:"src_instance_id"`

	// 原数据库版本。
	SrcDatabaseVersion string `json:"src_database_version"`

	// 目标实例ID。
	DstInstanceId string `json:"dst_instance_id"`

	// 目标数据库版本。
	DstDatabaseVersion string `json:"dst_database_version"`

	// 升级结果。 success，表示成功。 failed，表示失败。 running， 表示升级中。
	Result string `json:"result"`

	// 实例内网IP是否改变。 true，表示改变。 false，表示不改变。
	IsPrivateIpChanged bool `json:"is_private_ip_changed"`

	// 实例内网IP修改时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	PrivateIpChangeTime string `json:"private_ip_change_time"`

	// 统计信息收集模式。 before_change_private_ip，修改实例内网IP前收集。 after_change_private_ip，修改实例内网IP后收集。
	StatisticsCollectionMode string `json:"statistics_collection_mode"`

	// 升级报告详情。
	Detail string `json:"detail"`
}

func (o UpgradeReports) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeReports struct{}"
	}

	return strings.Join([]string{"UpgradeReports", string(data)}, " ")
}
