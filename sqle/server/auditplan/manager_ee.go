//go:build enterprise
// +build enterprise

package auditplan

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

func (sap *SyncFromAuditPlan) SyncSqlManager() error {
	if sap.AuditReport == nil {
		return fmt.Errorf("schedule to audit Task failed, audit report is nil")
	}

	sqlApSqlMap := make(map[string]*model.AuditPlanSQLV2, len(sap.FilterSqls))
	for _, sql := range sap.FilterSqls {
		sqlApSqlMap[sql.SQLContent] = sql
	}

	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanById(sap.AuditReport.AuditPlanID)
	if err != nil {
		return fmt.Errorf("get audit plan failed, error: %v", err)
	}
	if !exist {
		return fmt.Errorf("audit plan %s not exist", ap.Name)
	}

	var sqlManageList []*model.SqlManage
	for _, reportSQL := range sap.Task.ExecuteSQLs {
		apSql := sqlApSqlMap[reportSQL.BaseSQL.Content]
		fp := apSql.Fingerprint
		instName := ap.InstanceName
		schema := ap.InstanceDatabase
		source := model.SQLManageSourceAuditPlan
		apInfo := apSql.Info

		var info = struct {
			Counter              *uint   `json:"counter"`
			LastReceiveTimestamp *string `json:"last_receive_timestamp"`
			FirstQueryAt         *string `json:"first_query_at"`
		}{}
		if len(apInfo) > 0 { // not null
			err = json.Unmarshal(apInfo, &info)
			if err != nil {
				return fmt.Errorf("unmarshal info failed, error: %v", err)
			}
		}

		var fpCount uint
		if info.Counter != nil {
			fpCount = *info.Counter
		}

		var firstAppearTime *time.Time
		if info.FirstQueryAt != nil {
			firstQueryAt, err := time.Parse("2006-01-02T15:04:05-07:00", *info.FirstQueryAt)
			if err != nil {
				return fmt.Errorf("parse first query at failed, error: %v", err)
			}
			firstAppearTime = &firstQueryAt
		}

		var lastReceiveTime *time.Time
		if info.LastReceiveTimestamp != nil {
			lastReceiveTimeStamp, err := time.Parse("2006-01-02T15:04:05-07:00", *info.LastReceiveTimestamp)
			if err != nil {
				return fmt.Errorf("parse last receive timestamp failed, error: %v", err)
			}
			lastReceiveTime = &lastReceiveTimeStamp
		}

		sqlManage, err := NewSqlManage(fp, reportSQL.BaseSQL.Content, schema, instName, source, reportSQL.AuditLevel, ap.ProjectId, ap.ID, firstAppearTime, lastReceiveTime, 0, fpCount, reportSQL.AuditResults, nil)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		sqlManageList = append(sqlManageList, sqlManage)
	}

	if err := s.InsertOrUpdateSqlManage(sqlManageList); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	var sqlManageList []*model.SqlManage

	s := model.GetStorage()
	sqlManageList, err := s.GetAllSqlManageList()
	if err != nil {
		return fmt.Errorf("get all sql manage list failed, error: %v", err)
	}

	md5SqlManageMap := make(map[string]*model.SqlManage, len(sqlManageList))
	for _, sqlManage := range sqlManageList {
		md5SqlManageMap[sqlManage.ProjFpSourceInstSchemaMd5] = sqlManage
	}

	for _, executeSQL := range sa.Task.ExecuteSQLs {
		sql := executeSQL.Content
		fp := sa.SqlFpMap[sql]
		schemaName := sa.Task.Schema
		instName := sa.Task.InstanceName()
		source := model.SQLManageSourceSqlAuditRecord

		sqlManage, err := NewSqlManage(fp, sql, schemaName, instName, source, executeSQL.AuditLevel, sa.ProjectId, 0, &executeSQL.CreatedAt, &executeSQL.CreatedAt, sa.SqlAuditRecordID, 0, executeSQL.AuditResults, md5SqlManageMap)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		sqlManageList = append(sqlManageList, sqlManage)
	}

	if err := s.InsertOrUpdateSqlManage(sqlManageList); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func GetSqlMangeMd5(projectId uint, fp string, schema string, instName string, source string, apID uint) (string, error) {
	md5Json, err := json.Marshal(
		struct {
			ProjectId   uint
			Fingerprint string
			Schema      string
			InstName    string
			Source      string
			ApID        uint
		}{
			ProjectId:   projectId,
			Fingerprint: fp,
			Schema:      schema,
			InstName:    instName,
			Source:      source,
			ApID:        apID,
		},
	)
	if err != nil {
		return "", fmt.Errorf("marshal sql identity failed, error: %v", err)
	}

	return utils.Md5String(string(md5Json)), nil
}

func NewSqlManage(fp, sql, schemaName, instName, source, auditLevel string, projectId, apId uint, createAt, LastReceiveAt *time.Time, sqlAuditRecordID, fpCount uint, auditResult model.AuditResults, md5SqlManageMap map[string]*model.SqlManage) (*model.SqlManage, error) {
	md5Str, err := GetSqlMangeMd5(projectId, fp, schemaName, instName, source, apId)
	if err != nil {
		return nil, fmt.Errorf("get sql manage md5 failed, error: %v", err)
	}

	sqlManage := &model.SqlManage{
		SqlFingerprint:            fp,
		SqlText:                   sql,
		ProjFpSourceInstSchemaMd5: md5Str,
		Source:                    source,
		ProjectId:                 projectId,
		SchemaName:                schemaName,
		InstanceName:              instName,
		SqlAuditRecordId:          sqlAuditRecordID,
		AuditLevel:                auditLevel,
		AuditResults:              auditResult,
		AuditPlanId:               apId,
		FirstAppearTimestamp:      createAt,
		LastReceiveTimestamp:      LastReceiveAt,
	}

	if source == model.SQLManageSourceSqlAuditRecord {
		if sql, ok := md5SqlManageMap[md5Str]; ok {
			md5SqlManageMap[md5Str].FpCount += 1
			sqlManage.FpCount = sql.FpCount
		} else {
			sqlManage.FpCount = 1
		}
	} else if source == model.SQLManageSourceAuditPlan {
		sqlManage.FpCount = uint64(fpCount)
	}

	return sqlManage, nil
}

func SyncToSqlManage(sqls []*SQL, ap *model.AuditPlan) error {
	var sqlManageList []*model.SqlManage
	for _, sql := range sqls {
		var firstQueryAtPtrFormat *time.Time
		var err error
		firstQueryAt, ok := sql.Info["first_query_at"]
		if ok && firstQueryAt != nil {
			firstQueryAtFormat, err := time.Parse("2006-01-02 15:04:05", firstQueryAt.(string))
			if err != nil {
				return fmt.Errorf("parse first query at failed, error: %v", err)
			}
			firstQueryAtPtrFormat = &firstQueryAtFormat
		}

		var lastReceiveAtPtrFormat *time.Time
		lastReceiveAt, ok := sql.Info["last_receive_timestamp"]
		if ok && lastReceiveAt != nil {
			lastReceiveAtFormat, err := time.Parse("2006-01-02 15:04:05", lastReceiveAt.(string))
			if err != nil {
				return fmt.Errorf("parse last receive timestamp failed, error: %v", err)
			}
			lastReceiveAtPtrFormat = &lastReceiveAtFormat
		}

		var countFormat uint64
		count, ok := sql.Info["counter"]
		if ok && count != 0 {
			countFormat = count.(uint64)
		}

		// todo: 更新审核等级
		sqlManage, err := NewSqlManage(sql.Fingerprint, sql.SQLContent, sql.Schema, "", model.SQLManageSourceAuditPlan, "",
			ap.ProjectId, ap.ID, firstQueryAtPtrFormat, lastReceiveAtPtrFormat, 0,
			uint(countFormat), model.AuditResults{model.AuditResult{Message: "未审核"}}, nil)
		if err != nil {
			return err
		}

		sqlManageList = append(sqlManageList, sqlManage)
	}

	s := model.GetStorage()
	// todo 计算count值
	if err := s.InsertOrUpdateSqlManageWithNotUpdateFpCount(sqlManageList); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}
