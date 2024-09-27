package model

import (
	"time"
)

type SqlVersion struct {
	Model
	Version     string     `json:"version" gorm:"type:varchar(255) not null"`
	Description string     `json:"description" gorm:"type:varchar(512)"`
	Status      string     `json:"status" gorm:"type:varchar(255)"`
	LockTime    *time.Time `json:"lock_time" gorm:"type:datetime(3)"`
	ProjectId   ProjectUID `gorm:"index; not null"`

	SqlVersionStage []*SqlVersionStage
}

const (
	SqlVersionStatusReleased = "is_being_released"
	SqlVersionStatusLock     = "locked"
)

type SqlVersionStage struct {
	Model
	SqlVersionID uint   `json:"sql_version_id" gorm:"not null"`
	Name         string `json:"name" gorm:"type:varchar(255) not null"`
	// stage_sequence标识版本阶段的顺序，是一段从1开始连续的int值
	StageSequence int `json:"stage_sequence" gorm:"type:int not null"`

	SqlVersionStagesDependency []*SqlVersionStagesDependency
	WorkflowVersionStage       []*WorkflowVersionStage
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
	// workflow_sequence标识工单所处阶段的顺序及占位，当工单不需要发布到下一阶段时（如工单关闭），下一阶段的workflow_sequence可能是不连续的
	// 同一版本每个阶段之间的工单占位和顺序都相互对应，如：开发阶段发布到测试阶段，workflow_sequence为1的工单发布成功后workflow_sequence也为1
	WorkflowSequence      int        `json:"workflow_sequence" gorm:"type:int"`
	WorkflowReleaseStatus string     `json:"workflow_release_status" gorm:"type:varchar(255) not null"`
	WorkflowExecTime      *time.Time `json:"workflow_exec_time" gorm:"type:datetime(3)"`

	Workflow *Workflow `gorm:"foreignkey:WorkflowID"`
}

const (
	WorkflowReleaseStatusIsBingReleased   = "wait_for_release"
	WorkflowReleaseStatusHaveBeenReleased = "released"
	WorkflowReleaseStatusNotNeedReleased  = "not_need_release"
)
