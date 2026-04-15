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

// DynPerformanceSysStat V$SYSSTAT 采集结果
type DynPerformanceSysStat struct {
	LogonsCurrent int64 `json:"logons_current"`
	ExecuteCount  int64 `json:"execute_count"`
}

const DynPerformanceViewSysStatQuery = `SELECT name, value FROM V$SYSSTAT WHERE name IN ('logons current', 'execute count')`

// DynPerformanceSession V$SESSION 活跃会话
type DynPerformanceSession struct {
	SID        int64  `json:"sid"`
	Username   string `json:"username"`
	SchemaName string `json:"schema_name"`
	SQLText    string `json:"sql_text"`
	LastCallET int64  `json:"last_call_et"` // 当前状态持续秒数
}

const DynPerformanceViewActiveSessionQuery = `
SELECT
    s.sid,
    s.username,
    s.schemaname,
    sq.sql_text,
    s.last_call_et
FROM V$SESSION s
LEFT JOIN V$SQL sq ON s.sql_id = sq.sql_id AND s.sql_child_number = sq.child_number
WHERE s.status = 'ACTIVE'
  AND s.type = 'USER'
  AND s.username IS NOT NULL
  AND ROWNUM <= 1000
`

const DynPerformanceViewActiveSessionCountQuery = `
SELECT COUNT(*) FROM V$SESSION WHERE STATUS = 'ACTIVE' AND TYPE = 'USER' AND USERNAME IS NOT NULL
`
