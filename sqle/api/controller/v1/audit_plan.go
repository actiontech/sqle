package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"

	"github.com/labstack/echo/v4"
	"github.com/ungerik/go-dry"
)

var (
	errAuditPlanNotExist         = errors.New(errors.DataNotExist, fmt.Errorf("audit plan is not exist"))
	errAuditPlanExisted          = errors.New(errors.DataNotExist, fmt.Errorf("audit plan existed"))
	errAuditPlanInstanceConflict = errors.New(errors.DataConflict, fmt.Errorf("instance_name can not be empty while instance_database is not empty"))
	errAuditPlanCannotAccess     = errors.New(errors.DataInvalid, fmt.Errorf("you can not access this audit plan"))
)

type GetAuditPlanMetasReqV1 struct {
	FilterInstanceType string `json:"filter_instance_type" query:"filter_instance_type"`
}

type GetAuditPlanMetasResV1 struct {
	controller.BaseRes
	Data []AuditPlanMetaV1 `json:"data"`
}

type AuditPlanMetaV1 struct {
	Type         string                `json:"audit_plan_type"`
	Desc         string                `json:"audit_plan_type_desc"`
	InstanceType string                `json:"instance_type"`
	Params       []AuditPlanParamResV1 `json:"audit_plan_params,omitempty"`
}

type AuditPlanParamResV1 struct {
	Key   string `json:"key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	Type  string `json:"type" enums:"string,int,bool"`
}

// @Summary 获取审核任务元信息
// @Description get audit plan metas
// @Id getAuditPlanMetasV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param filter_instance_type query string false "filter instance type"
// @Success 200 {object} v1.GetAuditPlanMetasResV1
// @router /v1/audit_plan_metas [get]
func GetAuditPlanMetas(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetAuditPlanMetasResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    []AuditPlanMetaV1{},
	})
}

type CreateAuditPlanReqV1 struct {
	Name             string                `json:"audit_plan_name" form:"audit_plan_name" example:"audit_plan_for_java_repo_1" valid:"required,name"`
	Cron             string                `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *" valid:"required,cron"`
	InstanceType     string                `json:"audit_plan_instance_type" form:"audit_plan_instance_type" example:"mysql" valid:"required"`
	InstanceName     string                `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string                `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
	Type             string                `json:"audit_plan_type" form:"audit_plan_type" example:"slow log"`
	Params           []AuditPlanParamReqV1 `json:"audit_plan_params" valid:"dive,required"`
}

type AuditPlanParamReqV1 struct {
	Key   string `json:"key" form:"key" valid:"required"`
	Value string `json:"value" form:"value" valid:"required"`
}

// @Summary 添加审核计划
// @Description create audit plan
// @Id createAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Accept json
// @Param audit_plan body v1.CreateAuditPlanReqV1 true "create audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans [post]
func CreateAuditPlan(c echo.Context) error {
	s := model.GetStorage()
	manager := auditplan.GetManager()

	req := new(CreateAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !dry.StringInSlice(req.InstanceType, driver.AllDrivers()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DriverNotExist, &driver.DriverNotSupportedError{DriverTyp: req.InstanceType}))
	}

	if req.InstanceDatabase != "" && req.InstanceName == "" {
		return controller.JSONBaseErrorReq(c, errAuditPlanInstanceConflict)
	}

	_, exist, err := s.GetAuditPlanByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanExisted)
	}

	currentUserName := controller.GetUserName(c)

	if req.InstanceName == "" {
		err := manager.AddStaticAuditPlan(req.Name, req.Cron, req.InstanceType, currentUserName)
		return controller.JSONBaseErrorReq(c, err)
	}

	instance, exist, err := s.GetInstanceByName(req.InstanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, errInstanceNotExist))
	}

	if req.InstanceDatabase != "" {
		d, err := newDriverWithoutAudit(log.NewEntry(), instance, "")
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		defer d.Close(context.TODO())

		schemas, err := d.Schemas(context.TODO())
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !dry.StringInSlice(req.InstanceDatabase, schemas) {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("database %v is not exist in instance", req.InstanceDatabase)))
		}
	}
	err = manager.AddDynamicAuditPlan(req.Name, req.Cron, req.InstanceName, req.InstanceDatabase, currentUserName)
	return controller.JSONBaseErrorReq(c, err)
}

