package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListXelLogResponseResult 扩展日志信息
type ListXelLogResponseResult struct {

	// 文件名
	FileName string `json:"file_name"`

	// 日志大小，单位：KB
	FileSize string `json:"file_size"`
}

func (o ListXelLogResponseResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListXelLogResponseResult struct{}"
	}

	return strings.Join([]string{"ListXelLogResponseResult", string(data)}, " ")
}
