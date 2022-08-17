package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/labstack/echo/v4"
)

type DirectAuditReqV1 struct {
	InstanceType string `json:"instance_type" form:"instance_type" example:"mysql" valid:"required"`
	// 调用方不应该关心SQL是否被完美的拆分成独立的条目, 拆分SQL由SQLE实现
	SQLContent string `json:"sql_content" form:"sql_content" example:"select * from t1; select * from t2;" valid:"required"`
}

type DirectAuditResV1 struct {
	controller.BaseRes
	Data *AuditResDataV1 `json:"data"`
}

type AuditResDataV1 struct {
	AuditLevel string          `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score      int32           `json:"score"`
	PassRate   float64         `json:"pass_rate"`
	SQLResults []AuditSQLResV1 `json:"sql_results"`
}

type AuditSQLResV1 struct {
	Number      uint   `json:"number"`
	ExecSQL     string `json:"exec_sql"`
	AuditResult string `json:"audit_result"`
	AuditLevel  string `json:"audit_level"`
}

var ErrDirectAudit = errors.New(errors.GenericError, fmt.Errorf("audit failed, please confirm whether the type of audit plugin supports static audit, please check the log for details"))

// @Summary 直接审核SQL
// @Description Direct audit sql
// @Id directAuditV1
// @Tags sql_audit
// @Security ApiKeyAuth
// @Param req body v1.DirectAuditReqV1 true "sqls that should be audited"
// @Success 200 {object} v1.DirectAuditResV1
// @router /v1/sql_audit [post]
func DirectAudit(c echo.Context) error {
	req := new(DirectAuditReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	l := log.NewEntry().WithField("/v1/sql_audit", "direct audit failed")

	task, err := server.AuditSQLByDBType(l, req.SQLContent, req.InstanceType, "")
	if err != nil {
		l.Errorf("audit sqls failed: %v", err)
		return controller.JSONBaseErrorReq(c, ErrDirectAudit)
	}

	return c.JSON(http.StatusOK, DirectAuditResV1{
		BaseRes: controller.BaseRes{},
		Data:    convertTaskResultToAuditResV1(task),
	})
}

func convertTaskResultToAuditResV1(task *model.Task) *AuditResDataV1 {
	results := make([]AuditSQLResV1, len(task.ExecuteSQLs))
	for i, sql := range task.ExecuteSQLs {
		results[i] = AuditSQLResV1{
			Number:      sql.Number,
			ExecSQL:     sql.Content,
			AuditResult: sql.AuditResult,
			AuditLevel:  sql.AuditLevel,
		}
	}
	return &AuditResDataV1{
		AuditLevel: task.AuditLevel,
		Score:      task.Score,
		PassRate:   task.PassRate,
		SQLResults: results,
	}
}
