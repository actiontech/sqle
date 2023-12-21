package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateXelLogDownloadResponse Response Object
type CreateXelLogDownloadResponse struct {

	// 扩展日志文件返回实体
	List *[]CreateXelLogDownloadResult `json:"list,omitempty"`

	// 扩展日志文件信息数量。
	Count          *int32 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateXelLogDownloadResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateXelLogDownloadResponse struct{}"
	}

	return strings.Join([]string{"CreateXelLogDownloadResponse", string(data)}, " ")
}
