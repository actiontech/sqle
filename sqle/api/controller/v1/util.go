package v1

import (
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newDriverWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.Driver, error) {
	if inst == nil {
		return nil, errors.Errorf("instance is nil")
	}

	dsn, err := newDSN(inst, database)
	if err != nil {
		return nil, errors.Wrap(err, "new dsn")
	}

	cfg, err := driver.NewConfig(dsn, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new driver without audit")
	}

	return driver.NewDriver(l, inst.DbType, cfg)
}

func newDriverWithoutCfg(l *logrus.Entry, dbType string) (driver.Driver, error) {
	return driver.NewDriver(l, dbType, &driver.Config{})
}

func newDSN(instance *model.Instance, database string) (*driver.DSN, error) {
	if instance == nil {
		return nil, errors.Errorf("instance is nil")
	}

	return &driver.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     database,
	}, nil
}
