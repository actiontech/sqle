package v2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/actiontech/sqle/sqle/utils"
)

type GetAuditPlansReqV2 struct {
	FilterAuditPlanDBType       string `json:"filter_audit_plan_db_type" query:"filter_audit_plan_db_type"`
	FuzzySearchAuditPlanName    string `json:"fuzzy_search_audit_plan_name" query:"fuzzy_search_audit_plan_name"`
	FilterAuditPlanType         string `json:"filter_audit_plan_type" query:"filter_audit_plan_type"`
	FilterAuditPlanInstanceName string `json:"filter_audit_plan_instance_name" query:"filter_audit_plan_instance_name"`
	PageIndex                   uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                    uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlansResV2 struct {
	controller.BaseRes
	Data      []AuditPlanResV2 `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

type AuditPlanResV2 struct {
	Name             string             `json:"audit_plan_name" example:"audit_for_java_app1"`
	Cron             string             `json:"audit_plan_cron" example:"0 */2 * * *"`
	DBType           string             `json:"audit_plan_db_type" example:"mysql"`
	Token            string             `json:"audit_plan_token" example:"it's a JWT Token for scanner"`
	InstanceName     string             `json:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string             `json:"audit_plan_instance_database" example:"app1"`
	RuleTemplate     *RuleTemplateV2    `json:"rule_template"`
	Meta             v1.AuditPlanMetaV1 `json:"audit_plan_meta"`
}

// GetAuditPlans
// @Summary 获取扫描任务信息列表
// @Description get audit plan info list
// @Id getAuditPlansV2
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_audit_plan_db_type query string false "filter audit plan db type"
// @Param fuzzy_search_audit_plan_name query string false "fuzzy search audit plan name"
// @Param filter_audit_plan_type query string false "filter audit plan type"
// @Param filter_audit_plan_instance_name query string false "filter audit plan instance name"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} GetAuditPlansResV2
// @router /v2/projects/{project_name}/audit_plans [get]
func GetAuditPlans(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlansReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := v1.CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	currentUser, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	instances, err := s.GetUserCanOpInstancesFromProject(currentUser, projectName, []uint{model.OP_AUDIT_PLAN_VIEW_OTHERS})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	names := []string{}
	for _, instance := range instances {
		names = append(names, instance.Name)
	}

	isManager, err := s.IsProjectManager(currentUser.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_audit_plan_db_type":       req.FilterAuditPlanDBType,
		"fuzzy_search_audit_plan_name":    req.FuzzySearchAuditPlanName,
		"filter_audit_plan_type":          req.FilterAuditPlanType,
		"filter_audit_plan_instance_name": req.FilterAuditPlanInstanceName,
		"current_user_name":               currentUser.Name,
		"current_user_is_admin":           model.DefaultAdminUser == currentUser.Name || isManager,
		"filter_project_name":             projectName,
		"limit":                           req.PageSize,
		"offset":                          offset,
	}
	if len(names) > 0 {
		data["accessible_instances_name"] = fmt.Sprintf("'%s'", strings.Join(names, "', '"))
	}
	auditPlans, count, err := s.GetAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	templateNamesInProject, err := s.GetRuleTemplateNamesByProjectName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlansResV1 := make([]AuditPlanResV2, len(auditPlans))
	for i, ap := range auditPlans {
		meta, err := auditplan.GetMeta(ap.Type.String)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		meta.Params = ap.Params

		ruleTemplateName := ap.RuleTemplateName.String
		ruleTemplate := &RuleTemplateV2{
			Name: ruleTemplateName,
		}
		if !utils.StringsContains(templateNamesInProject, ruleTemplateName) {
			ruleTemplate.IsGlobalRuleTemplate = true
		}

		auditPlansResV1[i] = AuditPlanResV2{
			Name:             ap.Name,
			Cron:             ap.Cron,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			RuleTemplate:     ruleTemplate,
			Token:            ap.Token,
			Meta:             v1.ConvertAuditPlanMetaToRes(meta),
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlansResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlansResV1,
		TotalNums: count,
	})
}

type GetAuditPlanReportSQLsReqV2 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportSQLsResV2 struct {
	controller.BaseRes
	Data      []*AuditPlanReportSQLResV2 `json:"data"`
	TotalNums uint64                     `json:"total_nums"`
}

type AuditPlanReportSQLResV2 struct {
	SQL         string         `json:"audit_plan_report_sql" example:"select * from t1 where id = 1"`
	AuditResult []*AuditResult `json:"audit_plan_report_sql_audit_result"`
	Number      uint           `json:"number" example:"1"`
}

// @Summary 获取指定扫描任务的SQL扫描详情
// @Description get audit plan report SQLs
// @Id getAuditPlanReportsSQLs
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetAuditPlanReportSQLsResV2
// @router /v2/projects/{project_name}/audit_plans/{audit_plan_name}/reports/{audit_plan_report_id}/sqls [get]
func GetAuditPlanReportSQLs(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportSQLsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	apName := c.Param("audit_plan_name")

	ap, exist, err := v1.GetAuditPlanIfCurrentUserCanAccess(c, projectName, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_report_id": c.Param("audit_plan_report_id"),
		"audit_plan_id":        ap.ID,
		"limit":                req.PageSize,
		"offset":               offset,
	}
	auditPlanReportSQLs, count, err := s.GetAuditPlanReportSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportSQLsRes := make([]*AuditPlanReportSQLResV2, len(auditPlanReportSQLs))
	for i, auditPlanReportSQL := range auditPlanReportSQLs {
		auditPlanReportSQLsRes[i] = &AuditPlanReportSQLResV2{
			SQL:    auditPlanReportSQL.SQL,
			Number: auditPlanReportSQL.Number,
		}
		for j := range auditPlanReportSQL.AuditResults {
			ar := auditPlanReportSQL.AuditResults[j]
			auditPlanReportSQLsRes[i].AuditResult = append(auditPlanReportSQLsRes[i].AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.Message,
				RuleName: ar.RuleName,
			})
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportSQLsResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportSQLsRes,
		TotalNums: count,
	})
}

type GetAuditPlanAnalysisDataResV2 struct {
	controller.BaseRes
	Data *TaskAnalysisDataV2 `json:"data"`
}

// GetAuditPlanAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取task相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getTaskAnalysisData
// @Tags audit_plan
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param number path string true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v2.GetAuditPlanAnalysisDataResV2
// @router /v2/projects/{project_name}/audit_plans/{audit_plan_name}/reports/{audit_plan_report_id}/sqls/{number}/analysis [get]
func GetAuditPlanAnalysisData(c echo.Context) error {
	return getAuditPlanAnalysisData(c)
}
