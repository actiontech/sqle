package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateManualBackupResponse Response Object
type CreateManualBackupResponse struct {
	Backup         *BackupInfo `json:"backup,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o CreateManualBackupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateManualBackupResponse struct{}"
	}

	return strings.Join([]string{"CreateManualBackupResponse", string(data)}, " ")
}
