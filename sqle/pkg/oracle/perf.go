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

	// QueryActiveSessionCount queries the number of active user sessions from V$SESSION.
	QueryActiveSessionCount = `SELECT COUNT(*) AS active_sessions FROM V$SESSION WHERE STATUS = 'ACTIVE' AND TYPE = 'USER'`

	// QueryActiveSessions queries active user session details along with the executing SQL.
	// The %s placeholder is for an optional user filter clause, e.g. AND s.USERNAME NOT IN ('SYS','SYSTEM').
	QueryActiveSessions = `
SELECT s.sql_id, s.username, s.status, s.event, q.sql_fulltext
FROM V$SESSION s
LEFT JOIN V$SQL q ON s.sql_id = q.sql_id AND s.sql_child_number = q.child_number
WHERE s.STATUS = 'ACTIVE' AND s.TYPE = 'USER'
    AND s.sql_id IS NOT NULL
    AND s.EVENT != 'SQL*Net message from client'
    %s
`

	// QuerySysstatExecuteCount queries the cumulative execute count from V$SYSSTAT.
	QuerySysstatExecuteCount = `SELECT VALUE FROM V$SYSSTAT WHERE NAME = 'execute count'`
)

// ActiveSession represents an active Oracle session with its executing SQL information.
type ActiveSession struct {
	SQLID       string
	Username    string
	Status      string
	Event       string
	SQLFullText string
}
