//go:build !enterprise
// +build !enterprise

package v2

import "github.com/actiontech/sqle/sqle/model"

func createNotifyRecord(notifyType string, curTaskRecord *model.WorkflowInstanceRecord) error {
	// nothing
	return nil
}
