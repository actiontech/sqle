//go:build enterprise
// +build enterprise

package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func getSqlManageList(c echo.Context) error {
	req := new(v1.GetSqlManageListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

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
		"filter_business":                   req.FilterBusiness,
		"filter_priority":                   req.FilterPriority,
		"fuzzy_search_endpoint":             req.FuzzySearchEndpoint,
		"fuzzy_search_schema_name":          req.FuzzySearchSchemaName,
		"sort_field":                        req.SortField,
		"sort_order":                        req.SortOrder,
		"limit":                             req.PageSize,
		"offset":                            offset,
	}
	if req.FilterBusiness != nil && *req.FilterBusiness != "" {
		insts, err := dms.GetInstancesInProjectByTypeAndBusiness(c.Request().Context(), projectUid, "", *req.FilterBusiness)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		instIds := make([]string, len(insts))
		for i, v := range insts {
			instIds[i] = v.GetIDStr()
		}

		data["filter_business_instance_ids"] = fmt.Sprintf("\"%s\"", strings.Join(instIds, "\",\""))
	}

	s := model.GetStorage()
	sqlManage, err := s.GetSqlManageListByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlManageRet, err := convertToGetSqlManageListResp(sqlManage.SqlManageList)
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

func convertToGetSqlManageListResp(sqlManageList []*model.SqlManageDetail) ([]*SqlManage, error) {
	sqlManageRespList := make([]*SqlManage, 0, len(sqlManageList))
	users, err := dms.GetMapUsers(context.TODO(), nil, dms.GetDMSServerAddress())
	if err != nil {

		return nil, err
	}
	for _, sqlManage := range sqlManageList {
		sqlMgr := new(SqlManage)
		sqlMgr.Id = uint64(sqlManage.ID)
		sqlMgr.SqlFingerprint = sqlManage.SqlFingerprint
		sqlMgr.Sql = sqlManage.SqlText
		sqlMgr.InstanceName = dms.GetInstancesByIdWithoutError(sqlManage.InstanceID).Name
		sqlMgr.SchemaName = sqlManage.SchemaName

		for i := range sqlManage.AuditResults {
			ar := sqlManage.AuditResults[i]
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &v1.AuditResult{
				Level:    ar.Level,
				Message:  ar.Message,
				RuleName: ar.RuleName,
			})
		}

		source := &v1.Source{
			SqlSourceType: sqlManage.Source,
			SqlSourceID:   sqlManage.SourceID,
		}
		auditPlanDesc := v1.ConvertAuditPlanDescByType(sqlManage.Source)
		source.SqlSourceDesc = auditPlanDesc
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

		sqlMgr.Status = sqlManage.Status.String
		sqlMgr.Remark = sqlManage.Remark.String
		sqlMgr.Endpoints = sqlManage.Endpoints.String
		sqlMgr.Priority = sqlManage.Priority.String
		sqlManageRespList = append(sqlManageRespList, sqlMgr)
	}

	return sqlManageRespList, nil
}
