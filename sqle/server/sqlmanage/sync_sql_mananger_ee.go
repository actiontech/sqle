//go:build enterprise
// +build enterprise

package sqlmanage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/actiontech/sqle/sqle/utils"
)

func (sa *SyncFromSqlAuditRecord) SyncSqlManager(source string) error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync sql audit recover err: %v", err)
		}
	}()

	s := model.GetStorage()
	blacklist, err := s.GetBlackListByProjectID(model.ProjectUID(sa.ProjectId))
	if err != nil {
		return fmt.Errorf("get blacklist failed, error: %v", err)
	}
	var newSqlManageList []*model.SQLManageRecord
	matchedCount := make(map[uint] /*blacklist id*/ uint /*count*/)
	for _, executeSQL := range sa.Task.ExecuteSQLs {
		sql := executeSQL.Content
		fp := sa.SqlFpMap[sql]
		schemaName := sa.Task.Schema
		instId := sa.Task.InstanceId
		instName := sa.Task.InstanceName()
		instHost := sa.Task.InstanceHost()

		matchedID, isInBlacklist := auditplan.FilterSQLsByBlackList(instHost, sql, fp, instName, blacklist)
		if isInBlacklist {
			matchedCount[matchedID]++
			continue
		}
		sqlManage, err := NewSqlManageRecord(fp, sql, schemaName, strconv.FormatUint(instId, 10), instName, source, sa.SqlAuditRecordID, executeSQL.AuditLevel, sa.ProjectId, executeSQL.AuditResults)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		newSqlManageList = append(newSqlManageList, sqlManage)
	}

	if err := s.InsertOrUpdateSqlManageRecord(newSqlManageList, sa.SqlAuditRecordID); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func (sa *SyncFromSqlAuditRecord) UpdateSqlManageRecord(sourceId, source string) error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync sql audit recover err: %v", err)
		}
	}()
	s := model.GetStorage()
	sqlAuditRecords, exist, err := s.GetSqlManageSqlAuditRecordByRecordId(sourceId)
	if err != nil || !exist {
		return fmt.Errorf("record does not exist or get manage sql error, error: %v", err)
	}
	for _, record := range sqlAuditRecords {
		aggSourceId, err := AggregateSourceIdsBysqlId(record.SQLID, sourceId, false)
		if err != nil {
			return fmt.Errorf("aggregate source id failed, error: %v", err)
		}

		err = s.UpdateSqlManageRecord(record.SQLID, sourceId, aggSourceId, source)
		if err != nil {
			return fmt.Errorf("update sql manage record failed, error: %v", err)
		}
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

func NewSqlManageRecord(fp, sql, schemaName, instId, instName, source, sourceId, auditLevel, projectId string, auditResult model.AuditResults) (*model.SQLManageRecord, error) {
	sqlId := genSQLId(fp, schemaName, instName, source, projectId)
	sourceID, err := AggregateSourceIdsBysqlId(sqlId, sourceId, true)
	if err != nil {
		return nil, err
	}
	sqlManage := &model.SQLManageRecord{
		Source:         source,
		SourceId:       sourceID,
		ProjectId:      projectId,
		InstanceID:     instId,
		SchemaName:     schemaName,
		SqlFingerprint: fp,
		SqlText:        sql,
		AuditLevel:     auditLevel,
		AuditResults:   auditResult,
		SQLID:          sqlId,
	}
	return sqlManage, nil
}

func genSQLId(fp, schemaName, instName, source, projectId string) string {
	md5Json, err := json.Marshal(
		struct {
			ProjectId   string
			Fingerprint string
			Schema      string
			InstName    string
			Source      string
		}{
			ProjectId:   projectId,
			Fingerprint: fp,
			Schema:      schemaName,
			InstName:    instName,
			Source:      source,
		},
	)
	if err != nil {
		log.NewEntry().Errorf("new sql manage record gen sql id err: %v", err)
		return fp
	} else {
		return utils.Md5String(string(md5Json))
	}
}

func AggregateSourceIdsBysqlId(sqlId, sourceId string, isAddSourceId bool) (string, error) {
	s := model.GetStorage()
	sqlAuditRecord, exist, err := s.GetSqlManageSqlAuditRecordBySqlId(sqlId)
	if err != nil {
		return "", err
	}
	if !exist || len(sqlAuditRecord) == 0 {
		return sourceId, nil
	}
	sourceIds := make([]string, 0, len(sqlAuditRecord))
	for _, record := range sqlAuditRecord {
		if record.SqlAuditRecordId != sourceId {
			sourceIds = append(sourceIds, record.SqlAuditRecordId)
		}
	}
	if isAddSourceId {
		sourceIds = append(sourceIds, sourceId)
	}
	return strings.Join(sourceIds, ","), nil
}
