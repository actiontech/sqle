//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
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

	ctx := c.Request().Context()
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

		for i := range sqlManage.AuditResults {
			ar := sqlManage.AuditResults[i]
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.GetAuditMsgByLangTag(lang),
				RuleName: ar.RuleName,
			})
		}

		source := &Source{
			SqlSourceType: sqlManage.Source.String,
			SqlSourceIDs:  sqlManage.SourceIDs,
		}
		auditPlanDesc := ConvertSqlSourceDescByType(ctx, sqlManage.Source.String)
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

	err = s.BatchUpdateSqlManager(distinctSqlManageIDs, req.Status, req.Remark, req.Priority, req.Assignees)
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

	ctx := c.Request().Context()
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
		"filter_priority":                   req.FilterPriority,
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
		{locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportTotalSQLCount), strconv.FormatUint(sqlManageResp.SqlManageTotalNum, 10)},         // SQL总数
		{locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportProblemSQLCount), strconv.FormatUint(sqlManageResp.SqlManageBadNum, 10)},         // 问题SQL数
		{locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportOptimizedSQLCount), strconv.FormatUint(sqlManageResp.SqlManageOptimizedNum, 10)}, // 已优化SQL数
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := csvWriter.Write([]string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportSQLFingerprint), // "SQL指纹",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportSQL),            // "SQL",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportSource),         // "来源",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportDataSource),     // "数据源",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportSCHEMA),         // "SCHEMA",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportAuditResult),    // "审核结果",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportEndpoint),       // "端点信息",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportPersonInCharge), // "负责人",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportState),          // "状态",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.SMExportRemarks),        // "备注",
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
			sqlManage.SqlFingerprint.String,
			sqlManage.SqlText.String,
			ConvertSqlSourceDescByType(ctx, sqlManage.Source.String),
			dms.GetInstancesByIdWithoutError(sqlManage.InstanceID.String).Name,
			sqlManage.SchemaName.String,
			spliceAuditResults(ctx, sqlManage.AuditResults),
			sqlManage.Endpoints.String,
			strings.Join(assignees, ","),
			locale.Bundle.LocalizeMsgByCtx(ctx, model.SqlManageStatusMap[sqlManage.Status.String]),
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

	fileName := fmt.Sprintf("%s_sql_manager.csv", time.Now().Format("20060102150405"))
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
		Data:    convertRuleTipsToResp(c.Request().Context(), sqlManageRuleTips),
	})
}

func convertRuleTipsToResp(ctx context.Context, tips []*model.SqlManageRuleTips) []RuleTips {
	m := make(map[string] /*数据库类型*/ []RuleRespV1)
	for _, tip := range tips {
		m[tip.DbType] = append(m[tip.DbType], RuleRespV1{
			RuleName: tip.RuleName,
			Desc:     tip.I18nRuleInfo.GetRuleInfoByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)).Desc,
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
		Data:    convertSQLAnalysisResultToRes(c.Request().Context(), analysisResp, omg.SqlText),
	})
}

func getAuditPlanUnsolvedSQLCount(auditPlanId uint) (int64, error) {
	s := model.GetStorage()
	count, err := s.GetAuditPlanUnsolvedSQLCount(auditPlanId,
		[]string{model.SQLManageStatusIgnored,
			model.SQLManageStatusSolved,
			model.SQLManageStatusManualAudited})
	if err != nil {
		return count, err
	}
	return count, nil
}

func ConvertSqlSourceDescByType(ctx context.Context, source string) string {
	if source == model.SQLManageSourceSqlAuditRecord {
		return locale.Bundle.LocalizeMsgByCtx(ctx, model.SqlManageSourceMap[source])
	}
	for _, meta := range auditplan.Metas {
		if meta.Type == source {
			return locale.Bundle.LocalizeMsgByCtx(ctx, meta.Desc)
		}
	}
	return ""
}

