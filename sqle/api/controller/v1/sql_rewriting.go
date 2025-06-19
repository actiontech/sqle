package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type RewriteSQLReq struct {
	// @Description 是否启用结构化类型的重写
	EnableStructureType bool `json:"enable_structure_type" example:"false" default:"false"`
}

type RewriteSQLData struct {
	// @Description 重写前的SQL业务描述
	BusinessDesc string `json:"business_desc"`
	// @Description 重写前的SQL执行逻辑描述
	LogicDesc string `json:"logic_desc"`
	// @Description 重写建议列表
	Suggestions []*RewriteSuggestion `json:"suggestions"`
	// @Description 重写后的SQL
	RewrittenSQL string `json:"rewritten_sql"`
	// @Description 重写后的SQL业务描述
	RewrittenSQLBusinessDesc string `json:"rewritten_sql_business_desc"`
	// @Description 重写后的SQL执行逻辑描述
	RewrittenSQLLogicDesc string `json:"rewritten_sql_logic_desc"`
	// @Description 重写前后的业务不等价性描述，为空表示等价
	BusinessNonEquivalentDesc string `json:"business_non_equivalent_desc"`

	// TODO: 重写SQL的审核结果
}

type RewriteSuggestion struct {
	// @Description 审核规则名称
	// @Required
	RuleName string `json:"rule_name"`
	// @Description 审核规则等级
	// @Required
	AuditLevel string `json:"audit_level" enums:"normal,notice,warn,error"`
	// @Description 重写建议类型：语句级重写、结构级重写、其他
	// @Required
	Type string `json:"type" enums:"statement,structure,other"`
	// @Description 重写描述（适用于所有类型）
	// @Required
	Desc string `json:"desc"`
	// @Description 处理状态：初始化、已处理
	Status string `json:"status" enums:"initial,processed" default:"initial"`
	// @Description 重写后的SQL（适用于语句级重写和结构级重写）
	RewrittenSQL string `json:"rewritten_sql"`
	// @Description 数据库结构变更建议说明（例如：建议添加索引、修改表结构等优化建议）（适用于结构级重写）
	DDL_DCL_desc string `json:"ddl_dcl_desc"`
	// @Description 具体的数据库结构变更语句，需要在数据库中执行该变更语句之后再应用重写SQL（包含CREATE/ALTER/DROP等DDL语句，或SET等DCL语句）（适用于结构级重写）
	DDL_DCL string `json:"ddl_dcl"`
}

type RewriteSQLRes struct {
	controller.BaseRes
	Data *RewriteSQLData `json:"data"`
}

// AsyncRewriteTask 异步重写任务（Controller层定义）
type AsyncRewriteTask struct {
	// @Description 任务ID
	TaskID string `json:"task_id"`
	// @Description SQL编号
	SQLNumber string `json:"sql_number"`
	// @Description 任务状态
	Status string `json:"status" enums:"pending,running,completed,failed"`
	// @Description 错误信息
	ErrorMessage string `json:"error_message,omitempty"`
	// @Description 开始时间
	StartTime string `json:"start_time"`
	// @Description 结束时间
	EndTime *string `json:"end_time,omitempty"`
	// @Description 重写结果
	Result *RewriteSQLData `json:"result,omitempty"`
}

// RewriteSQL starts an asynchronous SQL rewriting task for a specific SQL statement within a task.
// @Summary 启动异步SQL重写任务
// @Description 启动特定任务中某条SQL语句的异步重写任务，返回任务状态
// @Id RewriteSQL
// @Accept json
// @Produce json
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Param req body v1.RewriteSQLReq true "request body"
// @Security ApiKeyAuth
// @Success 200 {object} v1.AsyncRewriteTaskStatusRes "成功启动异步重写任务"
// @router /v1/tasks/audits/{task_id}/sqls/{number}/rewrite [post]
func RewriteSQL(c echo.Context) error {
	return getRewriteSQLData(c)
}

// 异步重写任务状态响应
type AsyncRewriteTaskStatusRes struct {
	controller.BaseRes
	Data *AsyncRewriteTask `json:"data"`
}

// GetAsyncRewriteTaskStatus 获取异步重写任务状态接口
// @Summary 获取异步重写任务状态
// @Description 获取特定任务中某条SQL语句的异步重写任务状态
// @Id GetAsyncRewriteTaskStatus
// @Accept json
// @Produce json
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.AsyncRewriteTaskStatusRes "成功返回任务状态"
// @router /v1/tasks/audits/{task_id}/sqls/{number}/rewrite/status [get]
func GetAsyncRewriteTaskStatus(c echo.Context) error {
	return getAsyncRewriteTaskStatus(c)
}
