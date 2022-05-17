//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"net/http"
	"strconv"

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

	// todo: check DQL using driver

	// save db
	rawSQL := &model.SqlQueryHistory{
		CreateUserId: user.ID,
		InstanceId:   instance.ID,
		Database:     req.InstanceScheme,
		RawSql:       req.SQL,
	}
	for _, node := range nodes {
		rawSQL.ExecSQLs = append(rawSQL.ExecSQLs, &model.SqlQueryExecutionSql{
			Sql:         node.Text,
			ExecStartAt: nil,
			ExecEndAt:   nil,
			ExecResult:  "",
		})
	}
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
