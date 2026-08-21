//go:build !enterprise
// +build !enterprise

package workflowexport

import dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"

// opsTypeNameFromListDataExport is a CE stub: CE vendor ListDataExportWorkflow has no OpsType field.
func opsTypeNameFromListDataExport(_ *dmsV1.ListDataExportWorkflow) string {
	return ""
}
