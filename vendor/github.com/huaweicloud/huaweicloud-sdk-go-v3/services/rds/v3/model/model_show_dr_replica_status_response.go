package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDrReplicaStatusResponse Response Object
type ShowDrReplicaStatusResponse struct {

	// 同步状态，取值范围是0或-1，0表示正常，-1表示异常。
	ReplicaState *string `json:"replica_state,omitempty"`

	// 发送延迟大小（MB），即主实例当前wal日志写入位点与灾备实例当前接收wal日志位点的差值。
	WalWriteReceiveDelayInMb *string `json:"wal_write_receive_delay_in_mb,omitempty"`

	// 端到端延迟大小（MB），即主实例当前wal日志写入位点与灾备实例当前回放wal日志位点的差值。
	WalWriteReplayDelayInMb *string `json:"wal_write_replay_delay_in_mb,omitempty"`

	// 回放延迟时间（ms），即数据在灾备上回放的延迟时间。
	WalReceiveReplayDelayInMs *string `json:"wal_receive_replay_delay_in_ms,omitempty"`
	HttpStatusCode            int     `json:"-"`
}

func (o ShowDrReplicaStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDrReplicaStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowDrReplicaStatusResponse", string(data)}, " ")
}
