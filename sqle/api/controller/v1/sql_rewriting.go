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
	// @Description 重写建议列表
	Suggestions []*RewriteSuggestion `json:"suggestions"`
	// @Description 重写后的SQL
	RewrittenSQL string `json:"rewritten_sql"`
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

// RewriteSQL retrieves the rewritten SQL and associated suggestions for a specific SQL statement within a task.
// @Summary 获取任务中指定SQL的重写结果和建议
// @Description 获取特定任务中某条SQL语句的重写后的SQL及相关建议
// @Id RewriteSQL
// @Accept json
// @Produce json
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Param req body v1.RewriteSQLReq true "request body"
// @Security ApiKeyAuth
// @Success 200 {object} v1.RewriteSQLRes "成功返回重写结果"
// @router /v1/tasks/audits/{task_id}/sqls/{number}/rewrite [post]
func RewriteSQL(c echo.Context) error {
	return getRewriteSQLData(c)
}
