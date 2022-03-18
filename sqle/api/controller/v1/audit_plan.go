package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/labstack/echo/v4"
	"github.com/ungerik/go-dry"
)

var tokenExpire = 365 * 24 * time.Hour

var (
	errAuditPlanNotExist         = errors.New(errors.DataNotExist, fmt.Errorf("audit plan is not exist"))
	errAuditPlanExisted          = errors.New(errors.DataNotExist, fmt.Errorf("audit plan existed"))
	errAuditPlanInstanceConflict = errors.New(errors.DataConflict, fmt.Errorf("instance_name can not be empty while instance_database is not empty"))
	errAuditPlanCannotAccess     = errors.New(errors.DataInvalid, fmt.Errorf("you can not access this audit plan"))
)

type GetAuditPlanMetasReqV1 struct {
	FilterInstanceType *string `json:"filter_instance_type" query:"filter_instance_type"`
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

func convertAuditPlanMetaToRes(meta auditplan.Meta) AuditPlanMetaV1 {
	res := AuditPlanMetaV1{
		Type:         meta.Type,
		Desc:         meta.Desc,
		InstanceType: meta.InstanceType,
	}
	if meta.Params != nil && len(meta.Params) > 0 {
		paramsRes := make([]AuditPlanParamResV1, 0, len(meta.Params))
		for _, p := range meta.Params {
			paramRes := AuditPlanParamResV1{
				Key:   p.Key,
				Desc:  p.Desc,
				Type:  string(p.Type),
				Value: p.Value,
			}
			paramsRes = append(paramsRes, paramRes)
		}
		res.Params = paramsRes
	}
	return res
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
	req := new(GetAuditPlanMetasReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var metas []AuditPlanMetaV1
	for _, meta := range auditplan.Metas {
		// filter instance type
		if req.FilterInstanceType == nil ||
			meta.InstanceType == auditplan.InstanceTypeAll ||
			meta.InstanceType == *req.FilterInstanceType {
			metas = append(metas, convertAuditPlanMetaToRes(meta))
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanMetasResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    metas,
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

func checkAndGenerateAuditPlanParams(auditPlanType, instanceType string, paramsReq []AuditPlanParamReqV1) (params.Params, error) {
	meta, err := auditplan.GetMeta(auditPlanType)
	if err != nil {
		return nil, err
	}
	if meta.InstanceType != auditplan.InstanceTypeAll && meta.InstanceType != instanceType {
		return nil, fmt.Errorf("audit plan type %s not found", auditPlanType)
	}
	// check request params is equal params.
	if len(paramsReq) != len(meta.Params) {
		reqParamsKey := make([]string, 0, len(paramsReq))
		for _, p := range paramsReq {
			reqParamsKey = append(reqParamsKey, p.Key)
		}
		paramsKey := make([]string, 0, len(meta.Params))
		for _, p := range meta.Params {
			paramsKey = append(paramsKey, p.Key)
		}
		return nil, fmt.Errorf("request params key is [%s], but need [%s]",
			strings.Join(reqParamsKey, ", "), strings.Join(paramsKey, ", "))
	}
	for _, p := range paramsReq {
		// set and valid param.
		err := meta.Params.SetParamValue(p.Key, p.Value)
		if err != nil {
			return nil, fmt.Errorf("set param error: %s", err)
		}
	}
	return meta.Params, nil
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

	// check instance
	var instanceType string
	if req.InstanceName != "" {
		inst, exist, err := s.GetInstanceByName(req.InstanceName)
		if !exist {
			return controller.JSONBaseErrorReq(c, errInstanceNotExist)
		} else if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
		// check instance database
		if req.InstanceDatabase != "" {
			d, err := newDriverWithoutAudit(log.NewEntry(), inst, "")
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
		instanceType = inst.DbType
	} else {
		instanceType = req.InstanceType
	}

	// check params
	if req.Type == "" {
		req.Type = auditplan.TypeDefault
	}
	ps, err := checkAndGenerateAuditPlanParams(req.Type, instanceType, req.Params)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	// check user and generate token
	currentUserName := controller.GetUserName(c)
	user, exist, err := s.GetUserByName(currentUserName)
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("user is not exist")))
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	j := utils.NewJWT(utils.SecretKey)
	t, err := j.CreateToken(currentUserName, time.Now().Add(tokenExpire).Unix(),
		utils.WithAuditPlanName(req.Name))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	ap := &model.AuditPlan{
		Name:             req.Name,
		CronExpression:   req.Cron,
		Type:             req.Type,
		Params:           ps,
		CreateUserID:     user.ID,
		Token:            t,
		DBType:           instanceType,
		InstanceName:     req.InstanceName,
		InstanceDatabase: req.InstanceDatabase,
	}
	err = s.Save(ap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(ap.Name))
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
	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	ap, exist, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}
	err = s.Delete(ap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(apName))
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

	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
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
	if req.Params != nil {
		ps, err := checkAndGenerateAuditPlanParams(ap.Type, ap.DBType, req.Params)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateAttr["params"] = ps
	}

	err = storage.UpdateAuditPlanByName(apName, updateAttr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(apName))
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
	Name             string          `json:"audit_plan_name" example:"audit_for_java_app1"`
	Cron             string          `json:"audit_plan_cron" example:"0 */2 * * *"`
	DBType           string          `json:"audit_plan_db_type" example:"mysql"`
	Token            string          `json:"audit_plan_token" example:"it's a JWT Token for scanner"`
	InstanceName     string          `json:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string          `json:"audit_plan_instance_database" example:"app1"`
	Meta             AuditPlanMetaV1 `json:"audit_plan_meta"`
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
		"current_user_name":         currentUserName,
		"current_user_is_admin":     model.DefaultAdminUser == currentUserName,
		"limit":                     req.PageSize,
		"offset":                    offset,
	}
	auditPlans, count, err := s.GetAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlansResV1 := make([]AuditPlanResV1, len(auditPlans))
	for i, ap := range auditPlans {
		meta, err := auditplan.GetMeta(ap.Type.String)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		meta.Params = ap.Params
		auditPlansResV1[i] = AuditPlanResV1{
			Name:             ap.Name,
			Cron:             ap.Cron,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			Token:            ap.Token,
			Meta:             convertAuditPlanMetaToRes(meta),
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
	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
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
	meta, err := auditplan.GetMeta(ap.Type)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	meta.Params = ap.Params

	return c.JSON(http.StatusOK, &GetAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanResV1{
			Name:             ap.Name,
			Cron:             ap.CronExpression,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			Token:            ap.Token,
			Meta:             convertAuditPlanMetaToRes(meta),
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
	Id         string  `json:"audit_plan_report_id" example:"1"`
	AuditLevel string  `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score      int32   `json:"score"`
	PassRate   float64 `json:"pass_rate"`
	Timestamp  string  `json:"audit_plan_report_timestamp" example:"RFC3339"`
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
	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
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
			Id:         auditPlanReport.ID,
			AuditLevel: auditPlanReport.AuditLevel.String,
			Score:      auditPlanReport.Score.Int32,
			PassRate:   auditPlanReport.PassRate.Float64,
			Timestamp:  auditPlanReport.CreateAt,
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
	return nil
}

type FullSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list"`
}

type AuditPlanSQLReqV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" form:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?"`
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

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.UploadSQLs(apName, sqls, false))
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

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.UploadSQLs(apName, sqls, true))
}

func checkAndConvertToModelAuditPlanSQL(c echo.Context, apName string, reqSQLs []AuditPlanSQLReqV1) ([]*auditplan.SQL, error) {
	s := model.GetStorage()

	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return nil, err
	}

	ap, exist, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errAuditPlanNotExist
	}

	var driver driver.Driver
	// lazy load driver
	initDriver := func() error {
		if driver == nil {
			driver, err = newDriverWithoutCfg(log.NewEntry(), ap.DBType)
			if err != nil {
				return err
			}
		}
		return nil
	}
	defer func() {
		if driver != nil {
			driver.Close(context.TODO())
		}
	}()

	sqls := make([]*auditplan.SQL, len(reqSQLs))
	for i, reqSQL := range reqSQLs {
		fp := reqSQL.Fingerprint
		// the caller may be written in a different language, such as (Java, Bash, Python), so the fingerprint is
		// generated in different ways. In order to maintain th same fingerprint generation logic, we provide a way to
		// generate it by sqle, if the request fingerprint is empty.
		if fp == "" {
			err := initDriver()
			if err != nil {
				return nil, err
			}
			nodes, err := driver.Parse(context.TODO(), reqSQL.LastReceiveText)
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
		}
		sqls[i] = &auditplan.SQL{
			Fingerprint: fp,
			SQLContent:  reqSQL.LastReceiveText,
			Info:        info,
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
	return nil
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
	err := CheckCurrentUserCanAccessAuditPlan(c, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	manager := auditplan.GetManager()
	report, err := manager.Audit(apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &TriggerAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanReportResV1{
			Id:         fmt.Sprintf("%v", report.ID),
			AuditLevel: report.AuditLevel,
			Score:      report.Score,
			PassRate:   report.PassRate,
			Timestamp:  report.CreatedAt.Format(time.RFC3339),
		},
	})
}

func CheckCurrentUserCanAccessAuditPlan(c echo.Context, apName string) error {
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
