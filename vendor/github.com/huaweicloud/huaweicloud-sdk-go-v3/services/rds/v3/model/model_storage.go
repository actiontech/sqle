package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Storage 实例磁盘类型信息。
type Storage struct {

	// 磁盘类型名称，可能取值如下： - ULTRAHIGH：表示SSD。 - LOCALSSD：表示本地SSD。 - CLOUDSSD：表示SSD云盘，仅支持通用型和独享型规格实例。 - ESSD：表示极速型SSD，仅支持独享型规格实例。
	Name string `json:"name"`

	// 其中key是可用区编号，value是规格所在az的状态，包含以下状态： - normal，在售。 - unsupported，暂不支持该规格。 - sellout，售罄。
	AzStatus map[string]string `json:"az_status"`

	// 性能规格，包含以下状态： - normal：通用增强型。 - normal2：通用增强Ⅱ型。 - armFlavors：鲲鹏通用增强型。 - dedicicatenormal ：x86独享型。 - armlocalssd：鲲鹏通用型。 - normallocalssd：x86通用型。 - general：通用型。 - dedicated：独享型，仅云盘SSD支持。 - rapid：独享型，仅极速型SSD支持。 - bigmen：超大内存型。
	SupportComputeGroupType *[]string `json:"support_compute_group_type,omitempty"`
}

func (o Storage) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Storage struct{}"
	}

	return strings.Join([]string{"Storage", string(data)}, " ")
}
