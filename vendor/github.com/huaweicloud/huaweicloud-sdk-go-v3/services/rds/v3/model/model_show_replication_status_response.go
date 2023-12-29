package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReplicationStatusResponse Response Object
type ShowReplicationStatusResponse struct {

	// 复制状态。
	ReplicationStatus *string `json:"replication_status,omitempty"`

	// 复制异常原因。
	AbnormalReason *string `json:"abnormal_reason,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowReplicationStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReplicationStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowReplicationStatusResponse", string(data)}, " ")
}
