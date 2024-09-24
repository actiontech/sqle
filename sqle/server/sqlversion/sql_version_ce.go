//go:build !enterprise
// +build !enterprise

package sqlversion

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
)

func AssociateWorkflowToTheFirstStageOfSQLVersion(projectUID, workflowID string, versionID uint) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql version is enterprise version feature"))
}
