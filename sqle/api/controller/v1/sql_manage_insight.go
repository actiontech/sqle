package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlPerformanceInsightsResp struct {
	controller.BaseRes
	Data *SqlPerformanceInsights `json:"data"`
}

// SqlPerformanceInsights SQL性能洞察数据结构体
type SqlPerformanceInsights struct {
	TaskSupport bool   `json:"task_support"` // 是否支持任务
	TaskEnable  bool   `json:"task_enable"`  // 是否启用任务
	XInfo       string `json:"x_info"`       // X轴信息
	YInfo       string `json:"y_info"`       // Y轴信息
	Message     string `json:"message"`      // 提示信息
	Lines       []Line `json:"lines"`        // 图表线条数据
}

// Line 图表线条数据结构体
type Line struct {
	LineName string        `json:"line_name"` // 线条名称
	Points   *[]ChartPoint `json:"points"`    // 图表点数据
}

// MetricName SQL性能洞察指标类型
type MetricName string

const (
	MetricNameComprehensiveTrend MetricName = "comprehensive_trend"  // 数据源综合趋势
	MetricNameSlowSQLTrend       MetricName = "slow_sql_trend"       // 慢SQL趋势
	MetricNameTopSQLTrend        MetricName = "top_sql_trend"        // TopSQL趋势
	MetricNameActiveSessionTrend MetricName = "active_session_trend" // 活跃会话数趋势
)

type GetSqlPerformanceInsightsReq struct {
	ProjectName string     `param:"project_name" json:"project_name" valid:"required"`
	InstanceId  string     `query:"instance_id" json:"instance_id" valid:"required"`
	StartTime   string     `query:"start_time" json:"start_time" valid:"required"`
	EndTime     string     `query:"end_time" json:"end_time" valid:"required"`
	MetricName  MetricName `query:"metric_name" json:"metric_name" enums:"comprehensive_trend,slow_sql_trend,active_session_trend" valid:"required"`
}

// GetSqlPerformanceInsights
// @Summary 获取SQL管控SQL性能洞察图表数据
// @Description get sql manage sql performance insights
// @Id GetSqlPerformanceInsights
// @Tags SqlInsight
// @Param project_name path string true "project name"
// @Param metric_name query string true "metric name" Enums(comprehensive_trend,slow_sql_trend,top_sql_trend,active_session_trend)
// @Param start_time query string true "start time"
// @Param end_time query string true "end time"
// @Param instance_id query string true "instance id"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlPerformanceInsightsResp
// @Router /v1/projects/{project_name}/sql_performance_insights [get]
func GetSqlPerformanceInsights(c echo.Context) error {
	return getSqlPerformanceInsights(c)
}

type GetSqlPerformanceInsightsRelatedSQLResp struct {
	controller.BaseRes
	Data      []*RelatedSQLInfo `json:"data"`
	TotalNums uint32            `json:"total_nums"`
}

type SqlSourceTypeEnum string

const (
	SqlSourceTypeWorkflow  SqlSourceTypeEnum = "workflow"
	SqlSourceTypeSqlManage SqlSourceTypeEnum = "sql_manage"
	SqlSourceTypeWorkBench SqlSourceTypeEnum = "workbench"
)

type RelatedSQLInfo struct {
	SqlFingerprint     string                   `json:"sql_fingerprint"`
	Source             SqlSourceTypeEnum        `json:"source" enums:"workflow,sql_manage"`
	ExecuteTimeAvg     float64                  `json:"execute_time_avg"`               // 平均执行时间(s)
	ExecuteTimeMax     float64                  `json:"execute_time_max"`               // 最大执行时间(s)
	ExecuteTimeMin     float64                  `json:"execute_time_min"`               // 最小执行时间(s)
	ExecuteTimeSum     float64                  `json:"execute_time_sum"`               // 总执行时间(s)
	LockWaitTime       float64                  `json:"lock_wait_time"`                 // 锁等待时间(s)
	ExecutionTimeTrend *SqlAnalysisScatterChart `json:"execution_time_trend,omitempty"` // SQL 趋势图表
}

