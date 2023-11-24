package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EngineFlavorData struct {

	// CPU大小。例如：1表示1U。
	Vcpus *string `json:"vcpus,omitempty"`

	// 内存大小，单位为GB。
	Ram *string `json:"ram,omitempty"`

	// 资源规格编码。例如：rds.mysql.m1.xlarge.rr。  更多规格说明请参考数据库实例规格。  “rds”代表RDS产品。 “mysql”代表数据库引擎。 “m1.xlarge”代表性能规格，为高内存类型。 “rr”表示只读实例（“.ha”表示主备实例）。 “rha.rr”表示高可用只读实例，规格编码示例：rds.mysql.n1.large.4.rha.rr。 具有公测权限的用户才可选择高可用，您可联系华为云客服人员申请。 高可用只读功能介绍请参见高可用只读简介。
	SpecCode *string `json:"spec_code,omitempty"`

	// 是否支持ipv6。
	IsIpv6Supported *bool `json:"is_ipv6_supported,omitempty"`

	// 资源类型
	TypeCode *string `json:"type_code,omitempty"`

	// 规格所在az的状态，包含以下状态： normal：在售。 unsupported：暂不支持该规格。 sellout：售罄。 abandon：未启用
	AzStatus map[string]string `json:"az_status,omitempty"`

	// 性能规格，包含以下状态： normal：通用增强型。 normal2：通用增强Ⅱ型。 armFlavors：鲲鹏通用增强型。 dedicicateNormal（dedicatedNormalLocalssd）：x86独享型。 armLocalssd：鲲鹏通用型。 normalLocalssd：x86通用型。 general：通用型。 dedicated 对于MySQL引擎：独享型。 对于PostgreSQL和SQL Server引擎：独享型，仅云盘SSD支持。 rapid 对于MySQL引擎：独享型（已下线）。 对于PostgreSQL和SQL Server引擎：独享型，仅极速型SSD支持。 bigmem：超大内存型。 highPerformancePrivilegeEdition：超高IO 尊享版
	GroupType *string `json:"group_type,omitempty"`

	// 最大连接数
	MaxConnection *string `json:"max_connection,omitempty"`

	// 数据库每秒执行的事务数，每个事务中包含18条SQL语句。
	Tps *string `json:"tps,omitempty"`

	// 数据库每秒执行的SQL数，包含insert、select、update、delete等。
	Qps *string `json:"qps,omitempty"`

	// 最小磁盘容量，单位G
	MinVolumeSize *string `json:"min_volume_size,omitempty"`

	// 最大磁盘容量，单位G
	MaxVolumeSize *string `json:"max_volume_size,omitempty"`
}

func (o EngineFlavorData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EngineFlavorData struct{}"
	}

	return strings.Join([]string{"EngineFlavorData", string(data)}, " ")
}
