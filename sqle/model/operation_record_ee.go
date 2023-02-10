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
	OperationRecordTypeAuditPlan           = "audit_plan"
	OperationRecordTypeWorkflow            = "workflow"

	// operation record action
	OperationRecordActionCreateProject             = "create_project"
	OperationRecordActionCreateProjectRuleTemplate = "create_project_rule_template"
	OperationRecordActionDeleteProjectRuleTemplate = "delete_project_rule_template"
	OperationRecordActionUpdateProjectRuleTemplate = "update_project_rule_template"
	OperationRecordActionUpdateWorkflowTemplate    = "update_workflow_template"
	OperationRecordActionCreateAuditPlan           = "create_audit_plan"
	OperationRecordActionDeleteAuditPlan           = "delete_audit_plan"
	OperationRecordActionUpdateAuditPlan           = "update_audit_plan"

	// Status operation record status
	OperationRecordStatusSuccess = "success"
	OperationRecordStatusFail    = "fail"
)
