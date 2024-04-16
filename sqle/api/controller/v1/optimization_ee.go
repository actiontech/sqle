//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/optimization"
	"github.com/labstack/echo/v4"
)

func sqlOptimizate(c echo.Context) error {
	req := new(OptimizeSQLReq)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}

	// 获取入参中的SQL
	sqls := getSQLFromFileResp{}
	if req.SQLContent != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.SQLContent,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	optimizationId, err := optimization.Optimizate(c.Request().Context(), user.Name, projectUid, req.InstanceName, req.SchemaName, req.OptimizationName, sqls.SQLsFromFormData)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, OptimizeSQLRes{
		BaseRes: controller.NewBaseReq(nil),
		Data: &OptimizeSQLResData{
			OptimizationRecordId: optimizationId,
		},
	})
}

func getOptimizationRecord(c echo.Context) error {
	optimizationId := c.Param("optimization_record_id")
	record, err := model.GetStorage().GetOptimizationRecordId(optimizationId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ret := OptimizationDetail{
		OptimizationID:   optimizationId,
		OptimizationName: record.OptimizationName,
		InstanceNmae:     record.InstanceName,
		DBType:           record.DBType,
		CreatedTime:      record.CreatedAt,
		CreatedUser:      record.Creator,
		Optimizationsummary: Optimizationsummary{
			NumberOfQuery:          record.NumberOfQuery,
			NumberOfSyntaxError:    record.NumberOfSyntaxError,
			NumberOfRewrite:        record.NumberOfRewrite,
			NumberOfRewrittenQuery: record.NumberOfRewrittenQuery,
			NumberOfIndex:          record.NumberOfIndex,
			NumberOfQueryIndex:     record.NumberOfQuery,
			PerformanceGain:        record.PerformanceImprove,
		},
		IndexRecommendations: record.IndexRecommendations,
	}

	return c.JSON(http.StatusOK, GetOptimizationRecordRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    &ret,
	})
}

func getOptimizationRecords(c echo.Context) error {
	req := new(GetOptimizationRecordsReq)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), false)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("check project manager failed: %v", err))
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_project_id":       projectUid,
		"filter_instance_name":    req.FilterInstanceName,
		"fuzzy_search":            req.FuzzySearch,
		"filter_create_time_from": req.FilterCreateTimeFrom,
		"filter_create_time_to":   req.FilterCreateTimeTo,
		"limit":                   req.PageSize,
		"offset":                  offset,
		"check_user_can_access":   !up.IsProjectAdmin(),
		"current_user":            user.Name,
	}

	records, total, err := s.GetOptimizationRecordsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ret := make([]OptimizationRecord, 0)
	for _, v := range records {
		ret = append(ret, OptimizationRecord{
			OptimizationName: v.OptimizationName,
			OptimizationID:   v.OptimizationId,
			InstanceNmae:     v.InstanceName,
			DBType:           v.DBType,
			PerformanceGain:  v.PerformanceImprove,
			CreatedTime:      v.CreatedAt,
			CreatedUser:      v.Creator,
		})
	}
	return c.JSON(http.StatusOK, GetOptimizationRecordsRes{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      ret,
		TotalNums: total,
	})
}

func getOptimizationSQL(c echo.Context) error {
	optimizationId := c.Param("optimization_record_id")

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO 权限校验
	optimizationSQL, err := model.GetStorage().GetOptimizationSQLById(optimizationId, number)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	trs := make([]RewriteRule, 0)
	for _, tr := range optimizationSQL.TriggeredRules {
		trs = append(trs, RewriteRule{
			RuleName:            tr.RuleName,
			Message:             tr.Message,
			RewrittenQueriesStr: tr.RewrittenQueriesStr,
			ViolatedQueriesStr:  tr.ViolatedQueriesStr,
		})
	}
	return c.JSON(http.StatusOK, GetOptimizationSQLRes{
		BaseRes: controller.NewBaseReq(nil),
		Data: &OptimizationSQLDetail{
			OriginalSQL:              optimizationSQL.OriginalSQL,
			OptimizedSQL:             optimizationSQL.OptimizedSQL,
			TriggeredRule:            trs,
			IndexRecommendations:     optimizationSQL.IndexRecommendations,
			ExplainValidationDetails: ExplainValidationDetail(optimizationSQL.ExplainValidationDetails),
		},
	})
}

func getOptimizationSQLs(c echo.Context) error {
	req := new(GetOptimizationSQLsReq)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}
	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"limit":           req.PageSize,
		"offset":          offset,
		"optimization_id": c.Param("optimization_record_id"),
	}
	sqls, total, err := model.GetStorage().GetOptimizationSQLsByReq(data)
	// TODO 权限校验
	ret := make([]OptimizationSQL, 0)
	for _, s := range sqls {
		ret = append(ret, OptimizationSQL{
			Number:              uint64(s.ID),
			OriginalSQL:         s.OriginalSQL,
			NumberOfRewrite:     s.NumberOfRewrite,
			NumberOfSyntaxError: s.NumberOfSyntaxError,
			NumberOfIndex:       s.NumberOfIndex,
			NumberOfHitIndex:    s.NumberOfHitIndex,
			Performance:         s.Performance,
			ContributingIndices: s.ContributingIndices,
		})
	}
	return c.JSON(http.StatusOK, GetOptimizationSQLsRes{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: total,
		Data:      ret,
	})
}
