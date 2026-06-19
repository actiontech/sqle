package v1

import (
	"context"
	"fmt"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type CreateAuditWhitelistReqV1 struct {
	Value     string `json:"value" example:"create table" valid:"required"`
	MatchType string `json:"match_type" example:"exact_match" enums:"exact_match,fp_match" valid:"omitempty,oneof=exact_match fp_match"`
	Desc      string `json:"desc" example:"used for rapid release"`
}

type CreateSQLRuleExceptionReqV1 struct {
	InstanceID     uint64 `json:"instance_id" example:"1" valid:"required"`
	SQLFingerprint string `json:"sql_fingerprint" example:"select * from ?"`
	RuleName       string `json:"rule_name" example:"all_check_prepare_statement_placeholders" valid:"required"`
	RuleDesc       string `json:"rule_desc" example:"rule description"`
	RuleLevel      string `json:"rule_level" example:"error"`
	Reason         string `json:"reason" example:"业务确认该 SQL 可例外" valid:"required"`
}

type CreateSQLRuleExceptionResV1 struct {
	controller.BaseRes
	Data *SQLRuleExceptionResV1 `json:"data"`
}

type SQLRuleExceptionResV1 struct {
	Id             uint   `json:"sql_rule_exception_id"`
	ProjectId      string `json:"project_id"`
	InstanceID     uint64 `json:"instance_id"`
	InstanceName   string `json:"instance_name"`
	SQLFingerprint string `json:"sql_fingerprint"`
	RuleName       string `json:"rule_name"`
	RuleDesc       string `json:"rule_desc"`
	RuleLevel      string `json:"rule_level"`
	Reason         string `json:"reason"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      string `json:"created_at"`
}

// @Summary 添加单规则 SQL 审核例外
// @Description create a sql rule exception with project, instance, sql fingerprint and rule unique tuple
// @Accept json
// @Id createSQLRuleExceptionV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance body v1.CreateSQLRuleExceptionReqV1 true "add sql rule exception req"
// @Success 200 {object} v1.CreateSQLRuleExceptionResV1
// @router /v1/projects/{project_name}/audit_whitelist/rule_exceptions [post]
func CreateSQLRuleException(c echo.Context) error {
	req := new(CreateSQLRuleExceptionReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasManagePermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to create sql rule exception")))
	}

	s := model.GetStorage()
	sqlRuleException := &model.SQLRuleException{
		ProjectId:      model.ProjectUID(projectUid),
		InstanceID:     req.InstanceID,
		SQLFingerprint: req.SQLFingerprint,
		RuleName:       req.RuleName,
		RuleDesc:       req.RuleDesc,
		RuleLevel:      req.RuleLevel,
		Reason:         req.Reason,
		CreatedBy:      user.Name,
	}

	savedSQLRuleException, _, err := s.CreateSQLRuleExceptionIfNotExist(sqlRuleException)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &CreateSQLRuleExceptionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLRuleExceptionToRes(savedSQLRuleException),
	})
}

// @Summary 取消单规则 SQL 审核例外
// @Description delete a sql rule exception by id in project scope
// @Id deleteSQLRuleExceptionV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_rule_exception_id path string true "sql rule exception id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist/rule_exceptions/{sql_rule_exception_id} [delete]
func DeleteSQLRuleException(c echo.Context) error {
	s := model.GetStorage()
	sqlRuleExceptionID := c.Param("sql_rule_exception_id")
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasManagePermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to delete sql rule exception")))
	}
	sqlRuleException, exist, err := s.GetSQLRuleExceptionByIDAndProjectID(sqlRuleExceptionID, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("sql rule exception is not exist")))
	}
	if err := s.DeleteSQLRuleException(sqlRuleException); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 添加SQL白名单
// @Description create a sql whitelist
// @Accept json
// @Id createAuditWhitelistV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance body v1.CreateAuditWhitelistReqV1 true "add sql whitelist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist [post]
func CreateAuditWhitelist(c echo.Context) error {
	req := new(CreateAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasManagePermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to create audit whitelist")))
	}
	s := model.GetStorage()

	sqlWhitelist := &model.SqlWhitelist{
		ProjectId: model.ProjectUID(projectUid),
		Value:     req.Value,
		Desc:      req.Desc,
		MatchType: req.MatchType,
	}

	err = s.Save(sqlWhitelist)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateAuditWhitelistReqV1 struct {
	Value     *string `json:"value" example:"create table"`
	MatchType *string `json:"match_type" example:"exact_match" enums:"exact_match,fp_match"`
	Desc      *string `json:"desc" example:"used for rapid release"`
}

// @Summary 更新SQL白名单
// @Description update sql whitelist by id
// @Accept json
// @Id UpdateAuditWhitelistByIdV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_whitelist_id path string true "sql audit whitelist id"
// @Param instance body v1.UpdateAuditWhitelistReqV1 true "update sql whitelist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist/{audit_whitelist_id}/ [patch]
func UpdateAuditWhitelistById(c echo.Context) error {
	req := new(UpdateAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasManagePermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to update audit whitelist")))
	}

	s := model.GetStorage()

	whitelistId := c.Param("audit_whitelist_id")
	sqlWhitelist, exist, err := s.GetSqlWhitelistByIdAndProjectUID(whitelistId, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("sql audit whitelist is not exist")))
	}

	// nothing to update
	if req.Value == nil && req.Desc == nil && req.MatchType == nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
	}

	if req.Value != nil {
		sqlWhitelist.Value = *req.Value
	}
	if req.MatchType != nil {
		sqlWhitelist.MatchType = *req.MatchType
	}

	if req.Value != nil || req.MatchType != nil {
		sqlWhitelist.MatchedCount = 0
		sqlWhitelist.LastMatchedTime = nil
	}

	if req.Desc != nil {
		sqlWhitelist.Desc = *req.Desc
	}

	err = s.Save(sqlWhitelist)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 删除SQL白名单信息
// @Description remove sql white
// @Id deleteAuditWhitelistByIdV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_whitelist_id path string true "audit whitelist id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist/{audit_whitelist_id}/ [delete]
func DeleteAuditWhitelistById(c echo.Context) error {
	s := model.GetStorage()
	whitelistId := c.Param("audit_whitelist_id")
	// projectName := c.Param("project_name")

	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasManagePermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to delete audit whitelist")))
	}
	sqlWhitelist, exist, err := s.GetSqlWhitelistByIdAndProjectUID(whitelistId, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("sql audit whitelist is not exist")))
	}
	err = s.Delete(sqlWhitelist)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetAuditWhitelistReqV1 struct {
	FuzzySearchValue *string `json:"fuzzy_search_value" query:"fuzzy_search_value" valid:"omitempty"`
	FilterMatchType  *string `json:"filter_match_type" query:"filter_match_type" valid:"omitempty,oneof=exact_match fp_match" enums:"exact_match,fp_match"`
	PageIndex        uint32  `json:"page_index" query:"page_index" valid:"required"`
	PageSize         uint32  `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditWhitelistResV1 struct {
	controller.BaseRes
	Data      []*AuditWhitelistResV1 `json:"data"`
	TotalNums int64                  `json:"total_nums"`
}

type GetSQLRuleExceptionReqV1 struct {
	FuzzySearchValue      *string `json:"fuzzy_search_value" query:"fuzzy_search_value" valid:"omitempty"`
	FilterInstanceID      *uint64 `json:"filter_instance_id" query:"filter_instance_id" valid:"omitempty"`
	FilterRuleName        *string `json:"filter_rule_name" query:"filter_rule_name" valid:"omitempty"`
	FilterCreatedBy       *string `json:"filter_created_by" query:"filter_created_by" valid:"omitempty"`
	FilterCreatedTimeFrom *string `json:"filter_created_time_from" query:"filter_created_time_from" valid:"omitempty"`
	FilterCreatedTimeTo   *string `json:"filter_created_time_to" query:"filter_created_time_to" valid:"omitempty"`
	FilterSQLFingerprint  *string `json:"filter_sql_fingerprint" query:"filter_sql_fingerprint" valid:"omitempty"`
	PageIndex             uint32  `json:"page_index" query:"page_index" valid:"required"`
	PageSize              uint32  `json:"page_size" query:"page_size" valid:"required"`
}

type GetSQLRuleExceptionResV1 struct {
	controller.BaseRes
	Data      []*SQLRuleExceptionResV1 `json:"data"`
	TotalNums int64                    `json:"total_nums"`
}

type AuditWhitelistResV1 struct {
	Id            uint       `json:"audit_whitelist_id"`
	Value         string     `json:"value"`
	MatchType     string     `json:"match_type"`
	MatchedCount  uint       `json:"matched_count"`
	LastMatchTime *time.Time `json:"last_match_time"`
	Desc          string     `json:"desc"`
}

// @Summary 获取Sql审核白名单
// @Description get all whitelist
// @Id getAuditWhitelistV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_value query string false "fuzzy value"
// @Param filter_match_type query string false "match type"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetAuditWhitelistResV1
// @router /v1/projects/{project_name}/audit_whitelist [get]
func GetSqlWhitelist(c echo.Context) error {
	req := new(GetAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasViewPermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to select audit whitelist")))
	}
	s := model.GetStorage()
	sqlWhitelist, count, err := s.GetSqlWhitelistByProjectUID(req.PageIndex, req.PageSize, model.ProjectUID(projectUid), req.FuzzySearchValue, req.FilterMatchType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	whitelistRes := make([]*AuditWhitelistResV1, 0, len(sqlWhitelist))
	for _, v := range sqlWhitelist {
		whitelistRes = append(whitelistRes, &AuditWhitelistResV1{
			Id:            v.ID,
			Value:         v.Value,
			Desc:          v.Desc,
			MatchType:     v.MatchType,
			MatchedCount:  uint(v.MatchedCount),
			LastMatchTime: v.LastMatchedTime,
		})
	}
	return c.JSON(http.StatusOK, &GetAuditWhitelistResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      whitelistRes,
		TotalNums: count,
	})
}

// @Summary 获取单规则 SQL 审核例外列表
// @Description get sql rule exceptions with server side filters
// @Id getSQLRuleExceptionV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_value query string false "fuzzy value"
// @Param filter_instance_id query string false "instance id"
// @Param filter_rule_name query string false "rule name"
// @Param filter_created_by query string false "created by"
// @Param filter_created_time_from query string false "created time from"
// @Param filter_created_time_to query string false "created time to"
// @Param filter_sql_fingerprint query string false "sql fingerprint"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetSQLRuleExceptionResV1
// @router /v1/projects/{project_name}/audit_whitelist/rule_exceptions [get]
func GetSQLRuleException(c echo.Context) error {
	req := new(GetSQLRuleExceptionReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	hasPermission, err := hasViewPermission(user.GetIDStr(), projectUid, dmsV1.OpPermissionMangeAuditSQLWhiteList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !hasPermission {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you have no permission to select sql rule exception")))
	}

	sqlRuleExceptions, count, err := model.GetStorage().GetSQLRuleExceptions(req.PageIndex, req.PageSize, model.SQLRuleExceptionListFilter{
		ProjectID:        model.ProjectUID(projectUid),
		InstanceID:       req.FilterInstanceID,
		RuleName:         req.FilterRuleName,
		CreatedBy:        req.FilterCreatedBy,
		CreatedTimeFrom:  req.FilterCreatedTimeFrom,
		CreatedTimeTo:    req.FilterCreatedTimeTo,
		SQLFingerprint:   req.FilterSQLFingerprint,
		FuzzySearchValue: req.FuzzySearchValue,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ret := make([]*SQLRuleExceptionResV1, 0, len(sqlRuleExceptions))
	for _, sqlRuleException := range sqlRuleExceptions {
		ret = append(ret, convertSQLRuleExceptionListItemToRes(sqlRuleException))
	}
	return c.JSON(http.StatusOK, &GetSQLRuleExceptionResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      ret,
		TotalNums: count,
	})
}

func convertSQLRuleExceptionToRes(sqlRuleException *model.SQLRuleException) *SQLRuleExceptionResV1 {
	if sqlRuleException == nil {
		return nil
	}
	return &SQLRuleExceptionResV1{
		Id:             sqlRuleException.ID,
		ProjectId:      string(sqlRuleException.ProjectId),
		InstanceID:     sqlRuleException.InstanceID,
		SQLFingerprint: sqlRuleException.SQLFingerprint,
		RuleName:       sqlRuleException.RuleName,
		RuleDesc:       sqlRuleException.RuleDesc,
		RuleLevel:      sqlRuleException.RuleLevel,
		Reason:         sqlRuleException.Reason,
		CreatedBy:      sqlRuleException.CreatedBy,
		CreatedAt:      sqlRuleException.CreatedAt.Format(time.RFC3339),
	}
}

func convertSQLRuleExceptionListItemToRes(sqlRuleException *model.SQLRuleExceptionListItem) *SQLRuleExceptionResV1 {
	if sqlRuleException == nil {
		return nil
	}
	ret := convertSQLRuleExceptionToRes(&sqlRuleException.SQLRuleException)
	ret.InstanceName = sqlRuleException.InstanceName
	return ret
}
