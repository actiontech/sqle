package model

import (
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"
	"gorm.io/gorm"
)

type AuditFile struct {
	Model
	TaskId     uint   `json:"-" gorm:"index"`
	UniqueName string `json:"unique_name" gorm:"type:varchar(255)"`
	FileHost   string `json:"file_host" gorm:"type:varchar(255)"`
	FileName   string `json:"file_name" gorm:"type:varchar(255)"`
	/*
		默认顺序
			1. 单个文件的执行顺序都是1，在页面上默认显示执行顺序1
			2. 对zip文件来说，zip文件本身的执行顺序默认为0，因为不执行也不展示
			3. 对zip文件中的子文件，只保存sql文件，sql文件的执行顺序从1开始，按序递增，默认顺序为读取到的sql的先后顺序
	*/
	ExecOrder uint `json:"exec_order"`
	ParentID  uint `json:"parent_id"`
}

const FixFilePath string = "audit_files/"

func NewFileRecord(taskID, order uint, nickName, uniqueName string) *AuditFile {
	return &AuditFile{
		TaskId:     taskID,
		UniqueName: uniqueName,
		FileHost:   config.GetOptions().SqleOptions.ReportHost,
		FileName:   nickName,
		ExecOrder:  order,
	}
}

func DefaultFilePath(fileName string) string {
	return FixFilePath + fileName
}

func GenUniqueFileName() string {
	return time.Now().Format("2006-01-02") + "_" + utils.GenerateRandomString(5)
}

func (s *Storage) GetFileByTaskId(taskId string) ([]*AuditFile, error) {
	auditFiles := []*AuditFile{}
	result := s.db.Where("task_id = ?", taskId).Find(&auditFiles)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return auditFiles, errors.New(errors.ConnectStorageError, result.Error)
}

