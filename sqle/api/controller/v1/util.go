package v1

import (
	"fmt"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newDriverWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.Driver, error) {
	if inst == nil {
		return nil, errors.Errorf("instance is nil")
	}

	dsn := &driver.DSN{
		Host:             inst.Host,
		Port:             inst.Port,
		User:             inst.User,
		Password:         inst.Password,
		AdditionalParams: params.NewParamSliceFromMap(inst.GetAdditionalParams()),

		DatabaseName: database,
	}
	fmt.Println("dsn:", inst.GetAdditionalParams())
	cfg, err := driver.NewConfig(dsn, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new driver without audit")
	}

	return driver.NewDriver(l, inst.DbType, cfg)
}

func newDriverWithoutCfg(l *logrus.Entry, dbType string) (driver.Driver, error) {
	return driver.NewDriver(l, dbType, &driver.Config{})
}