// @Summary 删除审核计划
// @Description delete audit plan
// @Id deleteAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/ [delete]
func DeleteAuditPlan(c echo.Context) error {
	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.DeleteAuditPlan(apName))
}

type UpdateAuditPlanReqV1 struct {
	Cron             *string               `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *" valid:"omitempty,cron"`
	InstanceName     *string               `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase *string               `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
	Params           []AuditPlanParamReqV1 `json:"audit_plan_params" valid:"dive,required"`
}

// @Summary 更新审核计划
// @Description update audit plan
// @Id updateAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @param audit_plan body v1.UpdateAuditPlanReqV1 true "update audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/ [patch]
func UpdateAuditPlan(c echo.Context) error {
	req := new(UpdateAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")

	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	updateAttr := make(map[string]interface{})
	if req.Cron != nil {
		updateAttr["cron_expression"] = *req.Cron
	}
	if req.InstanceName != nil {
		updateAttr["instance_name"] = *req.InstanceName
	}
	if req.InstanceDatabase != nil {
		updateAttr["instance_database"] = *req.InstanceDatabase
	}
	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.UpdateAuditPlan(apName, updateAttr))
}

type GetAuditPlansReqV1 struct {
	FilterAuditPlanDBType string `json:"filter_audit_plan_db_type" query:"filter_audit_plan_db_type"`
	PageIndex             uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize              uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlansResV1 struct {
	controller.BaseRes
	Data      []AuditPlanResV1 `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

type AuditPlanResV1 struct {
	Name             string           `json:"audit_plan_name" example:"audit_for_java_app1"`
	Cron             string           `json:"audit_plan_cron" example:"0 */2 * * *"`
	DBType           string           `json:"audit_plan_db_type" example:"mysql"`
	Token            string           `json:"audit_plan_token" example:"it's a JWT Token for scanner"`
	InstanceName     string           `json:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string           `json:"audit_plan_instance_database" example:"app1"`
	Meta             *AuditPlanMetaV1 `json:"audit_plan_meta"`
}

// @Summary 获取审核计划信息列表
// @Description get audit plan info list
// @Id getAuditPlansV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param filter_audit_plan_db_type query string false "filter audit plan db type"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlansResV1
// @router /v1/audit_plans [get]
func GetAuditPlans(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlansReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	currentUserName := controller.GetUserName(c)
	data := map[string]interface{}{
		"filter_audit_plan_db_type": req.FilterAuditPlanDBType,
		"current_user_id":           currentUserName,
		"current_user_is_admin":     model.DefaultAdminUser == currentUserName,
		"limit":                     req.PageSize,
		"offset":                    offset,
	}
	auditPlans, count, err := s.GetAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlansResV1 := make([]AuditPlanResV1, len(auditPlans))
	for i, auditPlan := range auditPlans {
		auditPlansResV1[i] = AuditPlanResV1{
			Name:             auditPlan.Name,
			Cron:             auditPlan.Cron,
			DBType:           auditPlan.DBType,
			InstanceName:     auditPlan.InstanceName,
			InstanceDatabase: auditPlan.InstanceDatabase,

			Token: auditPlan.Token,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlansResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlansResV1,
		TotalNums: count,
	})
}

type GetAuditPlanResV1 struct {
	controller.BaseRes
	Data AuditPlanResV1 `json:"data"`
}

// @Summary 获取指定审核计划
// @Description get audit plan
// @Id getAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.GetAuditPlanResV1
// @router /v1/audit_plans/{audit_plan_name}/ [get]
func GetAuditPlan(c echo.Context) error {
	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	storage := model.GetStorage()

	ap, exist, err := storage.GetAuditPlanByName(apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	return c.JSON(http.StatusOK, &GetAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanResV1{
			Name:             ap.Name,
			Cron:             ap.CronExpression,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			Token:            ap.Token,
		},
	})
}

type GetAuditPlanReportsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanReportResV1 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type AuditPlanReportResV1 struct {
	Id        string `json:"audit_plan_report_id" example:"1"`
	Timestamp string `json:"audit_plan_report_timestamp" example:"RFC3339"`
}

// @Summary 获取指定审核计划的报告列表
// @Description get audit plan report list
// @Id getAuditPlanReportsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportsResV1
// @router /v1/audit_plans/{audit_plan_name}/reports [get]
func GetAuditPlanReports(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_name": apName,
		"limit":           req.PageSize,
		"offset":          offset,
	}
	auditPlanReports, count, err := s.GetAuditPlanReportsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportsResV1 := make([]AuditPlanReportResV1, len(auditPlanReports))
	for i, auditPlanReport := range auditPlanReports {
		auditPlanReportsResV1[i] = AuditPlanReportResV1{
			Id:        auditPlanReport.ID,
			Timestamp: auditPlanReport.CreateAt,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportsResV1,
		TotalNums: count,
	})
}

type GetAuditPlanReportSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportSQLsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanReportSQLResV1 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

type AuditPlanReportSQLResV1 struct {
	Fingerprint          string `json:"audit_plan_report_sql_fingerprint" example:"select * from t1 where id = ?"`
	LastReceiveText      string `json:"audit_plan_report_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_report_sql_last_receive_timestamp" example:"RFC3339"`
	AuditResult          string `json:"audit_plan_report_sql_audit_result" example:"same format as task audit result"`
}

// @Summary 获取指定审核计划的SQL审核详情
// @Description get audit plan report SQLs
// @Deprecated
// @Id getAuditPlanReportSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportSQLsResV1
// @router /v1/audit_plans/{audit_plan_name}/report/{audit_plan_report_id}/ [get]
func GetAuditPlanReportSQLs(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_name":      apName,
		"audit_plan_report_id": c.Param("audit_plan_report_id"),
		"limit":                req.PageSize,
		"offset":               offset,
	}
	auditPlanReportSQLs, count, err := s.GetAuditPlanReportSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportSQLsResV1 := make([]AuditPlanReportSQLResV1, len(auditPlanReportSQLs))
	for i, auditPlanReportSQL := range auditPlanReportSQLs {
		auditPlanReportSQLsResV1[i] = AuditPlanReportSQLResV1{
			Fingerprint:          auditPlanReportSQL.Fingerprint,
			LastReceiveText:      auditPlanReportSQL.LastReceiveText,
			LastReceiveTimestamp: auditPlanReportSQL.LastReceiveTimestamp,
			AuditResult:          auditPlanReportSQL.AuditResult,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportSQLsResV1,
		TotalNums: count,
	})
}

type FullSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list"`
}

type AuditPlanSQLReqV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" form:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?" valid:"required"`
	Counter              string `json:"audit_plan_sql_counter" form:"audit_plan_sql_counter" example:"6" valid:"required"`
	LastReceiveText      string `json:"audit_plan_sql_last_receive_text" form:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_sql_last_receive_timestamp" form:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
}

