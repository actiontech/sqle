//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sqlQuery "github.com/actiontech/sqle/sqle/server/sql_query"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/sirupsen/logrus"

	"github.com/actiontech/sqle/sqle/driver"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var errSqlQueryUserNoAccessToSql = errors.New(errors.DataNotExist, fmt.Errorf("current user has no access to this sql"))
var errSqlQueryNoSql = errors.New(errors.DataNotExist, fmt.Errorf("there is no sql"))

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
	drvMgr, err := newDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer drvMgr.Close(context.TODO())

	d, err := drvMgr.GetAuditDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err := d.Ping(context.TODO()); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nodes, err := d.Parse(context.TODO(), req.SQL)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if len(nodes) == 0 {
		return controller.JSONBaseErrorReq(c, errSqlQueryNoSql)
	}

	// audit
	if instance.SqlQueryConfig.AuditEnabled {
		singleSqls := make([]string, len(nodes))
		for i, node := range nodes {
			singleSqls[i] = node.Text
		}
		err = sqlQuery.Audit(singleSqls, req.InstanceSchema, instance)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	rawSQL := &model.SqlQueryHistory{
		CreateUserId: user.ID,
		InstanceId:   instance.ID,
		Schema:       req.InstanceSchema,
		RawSql:       req.SQL,
	}

	queryDriver, err := drvMgr.GetSQLQueryDriver()
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
			return controller.JSONBaseErrorReq(c, fmt.Errorf("the SQL[%s] is invalid: %v", node.Text, validateResult.Error))
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
	queryIdStr := c.Param("query_id")
	queryId, err := strconv.Atoi(queryIdStr)
	if err != nil {
		return err
	}
	req := new(GetSQLResultReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	// get data from db
	s := model.GetStorage()
	singleSql, err := s.GetSqlQueryExecSqlByQueryId(queryId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	history, err := s.GetSqlQueryHistoryById(singleSql.SqlQueryHistoryId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceId := strconv.FormatUint(uint64(history.InstanceId), 10)
	instance, exist, err := s.GetInstanceById(instanceId)
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
		if user.ID != history.CreateUserId {
			return controller.JSONBaseErrorReq(c, errSqlQueryUserNoAccessToSql)
		}
	}

	l := log.NewEntry().WithFields(logrus.Fields{
		"user":          user.Name,
		"host":          c.Request().Host,
		"time":          time.Now(),
		"instance_name": instance.Name,
		"instance_addr": fmt.Sprintf("%v:%v", instance.Host, instance.Port),
		"schema":        history.Schema,
		"raw_sql":       singleSql.Sql,
	})
	l.Infoln("SQL Query begin")

	// rewrite sql
	dsn, err := newDSN(instance, history.Schema)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	drvMgr, err := driver.NewDriverManger(log.NewEntry(), instance.DbType, &driver.Config{DSN: dsn})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer drvMgr.Close(context.TODO())

	queryDriver, err := drvMgr.GetSQLQueryDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	limit := uint32(instance.SqlQueryConfig.MaxPreQueryRows)
	if limit > req.PageSize {
		limit = req.PageSize
	}
	offset := limit * (req.PageIndex - 1)
	rewriteRes, err := queryDriver.QueryPrepare(context.TODO(), singleSql.Sql, &driver.QueryPrepareConf{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if rewriteRes.ErrorType != driver.ErrorTypeNotError {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("the SQL[%s] is invalid: %w", singleSql.Sql, err))
	}

	startTime := time.Now()
	singleSql.ExecStartAt = &startTime

	// execute sql
	queryRes, err := queryDriver.Query(context.TODO(), rewriteRes.NewSQL, &driver.QueryConf{TimeOutSecond: uint32(instance.SqlQueryConfig.QueryTimeoutSecond)})
	if err != nil {
		// update sql_query_execution_sqls table
		singleSql.ExecResult = err.Error()
		endAt := time.Now()
		singleSql.ExecEndAt = &endAt

		l.WithFields(logrus.Fields{
			"exec_start_time": startTime,
			"exec_sql":        rewriteRes.NewSQL,
			"elapsed_time":    endAt.Sub(startTime) / time.Millisecond,
		}).Errorln("SQL Query error")

		if err := s.Save(singleSql); err != nil {
			log.Logger().Errorf("update result to sql_query_execution_sqls failed: %v", err)
		}
		return controller.JSONBaseErrorReq(c, err)
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime) / time.Millisecond // ms

	l.WithFields(logrus.Fields{
		"exec_start_time":   startTime,
		"exec_sql":          rewriteRes.NewSQL,
		"result_rows_count": len(queryRes.Rows),
		"elapsed_time":      elapsedTime,
	}).Infoln("SQL Query end")

	// update sql_query_execution_sqls table
	singleSql.ExecEndAt = &endTime
	singleSql.ExecResult = "OK"
	if err := s.Save(singleSql); err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("update result to sql_query_execution_sqls failed: %v", err))
	}

	// build response
	data, err := buildSqlQueryRes(rewriteRes.NewSQL, int(req.PageIndex), int(limit), int(elapsedTime), queryRes)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSQLResultResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}

func buildSqlQueryRes(sql string, pageIndex, pageSize, elapsedTime int, queryRes *driver.QueryResult) (GetSQLResultResDataV1, error) {
	rows := make([]map[string]string, len(queryRes.Rows))
	for i, r := range queryRes.Rows {
		row := make(map[string]string)
		for j, v := range r.Values {
			column := queryRes.Column[j]
			row[column.Value] = v.Value
		}
		rows[i] = row
	}

	head := make([]SQLResultItemHeadResV1, len(queryRes.Column))
	for i, c := range queryRes.Column {
		head[i].FieldName = c.Value
	}

	// calculate start-line
	rowsCount := len(rows)
	var startLine, endLine int
	if rowsCount == 0 {
		startLine = 0
		endLine = 0
	} else if rowsCount == pageSize {
		startLine = pageSize*(pageIndex-1) + 1
		endLine = pageSize * pageIndex
	} else if rowsCount < pageSize {
		startLine = pageSize*(pageIndex-1) + 1
		endLine = pageSize*(pageIndex-1) + rowsCount
	}

	res := GetSQLResultResDataV1{
		SQL:         sql,
		StartLine:   startLine,
		EndLine:     endLine,
		CurrentPage: pageIndex,
		ExecuteTime: elapsedTime,
		Rows:        rows,
		Head:        head,
	}
	return res, nil
}

func getSQLQueryHistory(c echo.Context) error {
	instanceName := c.Param("instance_name")
	req := new(GetSQLQueryHistoryReqV1)
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

	sqlHistories, err := s.GetSqlQueryRawSqlByUserId(user.ID, instance.ID, req.PageIndex, req.PageSize, req.FilterFuzzySearch)
	items := make([]SQLHistoryItemResV1, len(sqlHistories))
	for i, h := range sqlHistories {
		items[i] = SQLHistoryItemResV1{
			SQL: h.RawSql,
		}
	}

	return c.JSON(http.StatusOK, &GetSQLQueryHistoryResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetSQLQueryHistoryResDataV1{
			SQLHistories: items,
		},
	})
}

