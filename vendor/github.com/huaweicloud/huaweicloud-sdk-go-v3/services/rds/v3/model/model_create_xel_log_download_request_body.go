package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateXelLogDownloadRequestBody struct {

	// 文件名。取值范围：不为null和空字符串，只为大小写字母，数字和下划线，以“.xel”结尾
	FileName string `json:"file_name"`
}

func (o CreateXelLogDownloadRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateXelLogDownloadRequestBody struct{}"
	}

	return strings.Join([]string{"CreateXelLogDownloadRequestBody", string(data)}, " ")
}
