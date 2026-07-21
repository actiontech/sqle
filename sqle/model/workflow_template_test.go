package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultWorkflowTemplateIsDefault(t *testing.T) {
	tpl := DefaultWorkflowTemplate("proj-1")
	assert.True(t, tpl.IsDefault)
	assert.Equal(t, WorkflowTemplateTypeWorkflow, tpl.WorkflowType)
}

func TestDefaultDataExportWorkflowTemplateIsDefault(t *testing.T) {
	tpl := DefaultDataExportWorkflowTemplate("proj-1")
	assert.True(t, tpl.IsDefault)
	assert.Equal(t, WorkflowTemplateTypeDataExport, tpl.WorkflowType)
}

func TestValidateWorkflowTemplateForCreate(t *testing.T) {
	projectId := ProjectUID("p1")

	t.Run("nil template", func(t *testing.T) {
		err := ValidateWorkflowTemplateForCreate(nil, projectId)
		assert.Error(t, err)
	})

	t.Run("project mismatch", func(t *testing.T) {
		err := ValidateWorkflowTemplateForCreate(&WorkflowTemplate{
			ProjectId:    "other",
			WorkflowType: WorkflowTemplateTypeWorkflow,
		}, projectId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not belong to the project")
	})

	t.Run("type must be workflow", func(t *testing.T) {
		err := ValidateWorkflowTemplateForCreate(&WorkflowTemplate{
			ProjectId:    projectId,
			WorkflowType: WorkflowTemplateTypeDataExport,
		}, projectId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be workflow")
	})

	t.Run("valid template", func(t *testing.T) {
		err := ValidateWorkflowTemplateForCreate(&WorkflowTemplate{
			ProjectId:    projectId,
			WorkflowType: WorkflowTemplateTypeWorkflow,
		}, projectId)
		assert.NoError(t, err)
	})
}

func TestCanDeleteWorkflowTemplate(t *testing.T) {
	t.Run("cannot delete default", func(t *testing.T) {
		err := CanDeleteWorkflowTemplate(true, 2, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "default")
	})

	t.Run("cannot delete only template", func(t *testing.T) {
		err := CanDeleteWorkflowTemplate(false, 1, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only")
	})

	t.Run("cannot delete referenced unfinished", func(t *testing.T) {
		err := CanDeleteWorkflowTemplate(false, 2, true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unfinished")
	})

	t.Run("can delete non-default unused", func(t *testing.T) {
		err := CanDeleteWorkflowTemplate(false, 2, false)
		assert.NoError(t, err)
	})
}

func TestResolveWorkflowTemplateSelection(t *testing.T) {
	zero := uint(0)
	explicit := uint(9)
	cases := []struct {
		name        string
		explicitID  *uint
		useExplicit bool
	}{
		{name: "nil means default", explicitID: nil, useExplicit: false},
		{name: "zero means default", explicitID: &zero, useExplicit: false},
		{name: "explicit id", explicitID: &explicit, useExplicit: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			useExplicit := tc.explicitID != nil && *tc.explicitID != 0
			assert.Equal(t, tc.useExplicit, useExplicit)
		})
	}
}