func getSQLExplain(c echo.Context) error {
	instanceName := c.Param("instance_name")
	req := new(GetSqlExplainReqV1)
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

	drvMgr, err := newDriverManagerWithoutAudit(log.NewEntry(), instance, req.InstanceSchema)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer drvMgr.Close(context.TODO())

	driverWithoutAudit, err := drvMgr.GetAuditDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := driverWithoutAudit.Ping(context.TODO()); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nodes, err := driverWithoutAudit.Parse(context.TODO(), req.Sql)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if len(nodes) == 0 {
		return controller.JSONBaseErrorReq(c, errSqlQueryNoSql)
	}

	var sqlExplains []SQLQuerySQLExplain
	analysisDriver, err := drvMgr.GetAnalysisDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for _, node := range nodes {
		explainResult, err := analysisDriver.Explain(context.TODO(), &driver.ExplainConf{Sql: node.Text})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		classicResult := convertExplainResultsToRes(explainResult)
		sqlExplains = append(sqlExplains, SQLQuerySQLExplain{SQL: node.Text, ClassicResult: classicResult})
	}

	return c.JSON(http.StatusOK, &GetSQLExplainResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    sqlExplains,
	})
}

func convertExplainResultsToRes(explainResult *driver.ExplainResult) ExplainClassicResult {
	explainClassicResult := ExplainClassicResult{
		Rows: make([]map[string]string, len(explainResult.ClassicResult.Rows)),
		Head: make([]TableMetaItemHeadResV1, len(explainResult.ClassicResult.Columns)),
	}

	for i, column := range explainResult.ClassicResult.Columns {
		explainClassicResult.Head[i].FieldName = column.Name
		explainClassicResult.Head[i].Desc = column.Desc
	}

	for k, rows := range explainResult.ClassicResult.Rows {
		explainClassicResult.Rows[k] = make(map[string]string)
		for i, row := range rows {
			columnName := explainResult.ClassicResult.Columns[i].Name
			explainClassicResult.Rows[k][columnName] = row
		}
	}

	return explainClassicResult
}
