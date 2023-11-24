package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOffSiteInstancesResponse Response Object
type ListOffSiteInstancesResponse struct {

	// 跨区域备份实例信息。
	OffsiteBackupInstances *[]OffsiteBackupInstance `json:"offsite_backup_instances,omitempty"`

	// 总记录数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListOffSiteInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOffSiteInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListOffSiteInstancesResponse", string(data)}, " ")
}
