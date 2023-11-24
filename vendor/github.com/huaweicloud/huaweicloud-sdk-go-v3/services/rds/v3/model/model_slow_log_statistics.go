package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SlowLogStatistics 慢日志信息。
type SlowLogStatistics struct {

	// 执行次数。
	Count string `json:"count"`

	// 平均执行时间。
	Time string `json:"time"`

	// 平均等待锁时间。
	LockTime string `json:"lockTime"`

	// 平均结果行数量。
	RowsSent int64 `json:"rowsSent"`

	// 平均扫描的行数量。
	RowsExamined int64 `json:"rowsExamined"`

	// 所属数据库。
	Database string `json:"database"`

	// 帐号。
	Users string `json:"users"`

	// 执行语法。
	QuerySample string `json:"querySample"`

	// 语句类型。
	Type string `json:"type"`

	// IP地址。
	ClientIP string `json:"clientIP"`
}

func (o SlowLogStatistics) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlowLogStatistics struct{}"
	}

	return strings.Join([]string{"SlowLogStatistics", string(data)}, " ")
}
