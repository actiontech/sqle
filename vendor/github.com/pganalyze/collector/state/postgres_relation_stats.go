package state

import "github.com/guregu/null"

type PostgresRelationStats struct {
	SizeBytes         int64     // On-disk size including FSM and VM, plus TOAST table if any, excluding indices
	ToastSizeBytes    int64     // TOAST table and TOAST index size (included in SizeBytes as well)
	SeqScan           int64     // Number of sequential scans initiated on this table
	SeqTupRead        int64     // Number of live rows fetched by sequential scans
	IdxScan           int64     // Number of index scans initiated on this table
	IdxTupFetch       int64     // Number of live rows fetched by index scans
	NTupIns           int64     // Number of rows inserted
	NTupUpd           int64     // Number of rows updated
	NTupDel           int64     // Number of rows deleted
	NTupHotUpd        int64     // Number of rows HOT updated (i.e., with no separate index update required)
	NLiveTup          int64     // Estimated number of live rows
	NDeadTup          int64     // Estimated number of dead rows
	NModSinceAnalyze  int64     // Estimated number of rows modified since this table was last analyzed
	NInsSinceVacuum   int64     // Estimated number of rows inserted since this table was last vacuumed
	LastVacuum        null.Time // Last time at which this table was manually vacuumed (not counting VACUUM FULL)
	LastAutovacuum    null.Time // Last time at which this table was vacuumed by the autovacuum daemon
	LastAnalyze       null.Time // Last time at which this table was manually analyzed
	LastAutoanalyze   null.Time // Last time at which this table was analyzed by the autovacuum daemon
	VacuumCount       int64     // Number of times this table has been manually vacuumed (not counting VACUUM FULL)
	AutovacuumCount   int64     // Number of times this table has been vacuumed by the autovacuum daemon
	AnalyzeCount      int64     // Number of times this table has been manually analyzed
	AutoanalyzeCount  int64     // Number of times this table has been analyzed by the autovacuum daemon
	HeapBlksRead      int64     // Number of disk blocks read from this table
	HeapBlksHit       int64     // Number of buffer hits in this table
	IdxBlksRead       int64     // Number of disk blocks read from all indexes on this table
	IdxBlksHit        int64     // Number of buffer hits in all indexes on this table
	ToastBlksRead     int64     // Number of disk blocks read from this table's TOAST table (if any)
	ToastBlksHit      int64     // Number of buffer hits in this table's TOAST table (if any)
	TidxBlksRead      int64     // Number of disk blocks read from this table's TOAST table indexes (if any)
	TidxBlksHit       int64     // Number of buffer hits in this table's TOAST table indexes (if any)
	FrozenXIDAge      int32     // Age of frozen XID for this table
	MinMXIDAge        int32     // Age of minimum multixact ID for this table
	Relpages          int32     // Size of the on-disk representation of this table in pages (of size BLCKSZ)
	Reltuples         float32   // Number of live rows in the table. -1 indicating that the row count is unknown
	Relallvisible     int32     // Number of pages that are marked all-visible in the table's visibility map
	ExclusivelyLocked bool      // Whether these statistics are zeroed out because the table was locked at collection time
	ToastReltuples    float32   // Number of live rows in the TOAST table. -1 indicating that the row count is unknown
	ToastRelpages     int32     // Size of the on-disk representation of the TOAST table in pages (of size BLCKSZ)
}

type PostgresIndexStats struct {
	SizeBytes         int64
	IdxScan           int64 // Number of index scans initiated on this index
	IdxTupRead        int64 // Number of index entries returned by scans on this index
	IdxTupFetch       int64 // Number of live table rows fetched by simple index scans using this index
	IdxBlksRead       int64 // Number of disk blocks read from this index
	IdxBlksHit        int64 // Number of buffer hits in this index
	ExclusivelyLocked bool  // Whether these statistics are zeroed out because the index was locked at collection time
}

type PostgresColumnStats struct {
	SchemaName  string
	TableName   string
	ColumnName  string
	Inherited   bool
	NullFrac    float64
	AvgWidth    int32
	NDistinct   float64
	Correlation null.Float
}

