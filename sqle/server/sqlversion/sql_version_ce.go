//go:build !enterprise
// +build !enterprise

package sqlversion

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
)

func CheckInstanceInWorkflowCanAssociateToTheFirstStageOfVersion(versionID uint, instanceId []uint64) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql version is enterprise version feature"))
}
