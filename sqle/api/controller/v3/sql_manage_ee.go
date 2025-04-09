//go:build enterprise
// +build enterprise

package v3

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func getSqlManageList(c echo.Context) error {
	req := new(GetSqlManageListReq)
 	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetProjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}
	ctx := c.Request().Context()

	searchSqlFingerprint := ""
	if req.FuzzySearchSqlFingerprint != nil {
		searchSqlFingerprint = strings.Replace(*req.FuzzySearchSqlFingerprint, "'", "\\'", -1)
	}

	data := map[string]interface{}{
		"fuzzy_search_sql_fingerprint":      searchSqlFingerprint,
		"filter_assignee":                   req.FilterAssignee,
		"filter_instance_id":                req.FilterInstanceID,
		"filter_source":                     req.FilterSource,
		"filter_audit_level":                req.FilterAuditLevel,
		"filter_last_audit_start_time_from": req.FilterLastAuditStartTimeFrom,
		"filter_last_audit_start_time_to":   req.FilterLastAuditStartTimeTo,
		"filter_status":                     req.FilterStatus,
		"project_id":                        projectUid,
		"filter_db_type":                    req.FilterDbType,
		"filter_rule_name":                  req.FilterRuleName,
		"filter_by_environment_tag":         req.FilterByEnvironmentTag,
		"filter_priority":                   req.FilterPriority,
		"fuzzy_search_endpoint":             req.FuzzySearchEndpoint,
		"fuzzy_search_schema_name":          req.FuzzySearchSchemaName,
		"sort_field":                        req.SortField,
		"sort_order":                        req.SortOrder,
		"limit":                             req.PageSize,
		"offset":                            offset,
	}
	if req.FilterByEnvironmentTag != nil && *req.FilterByEnvironmentTag != "" {
		instances, err := dms.GetInstancesInProjectByTypeAndEnvironmentTag(c.Request().Context(), projectUid, "", *req.FilterByEnvironmentTag)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		instIds := make([]string, len(instances))
		for i, v := range instances {
			instIds[i] = v.GetIDStr()
		}

		data["filter_business_instance_ids"] = fmt.Sprintf("\"%s\"", strings.Join(instIds, "\",\""))
	}

	s := model.GetStorage()
	sqlManage, err := s.GetSqlManageListByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlManageRet, err := convertToGetSqlManageListResp(ctx, sqlManage.SqlManageList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlManageListResp{
		BaseRes:               controller.NewBaseReq(nil),
		Data:                  sqlManageRet,
		SqlManageTotalNum:     sqlManage.SqlManageTotalNum,
		SqlManageBadNum:       sqlManage.SqlManageBadNum,
		SqlManageOptimizedNum: sqlManage.SqlManageOptimizedNum,
	})
}

func convertToGetSqlManageListResp(ctx context.Context, sqlManageList []*model.SqlManageDetail) ([]*SqlManage, error) {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	sqlManageRespList := make([]*SqlManage, 0, len(sqlManageList))
	users, err := dms.GetMapUsers(context.TODO(), nil, dms.GetDMSServerAddress())
	if err != nil {

		return nil, err
	}
	for _, sqlManage := range sqlManageList {
		sqlMgr := new(SqlManage)
		sqlMgr.Id = uint64(sqlManage.ID)
		sqlMgr.SqlFingerprint = sqlManage.SqlFingerprint.String
		sqlMgr.Sql = sqlManage.SqlText.String
		sqlMgr.InstanceName = dms.GetInstancesByIdWithoutError(sqlManage.InstanceID.String).Name
		sqlMgr.SchemaName = sqlManage.SchemaName.String
		sqlMgr.AuditStatus = sqlManage.AuditStatus.String
		for i := range sqlManage.AuditResults {
			ar := sqlManage.AuditResults[i]
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &v1.AuditResult{
				Level:           ar.Level,
				Message:         ar.GetAuditMsgByLangTag(lang),
				RuleName:        ar.RuleName,
				ExecutionFailed: ar.ExecutionFailed,
				ErrorInfo:       ar.GetAuditErrorMsgByLangTag(lang),
			})
		}

		source := &v1.Source{
			SqlSourceType: sqlManage.Source.String,
			SqlSourceIDs:  sqlManage.SourceIDs,
		}
		sqlSourceDesc := v1.ConvertSqlSourceDescByType(ctx, sqlManage.Source.String)
		source.SqlSourceDesc = sqlSourceDesc
		sqlMgr.Source = source

		if sqlManage.AppearTimestamp != nil {
			sqlMgr.FirstAppearTimeStamp = sqlManage.AppearTimestamp.Format("2006-01-02 15:04:05")
		}
		if sqlManage.LastReceiveTimestamp != nil {
			sqlMgr.LastReceiveTimeStamp = sqlManage.LastReceiveTimestamp.Format("2006-01-02 15:04:05")
		}
		sqlMgr.FpCount = sqlManage.FpCount

		if sqlManage.Assignees != nil {
			for _, assignees := range strings.Split(*sqlManage.Assignees, ",") {
				if v, ok := users[assignees]; ok {
					sqlMgr.Assignees = append(sqlMgr.Assignees, v.Name)
				}
			}
		}

		endpoints, err := sqlManage.Endpoints()
		if err != nil {
			return nil, err
		}

		sqlMgr.Status = sqlManage.Status.String
		sqlMgr.Remark = sqlManage.Remark.String
		sqlMgr.Endpoints = endpoints
		sqlMgr.Priority = sqlManage.Priority.String
		sqlManageRespList = append(sqlManageRespList, sqlMgr)
	}

	return sqlManageRespList, nil
}
