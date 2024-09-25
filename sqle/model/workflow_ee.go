//go:build enterprise
// +build enterprise
package model

import "github.com/actiontech/sqle/sqle/errors"


type WorkflowWithInstanceID struct {
	ID          uint    `json:"id"`
	WorkFlowID  string  `json:"workflow_id"`
	Subject     string  `json:"subject"`
	Desc        string  `json:"desc"`
	InstanceIDs RowList `json:"instance_ids"`
}

func (s *Storage) GetWorkflowsThatCanBeAssociatedToStage(instanceIdRange []uint64, excludeWorkflowIds []string) ([]*WorkflowWithInstanceID, error) {
	workflows := []*WorkflowWithInstanceID{}
	data := map[string]interface{}{
		"filter_instance_id": instanceIdRange,
	}
	if len(excludeWorkflowIds) > 0 {
		data["filter_workflow_id"] = excludeWorkflowIds
	}
	err := s.getListResult(getWorkflowByInstanceIDsBodyTpl, getWorkflowByInstanceIDsQueryBodyTpl, data, &workflows)
	if err != nil {
		return workflows, err
	}
	return workflows, errors.ConnectStorageErrWrapper(err)
}

var getWorkflowByInstanceIDsBodyTpl = `
SELECT 
	workflows.id, 
	workflows.workflow_id,
	workflows.subject, 
	workflows.desc, 
	GROUP_CONCAT(workflow_instance_records.instance_id, "") as instance_ids
FROM workflows 
LEFT JOIN workflow_records 
ON workflows.workflow_record_id = workflow_records.id
LEFT JOIN workflow_instance_records 
ON workflow_instance_records.workflow_record_id = workflow_records.id
{{- template "body" . -}}
`

var getWorkflowByInstanceIDsQueryBodyTpl = `
{{ define "body" }}
{{- if .filter_workflow_id }}
WHERE workflows.workflow_id NOT IN ( {{range $index, $element := .filter_workflow_id}}{{if $index}},{{end}}{{$element}}{{end}} )
{{- end }}
GROUP BY workflows.id
HAVING COUNT(CASE WHEN workflow_instance_records.instance_id NOT IN ( {{range $index, $element := .filter_instance_id}}{{if $index}},{{end}}{{$element}}{{end}} ) THEN 1 END) = 0
ORDER BY workflows.created_at DESC
LIMIT 100
{{ end }}
`
