package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteManualBackupRequest Request Object
type BatchDeleteManualBackupRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *BatchDeleteBackupRequestBody `json:"body,omitempty"`
}

func (o BatchDeleteManualBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteManualBackupRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteManualBackupRequest", string(data)}, " ")
}
