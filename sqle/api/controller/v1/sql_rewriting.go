package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type TaskSQLRewrittenData struct {
	// @Description 表结构
	TableMetas *TableMetas `json:"table_metas"`
	// @Description 重写建议列表
	Suggestions []*SQLRewrittenSuggestion `json:"suggestions"`
	// @Description 重写后的SQL
	RewrittenSQL string `json:"rewritten_sql"`
	// @Description 业务不等价性描述，为空表示等价
	BusinessNonEquivalentDesc string `json:"business_non_equivalent_desc"`

	// TODO: 重写SQL的审核结果
}

type SQLRewrittenSuggestion struct {
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
	// @Description 重写后的SQL（适用于语句级重写和结构级重写）
	RewrittenSQL string `json:"rewritten_sql"`
	// @Description ddl dcl 描述(适用于结构级重写)
	DDL_DCL_desc string `json:"ddl_dcl_desc"`
	// @Description ddl dcl(适用于结构级重写)
	DDL_DCL string `json:"ddl_dcl"`
}

type GetTaskSQLRewrittenRes struct {
	controller.BaseRes
	Data *TaskSQLRewrittenData `json:"data"`
}

// GetTaskSQLRewritten retrieves the rewritten SQL and associated suggestions for a specific SQL statement within a task.
// @Summary 获取任务中指定SQL的重写结果和建议
// @Description 获取特定任务中某条SQL语句的重写后的SQL及相关建议
// @Id getTaskSQLRewritten
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskSQLRewrittenRes "成功返回重写结果"
// @router /v1/tasks/audits/{task_id}/sqls/{number}/rewrite [get]
func GetTaskSQLRewritten(c echo.Context) error {
	return getTaskSQLRewrittenData(c)
}
