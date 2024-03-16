package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchRestoreDatabaseResponse Response Object
type BatchRestoreDatabaseResponse struct {

	// 表信息
	RestoreResult  *[]PostgreSqlRestoreResult `json:"restore_result,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o BatchRestoreDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestoreDatabaseResponse struct{}"
	}

	return strings.Join([]string{"BatchRestoreDatabaseResponse", string(data)}, " ")
}
