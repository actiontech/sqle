package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlHbaInfoResponse Response Object
type ListPostgresqlHbaInfoResponse struct {
	Body           *[]PostgresqlHbaConf `json:"body,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListPostgresqlHbaInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlHbaInfoResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlHbaInfoResponse", string(data)}, " ")
}
