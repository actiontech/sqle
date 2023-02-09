//go:build enterprise
// +build enterprise

package model

const (
	OperationRecordPlatform = "--"

	// operation record type
	OperationRecordTypeProjectManage = "project_manage"

	// operation record action
	OperationRecordActionCreateProject = "create_project"

	// Status operation record status
	OperationRecordStatusSuccess = "success"
	OperationRecordStatusFail    = "fail"
)
