package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newDriverWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.Driver, error) {
	if inst == nil {
		return nil, errors.Errorf("instance is nil")
	}

	dsn := &driver.DSN{
		Host:     inst.Host,
		Port:     inst.Port,
		User:     inst.User,
		Password: inst.Password,

		DatabaseName: database,
	}

	cfg, err := driver.NewConfig(dsn, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new driver without audit")
	}

	return driver.NewDriver(l, inst.DbType, cfg)
}

func JSONNewNotImplementedErr(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.New("not implemented yet"))
}
