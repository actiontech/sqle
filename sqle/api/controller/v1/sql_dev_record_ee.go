//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

func getSqlDEVRecordList(c echo.Context) error {
	req := new(GetSqlDEVRecordListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
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
		"fuzzy_search_sql_fingerprint":  searchSqlFingerprint,
		"filter_instance_name":          req.FilterInstanceName,
		"filter_source":                 req.FilterSource,
		"filter_creator":                req.FilterCreator,
		"filter_last_receive_time_from": req.FilterLastReceiveTimeFrom,
		"filter_last_receive_time_to":   req.FilterLastReceiveTimeTo,
		"project_id":                    projectID,
		"fuzzy_search_schema_name":      req.FuzzySearchSchemaName,
		"sort_field":                    req.SortField,
		"sort_order":                    req.SortOrder,
		"limit":                         req.PageSize,
		"offset":                        offset,
	}

	ctx := c.Request().Context()
	s := model.GetStorage()
	sqlDEVRecords, total, err := s.GetSqlDEVRecordListByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlDEVRecordRet, err := convertToGetSqlDEVRecordListResp(ctx, sqlDEVRecords)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlDEVRecordListResp{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      sqlDEVRecordRet,
		TotalNums: total,
	})
}

func convertToGetSqlDEVRecordListResp(ctx context.Context, sqlDEVRecordList []*model.SQLDevRecord) ([]*SqlDEVRecord, error) {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	sqlDEVRecordRespList := make([]*SqlDEVRecord, 0, len(sqlDEVRecordList))
	for _, sqlDEVRecord := range sqlDEVRecordList {
		sqlDEV := new(SqlDEVRecord)
		sqlDEV.Id = uint64(sqlDEVRecord.ID)
		sqlDEV.SqlFingerprint = sqlDEVRecord.SqlFingerprint
		sqlDEV.Sql = sqlDEVRecord.SqlText
		sqlDEV.InstanceName = sqlDEVRecord.InstanceName
		sqlDEV.SchemaName = sqlDEVRecord.SchemaName
		sqlDEV.Creator = sqlDEVRecord.Creator
		sqlDEV.FpCount = sqlDEVRecord.FpCount
		sqlDEV.Source = &RecordSource{Name: sqlDEVRecord.Source}

		for i := range sqlDEVRecord.AuditResults {
			ar := sqlDEVRecord.AuditResults[i]
			sqlDEV.AuditResult = append(sqlDEV.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.GetAuditMsgByLangTag(lang),
				RuleName: ar.RuleName,
			})
		}

		if sqlDEVRecord.FirstAppearTimestamp != nil {
			sqlDEV.FirstAppearTimeStamp = sqlDEVRecord.FirstAppearTimestamp.Format("2006-01-02 15:04:05")
		}
		if sqlDEVRecord.LastReceiveTimestamp != nil {
			sqlDEV.LastReceiveTimeStamp = sqlDEVRecord.LastReceiveTimestamp.Format("2006-01-02 15:04:05")
		}

		sqlDEVRecordRespList = append(sqlDEVRecordRespList, sqlDEV)
	}

	return sqlDEVRecordRespList, nil
}

func SyncSqlDevRecord(ctx context.Context, task *model.Task, projectId, creator string) {
	logger := log.NewEntry().WithField("sync_sql_dev_record", creator)

	err := syncSqlDevRecord(ctx, task, projectId, creator)
	if err != nil {
		logger.Errorf("sync sql dev record failed, err: %v", err)
	}
}
func syncSqlDevRecord(ctx context.Context, task *model.Task, projectId, creator string) error {
	if task == nil || task.ExecuteSQLs == nil {
		return fmt.Errorf("sql audit task is nil")
	}

	plugin, err := common.NewDriverManagerWithoutCfg(log.NewEntry(), task.DBType)
	if err != nil {
		return fmt.Errorf("open plugin failed: %v", err)
	}
	defer plugin.Close(ctx)

	s := model.GetStorage()

	var sqlDevRecords []*model.SQLDevRecord
	for _, executeSQL := range task.ExecuteSQLs {
		node, err := plugin.Parse(ctx, executeSQL.Content)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
		sql := executeSQL.Content
		schemaName := task.Schema
		instName := task.InstanceName()
		source := model.SQLDevRecordSourceIDEPlugin
		createdAt := time.Now()

		sqlDevRecord, err := NewSqlDevRecord(node[0].Fingerprint, sql, schemaName, instName, source, executeSQL.AuditLevel,
			projectId, creator, &createdAt, &createdAt, 1, executeSQL.AuditResults)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		sqlDevRecords = append(sqlDevRecords, sqlDevRecord)
	}

	if err := s.InsertOrUpdateSqlDevRecord(sqlDevRecords); err != nil {
		return fmt.Errorf("insert or update sql dev record failed, error: %v", err)
	}

	return nil
}

func GetSqlDevRecordMd5(projectId string, fp string, schema string, instName string, source string, creator string) (string, error) {
	md5Json, err := json.Marshal(
		struct {
			ProjectId   string
			Fingerprint string
			Schema      string
			InstName    string
			Source      string
			Creator     string
		}{
			ProjectId:   projectId,
			Fingerprint: fp,
			Schema:      schema,
			InstName:    instName,
			Source:      source,
			Creator:     creator,
		},
	)
	if err != nil {
		return "", fmt.Errorf("marshal sql identity failed, error: %v", err)
	}

	return utils.Md5String(string(md5Json)), nil
}

func NewSqlDevRecord(fp, sql, schemaName, instName, source, auditLevel, projectId string, creator string, createAt, LastReceiveAt *time.Time, fpCount uint, auditResult model.AuditResults) (*model.SQLDevRecord, error) {
	md5Str, err := GetSqlDevRecordMd5(projectId, fp, schemaName, instName, source, creator)
	if err != nil {
		return nil, fmt.Errorf("get sql manage md5 failed, error: %v", err)
	}

	SqlDevRecord := &model.SQLDevRecord{
		SqlFingerprint:            fp,
		SqlText:                   sql,
		ProjFpSourceInstSchemaMd5: md5Str,
		Source:                    source,
		ProjectId:                 projectId,
		SchemaName:                schemaName,
		InstanceName:              instName,
		AuditLevel:                auditLevel,
		AuditResults:              auditResult,
		Creator:                   creator,
		FirstAppearTimestamp:      createAt,
		LastReceiveTimestamp:      LastReceiveAt,
		FpCount:                   uint64(fpCount),
	}

	return SqlDevRecord, nil
}
