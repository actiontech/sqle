package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateRestoreInstanceRequestBody 实例信息。
type CreateRestoreInstanceRequestBody struct {

	// 实例名称。 用于表示实例的名称，同一租户下，同类型的实例名可重名。取值范围如下： - MySQL数据库支持的字符长度是4~64个字符，必须以字母开头，区分大小写，可以包含字母、数字、中文字符、中划线或者下划线，不能包含其他的特殊字符。 - PostgreSQL和SQL Server数据库支持的字符长度是4~64个字符，必须以字母开头，区分大小写，可以包含字母、数字、中划线或者下划线，不能包含其他的特殊字符。
	Name string `json:"name"`

	Datastore *Datastore `json:"datastore"`

	Ha *Ha `json:"ha,omitempty"`

	// 参数模板ID。可调用[获取参数模板列表](https://support.huaweicloud.com/api-rds/rds_09_0301.html)接口返回的ID获取。
	ConfigurationId *string `json:"configuration_id,omitempty"`

	// 数据库端口信息。  - MySQL数据库端口设置范围为1024～65535（其中12017和33071被RDS系统占用不可设置）。 - PostgreSQL数据库端口修改范围为2100～9500。 - Microsoft SQL Server实例的端口设置范围为1433和2100~9500（其中5355和5985不可设置。对于2017 EE、2017 SE、2017 Web版，5050、5353和5986不可设置。  当不传该参数时，默认端口如下：  - MySQL默认3306。 - PostgreSQL默认5432。 - Microsoft SQL Server默认1433。
	Port *string `json:"port,omitempty"`

	// 数据库密码。创建只读实例时不可选，其它场景必选。  取值范围：  非空，由大小写字母、数字和特殊符号~!@#%^*-_=+?组成，长度8~32个字符。  其中允许的特殊字符如下： - MySQL数据库允许输入~!@#$%^*-_=+?,特殊字符。 - SQL Server数据库允许输入~!@#$%^*-_+?,特殊字符。 - PostgreSQL数据库允许输入~!@#%^*-_=+?,特殊字符。  建议您输入高强度密码，以提高安全性，防止出现密码被暴力破解等安全风险。
	Password *string `json:"password,omitempty"`

	BackupStrategy *BackupStrategy `json:"backup_strategy,omitempty"`

	// 企业项目ID。 使用请参考《企业管理 API参考》的[查询企业项目列表](https://support.huaweicloud.com/api-em/zh-cn_topic_0121230880.html)响应消息表“enterprise_project字段数据结构说明”的“id”。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 用于磁盘加密的密钥ID，默认为空。
	DiskEncryptionId *string `json:"disk_encryption_id,omitempty"`

	// 规格码,取值范围：非空。 使用[查询数据库规格](https://support.huaweicloud.com/api-rds/rds_06_0002.html)接口响应消息中的 flavors字段中“spec_code”获取且对应az_status为“在售”状态。
	FlavorRef string `json:"flavor_ref"`

	Volume *Volume `json:"volume"`

	// 区域ID。创建主实例时必选，其它场景不可选。 取值参见[地区和终端节点](https://developer.huaweicloud.com/endpoint)。
	Region string `json:"region"`

	// 可用区ID。对于数据库实例类型不是单机的实例，需要分别为实例所有节点指定可用区，并用逗号隔开。 取值参见[地区和终端节点](https://developer.huaweicloud.com/endpoint)。
	AvailabilityZone string `json:"availability_zone"`

	// 虚拟私有云ID。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。 获取方法如下： - 方法1：登录虚拟私有云服务的控制台界面，在虚拟私有云的详情页面查找VPC ID。 - 方法2：通过虚拟私有云服务的API接口查询，具体操作可参考[查询VPC列表](https://support.huaweicloud.com/api-vpc/vpc_api01_0003.html)。
	VpcId string `json:"vpc_id"`

	// 子网的网络ID信息。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。 获取方法如下： - 方法1：登录虚拟私有云服务的控制台界面，单击VPC下的子网，进入子网详情页面，查找网络ID。 - 方法2：通过虚拟私有云服务的API接口查询，具体操作可参考[查询子网列表](https://support.huaweicloud.com/api-vpc/vpc_subnet01_0003.html)。
	SubnetId string `json:"subnet_id"`

	// 指定实例的内网IP,目前仅支持设置IPv4地址。 获取方法如下： - 方法1：登录虚拟私有云服务的控制台界面，单击VPC下的子网，进入子网详情页面，查找子网的网段，选择未被占用的IP。 - 方法2：通过虚拟私有云服务的API接口查询，具体操作可参考[查询私有IP列表](https://support.huaweicloud.com/api-vpc/vpc_privateip_0003.html),选择“device_owner”为空的私有IP。
	DataVip *string `json:"data_vip,omitempty"`

	// 安全组ID。创建只读实例时不可选（只读实例的网络属性默认和主实例相同），其它场景必选。 获取方法如下： - 方法1：登录虚拟私有云服务的控制台界面，在安全组的详情页面查找安全组ID。 - 方法2：通过虚拟私有云服务的API接口查询，具体操作可参考[查询安全组列表](https://support.huaweicloud.com/api-vpc/vpc_sg01_0003.html)。
	SecurityGroupId string `json:"security_group_id"`

	ChargeInfo *ChargeInfo `json:"charge_info,omitempty"`

	// 时区。  - 不选择时，各个引擎时区如下：   - MySQL国内站、国际站默认为UTC时间。   - PostgreSQL国内站、国际站默认为UTC时间。   - Microsoft SQL Server国内站默认为China Standard Time，国际站默认为UTC时间。 - MySQL和PostgreSQL引擎选择填写时，取值范围为UTC-12:00~UTC+12:00，且只支持整段时间，如UTC+08:00，不支持UTC+08:30。 - Microsoft SQL Server引擎选择填写时，请参见“创建数据库实例”接口[表10](https://support.huaweicloud.com/api-rds/rds_01_0002.html#rds_01_0002__table613473883617) 时区与UTC偏移量对照表，填写时区列字符串，如：China Standard Time。
	TimeZone *string `json:"time_zone,omitempty"`

	// Dec用户专属存储ID，每个az配置的专属存储不同，Dec用户创建实例时，对于数据库实例类型不是单机或只读的实例，需要分别为实例所有节点指定dsspoolId，并用逗号隔开。 获取方法如下： - 方法1：登录专属分布式存储服务DSS的控制台界面，查看专属存储列表，选择符合条件的az下的专属dss的ID。 - 方法2：通过专属分布式存储服务DSS的API接口查询，具体操作可参考[获取专属存储详情列表](https://support.huaweicloud.com/api-dss/dss_02_1002.html)。
	DsspoolId *string `json:"dsspool_id,omitempty"`

	// 只读实例的主实例ID。创建只读实例时必选，其它场景不可选。
	ReplicaOfId *string `json:"replica_of_id,omitempty"`

	RestorePoint *RestorePoint `json:"restore_point,omitempty"`

	// 仅限Microsoft SQL Server实例创建使用。对于MySQL和PostgreSQL实例，该参数无意义。取值范围：根据查询SQL Server可用字符集的字符集查询列表查询可设置的字符集。 取值范围：根据[查询SQL Server可用字符集](https://support.huaweicloud.com/api-rds/rds_05_0010.html)查询可设置的字符集。
	Collation *string `json:"collation,omitempty"`

	// 标签列表。单个实例总标签数上限20个。
	Tags *[]TagWithKeyValue `json:"tags,omitempty"`

	UnchangeableParam *UnchangeableParam `json:"unchangeable_param,omitempty"`

	// 是否只预检此次请求，仅支持MySQL。 - true：发送参数检查请求，不会创建实例。   - 检查通过：返回202状态码。   - 检查不通过：返回对应错误码，详情请参考错误码。 - false：发送正常请求，通过检查后，并且执行创建实例的请求。
	DryRun *bool `json:"dry_run,omitempty"`
}

func (o CreateRestoreInstanceRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRestoreInstanceRequestBody struct{}"
	}

	return strings.Join([]string{"CreateRestoreInstanceRequestBody", string(data)}, " ")
}
