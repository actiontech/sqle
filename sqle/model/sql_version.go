package model

import (
	"time"
)

type SqlVersion struct {
	Model
	Version   string     `json:"version" gorm:"type:varchar(255) not null"`
	Desc      string     `json:"desc" gorm:"type:varchar(512)"`
	Status    string     `json:"status" gorm:"type:varchar(255)"`
	LockTime  *time.Time `json:"lock_time" gorm:"type:datetime(3)"`
	ProjectId ProjectUID `gorm:"index; not null"`

	SqlVersionStage []*SqlVersionStage
}

const (
	SqlVersionStatusReleased = "is_being_released"
	SqlVersionStatusLock     = "locked"
)

type SqlVersionStage struct {
	Model
	SqlVersionID  uint   `json:"sql_version_id" gorm:"not null"`
	Name          string `json:"name" gorm:"type:varchar(255) not null"`
	StageSequence int    `json:"stage_sequence" gorm:"type:int not null"`

	SqlVersionStagesDependency []*SqlVersionStagesDependency
	WorkflowReleaseStage       []*WorkflowVersionStage
}

type SqlVersionStagesDependency struct {
	Model
	SqlVersionStageID   uint   `json:"sql_version_stage_id" gorm:"not null"`
	NextStageID         uint   `json:"next_stage_id"`
	StageInstanceID     uint64 `json:"stage_instance_id" gorm:"not null"`
	NextStageInstanceID uint64 `json:"next_stage_instance_id"`
}

type WorkflowVersionStage struct {
	Model
	WorkflowID        string `json:"workflow_id" gorm:"not null"`
	SqlVersionID      uint   `json:"sql_version_id"`
	SqlVersionStageID uint   `json:"sql_version_stage_id"`
	WorkflowSequence  int    `json:" workflow_sequence" gorm:"type:int"`

	Workflow *Workflow `gorm:"foreignkey:WorkflowID"`
}