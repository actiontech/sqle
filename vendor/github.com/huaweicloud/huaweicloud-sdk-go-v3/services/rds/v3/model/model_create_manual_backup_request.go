package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateManualBackupRequest Request Object
type CreateManualBackupRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *CreateManualBackupRequestBody `json:"body,omitempty"`
}

func (o CreateManualBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateManualBackupRequest struct{}"
	}

	return strings.Join([]string{"CreateManualBackupRequest", string(data)}, " ")
}
