package v1

import (
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newDriverManagerWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.DriverManager, error) {
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

	return driver.NewDriverManger(l, inst.DbType, cfg)
}

func newDriverManagerWithoutCfg(l *logrus.Entry, dbType string) (driver.DriverManager, error) {
	return driver.NewDriverManger(l, dbType, nil)
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
