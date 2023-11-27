package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DssPoolInfo 实例专属存储信息。
type DssPoolInfo struct {

	// 专属存储池所在az
	AzName string `json:"az_name"`

	// 专属存储池免费空间大小，单位GB
	FreeCapacityGb string `json:"free_capacity_gb"`

	// 专属存储池磁盘类型名称，可能取值如下：  - ULTRAHIGH，表示SSD。  - LOCALSSD：表示本地SSD。  - CLOUDSSD：表示SSD云盘，仅支持通用型和独享型规格实例。  - ESSD：表示极速型SSD，仅支持独享型规格实例。
	DsspoolVolumeType string `json:"dsspool_volume_type"`

	// 专属存储池ID
	DsspoolId string `json:"dsspool_id"`

	// 专属存储池当前状态，可能取值如下： - available，表示可用。 - deploying，表示正在部署。 - enlarging，表示正在扩容。 - frozen，表示冻结。 - sellout，表示售罄。
	DsspoolStatus string `json:"dsspool_status"`
}

func (o DssPoolInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DssPoolInfo struct{}"
	}

	return strings.Join([]string{"DssPoolInfo", string(data)}, " ")
}
