package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyPostgresqlHbaConfRequest Request Object
type ModifyPostgresqlHbaConfRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	Body *[]PostgresqlHbaConf `json:"body,omitempty"`
}

func (o ModifyPostgresqlHbaConfRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyPostgresqlHbaConfRequest struct{}"
	}

	return strings.Join([]string{"ModifyPostgresqlHbaConfRequest", string(data)}, " ")
}
