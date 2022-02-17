package oracle

type DynPerformanceSQLArea struct {
	SQLFullText string `json:"sql_fulltext"`
	Avg         string `json:"avg"`
}

const (
	DynPerformanceViewSQLAreaTpl                  = `SELECT * FROM (SELECT SQL_FULLTEXT, %v / EXECUTIONS AS avg  FROM V$SQLAREA  WHERE EXECUTIONS > 0 ORDER BY avg DESC ) WHERE rownum <= %v`
	DynPerformanceViewSQLAreaColumnElapsedTime    = "ELAPSED_TIME"
	DynPerformanceViewSQLAreaColumnCPUTime        = "CPU_TIME"
	DynPerformanceViewSQLAreaColumnDiskReads      = "DISK_READS"
	DynPerformanceViewSQLAreaColumnBufferGets     = "BUFFER_GETS"
	DynPerformanceViewSQLAreaColumnUserIOWaitTime = "USER_IO_WAIT_TIME"
)
