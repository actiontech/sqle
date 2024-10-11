//go:build !enterprise
// +build !enterprise

package sqlversion

import (
	"context"
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
)

func CheckInstanceInWorkflowCanAssociateToTheFirstStageOfVersion(versionID uint, instanceId []uint64) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql version is enterprise version feature"))
}

func CheckWorkflowExecutable(ctx context.Context, projectUid, workflowId string) (executable bool, reason string, err error) {
	return true, "", nil
}
