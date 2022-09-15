package common

import (
	"context"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

func CheckInstanceIsConnectable(instance *model.Instance) error {
	drvMgr, err := NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer drvMgr.Close(context.TODO())

	d, err := drvMgr.GetAuditDriver()
	if err != nil {
		return err
	}

	if err := d.Ping(context.TODO()); err != nil {
		return err
	}

	return nil
}
