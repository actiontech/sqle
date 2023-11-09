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
