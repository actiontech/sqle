package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchRestoreDatabaseRequest Request Object
type BatchRestoreDatabaseRequest struct {
	Body *PostgreSqlRestoreDatabaseRequest `json:"body,omitempty"`
}

func (o BatchRestoreDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestoreDatabaseRequest struct{}"
	}

	return strings.Join([]string{"BatchRestoreDatabaseRequest", string(data)}, " ")
}
