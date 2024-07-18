package state

import "time"

// PostgresVacuumProgress - PostgreSQL vacuum thats currently running
//
// See https://www.postgresql.org/docs/10/static/progress-reporting.html
type PostgresVacuumProgress struct {
	VacuumIdentity  uint64 // Combination of vacuum "query" start time and PID, used to identify a vacuum over time
	BackendIdentity uint64 // Combination of process start time and PID, used to identify a process over time

	DatabaseName string
	SchemaName   string
	RelationName string
	RoleName     string
	StartedAt    time.Time
	Autovacuum   bool
	Toast        bool

	Phase            string
	HeapBlksTotal    int64
	HeapBlksScanned  int64
	HeapBlksVacuumed int64
	IndexVacuumCount int64
	MaxDeadTuples    int64
	NumDeadTuples    int64
}
