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
	DynPerformanceViewSQLAreaSlowLogTpl = `
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
        last_active_time >= SYSDATE - INTERVAL '%v' MINUTE
        AND s.EXECUTIONS > 0
        AND s.elapsed_time / s.executions > %v
        %v
    ORDER BY s.elapsed_time / s.executions DESC
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

// ProcessListSession represents an active session from V$SESSION joined with V$SQL.
type ProcessListSession struct {
	SID         int64  `json:"sid"`
	Username    string `json:"username"`
	SchemaName  string `json:"schema_name"`
	SQLFullText string `json:"sql_fulltext"`
	LastCallET  int64  `json:"last_call_et"`
}

// DynPerformanceViewSessionTpl is a Go template for querying active sessions
// from V$SESSION joined with V$SQL. It filters user sessions that are active,
// have a SQL_ID, a USERNAME, and excludes the current session.
// When .MinSecond > 0, it additionally filters by LAST_CALL_ET >= .MinSecond.
const DynPerformanceViewSessionTpl = `
SELECT
    s.SID,
    s.USERNAME,
    s.SCHEMANAME,
    q.SQL_FULLTEXT,
    s.LAST_CALL_ET
FROM V$SESSION s
LEFT JOIN V$SQL q ON s.SQL_ID = q.SQL_ID AND s.SQL_CHILD_NUMBER = q.CHILD_NUMBER
WHERE s.TYPE = 'USER'
  AND s.STATUS = 'ACTIVE'
  AND s.SQL_ID IS NOT NULL
  AND s.USERNAME IS NOT NULL
  AND s.SID != SYS_CONTEXT('USERENV', 'SID')
{{- if gt .MinSecond 0 }}
  AND s.LAST_CALL_ET >= {{ .MinSecond }}
{{- end }}
`
