package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SlowLog 慢日志信息。
type SlowLog struct {

	// 执行次数。
	Count string `json:"count"`

	// 平均执行时间。
	Time string `json:"time"`

	// 平均等待锁时间。
	LockTime string `json:"lock_time"`

	// 平均结果行数量。
	RowsSent string `json:"rows_sent"`

	// 平均扫描的行数量。
	RowsExamined string `json:"rows_examined"`

	// 所属数据库。
	Database string `json:"database"`

	// 帐号。
	Users string `json:"users"`

	// 执行语法。
	QuerySample string `json:"query_sample"`

	// 语句类型。
	Type string `json:"type"`

	// 发生时间，UTC时间。
	StartTime string `json:"start_time"`

	// IP地址。
	ClientIp string `json:"client_ip"`
}

func (o SlowLog) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlowLog struct{}"
	}

	return strings.Join([]string{"SlowLog", string(data)}, " ")
}
