package obfororacle

type DynPerformanceObForOracleColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	CPUTime        float64 `json:"cpu_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewObForOracleTpl = `
select * from (
SELECT
t1.sql_fulltext as sql_fulltext,
sum(t1.EXECUTIONS) as executions,
sum(t1.ELAPSED_TIME) as elapsed_time,
sum(t1.CPU_TIME) as cpu_time,
sum(t1.DISK_READS) as disk_reads, // 所有执行物理读的次数
sum(t1.BUFFERS_GETS) as buffer_gets, // 所有执行逻辑读的次数
sum(t1.USER_IO_WAIT_TIME) as user_io_wait_time
FROM 
(select to_char(QUERY_SQL) sql_fulltext,EXECUTIONS,ELAPSED_TIME,CPU_TIME,DISK_READS,BUFFERS_GETS,USER_IO_WAIT_TIME
from GV$OB_PLAN_CACHE_PLAN_STAT
) t1 
where t1.sql_fulltext != 'null'
GROUP BY t1.sql_fulltext ORDER BY %v DESC
)
WHERE rownum <= %v`
	DynPerformanceViewObForOracleColumnExecutions     = "executions"
	DynPerformanceViewObForOracleColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewObForOracleColumnCPUTime        = "cpu_time"
	DynPerformanceViewObForOracleColumnDiskReads      = "disk_reads"
	DynPerformanceViewObForOracleColumnBufferGets     = "buffer_gets"
	DynPerformanceViewObForOracleColumnUserIOWaitTime = "user_io_wait_time"
)
