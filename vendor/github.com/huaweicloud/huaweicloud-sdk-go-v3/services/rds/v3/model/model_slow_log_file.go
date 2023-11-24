package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SlowLogFile 慢日志信息。
type SlowLogFile struct {

	// 文件名。
	FileName string `json:"file_name"`

	// 文件大小（单位Byte）
	FileSize string `json:"file_size"`
}

func (o SlowLogFile) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlowLogFile struct{}"
	}

	return strings.Join([]string{"SlowLogFile", string(data)}, " ")
}