// 散点图结构体，专用于SQL执行代价的散点图表示
type SqlAnalysisScatterChart struct {
	Points  *[]ScatterPoint `json:"points"`
	XInfo   *string         `json:"x_info"`
	YInfo   *string         `json:"y_info"`
	Message string          `json:"message"`
}

type ScatterPoint struct {
	Time            *string             `json:"time"`              // 时间点
	ExecuteTime     *float64            `json:"execute_time"`      // 执行时间
	SQL             *string             `json:"sql"`               // SQL语句
	Id              uint64              `json:"id"`                // SQL ID
	IsInTransaction bool                `json:"is_in_transaction"` // 是否在事务中
	Infos           []map[string]string `json:"info"`              // 额外信息
}

type GetSqlPerformanceInsightsRelatedSQLReq struct {
	ProjectName  string            `param:"project_name" json:"project_name" valid:"required"`
	InstanceId   string            `query:"instance_id" json:"instance_id" valid:"required"`
	StartTime    string            `query:"start_time" json:"start_time" valid:"required"`
	EndTime      string            `query:"end_time" json:"end_time" valid:"required"`
	FilterSource SqlSourceTypeEnum `query:"filter_source" json:"filter_source,omitempty" enums:"workflow,sql_manage,workbench" valid:"required"`
	OrderBy      *string           `query:"order_by" json:"order_by,omitempty" enums:"execute_time_avg,execute_time_max,execute_time_min,execute_time_sum,lock_wait_time"`
	IsAsc        *bool             `query:"is_asc" json:"is_asc,omitempty"`
	PageIndex    uint32            `query:"page_index" valid:"required" json:"page_index"`
	PageSize     uint32            `query:"page_size" valid:"required" json:"page_size"`
}

// fixme: 这个接口的设计上由于产品对于SQL指纹浮动的图表的展示情况还有疑问。所以对应的设计应该是缺失了部分。
// 1、由于后续还需要查询具体sql的管理事务等。所以可能需要设计一个 id 之类的字段以便后续查询。但是现在list里的内容是sql指纹，不是具体的sql，所以这个ID要放在哪还需要等产品的结论。
// 2、同样的，对于关联事物功能。目前产品的设计里，关联事物打开的是具体的SQL，而不是SQL指纹。但是又不是所有SQL都在事务里。
// 2.1、所以问题一为，点击按钮的时候要打开的是哪一条具体的SQL。
// 2.2、需要设计字段显示的告诉前端这条SQL是否在一个事务当中。以便前端展示 "关联事务" 按钮是否可用。
// 3、这个table还有一个跳转到SQL分析的功能。但是现有的SQL分析也是针对单条SQL。而这个地方也是SQL指纹。
// 3.1、而且看上去现有的SQL是基于生成了分析结果。然后这个结果有个类似于 sql_manage_id 的玩意来获取的。但是从这里跳过去就没这个id了。需要前后端讨论一下具体实现。
// GetSqlPerformanceInsightsRelatedSQL
// @Summary 获取sql洞察 时间选区 的关联SQL
// @Description Get related SQL for the selected time range in SQL performance insights
// @Id GetSqlPerformanceInsightsRelatedSQL
// @Tags SqlInsight
// @Param project_name path string true "project name"
// @Param instance_id query string true "instance id"
// @Param start_time query string true "start time"
// @Param end_time query string true "end time"
// @Param filter_source query string true "filter by SQL source" Enums(workflow,sql_manage,workbench)
// @Param order_by query string false "order by field"
// @Param is_asc query bool false "is ascending order"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlPerformanceInsightsRelatedSQLResp
// @Router /v1/projects/{project_name}/sql_performance_insights/related_sql [get]
func GetSqlPerformanceInsightsRelatedSQL(c echo.Context) error {
	return getSqlPerformanceInsightsRelatedSQL(c)
}

type LockType string

