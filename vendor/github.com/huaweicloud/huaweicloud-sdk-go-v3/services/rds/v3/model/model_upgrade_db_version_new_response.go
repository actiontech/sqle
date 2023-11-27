package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDbVersionNewResponse Response Object
type UpgradeDbVersionNewResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpgradeDbVersionNewResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbVersionNewResponse struct{}"
	}

	return strings.Join([]string{"UpgradeDbVersionNewResponse", string(data)}, " ")
}
