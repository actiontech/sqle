package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstancesSupportFastRestoreResponse Response Object
type ListInstancesSupportFastRestoreResponse struct {

	// 实例的极速恢复支持情况。
	SupportFastRestoreList *[]SupportFastRestoreList `json:"support_fast_restore_list,omitempty"`
	HttpStatusCode         int                       `json:"-"`
}

func (o ListInstancesSupportFastRestoreResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesSupportFastRestoreResponse struct{}"
	}

	return strings.Join([]string{"ListInstancesSupportFastRestoreResponse", string(data)}, " ")
}
