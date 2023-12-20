package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlExtensionResponse Response Object
type ListPostgresqlExtensionResponse struct {
	Extensions *[]ExtensionsResponse `json:"extensions,omitempty"`

	// 总插件数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListPostgresqlExtensionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlExtensionResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlExtensionResponse", string(data)}, " ")
}
