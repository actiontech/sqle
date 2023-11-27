package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddPostgresqlHbaConfRequest Request Object
type AddPostgresqlHbaConfRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	Body *[]PostgresqlHbaConf `json:"body,omitempty"`
}

func (o AddPostgresqlHbaConfRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddPostgresqlHbaConfRequest struct{}"
	}

	return strings.Join([]string{"AddPostgresqlHbaConfRequest", string(data)}, " ")
}
