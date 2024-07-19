package state

import "github.com/guregu/null"

// PostgresVacuumStatsEntry - One entry in the VACUUM statistics
type PostgresVacuumStatsEntry struct {
	SchemaName   string
	RelationName string

	LiveRowCount int32
	DeadRowCount int32
	Relfrozenxid int32
	Relminmxid   int32

	LastManualVacuumRun  null.Time
	LastAutoVacuumRun    null.Time
	LastManualAnalyzeRun null.Time
	LastAutoAnalyzeRun   null.Time

	AutovacuumEnabled               bool
	AutovacuumVacuumThreshold       int32
	AutovacuumAnalyzeThreshold      int32
	AutovacuumVacuumScaleFactor     float64
	AutovacuumAnalyzeScaleFactor    float64
	AutovacuumFreezeMaxAge          int32
	AutovacuumMultixactFreezeMaxAge int32
	AutovacuumVacuumCostDelay       int32
	AutovacuumVacuumCostLimit       int32

	Fillfactor int32
}

// PostgresVacuumStats - Details on VACUUM configuration and expected runs
type PostgresVacuumStats struct {
	DatabaseName string

	// Database-wide settings
	AutovacuumMaxWorkers     int32
	AutovacuumNaptimeSeconds int32

	// Defaults for per-table settings
	AutovacuumEnabled               bool
	AutovacuumVacuumThreshold       int32
	AutovacuumAnalyzeThreshold      int32
	AutovacuumVacuumScaleFactor     float64
	AutovacuumAnalyzeScaleFactor    float64
	AutovacuumFreezeMaxAge          int32
	AutovacuumMultixactFreezeMaxAge int32
	AutovacuumVacuumCostDelay       int32
	AutovacuumVacuumCostLimit       int32

	Relations []PostgresVacuumStatsEntry
}
