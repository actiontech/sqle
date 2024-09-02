//go:build enterprise
// +build enterprise

package model

import (
	"database/sql"
	e "errors"
	"fmt"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

const (
	SQLManageStatusUnhandled     = "unhandled"
	SQLManageStatusSolved        = "solved"
	SQLManageStatusIgnored       = "ignored"
	SQLManageStatusManualAudited = "manual_audited"

	SQLManageSourceAuditPlan      = "audit_plan"
	SQLManageSourceSqlAuditRecord = "sql_audit_record"

	CommonAuditLevel = "normal"
	NoticeAuditLevel = "notice"
	WarnAuditLevel   = "warn"
	ErrorAuditLevel  = "error"
)

var SqlManageSourceMap = map[string]*i18n.Message{
	SQLManageSourceSqlAuditRecord: locale.SQLManageSourceSqlAuditRecord,
	SQLManageSourceAuditPlan:      locale.SQLManageSourceAuditPlan,
}

var SqlManageStatusMap = map[string]*i18n.Message{
	SQLManageStatusUnhandled:     locale.SQLManageStatusUnhandled,
	SQLManageStatusSolved:        locale.SQLManageStatusSolved,
	SQLManageStatusIgnored:       locale.SQLManageStatusIgnored,
	SQLManageStatusManualAudited: locale.SQLManageStatusManualAudited,
}

func (s *Storage) UpdateSqlManage(auditRecordId uint) error {
	return s.Tx(func(tx *gorm.DB) error {
		err := tx.Exec(`DELETE sql_manages
FROM sql_manages,
     sql_manage_sql_audit_records smr
WHERE smr.sql_id = sql_manages.proj_fp_source_inst_schema_md5
  AND smr.sql_audit_record_id = ?
  AND sql_manages.fp_count = 1
  AND sql_manages.deleted_at IS NULL;`, auditRecordId).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`UPDATE sql_manages s,
    sql_manage_sql_audit_records smr
SET s.fp_count = s.fp_count - 1
WHERE s.proj_fp_source_inst_schema_md5 = smr.sql_id
  AND smr.sql_audit_record_id = ?
  AND s.fp_count > 1
  AND s.deleted_at IS NULL;`, auditRecordId).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`DELETE FROM sql_manage_sql_audit_records WHERE sql_audit_record_id = ?`, auditRecordId).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) UpdateSqlManageRecord(sqlId, sourceIds, source string) error {
	return s.Tx(func(tx *gorm.DB) error {
		if sourceIds == "" {
			err := tx.Exec(`UPDATE sql_manage_records oms 
			LEFT JOIN sql_manage_record_processes sm ON sm.sql_manage_record_id = oms.id
			SET oms.deleted_at = now(),
			sm.deleted_at = now()
			WHERE oms.source = ? AND oms.sql_id = ? `, source, sqlId).Error
			if err != nil {
				return err
			}
		} else {
			err := tx.Model(&SQLManageRecord{}).Where("source = ? AND sql_id = ?", source, sqlId).Update("source_id", sourceIds).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetSqlManageRuleTips(projectID string) ([]*SqlManageRuleTips, error) {
	sqlManageRuleTips := make([]*SqlManageRuleTips, 0)
	err := s.db.Raw(`SELECT DISTINCT t.db_type, r.name rule_name, r.desc
FROM sql_manages sm
         LEFT JOIN sql_manage_sql_audit_records msar
                   ON sm.proj_fp_source_inst_schema_md5 = msar.sql_id
         LEFT JOIN sql_audit_records sar ON msar.sql_audit_record_id = sar.id
         LEFT JOIN tasks t ON sar.task_id = t.id
         LEFT JOIN rules r ON r.db_type = t.db_type
WHERE sm.deleted_at IS NULL AND sm.project_id = ?
  AND sm.audit_results LIKE CONCAT('%"'
    , r.name
    , '"%')
UNION
SELECT DISTINCT ap.db_type, r.name rule_name, r.desc
FROM sql_manages sm
         LEFT JOIN audit_plans ap ON ap.id = sm.audit_plan_id
         LEFT JOIN rules r ON r.db_type = ap.db_type
WHERE sm.deleted_at IS NULL  AND sm.project_id = ?
  AND sm.audit_results LIKE CONCAT('%"'
    , r.name
    , '"%');`, projectID, projectID).Scan(&sqlManageRuleTips).Error
	if err != nil {
		return nil, err
	}

	return sqlManageRuleTips, nil
}

func (s *Storage) GetSqlManagerRuleTips(projectID string) ([]*SqlManageRuleTips, error) {
	sqlManageRuleTips := make([]*SqlManageRuleTips, 0)
	err := s.db.Table("sql_manage_records oms").
		Joins("LEFT JOIN audit_plans_v2 ap ON ap.id = oms.source_id").
		Joins("LEFT JOIN instance_audit_plans iap ON iap.id = ap.instance_audit_plan_id").
		Joins("LEFT JOIN rules ON rules.db_type = iap.db_type").
		Where("oms.audit_results LIKE CONCAT('%' , rules.name , '%') AND oms.project_id = ?", projectID).
		Select("DISTINCT iap.db_type, rules.name as rule_name, rules.desc").
		Scan(&sqlManageRuleTips).Error
	return sqlManageRuleTips, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetSqlManageByFingerprintSourceInstNameSchemaMd5(projFpSourceInstSchemaMd5 string) (*SqlManage, bool, error) {
	sqlManage := &SqlManage{}
	err := s.db.Where("proj_fp_source_inst_schema_md5 = ?", projFpSourceInstSchemaMd5).First(sqlManage).Error
	if e.Is(err, gorm.ErrRecordNotFound) {
		return sqlManage, false, nil
	}

	return sqlManage, true, errors.New(errors.ConnectStorageError, err)
}

type SqlManageResp struct {
	SqlManageList         []*SqlManageDetail
	SqlManageTotalNum     uint64 `json:"sql_manage_total_num"`
	SqlManageBadNum       uint64 `json:"sql_manage_bad_num"`
	SqlManageOptimizedNum uint64 `json:"sql_manage_optimized_num"`
}

type SqlManageDetail struct {
	ID                   uint           `json:"id"`
	SqlFingerprint       sql.NullString `json:"sql_fingerprint"`
	SqlText              sql.NullString `json:"sql_text"`
	Source               sql.NullString `json:"source"`
	SourceIDs            RowList        `json:"source_id"`
	AuditLevel           sql.NullString `json:"audit_level"`
	AuditResults         AuditResults   `json:"audit_results"`
	AuditStatus          sql.NullString `json:"audit_status"`
	FpCount              uint64         `json:"fp_count"`
	AppearTimestamp      *time.Time     `json:"first_appear_timestamp"`
	LastReceiveTimestamp *time.Time     `json:"last_receive_timestamp"`
	InstanceID           sql.NullString `json:"instance_id"`
	SchemaName           sql.NullString `json:"schema_name"`
	Status               sql.NullString `json:"status"`
	Remark               sql.NullString `json:"remark"`
	Assignees            *string        `json:"assignees"`
	Endpoints            sql.NullString `json:"endpoints"`
	Priority             sql.NullString `json:"priority"`
}

func (sm *SqlManageDetail) FirstAppearTime() string {
	if sm.AppearTimestamp == nil {
		return ""
	}
	return sm.AppearTimestamp.String()
}

func (sm *SqlManageDetail) LastReceiveTime() string {
	if sm.LastReceiveTimestamp == nil {
		return ""
	}
	return sm.LastReceiveTimestamp.String()
}

var sqlManageTotalCount = `
SELECT COUNT(DISTINCT sm.id)

{{- template "body" . -}}
`

var sqlManageQueryTpl = `
SELECT 
	sm.id,
	sm.sql_fingerprint,
	sm.sql_text,
	sm.source,
	sm.audit_level,
	sm.audit_results,
	sm.fp_count,
    sm.first_appear_timestamp,
	sm.last_receive_timestamp,
	sm.instance_id,
	sm.schema_name,
	sm.status,
	sm.remark,
	sm.assignees as assignees,
	ap.name as ap_name,
	GROUP_CONCAT(DISTINCT sar.audit_record_id) as sql_audit_record_ids,
	GROUP_CONCAT(DISTINCT all_sme.endpoint) as endpoints

{{- template "body" . -}} 

GROUP BY sm.id
ORDER BY 
{{- if and .sort_field .sort_order }}
	{{ .sort_field }} {{ .sort_order }}
{{- else }}
	sm.id desc
{{- end }}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var sqlManageBodyTpl = `
{{ define "body" }}

FROM sql_manages sm
         LEFT JOIN sql_manage_sql_audit_records msar ON sm.proj_fp_source_inst_schema_md5 = msar.sql_id
         LEFT JOIN sql_audit_records sar ON msar.sql_audit_record_id = sar.id
         LEFT JOIN sql_manage_endpoints sme ON sme.proj_fp_source_inst_schema_md5 = sm.proj_fp_source_inst_schema_md5
         LEFT JOIN sql_manage_endpoints all_sme ON all_sme.proj_fp_source_inst_schema_md5 = sm.proj_fp_source_inst_schema_md5
		 LEFT JOIN tasks t ON sar.task_id = t.id
         LEFT JOIN audit_plans ap ON ap.id = sm.audit_plan_id

WHERE sm.project_id = :project_id
  AND sm.deleted_at IS NULL

{{- if .fuzzy_search_sql_fingerprint }}
AND sm.sql_fingerprint LIKE '%{{ .fuzzy_search_sql_fingerprint }}%'
{{- end }}

{{- if .filter_assignee }}
AND sm.assignees REGEXP :filter_assignee
{{- end }}

{{- if .filter_instance_id }}
AND sm.instance_id = :filter_instance_id
{{- end }}

{{- if .filter_source }}
AND sm.source = :filter_source
{{- end }}

{{- if .filter_audit_level }}
AND sm.audit_level in ({{range $index, $element := .filter_audit_level}}{{if $index}}, {{end}}'{{$element}}'{{end}})
{{- end }}

{{- if .filter_db_type }}
AND (ap.db_type = :filter_db_type OR t.db_type = :filter_db_type)
{{- end }}

{{- if .filter_rule_name }}
AND JSON_CONTAINS(JSON_EXTRACT(sm.audit_results,'$[*].rule_name'), '"{{ .filter_rule_name }}"') > 0 
{{- end }}

{{- if .filter_last_audit_start_time_from }}
AND sm.last_receive_timestamp >= :filter_last_audit_start_time_from
{{- end }}

{{- if .filter_last_audit_start_time_to }}
AND sm.last_receive_timestamp <= :filter_last_audit_start_time_to
{{- end }}

{{- if .fuzzy_search_endpoint }}
AND sme.endpoint LIKE '%{{ .fuzzy_search_endpoint }}%'
{{- end }}

{{- if .fuzzy_search_schema_name }}
AND sm.schema_name LIKE '%{{ .fuzzy_search_schema_name }}%'
{{- end }}

{{- if .filter_status }}
AND sm.status = :filter_status
{{- end }}

{{- if .count_bad_sql }}
AND ( sm.audit_level = 'warn' OR sm.audit_level = 'error' )
AND sm.status = 'unhandled'
{{- end }}

{{- if .count_solved }}
AND sm.status = 'solved'
{{- end }}

{{ end }}
`

var sqlManagerTotalCount = `
SELECT COUNT(DISTINCT oms.id)

{{- template "body" . -}}
`

var sqlManagerQueryTpl = `
SELECT 
	oms.id,
	oms.sql_fingerprint,
	oms.sql_text,
	oms.source,
	oms.audit_level,
	IF(oms.audit_results IS NULL,'null',oms.audit_results) AS audit_results,
	IF(oms.audit_results IS NULL,'being_audited','') AS audit_status,
	oms.instance_id,
	oms.schema_name,
	oms.end_point as endpoints,
	sm.status,
	sm.remark,
	sm.assignees as assignees,
	oms.source_id as source_id,
	oms.priority
{{- template "body" . -}} 

GROUP BY oms.id
ORDER BY 
{{- if and .sort_field .sort_order }}
	{{ .sort_field }} {{ .sort_order }}
{{- else }}
	oms.id desc
{{- end }}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var sqlManagerBodyTpl = `
{{ define "body" }}

FROM sql_manage_records oms
         LEFT JOIN sql_manage_record_processes sm ON sm.sql_manage_record_id = oms.id

WHERE oms.project_id = :project_id
  AND oms.deleted_at IS NULL

{{- if .fuzzy_search_sql_fingerprint }}
AND oms.sql_fingerprint LIKE '%{{ .fuzzy_search_sql_fingerprint }}%'
{{- end }}

{{- if .filter_assignee }}
AND sm.assignees REGEXP :filter_assignee
{{- end }}

{{- if .filter_instance_id }}
AND oms.instance_id = :filter_instance_id
{{- end }}


{{- if .filter_business_instance_ids }}
AND oms.instance_id in ( {{ .filter_business_instance_ids}} )
{{- end }}


{{- if .filter_source }}
AND oms.source = :filter_source
{{- end }}

{{- if .filter_priority }}
AND oms.priority = :filter_priority
{{- end }}

{{- if .filter_audit_level }}
AND oms.audit_level in ({{range $index, $element := .filter_audit_level}}{{if $index}}, {{end}}'{{$element}}'{{end}})
{{- end }}

{{- if .filter_rule_name }}
AND JSON_CONTAINS(JSON_EXTRACT(oms.audit_results,'$[*].rule_name'), '"{{ .filter_rule_name }}"') > 0 
{{- end }}

{{- if .filter_last_audit_start_time_from }}
AND oms.updated_at >= :filter_last_audit_start_time_from
{{- end }}

{{- if .filter_last_audit_start_time_to }}
AND oms.updated_at <= :filter_last_audit_start_time_to
{{- end }}

{{- if .fuzzy_search_schema_name }}
AND oms.schema_name LIKE '%{{ .fuzzy_search_schema_name }}%'
{{- end }}

{{- if .filter_status }}
AND sm.status = :filter_status
{{- end }}

{{- if .count_bad_sql }}
AND ( oms.audit_level = 'warn' OR oms.audit_level = 'error' )
AND sm.status = 'unhandled'
{{- end }}

{{- if .count_solved }}
AND sm.status = 'solved'
{{- end }}

{{ end }}
`

// 获取大于等于传参的审核等级
// 比如：传参的值为warn，需要返回的审核等级为warn和error；传参的值为notice，需要返回notice，warn和error
func getAuditLevelsByLowestLevel(auditLevel string) []string {
	auditLevels := []string{CommonAuditLevel, NoticeAuditLevel, WarnAuditLevel, ErrorAuditLevel}
	switch auditLevel {
	case CommonAuditLevel:
		return auditLevels
	case NoticeAuditLevel:
		return auditLevels[1:]
	case WarnAuditLevel:
		return auditLevels[2:]
	case ErrorAuditLevel:
		return auditLevels[3:]
	default:
		return []string{}
	}
}

func (s *Storage) GetSqlManageListByReq(data map[string]interface{}) (list *SqlManageResp, err error) {
	sqlManageList := make([]*SqlManageDetail, 0)
	auditLevel := data["filter_audit_level"]
	auditLevelPointer, ok := auditLevel.(*string)
	if !ok {
		return nil, e.New("sql_manage: filter_audit_level convert to *string failed")
	}
	if auditLevelPointer != nil {
		auditLevelStr := *auditLevelPointer
		data["filter_audit_level"] = getAuditLevelsByLowestLevel(auditLevelStr)
	}

	err = s.getListResult(sqlManagerQueryTpl, sqlManagerBodyTpl, data, &sqlManageList)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.getCountResult(sqlManagerBodyTpl, sqlManagerTotalCount, data)
	if err != nil {
		return nil, err
	}

	fn := func(srcData map[string]interface{}, addSearchKey string) map[string]interface{} {
		newData := make(map[string]interface{})
		for k, v := range srcData {
			newData[k] = v
		}
		newData[addSearchKey] = true
		return newData
	}

	badSqlCount, err := s.getCountResult(sqlManagerBodyTpl, sqlManagerTotalCount, fn(data, "count_bad_sql"))
	if err != nil {
		return nil, err
	}

	solvedCount, err := s.getCountResult(sqlManagerBodyTpl, sqlManagerTotalCount, fn(data, "count_solved"))
	if err != nil {
		return nil, err
	}

	return &SqlManageResp{
		SqlManageList:         sqlManageList,
		SqlManageTotalNum:     totalCount,
		SqlManageBadNum:       badSqlCount,
		SqlManageOptimizedNum: solvedCount,
	}, nil
}

func (s *Storage) GetAllSqlManageList() ([]*SqlManage, error) {
	sqlManageList := make([]*SqlManage, 0)
	err := s.db.Model(&SqlManage{}).Find(&sqlManageList).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return sqlManageList, nil
}

type SqlManageWithEndpoint struct {
	*SqlManage
	Endpoints []string
}

// todo : 指纹count未累加
func (s *Storage) InsertOrUpdateSqlManageWithNotUpdateFpCount(sqlManageList []*SqlManageWithEndpoint) error {
	return s.Tx(func(tx *gorm.DB) error {
		args := make([]interface{}, 0, len(sqlManageList))
		pattern := make([]string, 0, len(sqlManageList))

		sqlManageEndpointArgs := make([]interface{}, 0, len(sqlManageList))
		sqlManageEndpointPattern := make([]string, 0, len(sqlManageList))

		now := time.Now().Format("2006-01-02 15:04:05")
		for _, sqlManage := range sqlManageList {
			pattern = append(pattern, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			args = append(args, sqlManage.SqlFingerprint, sqlManage.ProjFpSourceInstSchemaMd5, sqlManage.SqlText,
				sqlManage.Source, sqlManage.AuditLevel, sqlManage.AuditResults, sqlManage.FpCount, sqlManage.FirstAppearTimestamp,
				sqlManage.LastReceiveTimestamp, sqlManage.InstanceName, sqlManage.SchemaName, sqlManage.Remark,
				sqlManage.AuditPlanId, sqlManage.ProjectId)

			if len(sqlManage.Endpoints) > 0 {
				for _, endpoint := range sqlManage.Endpoints {
					sqlManageEndpointArgs = append(sqlManageEndpointArgs, now, now, sqlManage.ProjFpSourceInstSchemaMd5, endpoint)
					sqlManageEndpointPattern = append(sqlManageEndpointPattern, "(?, ?, ?, ?)")
				}
			}
		}

		if len(sqlManageEndpointArgs) > 0 {
			rawSql := fmt.Sprintf(`
				INSERT INTO sql_manage_endpoints (created_at, updated_at, proj_fp_source_inst_schema_md5, endpoint) 
				 	VALUES %s
				 	ON DUPLICATE KEY UPDATE updated_at = '%s'`, strings.Join(sqlManageEndpointPattern, ", "), now)

			err := tx.Exec(rawSql, sqlManageEndpointArgs...).Error
			if err != nil {
				return err
			}
		}

		raw := fmt.Sprintf(`
INSERT INTO sql_manages (sql_fingerprint, proj_fp_source_inst_schema_md5, sql_text, source, audit_level, audit_results,
                         fp_count, first_appear_timestamp, last_receive_timestamp, instance_name, schema_name,
                         remark, audit_plan_id, project_id)
		VALUES %s
		ON DUPLICATE KEY UPDATE sql_text       = VALUES(sql_text),
                        audit_plan_id          = VALUES(audit_plan_id),
                        audit_level            = VALUES(audit_level),
                        audit_results          = VALUES(audit_results),
                        first_appear_timestamp = VALUES(first_appear_timestamp),
                        last_receive_timestamp = VALUES(last_receive_timestamp);`,
			strings.Join(pattern, ", "))
		err := tx.Exec(raw, args...).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) InsertOrUpdateSqlManage(sqlManageList []*SqlManage, sqlAuditRecordID uint) error {
	return s.Tx(func(tx *gorm.DB) error {
		batchSize := 50 // 每批处理的大小
		total := len(sqlManageList)
		start := 0

		for start < total {
			end := start + batchSize
			if end > total {
				end = total
			}

			batchSqlManageList := sqlManageList[start:end]

			args := make([]interface{}, 0, len(batchSqlManageList))
			pattern := make([]string, 0, len(batchSqlManageList))
			for _, sqlManage := range batchSqlManageList {
				pattern = append(pattern, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
				args = append(args, sqlManage.SqlFingerprint, sqlManage.ProjFpSourceInstSchemaMd5, sqlManage.SqlText,
					sqlManage.Source, sqlManage.AuditLevel, sqlManage.AuditResults, sqlManage.FpCount, sqlManage.FirstAppearTimestamp,
					sqlManage.LastReceiveTimestamp, sqlManage.InstanceName, sqlManage.SchemaName, sqlManage.Remark,
					sqlManage.AuditPlanId, sqlManage.ProjectId)
			}

			raw := fmt.Sprintf(`
			INSERT INTO sql_manages (sql_fingerprint, proj_fp_source_inst_schema_md5, sql_text, source, audit_level, audit_results,
			                        fp_count, first_appear_timestamp, last_receive_timestamp, instance_name, schema_name,
			                        remark, audit_plan_id, project_id)
					VALUES %s
					ON DUPLICATE KEY UPDATE sql_text       = VALUES(sql_text),
			                       audit_plan_id          = VALUES(audit_plan_id),
			                       audit_level            = VALUES(audit_level),
			                       audit_results          = VALUES(audit_results),
			                       fp_count 			   = VALUES(fp_count),
			                       first_appear_timestamp = VALUES(first_appear_timestamp),
			                       last_receive_timestamp = VALUES(last_receive_timestamp);`,
				strings.Join(pattern, ", "))

			err := tx.Exec(raw, args...).Error
			if err != nil {
				return err
			}

			if sqlAuditRecordID != 0 {
				sqlAuditArgs := make([]interface{}, 0, len(batchSqlManageList))
				sqlAuditPattern := make([]string, 0, len(batchSqlManageList))

				for _, sqlManage := range batchSqlManageList {
					sqlAuditPattern = append(sqlAuditPattern, "(?, ?)")
					sqlAuditArgs = append(sqlAuditArgs, sqlManage.ProjFpSourceInstSchemaMd5, sqlAuditRecordID)
				}

				rawSql := fmt.Sprintf(`
				INSERT INTO sql_manage_sql_audit_records (sql_id, sql_audit_record_id) 
				 	VALUES %s`, strings.Join(sqlAuditPattern, ", "))

				err := tx.Exec(rawSql, sqlAuditArgs...).Error
				if err != nil {
					return err
				}
			}

			start += batchSize
		}

		return nil
	})
}

func (s *Storage) InsertOrUpdateSqlManageRecord(sqlManageList []*SQLManageRecord) error {
	return s.Tx(func(tx *gorm.DB) error {
		batchSize := 50 // 每批处理的大小
		total := len(sqlManageList)
		start := 0

		for start < total {
			end := start + batchSize
			if end > total {
				end = total
			}

			batchSqlManageList := sqlManageList[start:end]

			args := make([]interface{}, 0, len(batchSqlManageList))
			pattern := make([]string, 0, len(batchSqlManageList))
			for _, sqlManage := range batchSqlManageList {
				pattern = append(pattern, "(?,?,?,?,?,?,?,?,?,?,?)")
				args = append(args, sqlManage.SQLID, sqlManage.Source, sqlManage.SourceId, sqlManage.ProjectId, sqlManage.InstanceID,
					sqlManage.SchemaName, sqlManage.SqlFingerprint, sqlManage.SqlText, sqlManage.Info, sqlManage.AuditLevel, sqlManage.AuditResults)
			}

			raw := fmt.Sprintf(`INSERT INTO sql_manage_records (sql_id, source, source_id, project_id, instance_id, schema_name,
														sql_fingerprint, sql_text,  info, audit_level, audit_results) 
							VALUES %s
							ON DUPLICATE KEY UPDATE source = VALUES(source),
											source_id = VALUES(source_id), 
											project_id = VALUES(project_id), 
											instance_id = VALUES(instance_id), 
											schema_name = VALUES(schema_name), 
											sql_text = VALUES(sql_text), 
											sql_fingerprint = VALUES(sql_fingerprint), 
											info = VALUES(info), 
											audit_level= VALUES(audit_level), 
											audit_results = VALUES(audit_results),
											deleted_at = NULL;`,
				strings.Join(pattern, ", "))

			err := tx.Exec(raw, args...).Error
			if err != nil {
				return err
			}

			for _, sqlManage := range batchSqlManageList {
				const query = `INSERT INTO sql_manage_record_processes (sql_manage_record_id)
									SELECT oms.id FROM sql_manage_records oms WHERE oms.sql_id = ?
								ON DUPLICATE KEY UPDATE sql_manage_record_id = VALUES(sql_manage_record_id),
								deleted_at = NULL;`
				err := tx.Exec(query, sqlManage.SQLID).Error
				if err != nil {
					return err
				}
			}
			start += batchSize
		}

		return nil
	})
}

func (s *Storage) BatchUpdateSqlManage(idList []*uint64, status *string, remark *string, assignees []string) error {
	return s.Tx(func(tx *gorm.DB) error {
		data := map[string]interface{}{}
		if status != nil {
			data["status"] = *status
		}

		if remark != nil {
			data["remark"] = *remark
		}

		if len(assignees) != 0 {
			data["assignees"] = strings.Join(assignees, ",")
		}
		if len(data) > 0 {
			err := tx.Model(&SqlManage{}).Where("id in (?)", idList).Updates(data).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Storage) BatchUpdateSqlManager(idList []*uint64, status, remark, priority *string, assignees []string) error {
	return s.Tx(func(tx *gorm.DB) error {
		data := map[string]interface{}{}
		if status != nil {
			data["status"] = *status
		}

		if remark != nil {
			data["remark"] = *remark
		}

		if len(assignees) != 0 {
			data["assignees"] = strings.Join(assignees, ",")
		}
		if len(data) > 0 {
			err := tx.Model(&SQLManageRecordProcess{}).Where("sql_manage_record_id in (?)", idList).Updates(data).Error
			if err != nil {
				return err
			}
		}

		// update priority
		if priority != nil {
			err := tx.Model(&SQLManageRecord{}).Where("id in (?)", idList).Updates(map[string]interface{}{"priority": *priority}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Storage) GetSqlManageByID(id string) (*SqlManage, bool, error) {
	sqlManage := new(SqlManage)
	err := s.db.Where("id = ?", id).First(&sqlManage).Error
	if err != nil {
		if e.Is(gorm.ErrRecordNotFound, err) {
			return sqlManage, false, nil
		}
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return sqlManage, true, nil
}

func (s *Storage) GetOriginManageSqlByID(id string) (*SQLManageRecord, bool, error) {
	originManageSQL := new(SQLManageRecord)
	err := s.db.Where("id = ?", id).First(&originManageSQL).Error
	if err != nil {
		if e.Is(gorm.ErrRecordNotFound, err) {
			return originManageSQL, false, nil
		}
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return originManageSQL, true, nil
}

func (s *Storage) GetSqlManageListByIDs(ids []*uint64) ([]*SqlManage, error) {
	sqlManageList := []*SqlManage{}
	err := s.db.Model(SqlManage{}).Where("id IN (?)", ids).Find(&sqlManageList).Error
	if err != nil {
		return nil, err
	}
	return sqlManageList, nil
}

func (s *Storage) GetSqlManagerListByIDs(ids []*uint64) ([]*SQLManageRecordProcess, error) {
	sqlManagerList := []*SQLManageRecordProcess{}
	err := s.db.Model(SQLManageRecordProcess{}).Where("sql_manage_record_id IN (?)", ids).Find(&sqlManagerList).Error
	if err != nil {
		return nil, err
	}
	return sqlManagerList, nil
}

func (s *Storage) GetAuditPlanUnsolvedSQLCount(id uint, status []string) (int64, error) {
	query := `SELECT
					count(smr.id)
				FROM
					sql_manage_records AS smr
				LEFT JOIN sql_manage_record_processes AS sm ON sm.sql_manage_record_id = smr.id
				LEFT JOIN audit_plans_v2 AS ap ON ap.instance_audit_plan_id = smr.source_id AND ap.type = smr.source 
				WHERE
					ap.id = ?
					AND smr.deleted_at IS NULL
					AND JSON_TYPE(smr.audit_results) <> 'NULL'
					AND smr.audit_results IS NOT NULL 
					AND sm.status NOT IN(?);`
	var count int64
	err := s.db.Raw(query, id, status).Count(&count).Error
	if err != nil {
		return count, errors.New(errors.ConnectStorageError, err)
	}
	return count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetHighLevelSQLsByTime(projectId string, fromTime time.Time) ([]*SQLManageRecord, error) {
	sqlManageList := []*SQLManageRecord{}
	err := s.db.Model(SQLManageRecord{}).Where("project_id = ? AND updated_at > ? AND priority = 'high'", projectId, fromTime).Find(&sqlManageList).Error
	if err != nil {
		return nil, err
	}
	return sqlManageList, nil
}

func (s *Storage) GetSqlManageRecordsBySourceId(source, sourceId string) ([]*SQLManageRecord, error) {
	sqlManageRecors := []*SQLManageRecord{}
	err := s.db.Model(SQLManageRecord{}).Where("source = ? AND source_id LIKE ?", source, "%"+sourceId+"%").Find(&sqlManageRecors).Error
	if err != nil {
		return nil, err
	}
	return sqlManageRecors, nil
}
