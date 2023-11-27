package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlHbaInfoRequest Request Object
type ListPostgresqlHbaInfoRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`
}

func (o ListPostgresqlHbaInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlHbaInfoRequest struct{}"
	}

	return strings.Join([]string{"ListPostgresqlHbaInfoRequest", string(data)}, " ")
}