// PostgresColumnStatsKey - Information that uniquely identifies column stats
type PostgresColumnStatsKey struct {
	SchemaName string
	TableName  string
	ColumnName string
}

type PostgresRelationStatsMap map[Oid]PostgresRelationStats
type PostgresIndexStatsMap map[Oid]PostgresIndexStats
type PostgresColumnStatsMap map[PostgresColumnStatsKey][]PostgresColumnStats

type DiffedPostgresRelationStats PostgresRelationStats
type DiffedPostgresIndexStats PostgresIndexStats
type DiffedPostgresRelationStatsMap map[Oid]DiffedPostgresRelationStats
type DiffedPostgresIndexStatsMap map[Oid]DiffedPostgresIndexStats

func (curr PostgresRelationStats) DiffSince(prev PostgresRelationStats) DiffedPostgresRelationStats {
	return DiffedPostgresRelationStats{
		SizeBytes:        curr.SizeBytes,
		ToastSizeBytes:   curr.ToastSizeBytes,
		SeqScan:          curr.SeqScan - prev.SeqScan,
		SeqTupRead:       curr.SeqTupRead - prev.SeqTupRead,
		IdxScan:          curr.IdxScan - prev.IdxScan,
		IdxTupFetch:      curr.IdxTupFetch - prev.IdxTupFetch,
		NTupIns:          curr.NTupIns - prev.NTupIns,
		NTupUpd:          curr.NTupUpd - prev.NTupUpd,
		NTupDel:          curr.NTupDel - prev.NTupDel,
		NTupHotUpd:       curr.NTupHotUpd - prev.NTupHotUpd,
		NLiveTup:         curr.NLiveTup,
		NDeadTup:         curr.NDeadTup,
		NModSinceAnalyze: curr.NModSinceAnalyze,
		NInsSinceVacuum:  curr.NInsSinceVacuum,
		LastVacuum:       curr.LastVacuum,
		LastAutovacuum:   curr.LastAutovacuum,
		LastAnalyze:      curr.LastAnalyze,
		LastAutoanalyze:  curr.LastAutoanalyze,
		VacuumCount:      curr.VacuumCount - prev.VacuumCount,
		AutovacuumCount:  curr.AutovacuumCount - prev.AutovacuumCount,
		AnalyzeCount:     curr.AnalyzeCount - prev.AnalyzeCount,
		AutoanalyzeCount: curr.AutoanalyzeCount - prev.AutoanalyzeCount,
		HeapBlksRead:     curr.HeapBlksRead - prev.HeapBlksRead,
		HeapBlksHit:      curr.HeapBlksHit - prev.HeapBlksHit,
		IdxBlksRead:      curr.IdxBlksRead - prev.IdxBlksRead,
		IdxBlksHit:       curr.IdxBlksHit - prev.IdxBlksHit,
		ToastBlksRead:    curr.ToastBlksRead - prev.ToastBlksRead,
		ToastBlksHit:     curr.ToastBlksHit - prev.ToastBlksHit,
		TidxBlksRead:     curr.TidxBlksRead - prev.TidxBlksRead,
		TidxBlksHit:      curr.TidxBlksHit - prev.TidxBlksHit,
		FrozenXIDAge:     curr.FrozenXIDAge,
		MinMXIDAge:       curr.MinMXIDAge,
		Relpages:         curr.Relpages,
		Reltuples:        curr.Reltuples,
		Relallvisible:    curr.Relallvisible,
		ToastReltuples:   curr.ToastReltuples,
		ToastRelpages:    curr.ToastRelpages,
	}
}

func (curr PostgresIndexStats) DiffSince(prev PostgresIndexStats) DiffedPostgresIndexStats {
	return DiffedPostgresIndexStats{
		SizeBytes:   curr.SizeBytes,
		IdxScan:     curr.IdxScan - prev.IdxScan,
		IdxTupRead:  curr.IdxTupRead - prev.IdxTupRead,
		IdxTupFetch: curr.IdxTupFetch - prev.IdxTupFetch,
		IdxBlksRead: curr.IdxBlksRead - prev.IdxBlksRead,
		IdxBlksHit:  curr.IdxBlksHit - prev.IdxBlksHit,
	}
}
