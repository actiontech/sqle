package auditplan

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

type SQLV2 struct {
	SQLId string
	// from audit plan
	Source      string
	SourceId    string // 实例扫描任务id: instance audit plan id
	ProjectId   string
	InstanceID  string
	AuditPlanId string // 扫描任务id: audit plan id

	// from collect
	SQLContent  string
	Fingerprint string
	SchemaName  string
	Info        Metrics
}

func (s *SQLV2) GenSQLId() {
	md5Json, err := json.Marshal(
		struct {
			ProjectId   string
			Fingerprint string
			Schema      string
			InstID      string
			Source      string
			AuditPlanID string
		}{
			ProjectId:   s.ProjectId,
			Fingerprint: s.Fingerprint,
			Schema:      s.SchemaName,
			InstID:      s.InstanceID,
			Source:      s.Source,
			AuditPlanID: s.AuditPlanId,
		},
	)
	if err != nil {
		s.SQLId = s.Fingerprint
	} else {
		s.SQLId = utils.Md5String(string(md5Json))
	}
}

// Deprecated
func NewSQLV2FromSQL(ap *AuditPlan, sql *SQL) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(ap.Type)
	if err == nil {
		metrics = meta.Metrics
	}
	s := &SQLV2{
		Source:      ap.Type,
		SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
		ProjectId:   ap.ProjectId,
		InstanceID:  ap.InstanceID,
		SchemaName:  sql.Schema,
		SQLContent:  sql.SQLContent,
		Fingerprint: sql.Fingerprint,
	}
	s.Info = LoadMetrics(sql.Info, metrics)
	s.GenSQLId()
	return s
}

func ConvertMangerSQLQueueToSQLV2(sql *model.SQLManageQueue) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(sql.Source)
	if err == nil {
		metrics = meta.Metrics
	}
	// todo: 错误处理
	info, _ := sql.Info.OriginValue()
	s := &SQLV2{
		SQLId:       sql.SQLID,
		Source:      sql.Source,
		SourceId:    sql.SourceId,
		ProjectId:   sql.ProjectId,
		InstanceID:  sql.InstanceID,
		SchemaName:  sql.SchemaName,
		SQLContent:  sql.SqlText,
		Fingerprint: sql.SqlFingerprint,
		Info:        LoadMetrics(info, metrics),
	}
	return s
}

func ConvertMangerSQLToSQLV2(sql *model.SQLManageRecord) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(sql.Source)
	if err == nil {
		metrics = meta.Metrics
	}
	// todo: 错误处理
	info, _ := sql.Info.OriginValue()
	s := &SQLV2{
		SQLId:       sql.SQLID,
		Source:      sql.Source,
		SourceId:    sql.SourceId,
		ProjectId:   sql.ProjectId,
		InstanceID:  sql.InstanceID,
		SchemaName:  sql.SchemaName,
		SQLContent:  sql.SqlText,
		Fingerprint: sql.SqlFingerprint,
		Info:        LoadMetrics(info, metrics),
	}
	return s
}

func ConvertSQLV2ToMangerSQL(sql *SQLV2) *model.SQLManageRecord {
	data, _ := json.Marshal(sql.Info.ToMap()) // todo: 错误处理
	return &model.SQLManageRecord{
		SQLID:          sql.SQLId,
		Source:         sql.Source,
		SourceId:       sql.SourceId,
		ProjectId:      sql.ProjectId,
		InstanceID:     sql.InstanceID,
		SchemaName:     sql.SchemaName,
		SqlFingerprint: sql.Fingerprint,
		SqlText:        sql.SQLContent,
		Info:           data,
	}
}

func ConvertSQLV2ToMangerSQLQueue(sql *SQLV2) *model.SQLManageQueue {
	data, _ := json.Marshal(sql.Info.ToMap()) // todo: 错误处理
	return &model.SQLManageQueue{
		SQLID:          sql.SQLId,
		Source:         sql.Source,
		SourceId:       sql.SourceId,
		ProjectId:      sql.ProjectId,
		InstanceID:     sql.InstanceID,
		SchemaName:     sql.SchemaName,
		SqlFingerprint: sql.Fingerprint,
		SqlText:        sql.SQLContent,
		Info:           data,
	}
}

func ConvertSQLV2ToMangerRawSQL(sql *SQLV2) *model.SQLManageRawSQL {
	data, _ := json.Marshal(sql.Info.ToMap()) // todo: 错误处理
	execTimeStr := sql.Info.Get(MetricNameLastReceiveTimestamp).String()
	execTime, err := time.Parse(time.RFC3339, execTimeStr)
	if err != nil {
		execTime = time.Now()
	}
	return &model.SQLManageRawSQL{
		Source:         sql.Source,
		SourceId:       sql.SourceId,
		ProjectId:      sql.ProjectId,
		InstanceID:     sql.InstanceID,
		SchemaName:     sql.SchemaName,
		SqlFingerprint: sql.Fingerprint,
		SqlText:        sql.SQLContent,
		Info:           data,
		SQLID:          sql.SQLId,
		SqlExecTime:    execTime,
	}
}
