package oracle

// DynPerformanceSQLArea ref to https://docs.oracle.com/cd/E18283_01/server.112/e17110/dynviews_3064.htm
type DynPerformanceSQLArea struct {
	SQLFullText    string `json:"sql_fulltext"`
	Executions     int64  `json:"executions"`
	ElapsedTime    int64  `json:"elapsed_time"`
	UserIOWaitTime int64  `json:"user_io_wait_time"`
	CPUTime        int64  `json:"cpu_time"`
	DiskReads      int64  `json:"disk_reads"`
	BufferGets     int64  `json:"buffer_gets"`
	UserName       string `json:"username"`
}

// Note:
// I can not use Oracle to convert microseconds to seconds by "ROUND(cpu_time/1000/1000)";
// it should return float64 or string, but the driver return empty value, seem to be a bug;
// So I get the original cpu_time (microseconds) and convert it to seconds within the SQLE code logic.
const (
	DynPerformanceViewSQLAreaTpl = `
SELECT * FROM (
    SELECT
        s.sql_fulltext,
        s.executions,
        s.elapsed_time,
        s.user_io_wait_time,
        s.cpu_time,
        s.disk_reads,
        s.buffer_gets,
        u.username
    FROM
        V$SQLAREA s
    JOIN
        DBA_USERS u ON s.parsing_user_id = u.user_id
    WHERE
        last_active_time >= SYSDATE - INTERVAL '%v' MINUTE AND
        s.EXECUTIONS > 0
        %v
    ORDER BY %v DESC
)
WHERE
    rownum <= %v
`
	DynPerformanceViewSQLAreaColumnExecutions     = "executions"
	DynPerformanceViewSQLAreaColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewSQLAreaColumnCPUTime        = "cpu_time"
	DynPerformanceViewSQLAreaColumnDiskReads      = "disk_reads"
	DynPerformanceViewSQLAreaColumnBufferGets     = "buffer_gets"
	DynPerformanceViewSQLAreaColumnUserIOWaitTime = "user_io_wait_time"
)
