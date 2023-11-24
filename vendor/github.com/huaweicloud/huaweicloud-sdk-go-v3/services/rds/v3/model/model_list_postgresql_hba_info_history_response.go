package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlHbaInfoHistoryResponse Response Object
type ListPostgresqlHbaInfoHistoryResponse struct {
	Body           *[]PostgresqlHbaHistory `json:"body,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListPostgresqlHbaInfoHistoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlHbaInfoHistoryResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlHbaInfoHistoryResponse", string(data)}, " ")
}
