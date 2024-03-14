package dm

type DynPerformanceDmColumns struct {
	SQLFullText      string  `json:"sql_fulltext"`
	Executions       float64 `json:"executions"`
	TotalExecTime    float64 `json:"total_exec_time"`
	AverageExecTime  float64 `json:"average_exec_time"`
	CPUTime          float64 `json:"cpu_time"`
	PhyReadPageCnt   float64 `json:"phy_read_page_cnt"`
	LogicReadPageCnt float64 `json:"logic_read_page_cnt"`
}

const (
	DynPerformanceViewDmTpl = `
SELECT
    sql_fulltext,
    executions,
    total_exec_time,
    average_exec_time,
    cpu_time,
    phy_read_page_cnt,
    logic_read_page_cnt
FROM (
    SELECT
        SQL_TXT AS sql_fulltext,
        COUNT(*) AS executions,
        SUM(EXEC_TIME) AS total_exec_time,
        SUM(EXEC_TIME) / COUNT(*) OVER () AS average_exec_time,
        (SUM(EXEC_TIME) - SUM(PARSE_TIME) - SUM(IO_WAIT_TIME)) AS cpu_time,
        SUM(PHY_READ_CNT) AS phy_read_page_cnt,
        SUM(LOGIC_READ_CNT) AS logic_read_page_cnt,
        ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) AS row_num
    FROM V$SQL_STAT_HISTORY
    GROUP BY SQL_TXT
) t WHERE executions > 0 AND row_num <= %v ORDER BY %v DESC`
	DynPerformanceViewDmColumnExecutions       = "executions"
	DynPerformanceViewDmColumnTotalExecTime    = "total_exec_time"
	DynPerformanceViewDmColumnAverageExecTime  = "average_exec_time"
	DynPerformanceViewDmColumnCPUTime          = "cpu_time"
	DynPerformanceViewDmColumnPhyReadPageCnt   = "phy_read_page_cnt"
	DynPerformanceViewDmColumnLogicReadPageCnt = "logic_read_page_cnt"
)
