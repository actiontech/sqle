package common

import (
	"github.com/actiontech/sqle/sqle/driver"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewDriverManagerWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.Plugin, error) {
	if inst == nil {
		return nil, errors.Errorf("instance is nil")
	}

	dsn, err := NewDSN(inst, database)
	if err != nil {
		return nil, errors.Wrap(err, "new dsn")
	}

	cfg := &v2.Config{
		DSN: dsn,
	}
	plugin, err := driver.GetPluginManager().OpenPlugin(l, inst.DbType, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "open plugin")
	}

	return plugin, nil
}

func NewDriverManagerWithoutCfg(l *logrus.Entry, dbType string) (driver.Plugin, error) {
	return driver.GetPluginManager().OpenPlugin(l, dbType, &v2.Config{})
}

func NewDSN(instance *model.Instance, database string) (*v2.DSN, error) {
	if instance == nil {
		return nil, errors.Errorf("instance is nil")
	}

	return &v2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     database,
	}, nil
}
