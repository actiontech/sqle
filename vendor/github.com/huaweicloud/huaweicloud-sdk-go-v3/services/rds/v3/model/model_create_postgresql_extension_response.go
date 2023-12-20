package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostgresqlExtensionResponse Response Object
type CreatePostgresqlExtensionResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreatePostgresqlExtensionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostgresqlExtensionResponse struct{}"
	}

	return strings.Join([]string{"CreatePostgresqlExtensionResponse", string(data)}, " ")
}
