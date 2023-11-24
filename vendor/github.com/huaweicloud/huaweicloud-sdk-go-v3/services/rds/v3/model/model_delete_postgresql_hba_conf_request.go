package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlHbaConfRequest Request Object
type DeletePostgresqlHbaConfRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	Body *[]PostgresqlHbaConf `json:"body,omitempty"`
}

func (o DeletePostgresqlHbaConfRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlHbaConfRequest struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlHbaConfRequest", string(data)}, " ")
}