const (
	LockTypeShared                   LockType = "SHARED"                     // 共享锁
	LockTypeExclusive                LockType = "EXCLUSIVE"                  // 排他锁
	LockTypeIntentionShared          LockType = "INTENTION_SHARED"           // 意向共享锁
	LockTypeIntentionExclusive       LockType = "INTENTION_EXCLUSIVE"        // 意向排他锁
	LockTypeSharedIntentionExclusive LockType = "SHARED_INTENTION_EXCLUSIVE" // 共享意向排他锁
	LockTypeRowLock                  LockType = "ROW_LOCK"                   // 行锁
	LockTypeTableLock                LockType = "TABLE_LOCK"                 // 表锁
	LockTypeMetadataLock             LockType = "METADATA_LOCK"              // 元数据锁
)

type TransactionState string

const (
	TransactionStateRunning   TransactionState = "RUNNING"   // 执行中
	TransactionStateCompleted TransactionState = "COMPLETED" // 已完成
)

type TransactionInfo struct {
	TransactionId        string           `json:"transaction_id"`                                                                                                                       // 事务ID
	LockType             LockType         `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"` // 锁类型
	TransactionStartTime string           `json:"transaction_start_time"`                                                                                                               // 事务开始时间
	TransactionEndTime   string           `json:"transaction_end_time"`                                                                                                                 // 事务结束时间
	TransactionDuration  float64          `json:"transaction_duration"`                                                                                                                 // 事务持续时间
	TransactionState     TransactionState `json:"transaction_state" enums:"RUNNING,COMPLETED"`                                                                                          // 事务状态
}

type TransactionTimelineItem struct {
	StartTime   string `json:"start_time"`  // 开始时间
	Description string `json:"description"` // 描述信息
}

type TransactionTimeline struct {
	Timeline         []*TransactionTimelineItem `json:"timeline"`           // 时间线项目列表
	CurrentStepIndex int                        `json:"current_step_index"` // 当前步骤索引
}

// fixme： 这里同样有SQL分析功能。所以有跟上面3.1相同的问题。

type TransactionSQL struct {
	SQL             string   `json:"sql"`                                                                                                                                  // SQL语句
	ExecuteDuration float64  `json:"execute_duration"`                                                                                                                     // 执行时长
	LockType        LockType `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"` // 锁类型
}

type TransactionLockInfo struct {
	LockType      LockType `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"` // 锁类型
	TableName     string   `json:"table_name"`                                                                                                                           // 表名
	CreateLockSQL string   `json:"create_lock_sql"`                                                                                                                      // 创建锁的SQL语句
}

type RelatedTransactionInfo struct {
	TransactionInfo     *TransactionInfo       `json:"transaction_info"`      // 事务信息
	TransactionTimeline *TransactionTimeline   `json:"transaction_timeline"`  // 事务时间线
	TransactionSQLList  []*TransactionSQL      `json:"related_sql_info"`      // 相关SQL列表
	TransactionLockInfo []*TransactionLockInfo `json:"transaction_lock_info"` // 事务锁信息
}

type GetSqlPerformanceInsightsRelatedTransactionResp struct {
	controller.BaseRes
	Data *RelatedTransactionInfo `json:"data"`
}

// GetSqlPerformanceInsightsRelatedTransaction
// @Summary 获取sql洞察 相关SQL中具体一条SQL 的关联事务
// @Description Get related transaction for the selected SQL in SQL performance insights
// @Id GetSqlPerformanceInsightsRelatedTransaction
// @Tags SqlInsight
// @Param project_name path string true "project name"
// @Param instance_id query string true "instance id"
// @Param sql_id query string true "sql id"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlPerformanceInsightsRelatedTransactionResp
// @Router /v1/projects/{project_name}/sql_performance_insights/related_sql/related_transaction [get]
func GetSqlPerformanceInsightsRelatedTransaction(c echo.Context) error {
	return getSqlPerformanceInsightsRelatedTransaction(c)
}
