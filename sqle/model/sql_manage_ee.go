//go:build enterprise
// +build enterprise

package model

import (
	"github.com/jinzhu/gorm"
)

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
