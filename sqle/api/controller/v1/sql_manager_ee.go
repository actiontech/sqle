//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

func getSqlManageList(c echo.Context) error {
	req := new(GetSqlManageListReq)
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
		"filter_instance_name":              req.FilterInstanceID,
		"filter_source":                     req.FilterSource,
		"filter_audit_level":                req.FilterAuditLevel,
		"filter_last_audit_start_time_from": req.FilterLastAuditStartTimeFrom,
		"filter_last_audit_start_time_to":   req.FilterLastAuditStartTimeTo,
		"filter_status":                     req.FilterStatus,
		"project_id":                        projectUid,
		"filter_db_type":                    req.FilterDbType,
		"filter_rule_name":                  req.FilterRuleName,
		"fuzzy_search_endpoint":             req.FuzzySearchEndpoint,
		"fuzzy_search_schema_name":          req.FuzzySearchSchemaName,
		"sort_field":                        req.SortField,
		"sort_order":                        req.SortOrder,
		"limit":                             req.PageSize,
		"offset":                            offset,
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
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.Message,
				RuleName: ar.RuleName,
			})
		}

		source := &Source{
			SqlSourceType: sqlManage.Source,
			SqlSourceID:   sqlManage.SourceID,
		}
		auditPlanDesc := ConvertAuditPlanDescByType(sqlManage.Source)
		source.SqlSourceDesc = auditPlanDesc
		sqlMgr.Source = source

		if sqlManage.AppearTimestamp != nil {
			sqlMgr.FirstAppearTime = sqlManage.AppearTimestamp.Format("2006-01-02 15:04:05")
		}
		if sqlManage.LastReceiveTimestamp != nil {
			sqlMgr.LastAppearTime = sqlManage.LastReceiveTimestamp.Format("2006-01-02 15:04:05")
		}
		sqlMgr.AppearNum = sqlManage.FpCount
		if sqlManage.Assignees != nil {
			for _, assignees := range strings.Split(*sqlManage.Assignees, ",") {
				if v, ok := users[assignees]; ok {
					sqlMgr.Assignees = append(sqlMgr.Assignees, v.Name)
				}
			}
		}

		sqlMgr.Status = sqlManage.Status.String
		sqlMgr.Remark = sqlManage.Remark.String
		sqlMgr.Endpoint = sqlManage.Endpoints.String
		sqlManageRespList = append(sqlManageRespList, sqlMgr)
	}

	return sqlManageRespList, nil
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

	distinctSqlManageIDs := utils.RemoveDuplicatePtrUint64(req.SqlManageIdList)
	sqlManages, err := s.GetSqlManagerListByIDs(distinctSqlManageIDs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if len(sqlManages) != len(distinctSqlManageIDs) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("sql manage record not exist")))
	}

	err = s.BatchUpdateSqlManager(distinctSqlManageIDs, req.Status, req.Remark, req.Assignees)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

