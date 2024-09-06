package state

import (
	"time"

	"github.com/guregu/null"
)

type PostgresReplication struct {
	InRecovery bool

	// Data available on primary
	CurrentXlogLocation null.String
	Standbys            []PostgresReplicationStandby

	// Data available on standby
	IsStreaming        null.Bool
	ReceiveLocation    null.String
	ReplayLocation     null.String
	ApplyByteLag       null.Int
	ReplayTimestamp    null.Time
	ReplayTimestampAge null.Int
}

// PostgresReplicationStandby - Standby information as seen from the primary
type PostgresReplicationStandby struct {
	ClientAddr string

	RoleOid         Oid
	Pid             int64
	ApplicationName string
	ClientHostname  null.String
	ClientPort      int32
	BackendStart    time.Time
	SyncPriority    int32
	SyncState       string

	State          string
	SentLocation   null.String
	WriteLocation  null.String
	FlushLocation  null.String
	ReplayLocation null.String
	RemoteByteLag  null.Int
	LocalByteLag   null.Int
}
