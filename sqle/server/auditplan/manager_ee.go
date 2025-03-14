//go:build enterprise
// +build enterprise

package auditplan

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

func (sap *SyncFromAuditPlan) SyncSqlManager() error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync audit plan recover err: %v", err)
		}
	}()

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
		schema := apSql.Schema
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

		sqlManage, err := NewSqlManage(fp, reportSQL.BaseSQL.Content, schema, instName, source, reportSQL.AuditLevel, string(ap.ProjectId), ap.ID, firstAppearTime, lastReceiveTime, fpCount, reportSQL.AuditResults, nil)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		sqlManageList = append(sqlManageList, sqlManage)
	}

	if err := s.InsertOrUpdateSqlManage(sqlManageList, 0); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync sql audit recover err: %v", err)
		}
	}()

	s := model.GetStorage()
	sqlManageList, err := s.GetAllSqlManageList()
	if err != nil {
		return fmt.Errorf("get all sql manage list failed, error: %v", err)
	}

	md5SqlManageMap := make(map[string]*model.SqlManage, len(sqlManageList))
	for _, sqlManage := range sqlManageList {
		md5SqlManageMap[sqlManage.ProjFpSourceInstSchemaMd5] = sqlManage
	}

	blacklist, err := s.GetBlackListByProjectID(model.ProjectUID(sa.ProjectId))
	if err != nil {
		return fmt.Errorf("get blacklist failed, error: %v", err)
	}

	var newSqlManageList []*model.SqlManage
	matchedCount := make(map[uint] /*blacklist id*/ uint /*count*/)
	for _, executeSQL := range sa.Task.ExecuteSQLs {
		sql := executeSQL.Content
		fp := sa.SqlFpMap[sql]
		schemaName := sa.Task.Schema
		instName := sa.Task.InstanceName()
		source := model.SQLManageSourceSqlAuditRecord
		instHost := sa.Task.InstanceHost()

		matchedID, isInBlacklist := FilterSQLsByBlackList([]string{instHost}, sql, fp, instName, "", blacklist)
		if isInBlacklist {
			matchedCount[matchedID]++
			continue
		}

		sqlManage, err := NewSqlManage(fp, sql, schemaName, instName, source, executeSQL.AuditLevel, sa.ProjectId, 0, &executeSQL.CreatedAt, &executeSQL.CreatedAt, 0, executeSQL.AuditResults, md5SqlManageMap)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		newSqlManageList = append(newSqlManageList, sqlManage)
	}

	lastMatchedTime := time.Now()
	if err := s.BatchUpdateBlackListCount(matchedCount, lastMatchedTime); err != nil {
		return fmt.Errorf("update blacklist failed, error: %v", err)
	}

	if err := s.InsertOrUpdateSqlManage(newSqlManageList, sa.SqlAuditRecordID); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func GetSqlMangeMd5(projectId string, fp string, schema string, instName string, source string, apID uint) (string, error) {
	md5Json, err := json.Marshal(
		struct {
			ProjectId   string
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

func NewSqlManage(fp, sql, schemaName, instName, source, auditLevel, projectId string, apId uint, createAt, LastReceiveAt *time.Time, fpCount uint, auditResult model.AuditResults, md5SqlManageMap map[string]*model.SqlManage) (*model.SqlManage, error) {
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
