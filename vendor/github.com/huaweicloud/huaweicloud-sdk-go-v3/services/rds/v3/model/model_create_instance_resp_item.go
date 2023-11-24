package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateInstanceRespItem 实例信息。
type CreateInstanceRespItem struct {

	// 实例id
	Id *string `json:"id,omitempty"`

	// 实例名称。 用于表示实例的名称，同一租户下，同类型的实例名可重名，其中，SQL Server实例名唯一。 取值范围：4~64个字符之间，必须以字母开头，区分大小写，可以包含字母、数字、中划线或者下划线，不能包含其他的特殊字符。
	Name string `json:"name"`

	// 实例状态。如BUILD，表示创建中。 仅创建按需实例时会返回该参数。
	Status *string `json:"status,omitempty"`

	Datastore *Datastore `json:"datastore"`

	Ha *Ha `json:"ha,omitempty"`

	// 参数组ID。
	ConfigurationId *string `json:"configuration_id,omitempty"`

	// 数据库端口信息。  - MySQL数据库端口设置范围为1024～65535（其中12017和33071被RDS系统占用不可设置）。 - PostgreSQL数据库端口修改范围为2100～9500。 - Microsoft SQL Server实例的端口设置范围为1433和2100~9500（其中5355和5985不可设置。对于2017 EE、2017 SE、2017 Web版，5050、5353和5986不可设置。  当不传该参数时，默认端口如下：  - MySQL默认3306。 - PostgreSQL默认5432。 - Microsoft SQL Server默认1433。
	Port *string `json:"port,omitempty"`

	BackupStrategy *BackupStrategy `json:"backup_strategy,omitempty"`

	// 企业项目ID。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 用于磁盘加密的密钥ID。
	DiskEncryptionId *string `json:"disk_encryption_id,omitempty"`

	// 规格码。
	FlavorRef string `json:"flavor_ref"`

	Volume *Volume `json:"volume"`

	// 区域ID。创建主实例时必选，其它场景不可选。 取值参见[地区和终端节点](https://developer.huaweicloud.com/endpoint)。
	Region string `json:"region"`

	// 可用区ID。对于数据库实例类型不是单机的实例，需要分别为实例所有节点指定可用区，并用逗号隔开。 取值参见[地区和终端节点](https://developer.huaweicloud.com/endpoint)。
	AvailabilityZone string `json:"availability_zone"`

	// 虚拟私有云ID。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。
	VpcId string `json:"vpc_id"`

	// 子网ID。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。
	SubnetId string `json:"subnet_id"`

	// 安全组ID。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。
	SecurityGroupId string `json:"security_group_id"`

	ChargeInfo *ChargeInfo `json:"charge_info,omitempty"`

	// 仅限Microsoft SQL Server实例使用。取值范围：根据查询SQL Server可用字符集的字符集查询列表查询可设置的字符集。
	Collation *string `json:"collation,omitempty"`

	RestorePoint *RestorePoint `json:"restore_point,omitempty"`
}

func (o CreateInstanceRespItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceRespItem struct{}"
	}

	return strings.Join([]string{"CreateInstanceRespItem", string(data)}, " ")
}
