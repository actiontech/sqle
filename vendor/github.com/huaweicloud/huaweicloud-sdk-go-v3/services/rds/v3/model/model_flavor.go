package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Flavor 规格信息。
type Flavor struct {

	// 规格id
	Id string `json:"id"`

	// CPU个数。
	Vcpus string `json:"vcpus"`

	// 内存大小，单位为GB。
	Ram int32 `json:"ram"`

	// 资源规格编码。例如：rds.mysql.m1.xlarge.rr。  其中形如“xxx.xxx.mcs.i3.xxx.xxx.xxx”是超高性能型（尊享版），需要申请一定权限才可使用，更多规格说明请参考数据库实例规格。 - “rds”代表RDS产品。 - “mysql”代表数据库引擎。 - “m1.xlarge”代表性能规格，为高内存类型。
	SpecCode string `json:"spec_code"`

	// 实例模型，包括如下类型： - ha，主备实例。 - replica，只读实例。 - single，单实例。
	InstanceMode string `json:"instance_mode"`

	// 其中key是可用区编号，value是规格所在az的状态，包含以下状态： - normal，在售。 - unsupported，暂不支持该规格。 - sellout，售罄。
	AzStatus map[string]string `json:"az_status"`

	// 规格所在az的描述。
	AzDesc map[string]string `json:"az_desc"`

	// 数组形式版本号
	VersionName []string `json:"version_name"`

	// 性能规格，包含以下状态： - normal：通用增强型。 - normal2：通用增强Ⅱ型。 - armFlavors：鲲鹏通用增强型。 - dedicicatenormal ：x86独享型。 - armlocalssd：鲲鹏通用型。 - normallocalssd：x86通用型。 - general：通用型。 - dedicated：独享型，仅云盘SSD支持。 - rapid：独享型，仅极速型SSD支持。 - bigmen：超大内存型。
	GroupType string `json:"group_type"`
}

func (o Flavor) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Flavor struct{}"
	}

	return strings.Join([]string{"Flavor", string(data)}, " ")
}
