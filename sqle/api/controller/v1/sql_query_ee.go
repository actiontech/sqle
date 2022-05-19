//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/sirupsen/logrus"

	"github.com/actiontech/sqle/sqle/driver"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func prepareSQLQuery(c echo.Context) error {
	instanceName := c.Param("instance_name")
	req := new(PrepareSQLQueryReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}

	// check user auth
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if user.Name != model.DefaultAdminUser {
		exist, err = s.CheckUserHasOpToInstance(user, instance, []uint{model.OP_SQL_QUERY_QUERY})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
		}
	}

	// parse sql using driver
	d, err := newDriverWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer d.Close(context.TODO())
	if err := d.Ping(context.TODO()); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nodes, err := d.Parse(context.TODO(), req.SQL)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	rawSQL := &model.SqlQueryHistory{
		CreateUserId: user.ID,
		InstanceId:   instance.ID,
		Database:     req.InstanceScheme,
		RawSql:       req.SQL,
	}
	queryDriver, err := driver.NewSQLQueryDriver(log.NewEntry(), instance.DbType, &driver.DSN{
		Host:         instance.Host,
		Port:         instance.Port,
		User:         instance.User,
		Password:     instance.Password,
		DatabaseName: req.InstanceSchema,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for _, node := range nodes {
		// validate SQL
		validateResult, err := queryDriver.QueryPrepare(context.TODO(), node.Text, &driver.QueryPrepareConf{
			// these two parameters are used to rewrite sql, but the rewrite result is useless here
			// so fill them with smaller numeric values
			Limit:  1,
			Offset: 1,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if validateResult.ErrorType == driver.ErrorTypeNotQuery {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("the SQL[%s] is invalid: %w", node.Text, err))
		}

		rawSQL.ExecSQLs = append(rawSQL.ExecSQLs, &model.SqlQueryExecutionSql{
			Sql:         node.Text,
			ExecStartAt: nil,
			ExecEndAt:   nil,
			ExecResult:  "",
		})
	}

	// save db
	err = s.Save(rawSQL)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	queryIds := make([]PrepareSQLQueryResSQLV1, len(rawSQL.ExecSQLs))
	for i, sql := range rawSQL.ExecSQLs {
		queryIds[i].QueryId = strconv.FormatUint(uint64(sql.ID), 10)
		queryIds[i].SQL = sql.Sql
	}
	return c.JSON(http.StatusOK, &PrepareSQLQueryResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: PrepareSQLQueryResDataV1{
			QueryIds: queryIds,
		},
	})
}

func getSQLResult(c echo.Context) error {
	return nil
}

func getSQLQueryHistory(c echo.Context) error {
	return nil
}
