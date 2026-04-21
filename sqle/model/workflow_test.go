package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDataExportWorkflowTemplate(t *testing.T) {
	cases := []struct {
		name      string
		projectId string
		checks    func(t *testing.T, tmpl *WorkflowTemplate)
	}{
		{
			name:      "returns correct WorkflowType",
			projectId: "proj-123",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Equal(t, WorkflowTemplateTypeDataExport, tmpl.WorkflowType)
			},
		},
		{
			name:      "returns correct name format",
			projectId: "proj-456",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Equal(t, "proj-456-DataExportWorkflowTemplate", tmpl.Name)
			},
		},
		{
			name:      "has single export_review step",
			projectId: "proj-789",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Len(t, tmpl.Steps, 1)
				assert.Equal(t, WorkflowStepTypeExportReview, tmpl.Steps[0].Typ)
			},
		},
		{
			name:      "step has ApprovedByAuthorized true",
			projectId: "proj-abc",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.True(t, tmpl.Steps[0].ApprovedByAuthorized.Valid)
				assert.True(t, tmpl.Steps[0].ApprovedByAuthorized.Bool)
			},
		},
		{
			name:      "sets correct ProjectId",
			projectId: "proj-xyz",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Equal(t, ProjectUID("proj-xyz"), tmpl.ProjectId)
			},
		},
		{
			name:      "step number is 1",
			projectId: "proj-num",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Equal(t, uint(1), tmpl.Steps[0].Number)
			},
		},
		{
			name:      "AllowSubmitWhenLessAuditLevel is empty for data export",
			projectId: "proj-level",
			checks: func(t *testing.T, tmpl *WorkflowTemplate) {
				assert.Empty(t, tmpl.AllowSubmitWhenLessAuditLevel)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := DefaultDataExportWorkflowTemplate(tc.projectId)
			tc.checks(t, tmpl)
		})
	}
}

func TestWorkflowStepTypesByWorkflowType(t *testing.T) {
	cases := []struct {
		name         string
		workflowType string
		expected     []string
	}{
		{
			name:         "workflow allows sql_review and sql_execute",
			workflowType: WorkflowTemplateTypeWorkflow,
			expected:     []string{WorkflowStepTypeSQLReview, WorkflowStepTypeSQLExecute},
		},
		{
			name:         "data_export allows export_review and export_execute",
			workflowType: WorkflowTemplateTypeDataExport,
			expected:     []string{WorkflowStepTypeExportReview, WorkflowStepTypeExportExecute},
		},
		{
			name:         "unknown type returns nil",
			workflowType: "unknown",
			expected:     nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := WorkflowStepTypesByWorkflowType[tc.workflowType]
			assert.Equal(t, tc.expected, result)
		})
	}
}
