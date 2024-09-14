package model

import "time"

type SqlVersion struct {
	Model
	VersionNumber string     `json:"version_number" gorm:"type:varchar(255) not null"`
	Desc          string     `json:"desc" gorm:"type:varchar(512)"`
	Status        string     `json:"status" gorm:"type:varchar(255)"`
	LockTime      *time.Time `json:"lock_time" gorm:"type:datetime(3)"`
	IsLocked      bool       `json:"is_locked" gorm:"type:bool" example:"false"`
	ProjectId     ProjectUID `gorm:"index; not null"`

	SqlVersionStage []*SqlVersionStage
}

type SqlVersionStage struct {
	Model
	SqlVersionID  uint   `json:"sql_version_id" gorm:"not null"`
	Name          string `json:"name" gorm:"type:varchar(255) not null"`
	StageSequence int    `json:"stage_sequence" gorm:"type:int not null"`

	SqlVersionStagesDependency []*SqlVersionStagesDependency
	WorkflowReleaseStage       []*WorkflowReleaseStage
}

type SqlVersionStagesDependency struct {
	Model
	SqlVersionStageID   uint `json:"sql_version_stage_id" gorm:"not null"`
	NextStageID         uint `json:"next_stage_id"`
	StageInstanceID     uint `json:"stage_instance_id" gorm:"not null"`
	NextStageInstanceID uint `json:"next_stage_instance_id"`
}

type WorkflowReleaseStage struct {
	Model
	WorkflowID        string `json:"workflow_id" gorm:"not null"`
	SqlVersionStageID uint   `json:"sql_version_stage_id"`
	WorkflowSequence  int    `json:" workflow_sequence" gorm:"type:int"`

	Workflow *Workflow `gorm:"foreignkey:WorkflowID"`
}
