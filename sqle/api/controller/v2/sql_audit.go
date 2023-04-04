package v2

import (
	"github.com/labstack/echo/v4"

	"github.com/actiontech/sqle/sqle/api/controller"
)

type DirectAuditReqV2 struct {
	InstanceType string `json:"instance_type" form:"instance_type" example:"MySQL" valid:"required"`
	// 调用方不应该关心SQL是否被完美的拆分成独立的条目, 拆分SQL由SQLE实现
	SQLContent string `json:"sql_content" form:"sql_content" example:"select * from t1; select * from t2;" valid:"required"`
	SQLType    string `json:"sql_type" form:"sql_type" example:"sql" enums:"sql,mybatis," valid:"omitempty,oneof=sql mybatis"`
}

type AuditResDataV2 struct {
	AuditLevel string          `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score      int32           `json:"score"`
	PassRate   float64         `json:"pass_rate"`
	SQLResults []AuditSQLResV2 `json:"sql_results"`
}

type AuditSQLResV2 struct {
	Number      uint           `json:"number"`
	ExecSQL     string         `json:"exec_sql"`
	AuditResult []*AuditResult `json:"audit_result"`
	AuditLevel  string         `json:"audit_level"`
}

type DirectAuditResV2 struct {
	controller.BaseRes
	Data *AuditResDataV2 `json:"data"`
}

// @Summary 直接审核SQL
// @Description Direct audit sql
// @Id directAuditV2
// @Tags sql_audit
// @Security ApiKeyAuth
// @Param req body v2.DirectAuditReqV2 true "sqls that should be audited"
// @Success 200 {object} v2.DirectAuditResV2
// @router /v2/sql_audit [post]
func DirectAudit(c echo.Context) error {
	return controller.JSONNewNotImplementedErr(c)
}
