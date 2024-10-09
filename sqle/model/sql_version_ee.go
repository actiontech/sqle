//go:build enterprise
// +build enterprise

package model

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

type SqlVersionListDetail struct {
	Id        uint           `json:"id"`
	Version   sql.NullString `json:"version"`
	Desc      sql.NullString `json:"description"`
	Status    sql.NullString `json:"status"`
	LockTime  *time.Time     `json:"lock_time"`
	CreatedAt *time.Time     `json:"created_at"`
}

var sqlVersionQueryTpl = `
SELECT  
	sv.id AS id,
	sv.version AS version,
	sv.description AS description,
	sv.status AS status,
	sv.lock_time AS lock_time,
	sv.created_at AS created_at
 
{{- template "body" . -}} 

{{- if .order_by -}}
ORDER BY {{.order_by}}
{{- if .is_asc }}
ASC
{{- else}}
DESC
{{- end -}}
{{else}}
ORDER BY sv.created_at desc
{{- end -}}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var sqlVersionCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var sqlVersionBodyTpl = `
{{ define "body" }}

FROM 
    sql_versions sv
WHERE 
    sv.deleted_at IS NULL

{{- if .fuzzy_search }}
AND (
sv.version LIKE '%{{ .fuzzy_search }}%'
OR
sv.description LIKE '%{{ .fuzzy_search }}%'
)
{{- end }}

{{- if .filter_by_created_at_from }}
AND sv.created_at >= :filter_by_created_at_from
{{- end }}

{{- if .filter_by_created_at_to }}
AND sv.created_at <= :filter_by_created_at_to
{{- end }}

{{- if .filter_by_lock_time_from }}
AND sv.lock_time >= :filter_by_lock_time_from
{{- end }}

{{- if .filter_by_lock_time_to }}
AND sv.lock_time <= :filter_by_lock_time_to
{{- end }}

{{- if .filter_by_version_status }}
AND sv.status = :filter_by_version_status
{{- end }}

