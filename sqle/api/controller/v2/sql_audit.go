package v2

import (
	e "errors"
	"net/http"

	parser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
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
	Number      uint          `json:"number"`
	ExecSQL     string        `json:"exec_sql"`
	AuditResult []AuditResult `json:"audit_result"`
	AuditLevel  string        `json:"audit_level"`
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
				DbType:   task.DBType,
			}
		}

		results[i] = AuditSQLResV2{
			Number:      sql.Number,
			ExecSQL:     sql.Content,
			AuditResult: convertAuditResultToAuditResV2(sql.AuditResults),
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

type DirectAuditFileReqV2 struct {
	InstanceType string `json:"instance_type" form:"instance_type" example:"MySQL" valid:"required"`
	// 调用方不应该关心SQL是否被完美的拆分成独立的条目, 拆分SQL由SQLE实现
	// 每个数组元素是一个文件内容
	FileContents []string `json:"file_contents" form:"file_contents" example:"select * from t1; select * from t2;"`
	SQLType      string   `json:"sql_type" form:"sql_type" example:"sql" enums:"sql,mybatis," valid:"omitempty,oneof=sql mybatis"`
	ProjectName  string   `json:"project_name" form:"project_name" example:"project1" valid:"required"`
	InstanceName *string  `json:"instance_name" form:"instance_name" example:"instance1"`
	SchemaName   *string  `json:"schema_name" form:"schema_name" example:"schema1"`
}

// @Summary 直接从文件内容提取SQL并审核，SQL文件暂时只支持一次解析一个文件
// @Description Direct audit sql from SQL files and MyBatis files
// @Id directAuditFilesV2
// @Tags sql_audit
// @Security ApiKeyAuth
// @Param req body v2.DirectAuditFileReqV2 true "files that should be audited"
// @Success 200 {object} v2.DirectAuditResV2
// @router /v2/audit_files [post]
func DirectAuditFiles(c echo.Context) error {
	req := new(DirectAuditFileReqV2)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}

	userId := controller.GetUserID(c)
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), req.ProjectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !up.IsProjectMember() {
		return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, e.New("you are not the project member")))
	}

	if len(req.FileContents) <= 0 {
		return controller.JSONBaseErrorReq(c, e.New("file_contents is required"))
	}

	sqls := ""
	if req.SQLType == v1.SQLTypeMyBatis {
		sqls, err = v1.ConvertXmlFileContentToSQLs(req.FileContents)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		// sql文件暂时只支持一次解析一个文件
		sqls = req.FileContents[0]
	}

	l := log.NewEntry().WithField("api", "[post]/v2/audit_files")

	var instance *model.Instance
	var exist bool
	if req.InstanceName != nil {
		instance, exist, err = dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, *req.InstanceName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
		}
	}

	var schemaName string
	if req.SchemaName != nil {
		schemaName = *req.SchemaName
	}

	var task *model.Task
	if instance != nil && schemaName != "" {
		task, err = server.DirectAuditByInstance(l, sqls, schemaName, instance)
	} else {
		task, err = server.AuditSQLByDBType(l, sqls, req.InstanceType)
	}
	if err != nil {
		l.Errorf("audit sqls failed: %v", err)
		return controller.JSONBaseErrorReq(c, v1.ErrDirectAudit)
	}

	return c.JSON(http.StatusOK, DirectAuditResV2{
		BaseRes: controller.BaseRes{},
		Data:    convertFileAuditTaskResultToAuditResV2(task),
	})
}

func convertFileAuditTaskResultToAuditResV2(task *model.Task) *AuditResDataV2 {
	results := make([]AuditSQLResV2, len(task.ExecuteSQLs))
	for i, sql := range task.ExecuteSQLs {
		results[i] = AuditSQLResV2{
			Number:      sql.Number,
			ExecSQL:     sql.Content,
			AuditResult: convertAuditResultToAuditResV2(sql.AuditResults),
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

func convertAuditResultToAuditResV2(auditResults model.AuditResults) []AuditResult {
	ar := make([]AuditResult, len(auditResults))
	for i := range auditResults {
		ar[i] = AuditResult{
			Level:    auditResults[i].Level,
			Message:  auditResults[i].Message,
			RuleName: auditResults[i].RuleName,
		}
	}
	return ar
}
