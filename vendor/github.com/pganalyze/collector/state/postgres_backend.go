package state

import "github.com/guregu/null"

// PostgresBackend - PostgreSQL server backend thats currently working, waiting
// or idling (also known as an open connection)
//
// See https://www.postgresql.org/docs/9.5/static/monitoring-stats.html#PG-STAT-ACTIVITY-VIEW
type PostgresBackend struct {
	Identity        uint64      // Combination of process start time and PID, used to identify a process over time
	DatabaseOid     null.Int    // OID of the database this backend is connected to
	DatabaseName    null.String // Name of the database this backend is connected to
	RoleOid         null.Int    // OID of the user logged into this backend
	RoleName        null.String // Name of the user logged into this backend
	Pid             int32       // Process ID of this backend
	ApplicationName null.String // Name of the application that is connected to this backend
	ClientAddr      null.String // IP address of the client connected to this backend. If this field is null, it indicates either that the client is connected via a Unix socket on the server machine or that this is an internal process such as autovacuum.
	ClientPort      null.Int    // TCP port number that the client is using for communication with this backend, or -1 if a Unix socket is used
	BackendStart    null.Time   // Time when this process was started, i.e., when the client connected to the server
	XactStart       null.Time   // Time when this process' current transaction was started, or null if no transaction is active. If the current query is the first of its transaction, this column is equal to the query_start column.
	QueryStart      null.Time   // Time when the currently active query was started, or if state is not active, when the last query was started
	StateChange     null.Time   // Time when the state was last changed
	Waiting         null.Bool   // True if this backend is currently waiting on a lock
	BackendXid      null.Int    // Top-level transaction identifier of this backend, if any.
	BackendXmin     null.Int    // The current backend's xmin horizon.

	WaitEventType null.String // 9.6+ The type of event for which the backend is waiting, if any; otherwise NULL
	WaitEvent     null.String // 9.6+ Wait event name if backend is currently waiting, otherwise NULL

	BackendType null.String // 10+ The process type of this backend

	Query null.String // Text of this backend's most recent query

	// Current overall state of this backend. Possible values are:
	// - active: The backend is executing a query.
	// - idle: The backend is waiting for a new client command.
	// - idle in transaction: The backend is in a transaction, but is not currently executing a query.
	// - idle in transaction (aborted): This state is similar to idle in transaction, except one of the statements in the transaction caused an error.
	// - fastpath function call: The backend is executing a fast-path function.
	// - disabled: This state is reported if track_activities is disabled in this backend.
	State null.String

	BlockedByPids []int32 // The list of PIDs this backend is blocked by
}

type PostgresBackendCount struct {
	DatabaseOid    null.Int // OID of the database
	RoleOid        null.Int // OID of the user
	State          string   // Current overall state of this backend
	BackendType    string   // The process type of this backend
	WaitingForLock bool     // True if this backend is currently waiting on a heavyweight lock
	Count          int32    // Number of this kind of backends
}
