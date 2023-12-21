package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddPostgresqlHbaConfResponse Response Object
type AddPostgresqlHbaConfResponse struct {

	// 结果码
	Code *string `json:"code,omitempty"`

	// 结果描述
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddPostgresqlHbaConfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddPostgresqlHbaConfResponse struct{}"
	}

	return strings.Join([]string{"AddPostgresqlHbaConfResponse", string(data)}, " ")
}
