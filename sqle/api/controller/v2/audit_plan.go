package v2

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/labstack/echo/v4"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
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
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	userId := controller.GetUserID(c)

	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"filter_audit_plan_db_type":       req.FilterAuditPlanDBType,
		"fuzzy_search_audit_plan_name":    req.FuzzySearchAuditPlanName,
		"filter_audit_plan_type":          req.FilterAuditPlanType,
		"filter_audit_plan_instance_name": req.FilterAuditPlanInstanceName,
		"current_user_id":                 userId,
		"current_user_is_admin":           up.IsAdmin(),
		"filter_project_id":               projectUid,
		"limit":                           req.PageSize,
		"offset":                          offset,
	}
	if !up.IsAdmin() {
		instanceNames, err := dms.GetInstanceNamesInProjectByIds(c.Request().Context(), projectUid, up.GetInstancesByOP(dmsV1.OpPermissionTypeViewOtherAuditPlan))
		if err != nil {
			return err
		}
		data["accessible_instances_name"] = fmt.Sprintf("\"%s\"", strings.Join(instanceNames, "\",\""))
	}

	auditPlans, count, err := s.GetAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	templateNamesInProject, err := s.GetRuleTemplateNamesByProjectId(projectUid)
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
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	apName := c.Param("audit_plan_name")

	ap, exist, err := v1.GetAuditPlanIfCurrentUserCanAccess(c, projectUid, apName, dmsV1.OpPermissionTypeViewOtherAuditPlan)
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
				DbType:   ap.DBType,
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
// @Id getAuditPlantAnalysisDataV2
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

type AuditPlanSQLReqV2 struct {
	Fingerprint          string    `json:"audit_plan_sql_fingerprint" form:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?"`
	Counter              string    `json:"audit_plan_sql_counter" form:"audit_plan_sql_counter" example:"6" valid:"required"`
	LastReceiveText      string    `json:"audit_plan_sql_last_receive_text" form:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string    `json:"audit_plan_sql_last_receive_timestamp" form:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
	Schema               string    `json:"audit_plan_sql_schema" from:"audit_plan_sql_schema" example:"db1"`
	QueryTimeAvg         *float64  `json:"query_time_avg" from:"query_time_avg" example:"3.22"`
	QueryTimeMax         *float64  `json:"query_time_max" from:"query_time_max" example:"5.22"`
	FirstQueryAt         time.Time `json:"first_query_at" from:"first_query_at" example:"2023-09-12T02:48:01.317880Z"`
	DBUser               string    `json:"db_user" from:"db_user" example:"database_user001"`
	Endpoints            []string  `json:"endpoints" from:"endpoints"`
}

func filterSQLsByBlackList(sqls []*AuditPlanSQLReqV2, blackList []*model.BlackListAuditPlanSQL) []*AuditPlanSQLReqV2 {
	if len(blackList) == 0 {
		return sqls
	}
	filteredSQLs := []*AuditPlanSQLReqV2{}
	filter := v1.ConvertToBlackFilter(blackList)
	for _, sql := range sqls {
		if filter.HasEndpointInBlackList(sql.Endpoints) || filter.IsSqlInBlackList(sql.LastReceiveText) {
			continue
		}
		filteredSQLs = append(filteredSQLs, sql)
	}
	return filteredSQLs
}

