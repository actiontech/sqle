package auditplan

import (
	"encoding/json"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

// todo: 换名称
type SQLV2 struct {
	SQLId string
	// from audit plan
	Source       string
	SourceId     uint
	ProjectId    string
	InstanceName string

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
			InstName    string
			Source      string
			ApID        uint
		}{
			ProjectId:   s.ProjectId,
			Fingerprint: s.Fingerprint,
			Schema:      s.SchemaName,
			InstName:    s.InstanceName,
			Source:      s.Source,
			ApID:        s.SourceId,
		},
	)
	if err != nil {
		s.SQLId = s.Fingerprint
	} else {
		s.SQLId = utils.Md5String(string(md5Json))
	}
}

func NewSQLV2FromSQL(ap *AuditPlan, sql *SQL) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(ap.Type)
	if err == nil {
		metrics = meta.Metrics
	}
	s := &SQLV2{
		Source:       ap.Type,
		SourceId:     ap.ID,
		ProjectId:    ap.ProjectId,
		InstanceName: ap.InstanceName,
		SchemaName:   sql.Schema,
		SQLContent:   sql.SQLContent,
		Fingerprint:  sql.Fingerprint,
	}
	s.Info = LoadMetrics(sql.Info, metrics)
	s.GenSQLId()
	return s
}

func ConvertMangerSQLQueueToSQLV2(sql *model.OriginManageSQLQueue) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(sql.Source)
	if err == nil {
		metrics = meta.Metrics
	}
	// todo: 错误处理
	info, _ := sql.Info.OriginValue()
	s := &SQLV2{
		SQLId:        sql.SQLID,
		Source:       sql.Source,
		SourceId:     sql.SourceId,
		ProjectId:    sql.ProjectId,
		InstanceName: sql.InstanceName,
		SchemaName:   sql.SchemaName,
		SQLContent:   sql.SqlText,
		Fingerprint:  sql.SqlFingerprint,
		Info:         LoadMetrics(info, metrics),
	}
	return s
}

func ConvertMangerSQLToSQLV2(sql *model.OriginManageSQL) *SQLV2 {
	metrics := []string{}
	meta, err := GetMeta(sql.Source)
	if err == nil {
		metrics = meta.Metrics
	}
	// todo: 错误处理
	info, _ := sql.Info.OriginValue()
	s := &SQLV2{
		SQLId:        sql.SQLID,
		Source:       sql.Source,
		SourceId:     sql.SourceId,
		ProjectId:    sql.ProjectId,
		InstanceName: sql.InstanceName,
		SchemaName:   sql.SchemaName,
		SQLContent:   sql.SqlText,
		Fingerprint:  sql.SqlFingerprint,
		Info:         LoadMetrics(info, metrics),
	}
	return s
}

func ConvertSQLV2ToMangerSQL(sql *SQLV2) *model.OriginManageSQL {
	data, _ := json.Marshal(sql.Info.ToMap()) // todo: 错误处理
	return &model.OriginManageSQL{
		SQLID:          sql.SQLId,
		Source:         sql.Source,
		SourceId:       sql.SourceId,
		ProjectId:      sql.ProjectId,
		InstanceName:   sql.InstanceName,
		SchemaName:     sql.SchemaName,
		SqlFingerprint: sql.Fingerprint,
		SqlText:        sql.SQLContent,
		Info:           data,
		EndPoint:       sql.Info.Get("endpoints").String(),
	}
}

func ConvertSQLV2ToMangerSQLQueue(sql *SQLV2) *model.OriginManageSQLQueue {
	data, _ := json.Marshal(sql.Info.ToMap()) // todo: 错误处理
	return &model.OriginManageSQLQueue{
		SQLID:          sql.SQLId,
		Source:         sql.Source,
		SourceId:       sql.SourceId,
		ProjectId:      sql.ProjectId,
		InstanceName:   sql.InstanceName,
		SchemaName:     sql.SchemaName,
		SqlFingerprint: sql.Fingerprint,
		SqlText:        sql.SQLContent,
		Info:           data,
		EndPoint:       sql.Info.Get("endpoints").String(),
	}
}