func exportSqlManagesV1(c echo.Context) error {
	req := new(ExportSqlManagesReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

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
		"fuzzy_search_endpoint":             req.FuzzySearchEndpoint,
		"fuzzy_search_schema_name":          req.FuzzySearchSchemaName,
		"sort_field":                        req.SortField,
		"sort_order":                        req.SortOrder,
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

	sqlManageResp, err := s.GetSqlManageListByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	csvWriter := csv.NewWriter(buff)

	err = csvWriter.WriteAll([][]string{
		{"SQL总数", strconv.FormatUint(sqlManageResp.SqlManageTotalNum, 10)},
		{"问题SQL数", strconv.FormatUint(sqlManageResp.SqlManageBadNum, 10)},
		{"已优化SQL数", strconv.FormatUint(sqlManageResp.SqlManageOptimizedNum, 10)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := csvWriter.Write([]string{
		"SQL指纹",
		"SQL",
		"来源",
		"数据源",
		"SCHEMA",
		"审核结果",
		"端点信息",
		"负责人",
		"状态",
		"备注",
	}); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	users, err := dms.GetMapUsers(c.Request().Context(), nil, dms.GetDMSServerAddress())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for _, sqlManage := range sqlManageResp.SqlManageList {
		var assignees []string
		if sqlManage.Assignees != nil {
			for _, assignee := range strings.Split(*sqlManage.Assignees, ",") {
				if user, ok := users[assignee]; ok {
					assignees = append(assignees, user.Name)
				}
			}
		}
		var newRow []string
		newRow = append(
			newRow,
			sqlManage.SqlFingerprint,
			sqlManage.SqlText,
			ConvertAuditPlanDescByType(sqlManage.Source),
			dms.GetInstancesByIdWithoutError(sqlManage.InstanceID).Name,
			sqlManage.SchemaName,
			spliceAuditResults(sqlManage.AuditResults),
			sqlManage.Endpoints.String,
			strings.Join(assignees, ","),
			model.SqlManageStatusMap[sqlManage.Status.String],
			sqlManage.Remark.String,
		)

		if err := csvWriter.Write(newRow); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	fileName := fmt.Sprintf("%s_SQL管控.csv", time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{
		"filename": fileName,
	}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

func getSqlManageRuleTips(c echo.Context) error {
	s := model.GetStorage()

	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlManageRuleTips, err := s.GetSqlManagerRuleTips(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlManageRuleTipsResp{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRuleTipsToResp(sqlManageRuleTips),
	})
}

func convertRuleTipsToResp(tips []*model.SqlManageRuleTips) []RuleTips {
	m := make(map[string] /*数据库类型*/ []RuleRespV1)
	for _, tip := range tips {
		m[tip.DbType] = append(m[tip.DbType], RuleRespV1{
			RuleName: tip.RuleName,
			Desc:     tip.Desc,
		})
	}

	var ruleResp []RuleTips
	for dbType, rule := range m {
		ruleResp = append(ruleResp, RuleTips{
			DbType: dbType,
			Rule:   rule,
		})
	}

	return ruleResp
}

func getSqlManageSqlAnalysisV1(c echo.Context) error {
	mgID := c.Param("sql_manage_id")

	s := model.GetStorage()
	omg, exist, err := s.GetOriginManageSqlByID(mgID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr(fmt.Sprintf("sql manage id %v not exist", mgID)))
	}

	instance, exist, err := dms.GetInstancesById(c.Request().Context(), omg.InstanceID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr(fmt.Sprintf("sql manage id %v instance not exist", mgID)))
	}

	entry := log.NewEntry().WithField("sql_manage_analysis", mgID)
	analysisResp, err := GetSQLAnalysisResult(entry, instance, omg.SchemaName, omg.SqlText)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlManageSqlAnalysisResp{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(analysisResp, omg.SqlText),
	})
}

func convertSQLAnalysisResultToRes(res *AnalysisResult, rawSQL string) *SqlAnalysis {
	data := &SqlAnalysis{}

	// explain
	{
		data.SQLExplain = &SQLExplain{SQL: rawSQL}
		if res.ExplainResultErr != nil {
			data.SQLExplain.Message = res.ExplainResultErr.Error()
		} else {
			classicResult := ExplainClassicResult{
				Head: make([]TableMetaItemHeadResV1, len(res.ExplainResult.ClassicResult.Columns)),
				Rows: make([]map[string]string, len(res.ExplainResult.ClassicResult.Rows)),
			}

			// head
			for i := range res.ExplainResult.ClassicResult.Columns {
				col := res.ExplainResult.ClassicResult.Columns[i]
				classicResult.Head[i].FieldName = col.Name
				classicResult.Head[i].Desc = col.Desc
			}

			// rows
			for i := range res.ExplainResult.ClassicResult.Rows {
				row := res.ExplainResult.ClassicResult.Rows[i]
				classicResult.Rows[i] = make(map[string]string, len(row))
				for k := range row {
					colName := res.ExplainResult.ClassicResult.Columns[k].Name
					classicResult.Rows[i][colName] = row[k]
				}
			}
			data.SQLExplain.ClassicResult = classicResult
		}
	}

	// table_metas
	{
		data.TableMetas = &TableMetas{}
		if res.TableMetaResultErr != nil {
			data.TableMetas.ErrMessage = res.TableMetaResultErr.Error()
		} else {
			tableMetaItemsData := make([]*TableMeta, len(res.TableMetaResult.TableMetas))
			for i := range res.TableMetaResult.TableMetas {
				tableMeta := res.TableMetaResult.TableMetas[i]
				tableMetaColumnsInfo := tableMeta.ColumnsInfo
				tableMetaIndexInfo := tableMeta.IndexesInfo
				tableMetaItemsData[i] = &TableMeta{}
				tableMetaItemsData[i].Columns = TableColumns{
					Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
					Head: make([]TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Columns)),
				}

				tableMetaItemsData[i].Indexes = TableIndexes{
					Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
					Head: make([]TableMetaItemHeadResV1, len(tableMetaIndexInfo.Columns)),
				}

				tableMetaColumnData := tableMetaItemsData[i].Columns
				for j := range tableMetaColumnsInfo.Columns {
					col := tableMetaColumnsInfo.Columns[j]
					tableMetaColumnData.Head[j].FieldName = col.Name
					tableMetaColumnData.Head[j].Desc = col.Desc
				}

				for j := range tableMetaColumnsInfo.Rows {
					tableMetaColumnData.Rows[j] = make(map[string]string, len(tableMetaColumnsInfo.Rows[j]))
					for k := range tableMetaColumnsInfo.Rows[j] {
						colName := tableMetaColumnsInfo.Columns[k].Name
						tableMetaColumnData.Rows[j][colName] = tableMetaColumnsInfo.Rows[j][k]
					}
				}

				tableMetaIndexData := tableMetaItemsData[i].Indexes
				for j := range tableMetaIndexInfo.Columns {
					tableMetaIndexData.Head[j].FieldName = tableMetaIndexInfo.Columns[j].Name
					tableMetaIndexData.Head[j].Desc = tableMetaIndexInfo.Columns[j].Desc
				}

				for j := range tableMetaIndexInfo.Rows {
					tableMetaIndexData.Rows[j] = make(map[string]string, len(tableMetaIndexInfo.Rows[j]))
					for k := range tableMetaIndexInfo.Rows[j] {
						colName := tableMetaIndexInfo.Columns[k].Name
						tableMetaIndexData.Rows[j][colName] = tableMetaIndexInfo.Rows[j][k]
					}
				}

				tableMetaItemsData[i].Name = tableMeta.Name
				tableMetaItemsData[i].Schema = tableMeta.Schema
				tableMetaItemsData[i].CreateTableSQL = tableMeta.CreateTableSQL
				tableMetaItemsData[i].Message = tableMeta.Message
			}
			data.TableMetas.Items = tableMetaItemsData
		}
	}

	// performance_statistics
	{
		data.PerformanceStatistics = &PerformanceStatistics{}

		// affect_rows
		data.PerformanceStatistics.AffectRows = &AffectRows{}
		if res.AffectRowsResultErr != nil {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResultErr.Error()
		} else {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResult.ErrMessage
			data.PerformanceStatistics.AffectRows.Count = int(res.AffectRowsResult.Count)
		}

	}

	return data
}

func getUnsolvedSQLCountByManagerStatus(id uint) (int64, error) {
	s := model.GetStorage()
	count, err := s.GetUnsolvedSQLCount(id,
		[]string{model.SQLManageStatusIgnored,
			model.SQLManageStatusSolved,
			model.SQLManageStatusManualAudited})
	if err != nil {
		return count, err
	}
	return count, nil
}

func ConvertAuditPlanDescByType(auditPlanType string) string {
	for _, meta := range auditplan.Metas {
		if meta.Type == auditPlanType {
			return meta.Desc
		}
	}
	return ""
}
