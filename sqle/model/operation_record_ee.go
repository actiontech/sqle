//go:build enterprise
// +build enterprise

package model

const (
	OperationRecordPlatform = "--"

	// operation record type
	OperationRecordTypeProject             = "project"
	OperationRecordTypeInstance            = "instance"
	OperationRecordTypeProjectRuleTemplate = "project_rule_template"
	OperationRecordTypeWorkflowTemplate    = "workflow_template"
	OperationRecordTypeWhiteList           = "white_list"
	OperationRecordTypeAuditPlan           = "audit_plan"
	OperationRecordTypeWorkflow            = "workflow"

	// operation record action
	OperationRecordActionCreateProject             = "create_project"
	OperationRecordActionCreateProjectRuleTemplate = "create_project_rule_template"
	OperationRecordActionDeleteProjectRuleTemplate = "delete_project_rule_template"
	OperationRecordActionUpdateProjectRuleTemplate = "update_project_rule_template"

	// Status operation record status
	OperationRecordStatusSuccess = "success"
	OperationRecordStatusFail    = "fail"
)
