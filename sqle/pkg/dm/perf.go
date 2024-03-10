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
	DynPerformanceViewDmTpl = `SELECT TOP %v * FROM (
SELECT SQL_TXT sql_fulltext
       ,count(*) executions
       ,sum(EXEC_TIME) total_exec_time
       ,sum(EXEC_TIME)/count(*) average_exec_time
       ,ABS(sum(EXEC_TIME) - sum(PARSE_TIME) - sum(IO_WAIT_TIME)) cpu_time
       ,sum(PHY_READ_CNT) phy_read_page_cnt
       ,sum(LOGIC_READ_CNT) logic_read_page_cnt
FROM V$SQL_STAT_HISTORY GROUP BY SQL_TXT
) t WHERE executions > 0 ORDER BY %v DESC`
	DynPerformanceViewDmColumnExecutions       = "executions"
	DynPerformanceViewDmColumnTotalExecTime    = "total_exec_time"
	DynPerformanceViewDmColumnAverageExecTime  = "average_exec_time"
	DynPerformanceViewDmColumnCPUTime          = "cpu_time"
	DynPerformanceViewDmColumnPhyReadPageCnt   = "phy_read_page_cnt"
	DynPerformanceViewDmColumnLogicReadPageCnt = "logic_read_page_cnt"
)
