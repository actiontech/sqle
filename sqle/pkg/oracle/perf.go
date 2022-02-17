package oracle

// DynPerformanceSQLArea ref to https://docs.oracle.com/cd/E18283_01/server.112/e17110/dynviews_3064.htm
type DynPerformanceSQLArea struct {
	SQLFullText    string `json:"sql_fulltext"`
	Executions     string `json:"executions"`
	ElapsedTime    string `json:"elapsed_time"`
	UserIOWaitTime string `json:"user_io_wait_time"`
	CPUTime        string `json:"cpu_time"`
	DiskReads      string `json:"disk_reads"`
	BufferGets     string `json:"buffer_gets"`
	Avg            string `json:"avg"`
}

const (
	DynPerformanceViewSQLAreaTpl = `
	SELECT * FROM (
		SELECT sql_fulltext
			, executions
			, elapsed_time
			, user_io_wait_time
			, cpu_time
			, disk_reads
			, buffer_gets
			, %v / EXECUTIONS AS avg
		FROM V$SQLAREA  WHERE EXECUTIONS > 0 ORDER BY avg DESC 
	) WHERE rownum <= %v
	`
	DynPerformanceViewSQLAreaColumnExecutions     = "executions"
	DynPerformanceViewSQLAreaColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewSQLAreaColumnCPUTime        = "cpu_time"
	DynPerformanceViewSQLAreaColumnDiskReads      = "disk_reads"
	DynPerformanceViewSQLAreaColumnBufferGets     = "buffer_gets"
	DynPerformanceViewSQLAreaColumnUserIOWaitTime = "user_io_wait_time"
)
