package common

import (
	"github.com/actiontech/sqle/sqle/dms"
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
	// 走 ParseInstanceDBType 兜底：dms.convertInstance 入口已经做过同名规范化，
	// 但 HTTP 直入口（/v1/.../check_instance_is_connectable 等）携带的原始 DBType
	// 字面值未经过 dms.convertInstance，仍可能是 "GaussDB / openGauss" 等
	// DMS 枚举形态。本处再做一次相同的规范化，保证 OpenPlugin 拿到的永远是
	// SQLE 后端契约值（如 "GaussDB"）。详见 sqle-ee #2877 compat_risks R5。
	plugin, err := driver.GetPluginManager().OpenPlugin(l, dms.ParseInstanceDBType(inst.DbType), cfg)
	if err != nil {
		return nil, errors.Wrap(err, "open plugin")
	}

	return plugin, nil
}

func NewDriverManagerWithoutCfg(l *logrus.Entry, dbType string) (driver.Plugin, error) {
	return driver.GetPluginManager().OpenPlugin(l, dms.ParseInstanceDBType(dbType), &v2.Config{})
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