func convertToModelAuditPlanSQL(c echo.Context, auditPlan *model.AuditPlan, reqSQLs []*AuditPlanSQLReqV2) ([]*auditplan.SQL, error) {
	var p driver.Plugin
	var err error

	// lazy load driver
	initDriver := func() error {
		if p == nil {
			p, err = common.NewDriverManagerWithoutCfg(log.NewEntry(), auditPlan.DBType)
			if err != nil {
				return err
			}
		}
		return nil
	}
	defer func() {
		if p != nil {
			p.Close(context.TODO())
		}
	}()

	sqls := make([]*auditplan.SQL, 0, len(reqSQLs))
	for _, reqSQL := range reqSQLs {
		if reqSQL.LastReceiveText == "" {
			continue
		}
		fp := reqSQL.Fingerprint
		// the caller may be written in a different language, such as (Java, Bash, Python), so the fingerprint is
		// generated in different ways. In order to maintain th same fingerprint generation logic, we provide a way to
		// generate it by sqle, if the request fingerprint is empty.
		if fp == "" {
			err := initDriver()
			if err != nil {
				return nil, err
			}
			nodes, err := p.Parse(context.TODO(), reqSQL.LastReceiveText)
			if err != nil {
				return nil, err
			}
			if len(nodes) > 0 {
				fp = nodes[0].Fingerprint
			} else {
				fp = reqSQL.LastReceiveText
			}
		}
		counter, err := strconv.ParseUint(reqSQL.Counter, 10, 64)
		if err != nil {
			return nil, err
		}
		info := map[string]interface{}{
			"counter":                counter,
			"last_receive_timestamp": reqSQL.LastReceiveTimestamp,
			server.AuditSchema:       reqSQL.Schema,
			"endpoints":              reqSQL.Endpoints,
		}
		// 兼容老版本的Scannerd
		// 老版本Scannerd不传输这两个字段，不记录到数据库中
		// 并且这里避免记录0值到数据库中，导致后续计算出的平均时间出错
		if reqSQL.QueryTimeAvg != nil {
			info["query_time_avg"] = utils.Round(*reqSQL.QueryTimeAvg, 4)
		}
		if reqSQL.QueryTimeMax != nil {
			info["query_time_max"] = utils.Round(*reqSQL.QueryTimeMax, 4)
		}
		if !reqSQL.FirstQueryAt.IsZero() {
			info["first_query_at"] = reqSQL.FirstQueryAt
		}
		if reqSQL.DBUser != "" {
			info["db_user"] = reqSQL.DBUser
		}
		sqls = append(sqls, &auditplan.SQL{
			Fingerprint: fp,
			SQLContent:  reqSQL.LastReceiveText,
			Info:        info,
			Schema:      reqSQL.Schema,
		})
	}

	return sqls, nil
}

type PartialSyncAuditPlanSQLsReqV2 struct {
	SQLs []*AuditPlanSQLReqV2 `json:"audit_plan_sql_list" form:"audit_plan_sql_list" valid:"dive"`
}

// PartialSyncAuditPlanSQLs
// @Summary 增量同步SQL到扫描任务
// @Description partial sync audit plan SQLs
// @Id partialSyncAuditPlanSQLsV2
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v2.PartialSyncAuditPlanSQLsReqV2 true "partial sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/audit_plans/{audit_plan_name}/sqls/partial [post]
func PartialSyncAuditPlanSQLs(c echo.Context) error {
	req := new(PartialSyncAuditPlanSQLsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	apName := c.Param("audit_plan_name")

	s := model.GetStorage()
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ap, exist, err := dms.GetAuditPlanWithInstanceFromProjectByName(projectUid, apName, s.GetAuditPlanFromProjectByName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}

	l := log.NewEntry()
	reqSQLs := req.SQLs
	blackList, err := s.GetBlackListAuditPlanSQLs()
	if err == nil {
		reqSQLs = filterSQLsByBlackList(reqSQLs, blackList)
	} else {
		l.Warnf("blacklist is not used, err:%v", err)
	}
	if len(reqSQLs) == 0 {
		return controller.JSONBaseErrorReq(c, nil)
	}
	sqls, err := convertToModelAuditPlanSQL(c, ap, reqSQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, auditplan.UploadSQLs(l, ap, sqls, true))
}

type FullSyncAuditPlanSQLsReqV2 struct {
	SQLs []*AuditPlanSQLReqV2 `json:"audit_plan_sql_list" form:"audit_plan_sql_list" valid:"dive"`
}

// @Summary 全量同步SQL到扫描任务
// @Description full sync audit plan SQLs
// @Id fullSyncAuditPlanSQLsV2
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v2.FullSyncAuditPlanSQLsReqV2 true "full sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/audit_plans/{audit_plan_name}/sqls/full [post]
func FullSyncAuditPlanSQLs(c echo.Context) error {
	req := new(FullSyncAuditPlanSQLsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	apName := c.Param("audit_plan_name")

	s := model.GetStorage()

	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ap, exist, err := s.GetAuditPlanFromProjectByName(projectUid, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}

	l := log.NewEntry()
	reqSQLs := req.SQLs
	blackList, err := s.GetBlackListAuditPlanSQLs()
	if err == nil {
		reqSQLs = filterSQLsByBlackList(reqSQLs, blackList)
	} else {
		l.Warnf("blacklist is not used, err:%v", err)
	}
	if len(reqSQLs) == 0 {
		return controller.JSONBaseErrorReq(c, nil)
	}
	sqls, err := convertToModelAuditPlanSQL(c, ap, reqSQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, auditplan.UploadSQLs(l, ap, sqls, false))
}
