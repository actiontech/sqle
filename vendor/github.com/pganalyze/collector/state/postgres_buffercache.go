package state

// PostgresBuffercacheEntry - One entry in the buffercache statistics (already aggregated)
type PostgresBuffercacheEntry struct {
	Bytes        int64
	DatabaseName string
	SchemaName   *string
	ObjectName   *string
	ObjectKind   *string
	Toast        bool
}

// PostgresBuffercache - Details on whats contained in the Postgres buffer cache
type PostgresBuffercache struct {
	TotalBytes int64
	FreeBytes  int64

	Entries []PostgresBuffercacheEntry
}
