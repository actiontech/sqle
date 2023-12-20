package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostgresqlPreCheckUpgradeMajorVersionReq struct {

	// 目标版本。
	TargetVersion string `json:"target_version"`
}

func (o PostgresqlPreCheckUpgradeMajorVersionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlPreCheckUpgradeMajorVersionReq struct{}"
	}

	return strings.Join([]string{"PostgresqlPreCheckUpgradeMajorVersionReq", string(data)}, " ")
}
