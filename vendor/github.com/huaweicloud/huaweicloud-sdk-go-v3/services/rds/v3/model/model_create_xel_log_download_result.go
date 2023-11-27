package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateXelLogDownloadResult struct {

	// 下载文件名
	FileName string `json:"file_name"`

	// 生成链接的生成状态。FINISH，表示下载链接已经生成完成。EXPORTING，，表示正在生成文件。FAILED，表示存在日志文件准备失败。
	Status string `json:"status"`

	// 日志大小，单位：KB
	FileSize string `json:"file_size"`

	// 下载链接,链接的生成状态为EXPORTING，或者FAILED不予返回
	FileLink string `json:"file_link"`

	// 生成时间
	CreateAt string `json:"create_at"`

	// 更新时间
	UpdateAt string `json:"update_at"`
}

func (o CreateXelLogDownloadResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateXelLogDownloadResult struct{}"
	}

	return strings.Join([]string{"CreateXelLogDownloadResult", string(data)}, " ")
}
