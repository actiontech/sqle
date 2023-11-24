package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOffSiteBackupsResponse Response Object
type ListOffSiteBackupsResponse struct {

	// 跨区域备份信息。
	Backups *[]OffSiteBackupForList `json:"backups,omitempty"`

	// 总记录数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListOffSiteBackupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOffSiteBackupsResponse struct{}"
	}

	return strings.Join([]string{"ListOffSiteBackupsResponse", string(data)}, " ")
}
