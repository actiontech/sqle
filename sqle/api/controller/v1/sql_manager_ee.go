//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func getSqlManageList(c echo.Context) error {
	req := new(GetSqlManageListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}
	data := map[string]interface{}{
		"fuzzy_search_sql_fingerprint":      req.FuzzySearchSqlFingerprint,
		"filter_assignee":                   req.FilterAssignee,
		"filter_instance_name":              req.FilterInstanceName,
		"filter_source":                     req.FilterSource,
		"filter_audit_level":                req.FilterAuditLevel,
		"filter_last_audit_start_time_from": req.FilterLastAuditStartTimeFrom,
		"filter_last_audit_start_time_to":   req.FilterLastAuditStartTimeTo,
		"filter_status":                     req.FilterStatus,
		"project_name":                      projectName,
		"limit":                             req.PageSize,
		"offset":                            offset,
	}

	s := model.GetStorage()
	sqlManage, err := s.GetSqlManageListByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlManageListResp{
		BaseRes:               controller.NewBaseReq(nil),
		Data:                  convertToGetSqlManageListResp(sqlManage.SqlManageList),
		SqlManageTotalNum:     sqlManage.SqlManageTotalNum,
		SqlManageBadNum:       sqlManage.SqlManageBadNum,
		SqlManageOptimizedNum: sqlManage.SqlManageOptimizedNum,
	})
}

func convertToGetSqlManageListResp(sqlManageList []*model.SqlManageDetail) []*SqlManage {
	sqlManageRespList := make([]*SqlManage, 0, len(sqlManageList))
	for _, sqlManage := range sqlManageList {
		sqlMgr := new(SqlManage)
		sqlMgr.Id = uint64(sqlManage.ID)
		sqlMgr.SqlFingerprint = sqlManage.SqlFingerprint
		sqlMgr.Sql = sqlManage.SqlText
		sqlMgr.InstanceName = sqlManage.InstanceName
		sqlMgr.SchemaName = sqlManage.SchemaName

		for i := range sqlManage.AuditResults {
			ar := sqlManage.AuditResults[i]
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.Message,
				RuleName: ar.RuleName,
			})
		}

		source := &Source{Type: sqlManage.Source}
		if sqlManage.ApName != nil {
			source.AuditPlanName = *sqlManage.ApName
		}
		if sqlManage.SqlAuditRecordID != nil {
			source.SqlAuditRecordId = *sqlManage.SqlAuditRecordID
		}
		sqlMgr.Source = source

		if sqlManage.AppearTimestamp != nil {
			sqlMgr.FirstAppearTime = sqlManage.AppearTimestamp.Format("2006-01-02 15:04:05")
		}
		if sqlManage.LastReceiveTimestamp != nil {
			sqlMgr.LastAppearTime = sqlManage.LastReceiveTimestamp.Format("2006-01-02 15:04:05")
		}
		sqlMgr.AppearNum = sqlManage.FpCount

		for _, assignee := range sqlManage.Assignees {
			sqlMgr.Assignees = append(sqlMgr.Assignees, assignee)
		}

		sqlMgr.Status = sqlManage.Status
		sqlMgr.Remark = sqlManage.Remark

		sqlManageRespList = append(sqlManageRespList, sqlMgr)
	}

	return sqlManageRespList
}

func batchUpdateSqlManage(c echo.Context) error {
	req := new(BatchUpdateSqlManageReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if len(req.SqlManageIdList) == 0 {
		return controller.JSONBaseErrorReq(c, nil)
	}

	s := model.GetStorage()

	err := s.BatchUpdateSqlManage(req.SqlManageIdList, req.Status, req.Remark, req.Assignees)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}
