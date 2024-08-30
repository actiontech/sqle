//go:build enterprise
// +build enterprise

package sqlmanage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/actiontech/sqle/sqle/utils"
)

// 同步sql至sql管控
func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync sql audit recover err: %v", err)
		}
	}()

	plugin, err := common.NewDriverManagerWithoutCfg(log.NewEntry(), sa.Task.DBType)
	if err != nil {
		return fmt.Errorf("open plugin failed: %v", err)
	}
	defer plugin.Close(context.TODO())

	s := model.GetStorage()
	blacklist, err := s.GetBlackListByProjectID(model.ProjectUID(sa.ProjectId))
	if err != nil {
		return fmt.Errorf("get blacklist failed, error: %v", err)
	}
	var newSqlManageList []*model.SQLManageRecord
	matchedCount := make(map[uint] /*blacklist id*/ uint /*count*/)
	for _, executeSQL := range sa.Task.ExecuteSQLs {
		node, err := plugin.Parse(context.TODO(), executeSQL.Content)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
		fp := node[0].Fingerprint
		sql := executeSQL.Content
		schemaName := sa.Task.Schema
		instId := sa.Task.InstanceId
		instName := sa.Task.InstanceName()
		instHost := sa.Task.InstanceHost()

		matchedID, isInBlacklist := auditplan.FilterSQLsByBlackList(instHost, sql, fp, instName, blacklist)
		if isInBlacklist {
			matchedCount[matchedID]++
			continue
		}
		sqlManage, err := NewSqlManageRecord(fp, sql, schemaName, strconv.FormatUint(instId, 10), instName, sa.Source, sa.SqlAuditRecordID, executeSQL.AuditLevel, sa.ProjectId, executeSQL.AuditResults)
		if err != nil {
			return fmt.Errorf("create or update sql manage failed, error: %v", err)
		}

		newSqlManageList = append(newSqlManageList, sqlManage)
	}
	lastMatchedTime := time.Now()
	if err := s.BatchUpdateBlackListCount(matchedCount, lastMatchedTime); err != nil {
		return fmt.Errorf("update blacklist failed, error: %v", err)
	}
	if err := s.InsertOrUpdateSqlManageRecord(newSqlManageList); err != nil {
		return fmt.Errorf("insert or update sql manage failed, error: %v", err)
	}

	return nil
}

func (sa *SyncFromSqlAuditRecord) UpdateSqlManageRecord() error {
	defer func() {
		if err := recover(); err != nil {
			log.NewEntry().Errorf("sync sql audit recover err: %v", err)
		}
	}()
	s := model.GetStorage()
	records, err := s.GetSqlManageRecordsBySourceId(sa.Source, sa.SqlAuditRecordID)
	if err != nil {
		return fmt.Errorf("get sql manage list by source id error, error: %v", err)
	}
	for _, record := range records {
		aggSourceId, err := AggregateSourceIdsBysqlId(record.SQLID, sa.SqlAuditRecordID, false)
		if err != nil {
			return fmt.Errorf("aggregate source id failed, error: %v", err)
		}

		err = s.UpdateSqlManageRecord(record.SQLID, aggSourceId, sa.Source)
		if err != nil {
			return fmt.Errorf("update sql manage record failed, error: %v", err)
		}
	}

	return nil
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

// 根据sql id聚合source id（删除或者新增一条记录中的source id），结果以逗号隔开
func AggregateSourceIdsBysqlId(sqlId, sourceId string, isAddSourceId bool) (string, error) {
	s := model.GetStorage()
	sqlManageRecord, exist, err := s.GetManageSQLBySQLId(sqlId)
	if err != nil {
		return "", err
	}
	if !exist {
		return sourceId, nil
	}
	recordSourceIds := strings.Split(sqlManageRecord.SourceId, ",")
	sourceIds := make([]string, 0)
	for _, recordsSourceId := range recordSourceIds {
		if recordsSourceId != sourceId {
			sourceIds = append(sourceIds, recordsSourceId)
		}
	}
	if isAddSourceId {
		sourceIds = append(sourceIds, sourceId)
	}
	return strings.Join(sourceIds, ","), nil
}
