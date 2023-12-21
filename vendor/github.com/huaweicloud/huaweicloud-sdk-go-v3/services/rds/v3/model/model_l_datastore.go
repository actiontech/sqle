package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// LDatastore 数据库版本信息。
type LDatastore struct {

	// 数据库版本ID。
	Id string `json:"id"`

	// 数据库版本号。 - 对于MySQL引擎可以返回小版本号，例如MySQL 5.6.51版本，将返回5.6.51。 - 对于PostgreSQL和SQL Server引擎，只返回两位数的大版本号，例如PostgreSQL 9.6.X版本，仅返回9.6。
	Name string `json:"name"`
}

func (o LDatastore) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LDatastore struct{}"
	}

	return strings.Join([]string{"LDatastore", string(data)}, " ")
}
