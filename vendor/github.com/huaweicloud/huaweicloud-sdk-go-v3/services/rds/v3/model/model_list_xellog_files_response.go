package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListXellogFilesResponse Response Object
type ListXellogFilesResponse struct {

	// 扩展日志文件返回体
	List *[]ListXelLogResponseResult `json:"list,omitempty"`

	// 扩展日志文件数量
	Count          *int32 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListXellogFilesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListXellogFilesResponse struct{}"
	}

	return strings.Join([]string{"ListXellogFilesResponse", string(data)}, " ")
}
