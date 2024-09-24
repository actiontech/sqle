package model

import (
	"database/sql"
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

func (s *Storage) GetSqlVersionDetailByVersionId(versionId string) (*SqlVersion, bool, error) {
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

func (s *Storage) BatchSaveSqlVersion(reqSqlVersion *SqlVersion) error {
	err := s.Tx(func(txDB *gorm.DB) error {
		sqlVersion := &SqlVersion{
			Version:     reqSqlVersion.Version,
			Description: reqSqlVersion.Description,
			Status:      SqlVersionStatusReleased,
			ProjectId:   reqSqlVersion.ProjectId,
		}
		err := txDB.Save(sqlVersion).Error
		if err != nil {
			return err
		}
		// 保存版本阶段
		versionStages := make([]*SqlVersionStage, 0, len(reqSqlVersion.SqlVersionStage))
		stageDepMap := make(map[int][]SqlVersionStagesDependency)
		for _, stage := range reqSqlVersion.SqlVersionStage {
			versionStages = append(versionStages, &SqlVersionStage{
				SqlVersionID:  sqlVersion.ID,
				Name:          stage.Name,
				StageSequence: stage.StageSequence,
			})
			deps := make([]SqlVersionStagesDependency, 0)
			for _, stageDep := range stage.SqlVersionStagesDependency {
				deps = append(deps, SqlVersionStagesDependency{
					StageInstanceID:     stageDep.StageInstanceID,
					NextStageInstanceID: stageDep.NextStageInstanceID,
				})
			}
			stageDepMap[stage.StageSequence] = deps
		}
		err = txDB.Save(versionStages).Error
		if err != nil {
			return err
		}

		// 保存阶段依赖关系
		stageDeps := make([]*SqlVersionStagesDependency, 0)
		for _, versionStage := range versionStages {
			nextStage, exist, err := func(versionId uint, sequence int) (*SqlVersionStage, bool, error) {
				stage := &SqlVersionStage{}
				// next stage sequence
				next := sequence + 1
				err := txDB.Where("sql_version_id = ? AND stage_sequence = ?", versionId, next).First(stage).Error
				if err == gorm.ErrRecordNotFound {
					return nil, false, nil
				} else if err != nil {
					return nil, false, err
				}
				return stage, true, nil
			}(versionStage.SqlVersionID, versionStage.StageSequence)
			if err != nil {
				return err
			}

			for _, dep := range stageDepMap[versionStage.StageSequence] {
				sqlVersionStagesDep := &SqlVersionStagesDependency{}
				if exist {
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
		err = txDB.Save(stageDeps).Error
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

func (s *Storage) GetFirstStageOfSQLVersion(sqlVersionID uint) (*SqlVersionStage, error) {
	firstStage := &SqlVersionStage{}
	err := s.db.Model(&SqlVersionStage{}).Preload("SqlVersionStagesDependency").Preload("WorkflowVersionStage").Where("sql_version_id = ?", sqlVersionID).Order("stage_sequence ASC").First(firstStage).Error
	if err != nil {
		return nil, err
	}
	return firstStage, nil
}