// @Summary 全量同步SQL到审核计划
// @Description full sync audit plan SQLs
// @Id fullSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.FullSyncAuditPlanSQLsReqV1 true "full sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/sqls/full [post]
func FullSyncAuditPlanSQLs(c echo.Context) error {
	req := new(FullSyncAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	sqls, err := checkAndConvertToModelAuditPlanSQL(c, apName, req.SQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	err = s.OverrideAuditPlanSQLs(apName, sqls)
	return controller.JSONBaseErrorReq(c, err)
}

type PartialSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list"`
}

// @Summary 增量同步SQL到审核计划
// @Description partial sync audit plan SQLs
// @Id partialSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.PartialSyncAuditPlanSQLsReqV1 true "partial sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/sqls/partial [post]
func PartialSyncAuditPlanSQLs(c echo.Context) error {
	req := new(PartialSyncAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	sqls, err := checkAndConvertToModelAuditPlanSQL(c, apName, req.SQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	err = s.UpdateAuditPlanSQLs(apName, sqls)
	return controller.JSONBaseErrorReq(c, err)
}

func checkAndConvertToModelAuditPlanSQL(c echo.Context, apName string, reqSQLs []AuditPlanSQLReqV1) ([]*model.AuditPlanSQL, error) {
	s := model.GetStorage()

	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return nil, err
	}

	_, exist, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errAuditPlanNotExist
	}

	sqls := make([]*model.AuditPlanSQL, len(reqSQLs))
	for i, reqSQL := range reqSQLs {
		counter, err := strconv.ParseInt(reqSQL.Counter, 10, 64)
		if err != nil {
			return nil, err
		}
		sqls[i] = &model.AuditPlanSQL{
			Fingerprint:          reqSQL.Fingerprint,
			Counter:              int(counter),
			LastSQL:              reqSQL.LastReceiveText,
			LastReceiveTimestamp: reqSQL.LastReceiveTimestamp,
		}
	}
	return sqls, nil
}

type GetAuditPlanSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanSQLsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanSQLResV1 `json:"data"`
	TotalNums uint64              `json:"total_nums"`
}

type AuditPlanSQLResV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?"`
	Counter              string `json:"audit_plan_sql_counter" example:"6"`
	LastReceiveText      string `json:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
}

// @Summary 获取指定审核计划的SQLs信息(不包括审核结果)
// @Description get audit plan SQLs
// @Deprecated
// @Id getAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanSQLsResV1
// @router /v1/audit_plans/{audit_plan_name}/sqls [get]
func GetAuditPlanSQLs(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_name": apName,
		"limit":           req.PageSize,
		"offset":          offset,
	}
	auditPlanSQLs, count, err := s.GetAuditPlanSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanSQLsResV1 := make([]AuditPlanSQLResV1, len(auditPlanSQLs))
	for i, auditPlanSQL := range auditPlanSQLs {
		auditPlanSQLsResV1[i] = AuditPlanSQLResV1{
			Fingerprint:          auditPlanSQL.Fingerprint,
			LastReceiveText:      auditPlanSQL.LastReceiveText,
			LastReceiveTimestamp: auditPlanSQL.LastReceiveTimestamp,
			Counter:              auditPlanSQL.Counter,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanSQLsResV1,
		TotalNums: count,
	})
}

type TriggerAuditPlanResV1 struct {
	controller.BaseRes
	Data AuditPlanReportResV1 `json:"data"`
}

// @Summary 触发审核计划
// @Description trigger audit plan
// @Id triggerAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.TriggerAuditPlanResV1
// @router /v1/audit_plans/{audit_plan_name}/trigger [post]
func TriggerAuditPlan(c echo.Context) error {
	apName := c.Param("audit_plan_name")
	err := checkCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	manager := auditplan.GetManager()
	report, err := manager.TriggerAuditPlan(apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &TriggerAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanReportResV1{
			Id:        fmt.Sprintf("%v", report.ID),
			Timestamp: report.CreatedAt.Format(time.RFC3339),
		},
	})
}

func checkCurrentUserCanAccessAuditPlan(c echo.Context, apName string) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}

	storage := model.GetStorage()

	ap, exist, err := storage.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}
	if !exist {
		return errAuditPlanNotExist
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	if user.ID != ap.CreateUserID {
		return errAuditPlanCannotAccess
	}
	return nil
}
