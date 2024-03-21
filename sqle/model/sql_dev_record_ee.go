//go:build enterprise
// +build enterprise

package model

import (
	"fmt"
	"strings"
)

const (
	SQLDevRecordSourceIDEPlugin = "ide_plugin"
)

var SqlDEVManageSourceMap = map[string]string{
	SQLDevRecordSourceIDEPlugin: "IDE插件",
}

func (s *Storage) GetSqlDEVRecordListByReq(data map[string]interface{}) (list []*SQLDevRecord, total uint64, err error) {
	sqlDEVRecords := make([]*SQLDevRecord, 0)

	err = s.getListResult(sqlDEVRecordQueryTpl, sqlDEVRecordBodyTpl, data, &sqlDEVRecords)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.getCountResult(sqlDEVRecordBodyTpl, sqlDEVRecordTotalCount, data)
	if err != nil {
		return nil, 0, err
	}

	return sqlDEVRecords, total, nil
}

var sqlDEVRecordTotalCount = `
SELECT COUNT(DISTINCT sdr.id)

{{- template "body" . -}}
`

var sqlDEVRecordQueryTpl = `
SELECT 
	sdr.id,
	sdr.sql_fingerprint,
	sdr.sql_text,
	sdr.source,
	sdr.audit_level,
	sdr.audit_results,
	sdr.fp_count,
    sdr.first_appear_timestamp,
	sdr.last_receive_timestamp,
	sdr.instance_name,
	sdr.schema_name,
	sdr.creator

{{- template "body" . -}} 

ORDER BY 
{{- if and .sort_field .sort_order }}
	{{ .sort_field }} {{ .sort_order }}
{{- else }}
	sdr.id desc
{{- end }}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var sqlDEVRecordBodyTpl = `
{{ define "body" }}

FROM sql_dev_records sdr

WHERE sdr.project_id = :project_id
  AND sdr.deleted_at IS NULL

{{- if .fuzzy_search_sql_fingerprint }}
AND sdr.sql_fingerprint LIKE '%{{ .fuzzy_search_sql_fingerprint }}%'
{{- end }}

{{- if .filter_creator }}
AND sdr.creator = :filter_creator
{{- end }}

{{- if .filter_instance_name }}
AND sdr.instance_name = :filter_instance_name
{{- end }}

{{- if .filter_source }}
AND sdr.source = :filter_source
{{- end }}

{{- if .filter_last_receive_time_from }}
AND sdr.last_receive_timestamp >= :filter_last_receive_time_from
{{- end }}

{{- if .filter_last_receive_time_to }}
AND sdr.last_receive_timestamp <= :filter_last_receive_time_to
{{- end }}


{{ end }}
`

func (s *Storage) InsertOrUpdateSqlDevRecord(sqlDevRecordList []*SQLDevRecord) error {

	// 聚合记录，减少数据库操作
	sqlDevRecordMap := make(map[string]*SQLDevRecord)
	for _, v := range sqlDevRecordList {
		sdr, ok := sqlDevRecordMap[v.ProjFpSourceInstSchemaMd5]
		if ok {
			v.FpCount += sdr.FpCount
		}
		sqlDevRecordMap[v.ProjFpSourceInstSchemaMd5] = v
	}
	mergeSQLDevRecord := make([]*SQLDevRecord, 0)
	for _, v := range sqlDevRecordMap {
		mergeSQLDevRecord = append(mergeSQLDevRecord, v)
	}

	// 分片提交
	for batchSize, start := 50, 0; start < len(mergeSQLDevRecord); start += batchSize {
		end := start + batchSize
		if end > len(mergeSQLDevRecord) {
			end = len(mergeSQLDevRecord)
		}
		batchSqlDevRecordList := mergeSQLDevRecord[start:end]

		args := make([]interface{}, 0)
		pattern := make([]string, 0)
		for _, sqlDevRecord := range batchSqlDevRecordList {
			pattern = append(pattern, "( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			args = append(args, sqlDevRecord.SqlFingerprint, sqlDevRecord.ProjFpSourceInstSchemaMd5, sqlDevRecord.SqlText,
				sqlDevRecord.Source, sqlDevRecord.AuditLevel, sqlDevRecord.AuditResults, sqlDevRecord.FpCount, sqlDevRecord.FirstAppearTimestamp,
				sqlDevRecord.LastReceiveTimestamp, sqlDevRecord.InstanceName, sqlDevRecord.SchemaName, sqlDevRecord.Creator,
				sqlDevRecord.ProjectId)
		}

		raw := fmt.Sprintf(`
			INSERT INTO sql_dev_records (sql_fingerprint, proj_fp_source_inst_schema_md5, sql_text, source, audit_level, audit_results,
			                        fp_count, first_appear_timestamp, last_receive_timestamp, instance_name, schema_name,
			                        creator, project_id)
					VALUES %s
					ON DUPLICATE KEY UPDATE sql_text       = VALUES(sql_text),
			                       audit_level            = VALUES(audit_level),
			                       audit_results          = VALUES(audit_results),
			                       fp_count 			   = fp_count + VALUES(fp_count),
			                       first_appear_timestamp = VALUES(first_appear_timestamp),
			                       last_receive_timestamp = VALUES(last_receive_timestamp);`,
			strings.Join(pattern, ", "))

		err := s.db.Exec(raw, args...).Error
		if err != nil {
			return err
		}

	}

	return nil
}
