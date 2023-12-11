//go:build enterprise
// +build enterprise

package model

import (
	e "errors"
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
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

var SqlManageSourceMap = map[string]string{
	SQLManageSourceSqlAuditRecord: "SQL审核",
	SQLManageSourceAuditPlan:      "智能扫描",
}

var SqlManageStatusMap = map[string]string{
	SQLManageStatusUnhandled:     "未处理",
	SQLManageStatusSolved:        "已解决",
	SQLManageStatusIgnored:       "已忽略",
	SQLManageStatusManualAudited: "已人工审核",
}

func (s *Storage) UpdateSqlManage(auditRecordId uint) error {
	return s.Tx(func(tx *gorm.DB) error {
		err := tx.Exec(`DELETE sql_manages
FROM sql_manages,
     sql_manage_sql_audit_records smr
WHERE smr.proj_fp_source_inst_schema_md5 = sql_manages.proj_fp_source_inst_schema_md5
  AND smr.sql_audit_record_id = ?
  AND sql_manages.fp_count = 1
  AND sql_manages.deleted_at IS NULL;`, auditRecordId).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`UPDATE sql_manages s,
    sql_manage_sql_audit_records smr
SET s.fp_count = s.fp_count - 1
WHERE s.proj_fp_source_inst_schema_md5 = smr.proj_fp_source_inst_schema_md5
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

type SqlManageRuleTips struct {
	DbType   string `json:"db_type"`
	RuleName string `json:"rule_name"`
	Desc     string `json:"desc"`
}

func (s *Storage) GetSqlManageRuleTips(projectID uint) ([]*SqlManageRuleTips, error) {
	sqlManageRuleTips := make([]*SqlManageRuleTips, 0)
	err := s.db.Raw(`SELECT DISTINCT t.db_type, r.name rule_name, r.desc
FROM sql_manages sm
         LEFT JOIN sql_manage_sql_audit_records msar
                   ON sm.proj_fp_source_inst_schema_md5 = msar.proj_fp_source_inst_schema_md5
         LEFT JOIN sql_audit_records sar ON msar.sql_audit_record_id = sar.id
         LEFT JOIN tasks t ON sar.task_id = t.id
         LEFT JOIN rules r ON r.db_type = t.db_type
         JOIN projects p ON sm.project_id = p.id
WHERE sm.deleted_at IS NULL
  AND p.id = ?
  AND sm.audit_results LIKE CONCAT('%"'
    , r.name
    , '"%')
UNION
SELECT DISTINCT ap.db_type, r.name rule_name, r.desc
FROM sql_manages sm
         LEFT JOIN audit_plans ap ON ap.id = sm.audit_plan_id
         LEFT JOIN rules r ON r.db_type = ap.db_type
         JOIN projects p ON sm.project_id = p.id
WHERE sm.deleted_at IS NULL
  AND p.id = ?
  AND sm.audit_results LIKE CONCAT('%"'
    , r.name
    , '"%');`, projectID, projectID).Scan(&sqlManageRuleTips).Error
	if err != nil {
		return nil, err
	}

	return sqlManageRuleTips, nil
}

func (s *Storage) GetSqlManageByFingerprintSourceInstNameSchemaMd5(projFpSourceInstSchemaMd5 string) (*SqlManage, bool, error) {
	sqlManage := &SqlManage{}
	err := s.db.Where("proj_fp_source_inst_schema_md5 = ?", projFpSourceInstSchemaMd5).Find(sqlManage).Error
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
	ID                   uint         `json:"id"`
	SqlFingerprint       string       `json:"sql_fingerprint"`
	SqlText              string       `json:"sql_text"`
	Source               string       `json:"source"`
	AuditLevel           string       `json:"audit_level"`
	AuditResults         AuditResults `json:"audit_results"`
	FpCount              uint64       `json:"fp_count"`
	AppearTimestamp      *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp *time.Time   `json:"last_receive_timestamp"`
	InstanceName         string       `json:"instance_name"`
	SchemaName           string       `json:"schema_name"`
	Status               string       `json:"status"`
	Remark               string       `json:"remark"`
	Assignees            RowList      `json:"assignees"`
	ApName               *string      `json:"ap_name"`
	SqlAuditRecordIDs    RowList      `json:"sql_audit_record_ids"`
	Endpoints            RowList      `json:"endpoints"`
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
	sm.instance_name,
	sm.schema_name,
	sm.status,
	sm.remark,
	GROUP_CONCAT(DISTINCT all_users.login_name) as assignees,
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
         LEFT JOIN sql_manage_sql_audit_records msar ON sm.proj_fp_source_inst_schema_md5 = msar.proj_fp_source_inst_schema_md5
         LEFT JOIN sql_audit_records sar ON msar.sql_audit_record_id = sar.id
         LEFT JOIN sql_manage_endpoints sme ON sme.proj_fp_source_inst_schema_md5 = sm.proj_fp_source_inst_schema_md5
         LEFT JOIN sql_manage_endpoints all_sme ON all_sme.proj_fp_source_inst_schema_md5 = sm.proj_fp_source_inst_schema_md5	
		 LEFT JOIN tasks t ON sar.task_id = t.id
         LEFT JOIN audit_plans ap ON ap.id = sm.audit_plan_id
         LEFT JOIN sql_manage_assignees sma ON sma.sql_manage_id = sm.id
         LEFT JOIN users u ON u.id = sma.user_id
         LEFT JOIN sql_manage_assignees all_sma ON all_sma.sql_manage_id = sm.id
         LEFT JOIN users all_users ON all_users.id = all_sma.user_id

WHERE sm.project_id = :project_id
  AND sm.deleted_at IS NULL

{{- if .fuzzy_search_sql_fingerprint }}
AND sm.sql_fingerprint LIKE '%{{ .fuzzy_search_sql_fingerprint }}%'
{{- end }}

{{- if .filter_assignee }}
AND u.login_name = :filter_assignee
{{- end }}

{{- if .filter_instance_name }}
AND sm.instance_name = :filter_instance_name
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

	err = s.getListResult(sqlManageQueryTpl, sqlManageBodyTpl, data, &sqlManageList)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.getCountResult(sqlManageBodyTpl, sqlManageTotalCount, data)
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

	badSqlCount, err := s.getCountResult(sqlManageBodyTpl, sqlManageTotalCount, fn(data, "count_bad_sql"))
	if err != nil {
		return nil, err
	}

	solvedCount, err := s.getCountResult(sqlManageBodyTpl, sqlManageTotalCount, fn(data, "count_solved"))
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

		now := time.Now().Format(time.RFC3339)
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
				INSERT INTO sql_manage_sql_audit_records (proj_fp_source_inst_schema_md5, sql_audit_record_id) 
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

func (s *Storage) BatchUpdateSqlManage(idList []*uint64, status *string, remark *string, assignees []*string) error {
	return s.Tx(func(tx *gorm.DB) error {
		data := map[string]interface{}{}
		if status != nil {
			data["status"] = *status
		}

		if remark != nil {
			data["remark"] = *remark
		}

		if len(data) > 0 {
			err := tx.Model(&SqlManage{}).Where("id in (?)", idList).Update(data).Error
			if err != nil {
				return err
			}
		}

		if assignees != nil {
			userList := []*User{}
			err := tx.Where("login_name in (?)", assignees).Find(&userList).Error
			if err != nil {
				return err
			}

			if len(userList) > 0 {
				err := tx.Exec("DELETE FROM sql_manage_assignees WHERE sql_manage_id IN (?)", idList).Error
				if err != nil {
					return err
				}

				pattern := make([]string, 0, len(userList))
				args := make([]interface{}, 0)
				for _, id := range idList {
					for _, user := range userList {
						pattern = append(pattern, "(?,?)")
						args = append(args, *id, user.ID)
					}
				}

				raw := fmt.Sprintf("INSERT INTO `sql_manage_assignees` (`sql_manage_id`, `user_id`) VALUES %s",
					strings.Join(pattern, ", "))

				err = tx.Exec(raw, args...).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *Storage) GetSqlManageListByIDs(ids []*uint64) ([]*SqlManage, error) {
	sqlManageList := []*SqlManage{}
	err := s.db.Model(SqlManage{}).Where("id IN (?)", ids).Find(&sqlManageList).Error
	if err != nil {
		return nil, err
	}
	return sqlManageList, nil
}