{{ end }}
`

func (s *Storage) GetSqlVersionByReq(data map[string]interface{}) (
	list []*SqlVersionListDetail, count uint64, err error) {
	err = s.getListResult(sqlVersionBodyTpl, sqlVersionQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(sqlVersionBodyTpl, sqlVersionCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (s *Storage) GetSqlVersionDetailByVersionId(versionId uint) (*SqlVersion, bool, error) {
	version := &SqlVersion{}
	err := s.db.Preload("SqlVersionStage").Preload("SqlVersionStage.SqlVersionStagesDependency").Preload("SqlVersionStage.WorkflowVersionStage").Where("id=?", versionId).First(version).Error
	if err == gorm.ErrRecordNotFound {
		return version, false, nil
	}
	return version, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetStageDependenciesByStageId(stageId string) ([]*SqlVersionStagesDependency, error) {
	var dependencies []*SqlVersionStagesDependency
	err := s.db.Model(SqlVersionStagesDependency{}).Where("sql_version_stage_id = ?", stageId).Find(&dependencies).Error
	return dependencies, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) SaveSqlVersion(sqlVersion *SqlVersion) (uint, error) {
	var sqlVersionId uint
	err := s.Tx(func(txDB *gorm.DB) error {
		err := txDB.Model(&SqlVersion{}).Omit("SqlVersionStage.SqlVersionStagesDependency").Save(sqlVersion).Error
		if err != nil {
			return err
		}
		sqlVersionId = sqlVersion.ID
		err = s.SaveVersionStageDependency(txDB, sqlVersion.SqlVersionStage)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return sqlVersionId, nil
}

func (s *Storage) SaveVersionStageDependency(txDB *gorm.DB, stages []*SqlVersionStage) error {
	stageDepMap := make(map[int][]SqlVersionStagesDependency)
	for _, stage := range stages {
		deps := make([]SqlVersionStagesDependency, 0)
		for _, stageDep := range stage.SqlVersionStagesDependency {
			deps = append(deps, SqlVersionStagesDependency{
				StageInstanceID:     stageDep.StageInstanceID,
				NextStageInstanceID: stageDep.NextStageInstanceID,
			})
		}
		stageDepMap[stage.StageSequence] = deps
	}
	// 保存阶段依赖关系
	stageDeps := make([]*SqlVersionStagesDependency, 0)
	for _, versionStage := range stages {
		nextStage := GetNextStageBySqlVersionStage(stages, versionStage.StageSequence)
		for _, dep := range stageDepMap[versionStage.StageSequence] {
			sqlVersionStagesDep := &SqlVersionStagesDependency{}
			if nextStage != nil {
				sqlVersionStagesDep.SqlVersionStageID = versionStage.ID
				sqlVersionStagesDep.NextStageID = nextStage.ID
				sqlVersionStagesDep.StageInstanceID = dep.StageInstanceID
				sqlVersionStagesDep.NextStageInstanceID = dep.NextStageInstanceID
			} else {
				sqlVersionStagesDep.SqlVersionStageID = versionStage.ID
				sqlVersionStagesDep.StageInstanceID = dep.StageInstanceID
			}
			stageDeps = append(stageDeps, sqlVersionStagesDep)
		}
	}
	err := txDB.Save(stageDeps).Error
	if err != nil {
		return err
	}
	return nil
}

func GetNextStageBySqlVersionStage(stages []*SqlVersionStage, currentSequence int) *SqlVersionStage {
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].StageSequence < stages[j].StageSequence
	})
	for i, stage := range stages {
		if stage.StageSequence == currentSequence {
			if i+1 < len(stages) {
				return stages[i+1]
			}
			break
		}
	}
	return nil
}

func (s *Storage) GetStageWorkflowsByWorkflowIds(sqlVersionId uint, workflowIds []string) ([]*WorkflowVersionStage, error) {
	var stagesWorkflows []*WorkflowVersionStage
	err := s.db.Model(WorkflowVersionStage{}).Where("sql_version_id = ? AND workflow_id in (?)", sqlVersionId, workflowIds).Find(&stagesWorkflows).Error
	return stagesWorkflows, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateStageWorkflowExecTimeIfNeed(workflowId string) error {
	stage, exist, err := s.GetStageOfTheWorkflow(workflowId)
	if err != nil {
		return err
	}
	if !exist {
		// 工单没有关联版本阶段信息，不需要更新上线时间
		return nil
	}
	stagesWorkflows, err := s.GetStageWorkflowsByWorkflowIds(stage.SqlVersionID, []string{workflowId})
	if err != nil {
		return err
	}
	// 若上线时间已经有，则不进行更新，记录工单中第一个task的上线时间
	if stagesWorkflows[0].WorkflowExecTime == nil {
		err = s.db.Model(WorkflowVersionStage{}).Where("workflow_id = ?", workflowId).Update("workflow_exec_time", time.Now()).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetStageOfTheWorkflow(workflowId string) (*SqlVersionStage, bool, error) {
	stage := &SqlVersionStage{}
	err := s.db.Model(&SqlVersionStage{}).
		Joins("JOIN workflow_version_stages ON sql_version_stages.id = workflow_version_stages.sql_version_stage_id").
		Where("workflow_version_stages.workflow_id = ?", workflowId).First(stage).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return stage, true, nil
}

func (s *Storage) GetStageWorkflowByWorkflowId(sqlVersionId uint, workflowId string) (*WorkflowVersionStage, error) {
	var stagesWorkflow *WorkflowVersionStage
	err := s.db.Model(WorkflowVersionStage{}).Where("sql_version_id = ? AND workflow_id = ?", sqlVersionId, workflowId).Find(&stagesWorkflow).Error
	return stagesWorkflow, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowOfFirstStage(sqlVersionID uint, workflowId string) (*Workflow, error) {
	workflow := &Workflow{}
	err := s.db.Model(&Workflow{}).
		Joins("JOIN workflow_version_stages ON workflows.workflow_id = workflow_version_stages.workflow_id").
		Joins("JOIN sql_version_stages ON sql_version_stages.sql_version_id = workflow_version_stages.sql_version_id ").
		Where("workflow_version_stages.sql_version_id = ? AND workflow_version_stages.workflow_sequence IN "+
			"(SELECT workflow_sequence from workflow_version_stages WHERE workflow_id = ?)", sqlVersionID, workflowId).
		Order("sql_version_stages.stage_sequence ASC").First(workflow).Error
	if err != nil {
		return nil, err
	}
	return workflow, nil
}

func (s *Storage) GetWorkflowOfNextStage(versionId uint, workflowId string) (*SqlVersionStage, error) {
	stage, exist, err := s.GetStageOfTheWorkflow(workflowId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist, fmt.Errorf("workflow current stage not found"))
	}
	nextStage, exist, err := s.GetNextStageByStageSequence(versionId, stage.StageSequence)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist, fmt.Errorf("workflow next stage not found"))
	}
	return nextStage, nil
}

func (s *Storage) GetNextStageByStageSequence(versionId uint, sequence int) (*SqlVersionStage, bool, error) {
	stage := &SqlVersionStage{}
	// next stage sequence
	next := sequence + 1
	err := s.db.Where("sql_version_id = ? AND stage_sequence = ?", versionId, next).First(stage).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return stage, true, nil
}
func (s *Storage) UpdateWorkflowReleaseStatus(sqlVersionId uint, workflowId, status string) error {
	err := s.db.Model(WorkflowVersionStage{}).Where("sql_version_id = ? AND workflow_id = ?", sqlVersionId, workflowId).Update("workflow_release_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (stage SqlVersionStage) InitialStatusOfWorkflow() string {
	if len(stage.SqlVersionStagesDependency) > 0 && stage.SqlVersionStagesDependency[0].NextStageID == 0 {
		return WorkflowReleaseStatusNotNeedReleased
	}
	return WorkflowReleaseStatusIsBingReleased
}

func (s *Storage) UpdateSQLVersionById(versionId uint, sqlVersion map[string]interface{}) error {
	err := s.db.Model(&SqlVersion{}).Where("id = ?", versionId).
		Updates(sqlVersion).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteSqlVersionById(versionId uint) error {
	version, exist, err := s.GetSqlVersionDetailByVersionId(versionId)
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewDataNotExistErr("sql version not found")
	}
	err = s.Tx(func(txDB *gorm.DB) error {
		err := txDB.Delete(version).Error
		if err != nil {
			return err
		}
		err = txDB.Select("SqlVersionStagesDependency").Delete(version.SqlVersionStage).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateSQLVersionStageByVersionId(versionId uint, deleteStageIds []uint, addVersionStages []*SqlVersionStage) error {
	// 因为只有未关联工单的版本才能修改阶段，所以这里可以覆盖式更新阶段及数据源的依赖关系
	err := s.Tx(func(txDB *gorm.DB) error {
		// 删除阶段
		err := s.db.Unscoped().Where("sql_version_id = ?", versionId).Delete(&SqlVersionStage{}).Error
		if err != nil {
			return err
		}
		// 删除阶段数据源依赖关系
		err = s.db.Unscoped().Where("sql_version_stage_id IN (?)", deleteStageIds).Delete(&SqlVersionStagesDependency{}).Error
		if err != nil {
			return err
		}
		err = s.db.Model(&SqlVersionStage{}).Omit("SqlVersionStagesDependency").Save(addVersionStages).Error
		if err != nil {
			return err
		}
		err = s.SaveVersionStageDependency(txDB, addVersionStages)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 根据工单id获取在版本中关联阶段的工单
func (s *Storage) GetAssociatedStageWorkflows(workflowId string) ([]*AssociatedStageWorkflow, error) {
	var stageWorkflows []*AssociatedStageWorkflow
	err := s.db.Model(&WorkflowVersionStage{}).Select("workflow_version_stages.workflow_id,"+
		"refs.sql_version_stage_id,"+
		"svs.stage_sequence ,"+
		"w.subject AS workflow_name,"+
		"wr.status").
		Joins("INNER JOIN ( "+
			"SELECT sql_version_id, workflow_sequence,workflow_id,sql_version_stage_id "+
			"FROM workflow_version_stages "+
			"WHERE workflow_id = ?"+
			") AS refs ON workflow_version_stages.sql_version_id = refs.sql_version_id "+
			"AND workflow_version_stages.workflow_sequence = refs.workflow_sequence", workflowId).
		Joins("INNER JOIN sql_version_stages svs ON svs.id = workflow_version_stages.sql_version_stage_id").
		Joins("INNER JOIN workflows w ON workflow_version_stages.workflow_id = w.workflow_id").
		Joins("INNER JOIN workflow_records wr ON w.workflow_record_id = wr.id ").
		Scan(&stageWorkflows).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	return stageWorkflows, nil
}

func (s *Storage) GetFirstStageOfSQLVersion(sqlVersionID uint) (*SqlVersionStage, error) {
	firstStage := &SqlVersionStage{}
	err := s.db.Model(&SqlVersionStage{}).Preload("SqlVersionStagesDependency").Preload("WorkflowVersionStage").Where("sql_version_id = ?", sqlVersionID).Order("stage_sequence ASC").First(firstStage).Error
	if err != nil {
		return nil, err
	}
	return firstStage, nil
}

func (s *Storage) GetStageOfSQLVersion(sqlVersionID, stageID uint) (*SqlVersionStage, error) {
	stage := &SqlVersionStage{}
	err := s.db.Model(&SqlVersionStage{}).
		Preload("SqlVersionStagesDependency").
		Preload("WorkflowVersionStage").
		Where("id = ? AND sql_version_id = ? ", stageID, sqlVersionID).
		First(stage).Error
	if err != nil {
		return nil, err
	}
	return stage, nil
}

func (s *Storage) BatchCreateWorkflowVerionRelation(stage *SqlVersionStage, workflowIds []string) error {
	workflowVersionModels := make([]*WorkflowVersionStage, 0, len(workflowIds))
	for index, woworkflowId := range workflowIds {
		workflowVersionModels = append(workflowVersionModels, &WorkflowVersionStage{
			WorkflowID:            woworkflowId,
			SqlVersionID:          stage.SqlVersionID,
			SqlVersionStageID:     stage.ID,
			WorkflowSequence:      len(stage.WorkflowVersionStage) + index + 1,
			WorkflowReleaseStatus: WorkflowReleaseStatusIsBingReleased,
			WorkflowExecTime:      nil,
		})
	}
	return s.db.Model(&WorkflowVersionStage{}).CreateInBatches(&workflowVersionModels, 100).Error
}

func (s *Storage) GetWorkflowVersionRelationByWorkflowId(workflowId string) (relation *WorkflowVersionStage, exist bool, err error) {
	relation = &WorkflowVersionStage{}
	err = s.db.Model(&WorkflowVersionStage{}).
		Where("workflow_id = ? ", workflowId).
		Find(relation).Error
	if err != nil {
		return nil, false, err
	}
	if relation.ID == 0 {
		return nil, false, nil
	}
	return relation, true, nil
}
