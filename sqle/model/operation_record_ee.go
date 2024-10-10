//go:build enterprise
// +build enterprise

package model

const (
	// operation record type
	OperationRecordTypeProject             = "project"
	OperationRecordTypeInstance            = "instance"
	OperationRecordTypeProjectRuleTemplate = "project_rule_template"
	OperationRecordTypeWorkflowTemplate    = "workflow_template"
	OperationRecordTypeAuditPlan           = "audit_plan"
	OperationRecordTypeWorkflow            = "workflow"
	OperationRecordTypeGlobalUser          = "global_user"
	OperationRecordTypeGlobalRuleTemplate  = "global_rule_template"
	OperationRecordTypeSystemConfiguration = "system_configuration"
	OperationRecordTypeProjectMember       = "project_member"

	// operation record action
	OperationRecordActionCreateProject               = "create_project"
	OperationRecordActionDeleteProject               = "delete_project"
	OperationRecordActionUpdateProject               = "update_project"
	OperationRecordActionArchiveProject              = "archive_project"
	OperationRecordActionUnarchiveProject            = "unarchive_project"
	OperationRecordActionCreateInstance              = "create_instance"
	OperationRecordActionUpdateInstance              = "update_instance"
	OperationRecordActionDeleteInstance              = "delete_instance"
	OperationRecordActionCreateProjectRuleTemplate   = "create_project_rule_template"
	OperationRecordActionDeleteProjectRuleTemplate   = "delete_project_rule_template"
	OperationRecordActionUpdateProjectRuleTemplate   = "update_project_rule_template"
	OperationRecordActionUpdateWorkflowTemplate      = "update_workflow_template"
	OperationRecordActionCreateAuditPlan             = "create_audit_plan"
	OperationRecordActionDeleteAuditPlan             = "delete_audit_plan"
	OperationRecordActionUpdateAuditPlan             = "update_audit_plan"
	OperationRecordActionCreateWorkflow              = "create_workflow"
	OperationRecordActionCancelWorkflow              = "cancel_workflow"
	OperationRecordActionApproveWorkflow             = "approve_workflow"
	OperationRecordActionRejectWorkflow              = "reject_workflow"
	OperationRecordActionExecuteWorkflow             = "execute_workflow"
	OperationRecordActionScheduleWorkflow            = "schedule_workflow"
	OperationRecordActionUpdateWorkflow              = "update_workflow"
	OperationRecordActionCreateUser                  = "create_user"
	OperationRecordActionUpdateUser                  = "update_user"
	OperationRecordActionDeleteUser                  = "delete_user"
	OperationRecordActionCreateGlobalRuleTemplate    = "create_global_rule_template"
	OperationRecordActionUpdateGlobalRuleTemplate    = "update_global_rule_template"
	OperationRecordActionDeleteGlobalRuleTemplate    = "delete_global_rule_template"
	OperationRecordActionUpdateDingTalkConfiguration = "update_ding_talk_configuration"
	OperationRecordActionUpdateSMTPConfiguration     = "update_smtp_configuration"
	OperationRecordActionUpdateWechatConfiguration   = "update_wechat_configuration"
	OperationRecordActionUpdateSystemVariables       = "update_system_variables"
	OperationRecordActionUpdateLDAPConfiguration     = "update_ldap_configuration"
	OperationRecordActionUpdateOAuth2Configuration   = "update_oauth2_configuration"
	OperationRecordActionCreateMember                = "create_member"
	OperationRecordActionCreateMemberGroup           = "create_member_group"
	OperationRecordActionDeleteMember                = "delete_member"
	OperationRecordActionDeleteMemberGroup           = "delete_member_group"
	OperationRecordActionUpdateMember                = "update_member"
	OperationRecordActionUpdateMemberGroup           = "update_member_group"

	// Status operation record status
	OperationRecordStatusSucceeded = "succeeded"
	OperationRecordStatusFailed    = "failed"
)

var operationRecordQueryTpl = `
SELECT o.id,
       o.operation_time,
       o.operation_user_name,
       o.operation_req_ip,
       o.operation_type_name,
       o.operation_action,
       o.operation_project_name,
       o.operation_status,
       o.operation_i18n_content
{{- template "body" . -}}
ORDER BY o.operation_time DESC
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var operationRecordExportTpl = `
SELECT o.operation_time,
	   o.operation_project_name,
       o.operation_user_name,
       o.operation_action,
       o.operation_status,
       o.operation_i18n_content
{{- template "body" . -}}
ORDER BY o.operation_time DESC
`

var operationRecordWorkflowsCountTpl = `SELECT COUNT(*)

{{- template "body" . -}}
`

var operationRecordQueryBodyTpl = `
{{ define "body" }}
FROM operation_records o

WHERE o.deleted_at IS NULL 

{{- if .filter_operate_time_from }}
AND o.operation_time > :filter_operate_time_from
{{- end }}

{{- if .filter_operate_time_to }}
AND o.operation_time < :filter_operate_time_to
{{- end }}

{{- if .filter_operate_project_name }}
AND o.operation_project_name = :filter_operate_project_name
{{- end }}

{{- if .fuzzy_search_operate_user_name }}
AND o.operation_user_name LIKE '%{{ .fuzzy_search_operate_user_name }}%'
{{- end }}

{{- if .filter_operate_type_name }}
AND o.operation_type_name = :filter_operate_type_name
{{- end }}

{{- if .filter_operate_action }}
AND o.operation_action = :filter_operate_action
{{- end }}

{{ end }}

`

func (s *Storage) GetOperationRecordList(data map[string]interface{}) (result []*OperationRecord, count uint64, err error) {
	err = s.getListResult(operationRecordQueryBodyTpl, operationRecordQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}

	count, err = s.getCountResult(operationRecordQueryBodyTpl, operationRecordWorkflowsCountTpl, data)

	return result, count, err
}

func (s *Storage) GetOperationRecordExportList(data map[string]interface{}) (result []*OperationRecord, err error) {
	err = s.getListResult(operationRecordQueryBodyTpl, operationRecordExportTpl, data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
