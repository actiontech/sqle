package v2

import (
	"net/http"

	parser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/labstack/echo/v4"
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

// @Deprecated
// @Summary 直接审核SQL
// @Description Direct audit sql
// @Id directAuditV2
// @Tags sql_audit
// @Security ApiKeyAuth
// @Param req body v2.DirectAuditReqV2 true "sqls that should be audited"
// @Success 200 {object} v2.DirectAuditResV2
// @router /v2/sql_audit [post]
func DirectAudit(c echo.Context) error {
	req := new(DirectAuditReqV2)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}

	sql := req.SQLContent
	if req.SQLType == v1.SQLTypeMyBatis {
		sql, err = parser.ParseXML(req.SQLContent)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	l := log.NewEntry().WithField(c.Path(), "direct audit failed")

	task, err := server.AuditSQLByDBType(l, sql, req.InstanceType)
	if err != nil {
		l.Errorf("audit sqls failed: %v", err)
		return controller.JSONBaseErrorReq(c, v1.ErrDirectAudit)
	}

	return c.JSON(http.StatusOK, DirectAuditResV2{
		BaseRes: controller.BaseRes{},
		Data:    convertTaskResultToAuditResV2(task),
	})
}

func convertTaskResultToAuditResV2(task *model.Task) *AuditResDataV2 {
	results := make([]AuditSQLResV2, len(task.ExecuteSQLs))
	for i, sql := range task.ExecuteSQLs {

		ar := make([]*AuditResult, len(sql.AuditResults))
		for j := range sql.AuditResults {
			ar[j] = &AuditResult{
				Level:    sql.AuditResults[j].Level,
				Message:  sql.AuditResults[j].Message,
				RuleName: sql.AuditResults[j].RuleName,
			}
		}

		results[i] = AuditSQLResV2{
			Number:      sql.Number,
			ExecSQL:     sql.Content,
			AuditResult: ar,
			AuditLevel:  sql.AuditLevel,
		}

	}
	return &AuditResDataV2{
		AuditLevel: task.AuditLevel,
		Score:      task.Score,
		PassRate:   task.PassRate,
		SQLResults: results,
	}
}