func getGlobalSqlManageList(c echo.Context) error {
	req := new(GetGlobalSqlManageListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var err error
	// 1. 权限控制
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 1.1. 获取用户的所有权限信息
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(c.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)

	// 2.2 页面筛选项：如果根据项目优先级筛选，则先筛选出对应优先级下的项目
	filter, err := constructGlobalSqlManageBasicFilter(c.Request().Context(), user, userVisibility, permissions,
		&globalSqlManageBasicFilter{
			FilterProjectUid:      req.FilterProjectUid,
			FilterInstanceId:      req.FilterInstanceId,
			FilterProjectPriority: req.FilterProjectPriority,
		},
	)
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	filter["limit"] = limit
	filter["offset"] = offset

	// 3. 根据筛选项筛选SQL管控的SQL信息
	s := model.GetStorage()
	modelGlobalSqlManages, total, err := s.GetGlobalSqlManageList(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var projectMap = make(map[string]*dmsV1.ListProject)
	// 3.1. 若未根据项目优先级筛选，需要根据SQL拉取对应的项目信息
	if req.FilterProjectPriority == nil {
		projectMap, err = loadProjectsBySqlManage(c.Request().Context(), modelGlobalSqlManages)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		_, projectMap, err = loadProjectsByPriority(c.Request().Context(), *req.FilterProjectPriority)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	// 3.2. 需要根据SQL拉取对应的数据源信息
	instanceMap, err := loadInstancesBySqlManage(c.Request().Context(), modelGlobalSqlManages)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetGlobalSqlManageListResp{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      toGlobalSqlManage(c.Request().Context(), modelGlobalSqlManages, projectMap, instanceMap),
		TotalNums: total,
	})
}

func getGlobalSqlManageStatistics(c echo.Context) error {
	req := new(GetGlobalSqlManageStatisticsReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var err error
	// 1. 权限控制
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 1.1. 获取用户的所有权限信息
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(c.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)
	// 1.3. 若用户不能查看所有待关注SQL，则需要判断用户是否拥有多项目查看待关注SQL的权限，可以看到所有待关注SQL的用户不需要判断项目范围
	filter, err := constructGlobalSqlManageBasicFilter(c.Request().Context(), user, userVisibility, permissions,
		&globalSqlManageBasicFilter{
			FilterProjectUid:      req.FilterProjectUid,
			FilterInstanceId:      req.FilterInstanceId,
			FilterProjectPriority: req.FilterProjectPriority,
		})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 3. 根据筛选项筛选SQL管控的SQL信息
	s := model.GetStorage()
	total, err := s.GetGlobalSqlManageStatics(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetGlobalSqlManageStatisticsResp{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: total,
	})
}

type globalSqlManageBasicFilter struct {
	FilterProjectUid      *string                `query:"filter_project_uid" json:"filter_project_uid,omitempty"`
	FilterInstanceId      *string                `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterProjectPriority *dmsV1.ProjectPriority `query:"filter_project_priority" json:"filter_project_priority,omitempty" enums:"high,medium,low"`
}

func constructGlobalSqlManageBasicFilter(ctx context.Context, user *model.User, userVisibility GlobalDashBoardVisibility, permissions []dmsV1.OpPermissionItem, req *globalSqlManageBasicFilter) (map[string]interface{}, error) {
	// 角色：多项目管理者 canViewProjects
	// 1.3. 若用户不能查看所有待关注SQL，则需要判断用户是否拥有多项目查看待关注SQL的权限，可以看到所有待关注SQL的用户不需要判断项目范围
	var projectIdsOfProjectAdmin []string
	if userVisibility == GlobalDashBoardVisibilityProjects {
		for _, permission := range permissions {
			if permission.OpPermissionType == dmsV1.OpPermissionTypeProjectAdmin {
				projectIdsOfProjectAdmin = append(projectIdsOfProjectAdmin, permission.RangeUids...)
			}
		}
	}
	var canViewProjects bool = len(projectIdsOfProjectAdmin) > 0
	// 角色：SQL被分配者 !canViewGlobal && !canViewProjects

	// 2. 组织筛选项
	// 2.1. 基本筛选项

	data := map[string]interface{}{
		"filter_instance_id": req.FilterInstanceId, // 页面筛选项
		"filter_project_uid": req.FilterProjectUid, // 页面筛选项
	}

	// 2.2 页面筛选项：如果根据项目优先级筛选，则先筛选出对应优先级下的项目
	var projectIdsByPriority []string
	var err error
	if req.FilterProjectPriority != nil {
		projectIdsByPriority, _, err = loadProjectsByPriority(ctx, *req.FilterProjectPriority)
		if err != nil {
			return nil, err
		}
	}

	if req.FilterProjectPriority != nil {
		if userVisibility == GlobalDashBoardVisibilityProjects {
			// 2.2.1 若根据项目优先级筛选，且可以查看多项目待关注SQL，则将可查看的项目和项目优先级筛选后的项目的集合取交集
			data["filter_project_id_list"] = utils.IntersectionStringSlice(projectIdsByPriority, projectIdsOfProjectAdmin)
		} else {
			// 2.2.2 若只根据项目优先级筛选，则根据优先级对应的项目筛选
			data["filter_project_id_list"] = projectIdsByPriority

		}
	}
	// 2.3 若不根据项目优先级筛选
	if req.FilterProjectPriority == nil && canViewProjects {
		// 2.3.1 若可以查看多项目待关注SQL，则通过用户的有权限的项目进行筛选
		data["filter_project_id_list"] = projectIdsOfProjectAdmin
	}
	// 角色：SQL被分配者
	// 2.4. 若用户既不能查看所有待关注SQL也不能查看多项目待关注SQL，则可以查看在SQL管控中分配给他的SQL
	if userVisibility == GlobalDashBoardVisibilityAssignee {
		data["filter_assignees_id"] = user.GetIDStr()
	}
	return data, nil
}

func loadProjectsBySqlManage(ctx context.Context, modelGlobalSqlManages []*model.GlobalSqlManage) (projectMap map[string]*dmsV1.ListProject, err error) {
	if len(modelGlobalSqlManages) == 0 {
		return make(map[string]*dmsV1.ListProject), nil
	}
	projectMap = make(map[string]*dmsV1.ListProject)
	projectIds := []string{}
	for _, modelSqlManage := range modelGlobalSqlManages {
		if _, exist := projectMap[modelSqlManage.ProjectUid]; !exist {
			projectIds = append(projectIds, modelSqlManage.ProjectUid)
			projectMap[modelSqlManage.ProjectUid] = nil
		}
	}
	projectMap, err = loadProjectsByProjectIds(ctx, projectIds)
	if err != nil {
		return nil, err
	}
	return projectMap, nil
}

func loadInstancesBySqlManage(ctx context.Context, modelGlobalSqlManages []*model.GlobalSqlManage) (instanceMap map[string]*dmsV1.ListDBService, err error) {
	if len(modelGlobalSqlManages) == 0 {
		return make(map[string]*dmsV1.ListDBService), nil
	}
	instanceMap = make(map[string]*dmsV1.ListDBService)
	instanceIds := []string{}
	for _, modelSqlManage := range modelGlobalSqlManages {
		if modelSqlManage.InstanceID.Valid {
			if _, exist := instanceMap[modelSqlManage.InstanceID.String]; !exist {
				instanceIds = append(instanceIds, modelSqlManage.InstanceID.String)
				instanceMap[modelSqlManage.InstanceID.String] = nil
			}
		}
	}
	instanceMap, err = loadInstanceByInstanceIds(ctx, instanceIds)
	if err != nil {
		return nil, err
	}
	return instanceMap, nil
}

func toGlobalSqlManage(ctx context.Context, modelGlobalSqlManages []*model.GlobalSqlManage, projectMap map[string]*dmsV1.ListProject, instanceMap map[string]*dmsV1.ListDBService) []*GlobalSqlManage {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	ret := make([]*GlobalSqlManage, 0, len(modelGlobalSqlManages))
	for _, mg := range modelGlobalSqlManages {
		sqlMgr := &GlobalSqlManage{
			Id:                   mg.Id,
			Sql:                  mg.SqlText.String,
			ProjectName:          projectMap[mg.ProjectUid].Name,
			ProjectUid:           mg.ProjectUid,
			Status:               mg.Status.String,
			FirstAppearTimeStamp: mg.FirstAppearTimestamp.Format(time.RFC3339),
			ProjectPriority:      projectMap[mg.ProjectUid].ProjectPriority,
			ProblemDescriptions:  problemDescriptions(mg.Info),
		}
		if mg.InstanceID.Valid {
			sqlMgr.InstanceName = instanceMap[mg.InstanceID.String].Name
			sqlMgr.InstanceId = mg.InstanceID.String
		}
		for i := range mg.AuditResults {
			ar := mg.AuditResults[i]
			sqlMgr.AuditResult = append(sqlMgr.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.GetAuditMsgByLangTag(lang),
				RuleName: ar.RuleName,
			})
		}
		sqlMgr.Source = &Source{
			SqlSourceType: mg.Source.String,
			SqlSourceIDs:  mg.SourceIDs,
			SqlSourceDesc: ConvertSqlSourceDescByType(ctx, mg.Source.String),
		}
		ret = append(ret, sqlMgr)
	}
	return ret
}

func problemDescriptions(info string) []string {
	var infoJson struct {
		Counter         uint    `json:"counter"`
		QueryTimeAvg    float32 `json:"query_time_avg"`
		QueryTimeMax    float32 `json:"query_time_max"`
		RowsExaminedAvg float32 `json:"row_examined_avg"`
	}

	err := json.Unmarshal([]byte(info), &infoJson)
	if err != nil {
		return []string{info}
	}
	return []string{
		fmt.Sprintf("SQL出现次数 %v", infoJson.Counter),
		fmt.Sprintf("平均执行时间 %v", infoJson.QueryTimeAvg),
		fmt.Sprintf("最大执行时间 %v", infoJson.QueryTimeMax),
		fmt.Sprintf("平均扫描行数 %v", infoJson.RowsExaminedAvg),
	}
}

func getGlobalSqlManageStatistics(c echo.Context) error {
	return nil
}
