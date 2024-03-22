package postgresql

type DynPerformancePgColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	CPUTime        float64 `json:"cpu_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewPgTpl = `
SELECT query as sql_fulltext,
sum(calls) as executions,
sum(total_exec_time) AS elapsed_time,
sum(cpu_user_time) as cpu_time,
sum(shared_blks_read) AS disk_reads, // 表示从共享缓冲区中读取的块数。这个值表示数据库系统从磁盘或其他存储介质中读取的数据块数量，而不是从内存中读取的数据。
sum(shared_blks_hit) AS buffer_gets, // 表示从共享缓冲区中命中的块数。这个值表示数据库系统从内存中读取的数据块数量，而不是从磁盘或其他存储介质中读取的数据。
sum(blk_read_time) as user_io_wait_time
FROM pg_stat_monitor
WHERE calls > 0
group by query
ORDER BY %v DESC limit %v`
	DynPerformanceViewPgSQLColumnExecutions     = "executions"
	DynPerformanceViewPgSQLColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewPgSQLColumnCPUTime        = "cpu_time"
	DynPerformanceViewPgSQLColumnDiskReads      = "disk_reads"
	DynPerformanceViewPgSQLColumnBufferGets     = "buffer_gets"
	DynPerformanceViewPgSQLColumnUserIOWaitTime = "user_io_wait_time"
)
