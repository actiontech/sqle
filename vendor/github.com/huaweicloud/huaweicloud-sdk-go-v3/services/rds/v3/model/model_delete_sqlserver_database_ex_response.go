package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlserverDatabaseExResponse Response Object
type DeleteSqlserverDatabaseExResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteSqlserverDatabaseExResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlserverDatabaseExResponse struct{}"
	}

	return strings.Join([]string{"DeleteSqlserverDatabaseExResponse", string(data)}, " ")
}