func (s *Storage) GetParentFileByTaskId(taskId string) (*AuditFile, bool, error) {
	auditFile := &AuditFile{}
	err := s.db.Where("parent_id = 0 AND task_id = ?", taskId).First(auditFile).Error
	if err == gorm.ErrRecordNotFound {
		return auditFile, false, nil
	}
	return auditFile, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetFileByIds(fileIds []uint) ([]*AuditFile, error) {
	auditFiles := []*AuditFile{}
	result := s.db.Where("id in (?)", fileIds).Find(&auditFiles)
	if result.RowsAffected == 0 {
		return auditFiles, nil
	}
	return auditFiles, errors.New(errors.ConnectStorageError, result.Error)
}

// expired time 24hour
func (s *Storage) GetExpiredFileWithNoWorkflow() ([]AuditFile, error) {
	auditFiles := []AuditFile{}
	err := s.db.Model(&AuditFile{}).
		Joins("LEFT JOIN workflow_instance_records wir ON audit_files.task_id = wir.task_id").
		Where("wir.task_id IS NULL AND audit_files.deleted_at IS NULL").                // 删除没有提交为工单的SQL文件
		Where("audit_files.file_host = ?", config.GetOptions().SqleOptions.ReportHost). // 删除本机文件
		Where("audit_files.created_at < ?", time.Now().Add(-24*time.Hour)).             // 减少提交前文件就被删除的几率
		Find(&auditFiles).Error
	if len(auditFiles) == 0 {
		return nil, nil
	}
	return auditFiles, errors.New(errors.ConnectStorageError, err)
}

// expired time 7*24hour
func (s *Storage) GetExpiredFile() ([]AuditFile, error) {
	auditFiles := []AuditFile{}
	err := s.db.Model(&AuditFile{}).
		Where("audit_files.file_host = ?", config.GetOptions().SqleOptions.ReportHost). // 删除本机文件
		Where("audit_files.created_at < ?", time.Now().Add(-7*24*time.Hour)).           // 减少提交前文件就被删除的几率
		Find(&auditFiles).Error
	if len(auditFiles) == 0 {
		return nil, nil
	}
	return auditFiles, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) BatchCreateFileRecords(records []*AuditFile) error {
	return s.Tx(func(txDB *gorm.DB) error {
		for _, record := range records {
			record.ID = 0
			if err := txDB.Create(record).Error; err != nil {
				txDB.Rollback()
				return err
			}
		}
		return nil
	})
}

func (s *Storage) BatchSaveFileRecords(records []*AuditFile) error {
	return s.Tx(func(txDB *gorm.DB) error {
		for _, record := range records {
			if err := txDB.Save(record).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 单个文件中触发XX审核等级的SQL的数量统计
type FileAuditStatistic struct {
	ErrorCount   uint `json:"error_count"`
	WarningCount uint `json:"warning_count"`
	NoticeCount  uint `json:"notice_count"`
	NormalCount  uint `json:"normal_count"`
}

// 单个文件中SQL执行状态的数量统计
type FileExecStatistic struct {
	FailedCount             uint `json:"failed_count"`
	SucceededCount          uint `json:"succeeded_count"`
	InitializedCount        uint `json:"initialized_count"`
	DoingCount              uint `json:"doing_count"`
	ManuallyExecutedCount   uint `json:"manually_executed_count"`
	TerminateSucceededCount uint `json:"terminate_succeeded_count"`
	TerminateFailedCount    uint `json:"terminate_failed_count"`
}

func (a FileExecStatistic) FileExecStatus() string {
	// 一旦存在执行失败的SQL，则执行失败
	if a.FailedCount > 0 {
		return SQLExecuteStatusFailed
	}
	// 手工执行后，手工执行和SQLE执行互斥，因此先判断
	if a.ManuallyExecutedCount > 0 {
		return SQLExecuteStatusManuallyExecuted
	}
	// 终止上线后，终止上线操作在上线前或者上线中，若有终止失败，则为失败
	if a.TerminateFailedCount > 0 {
		return SQLExecuteStatusTerminateFailed
	}
	if a.TerminateSucceededCount > 0 {
		return SQLExecuteStatusTerminateSucc
	}
	// 执行上线中
	if a.InitializedCount > 0 {
		// 若仅包含初始化的SQL，则状态为初始化
		if a.DoingCount == 0 && a.FailedCount == 0 && a.SucceededCount == 0 {
			return SQLExecuteStatusInitialized
		}
		// 若有其他状态，则为执行中
		return SQLExecuteStatusDoing
	}
	if a.DoingCount > 0 {
		// 若不包含初始化的SQL，但存在正在执行的SQL，则文件状态为执行中
		return SQLExecuteStatusDoing
	}
	// 执行完毕后 程序执行到这里，其他所有状态数量均等于0，若存在成功的SQL，则执行成功
	if a.SucceededCount > 0 {
		return SQLExecuteStatusSucceeded
	}
	// 如果没有执行成功或者失败的SQL，则说明没有SQL，则为初始化
	return SQLExecuteStatusInitialized
}

type AuditResultStatistic struct {
	ExecFileID   string `json:"exec_file_id"`
	ExecFileName string `json:"exec_file_name"`
	ExecOrder    uint   `json:"exec_order"`
	FileAuditStatistic
	FileExecStatistic
}

func (s *Storage) GetAuditStatisticByTaskId(data map[string]interface{}) (
	result []*AuditResultStatistic, count uint64, err error) {
	// add key value because suspension(:) will be considered as input variables
	data["key_error"] = "{\"level\": \"error\"}"
	data["key_warn"] = "{\"level\": \"warn\"}"
	data["key_normal"] = "{\"level\": \"normal\"}"
	data["key_notice"] = "{\"level\": \"notice\"}"
	// 执行报错的规规则计入warn中
	data["execution_failed_true"] = "{\"execution_failed\": true}"
	data["execution_failed_false"] = "{\"execution_failed\": false}"
	err = s.getListResult(auditFileStatisticQueryBodyTpl, auditFileStatisticQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult("", auditFileStatisticCountTpl, data)
	return result, count, err
}

var auditFileStatisticQueryTpl string = `
	SELECT 
		a_file.id AS exec_file_id,
		a_file.exec_order,
		a_file.file_name AS exec_file_name,
		SUM(CASE WHEN e_sql.exec_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
		SUM(CASE WHEN e_sql.exec_status = 'succeeded' THEN 1 ELSE 0 END) AS succeeded_count,
		SUM(CASE WHEN e_sql.exec_status = 'initialized' THEN 1 ELSE 0 END) AS initialized_count,
		SUM(CASE WHEN e_sql.exec_status = 'doing' THEN 1 ELSE 0 END) AS doing_count,
		SUM(CASE WHEN e_sql.exec_status = 'manually_executed' THEN 1 ELSE 0 END) AS manually_executed_count,
		SUM(CASE WHEN e_sql.exec_status = 'terminate_succeeded' THEN 1 ELSE 0 END) AS terminate_succeeded_count,
		SUM(CASE WHEN e_sql.exec_status = 'terminate_failed' THEN 1 ELSE 0 END) AS terminate_failed_count,
		SUM(JSON_CONTAINS(e_sql.audit_results, :key_error) AND JSON_CONTAINS(e_sql.audit_results, :execution_failed_false)) AS error_count,
		SUM(JSON_CONTAINS(e_sql.audit_results, :key_warn) OR JSON_CONTAINS(e_sql.audit_results, :execution_failed_true)) AS warning_count,
		SUM(JSON_CONTAINS(e_sql.audit_results, :key_notice) AND JSON_CONTAINS(e_sql.audit_results, :execution_failed_false)) AS notice_count,
		SUM(JSON_CONTAINS(e_sql.audit_results, :key_normal) AND JSON_CONTAINS(e_sql.audit_results, :execution_failed_false)) AS normal_count
	{{- template "body" . -}}
	{{- if .limit }}
		LIMIT :limit OFFSET :offset	
	{{- end -}}
`

var auditFileStatisticQueryBodyTpl = `
	{{ define "body" }}
		FROM 
			audit_files AS a_file
		LEFT JOIN
			execute_sql_detail AS e_sql
		ON 
			a_file.task_id = e_sql.task_id
		AND 
			a_file.file_name = e_sql.source_file
		WHERE 
			a_file.task_id = :task_id
		AND
			e_sql.exec_status IS NOT NULL
		GROUP BY 
			a_file.id,a_file.file_name,a_file.exec_order
		ORDER BY
			a_file.exec_order
	{{- end }}
`

var auditFileStatisticCountTpl = `
	SELECT 
		COUNT(DISTINCT a_file.id)
	FROM 
		audit_files AS a_file
	LEFT JOIN
		execute_sql_detail AS e_sql
	ON 
		a_file.task_id = e_sql.task_id
	AND 
		a_file.file_name = e_sql.source_file
	WHERE 
		a_file.task_id = :task_id
	AND
		e_sql.exec_status IS NOT NULL
`

type AuditFileExecStatistic struct {
	ExecFileID   string `json:"exec_file_id"`
	ExecFileName string `json:"exec_file_name"`
	FileExecStatistic
}

func (s *Storage) GetAuditFileExecStatisticByFileId(data map[string]interface{}) (
	overview *AuditFileExecStatistic, err error) {
	result := make([]*AuditFileExecStatistic, 0, 1)
	err = s.getListResult(auditFileExecStatisticQueryBodyTpl, auditFileExecStatisticQueryTpl, data, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("can not find any file execute overview match with this task")
	}
	return result[0], err
}

var auditFileExecStatisticQueryTpl string = `
	SELECT 
		a_file.id AS exec_file_id,
		a_file.file_name AS exec_file_name,
		SUM(CASE WHEN e_sql.exec_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
		SUM(CASE WHEN e_sql.exec_status = 'succeeded' THEN 1 ELSE 0 END) AS succeeded_count,
		SUM(CASE WHEN e_sql.exec_status = 'initialized' THEN 1 ELSE 0 END) AS initialized_count,
		SUM(CASE WHEN e_sql.exec_status = 'doing' THEN 1 ELSE 0 END) AS doing_count,
		SUM(CASE WHEN e_sql.exec_status = 'manually_executed' THEN 1 ELSE 0 END) AS manually_executed_count,
		SUM(CASE WHEN e_sql.exec_status = 'terminate_succeeded' THEN 1 ELSE 0 END) AS terminate_succeeded_count,
		SUM(CASE WHEN e_sql.exec_status = 'terminate_failed' THEN 1 ELSE 0 END) AS terminate_failed_count
	{{- template "body" . -}}
`

var auditFileExecStatisticQueryBodyTpl = `
	{{ define "body" }}
		FROM 
			audit_files AS a_file
		LEFT JOIN
			execute_sql_detail AS e_sql
		ON 
			a_file.task_id = e_sql.task_id
		AND 
			a_file.file_name = e_sql.source_file
		WHERE 
			a_file.task_id = :task_id
		AND
			e_sql.exec_status IS NOT NULL
		AND 
		a_file.id = :file_id
		GROUP BY 
			a_file.id,a_file.file_name
		LIMIT 1
	{{- end }}
`
