package dms

import (
	"context"
	"fmt"

	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

type checkDataExportWorkflowTemplateUsedReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		IsUsed bool  `json:"is_used"`
		Count  int64 `json:"count"`
	} `json:"data"`
}

// CheckDataExportWorkflowTemplateUsed asks DMS whether unfinished data-export
// workflows still reference the given SQLE template.
func CheckDataExportWorkflowTemplateUsed(ctx context.Context, dmsAddr, projectUID string, templateID uint) (bool, int64, error) {
	if dmsAddr == "" {
		return false, 0, nil
	}
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}
	reply := &checkDataExportWorkflowTemplateUsedReply{}
	url := fmt.Sprintf(
		"%s/v1/dms/projects/%s/data_export_workflows/workflow_template_used?workflow_template_id=%d",
		dmsAddr, projectUID, templateID,
	)
	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return false, 0, fmt.Errorf("failed to check data export workflow template usage: %v", err)
	}
	if reply.Code != 0 {
		return false, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}
	return reply.Data.IsUsed, reply.Data.Count, nil
}
