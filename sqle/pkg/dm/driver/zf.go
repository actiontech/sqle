/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"context"
	"database/sql/driver"
	"dm/util"
	"reflect"
	"strconv"
	"time"
)

type logFilter struct{}

func (filter *logFilter) DmDriverOpen(filterChain *filterChain, d *DmDriver, dsn string) (ret *DmConnection, err error) {
	var logRecord = d.logInfo.logRecord
	logRecord.Set(d, "open", dsn)
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmDriverOpen(d, dsn)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmDriverOpenConnector(filterChain *filterChain, d *DmDriver, dsn string) (ret *DmConnector, err error) {
	var logRecord = d.logInfo.logRecord
	logRecord.Set(d, "openConnector", dsn)
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmDriverOpenConnector(d, dsn)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectorConnect(filterChain *filterChain, c *DmConnector, ctx context.Context) (ret *DmConnection, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "connect")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmConnectorConnect(c, ctx)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectorDriver(filterChain *filterChain, c *DmConnector) (ret *DmDriver) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "driver")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmConnectorDriver(c)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionBegin(filterChain *filterChain, c *DmConnection) (ret *DmConnection, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "begin")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmConnectionBegin(c)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionBeginTx(filterChain *filterChain, c *DmConnection, ctx context.Context, opts driver.TxOptions) (ret *DmConnection, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "beginTx", opts)
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmConnectionBeginTx(c, ctx, opts)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionCommit(filterChain *filterChain, c *DmConnection) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "commit")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionCommit(c)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmConnectionRollback(filterChain *filterChain, c *DmConnection) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "rollback")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionRollback(c)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmConnectionClose(filterChain *filterChain, c *DmConnection) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "close")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionClose(c)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmConnectionPing(filterChain *filterChain, c *DmConnection, ctx context.Context) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "ping")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionPing(c, ctx)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmConnectionExec(filterChain *filterChain, c *DmConnection, query string, args []driver.Value) (ret *DmResult, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "executeByStmt", query, args)
	defer func() {
		filter.executeAfter(c.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	filter.executeBefore(c.logInfo)
	ret, err = filterChain.DmConnectionExec(c, query, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionExecContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string, args []driver.NamedValue) (ret *DmResult, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "executeCtx", query, args)
	defer func() {
		filter.executeAfter(c.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	filter.executeBefore(c.logInfo)
	ret, err = filterChain.DmConnectionExecContext(c, ctx, query, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionQuery(filterChain *filterChain, c *DmConnection, query string, args []driver.Value) (ret *DmRows, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "query", query, args)
	defer func() {
		filter.executeAfter(c.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	filter.executeBefore(c.logInfo)
	ret, err = filterChain.DmConnectionQuery(c, query, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionQueryContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string, args []driver.NamedValue) (ret *DmRows, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "queryCtx", query, args)
	defer func() {
		filter.executeAfter(c.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	filter.executeBefore(c.logInfo)
	ret, err = filterChain.DmConnectionQueryContext(c, ctx, query, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionPrepare(filterChain *filterChain, c *DmConnection, query string) (ret *DmStatement, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "prepareStatement", query)
	defer func() {
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	ret, err = filterChain.DmConnectionPrepare(c, query)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionPrepareContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string) (ret *DmStatement, err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "prepareStatementCtx", query)
	defer func() {
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(query)
	ret, err = filterChain.DmConnectionPrepareContext(c, ctx, query)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmConnectionResetSession(filterChain *filterChain, c *DmConnection, ctx context.Context) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "resetSession")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionResetSession(c, ctx)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmConnectionCheckNamedValue(filterChain *filterChain, c *DmConnection, nv *driver.NamedValue) (err error) {
	var logRecord = c.logInfo.logRecord
	logRecord.Set(c, "checkNamedValue", nv)
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmConnectionCheckNamedValue(c, nv)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmStatementClose(filterChain *filterChain, s *DmStatement) (err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "close")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmStatementClose(s)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmStatementNumInput(filterChain *filterChain, s *DmStatement) (ret int) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "numInput")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmStatementNumInput(s)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmStatementExec(filterChain *filterChain, s *DmStatement, args []driver.Value) (ret *DmResult, err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "executeByStmt", args)
	defer func() {
		filter.executeAfter(s.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(s.nativeSql)
	filter.executeBefore(s.logInfo)
	ret, err = filterChain.DmStatementExec(s, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmStatementExecContext(filterChain *filterChain, s *DmStatement, ctx context.Context, args []driver.NamedValue) (ret *DmResult, err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "executeCtx", args)
	defer func() {
		filter.executeAfter(s.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(s.nativeSql)
	filter.executeBefore(s.logInfo)
	ret, err = filterChain.DmStatementExecContext(s, ctx, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmStatementQuery(filterChain *filterChain, s *DmStatement, args []driver.Value) (ret *DmRows, err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "query", args)
	defer func() {
		filter.executeAfter(s.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(s.nativeSql)
	filter.executeBefore(s.logInfo)
	ret, err = filterChain.DmStatementQuery(s, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmStatementQueryContext(filterChain *filterChain, s *DmStatement, ctx context.Context, args []driver.NamedValue) (ret *DmRows, err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "queryCtx", args)
	defer func() {
		filter.executeAfter(s.logInfo, logRecord)
		filter.doLog(logRecord)
	}()
	logRecord.SetSql(s.nativeSql)
	filter.executeBefore(s.logInfo)
	ret, err = filterChain.DmStatementQueryContext(s, ctx, args)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmStatementCheckNamedValue(filterChain *filterChain, s *DmStatement, nv *driver.NamedValue) (err error) {
	var logRecord = s.logInfo.logRecord
	logRecord.Set(s, "checkNamedValue", nv)
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmStatementCheckNamedValue(s, nv)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmResultLastInsertId(filterChain *filterChain, r *DmResult) (ret int64, err error) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "lastInsertId")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmResultLastInsertId(r)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmResultRowsAffected(filterChain *filterChain, r *DmResult) (ret int64, err error) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "rowsAffected")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret, err = filterChain.DmResultRowsAffected(r)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmRowsColumns(filterChain *filterChain, r *DmRows) (ret []string) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columns")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmRowsColumns(r)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmRowsClose(filterChain *filterChain, r *DmRows) (err error) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "close")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmRowsClose(r)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmRowsNext(filterChain *filterChain, r *DmRows, dest []driver.Value) (err error) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "next", dest)
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmRowsNext(r, dest)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmRowsHasNextResultSet(filterChain *filterChain, r *DmRows) (ret bool) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "hasNextResultSet")
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmRowsHasNextResultSet(r)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmRowsNextResultSet(filterChain *filterChain, r *DmRows) (err error) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "nextResultSet")
	defer func() {
		filter.doLog(logRecord)
	}()
	err = filterChain.DmRowsNextResultSet(r)
	if err != nil {
		logRecord.SetError(err)
		return
	}
	return
}

func (filter *logFilter) DmRowsColumnTypeScanType(filterChain *filterChain, r *DmRows, index int) (ret reflect.Type) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columnTypeScanType", index)
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmRowsColumnTypeScanType(r, index)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmRowsColumnTypeDatabaseTypeName(filterChain *filterChain, r *DmRows, index int) (ret string) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columnTypeDatabaseTypeName", index)
	defer func() {
		filter.doLog(logRecord)
	}()
	ret = filterChain.DmRowsColumnTypeDatabaseTypeName(r, index)
	logRecord.SetReturnValue(ret)
	return
}

func (filter *logFilter) DmRowsColumnTypeLength(filterChain *filterChain, r *DmRows, index int) (length int64, ok bool) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columnTypeLength", index)
	defer func() {
		filter.doLog(logRecord)
	}()
	length, ok = filterChain.DmRowsColumnTypeLength(r, index)
	if ok {
		logRecord.SetReturnValue(length)
	} else {
		logRecord.SetReturnValue(-1)
	}
	return
}

func (filter *logFilter) DmRowsColumnTypeNullable(filterChain *filterChain, r *DmRows, index int) (nullable, ok bool) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columnTypeNullable", index)
	defer func() {
		filter.doLog(logRecord)
	}()
	nullable, ok = filterChain.DmRowsColumnTypeNullable(r, index)
	if ok {
		logRecord.SetReturnValue(nullable)
	} else {
		logRecord.SetReturnValue(false)
	}
	return
}

func (filter *logFilter) DmRowsColumnTypePrecisionScale(filterChain *filterChain, r *DmRows, index int) (precision, scale int64, ok bool) {
	var logRecord = r.logInfo.logRecord
	logRecord.Set(r, "columnTypePrecisionScale", index)
	defer func() {
		filter.doLog(logRecord)
	}()
	precision, scale, ok = filterChain.DmRowsColumnTypePrecisionScale(r, index)
	if ok {
		logRecord.SetReturnValue(strconv.FormatInt(precision, 10) + "&" + strconv.FormatInt(scale, 10))
	} else {
		logRecord.SetReturnValue("-1&-1")
	}
	return
}

func (filter *logFilter) executeBefore(logInfo *logInfo) {
	if LogFilterLogger.IsSqlEnabled() {
		logInfo.lastExecuteStartNano = time.Now()
	}
}

func (filter *logFilter) executeAfter(logInfo *logInfo, record *LogRecord) {
	if LogFilterLogger.IsSqlEnabled() {
		record.SetUsedTime(time.Since(logInfo.lastExecuteStartNano))
	}
}

func (filter *logFilter) doLog(record *LogRecord) {

	if record == nil {
		return
	}
	if record.GetError() != nil {
		LogFilterLogger.ErrorWithErr(record.ToString(), record.GetError())
	} else if record.GetSql() != "" && LogFilterLogger.IsSqlEnabled() {
		LogFilterLogger.Sql(record.ToString())
	} else {
		LogFilterLogger.Info(record.ToString())
	}
}

/************************************************************************************************************/
type Logger struct {
}

var LogFilterLogger = &Logger{}
var ConnLogger = &Logger{}
var AccessLogger = &Logger{}

func (logger Logger) IsDebugEnabled() bool {
	return LogLevel >= LOG_DEBUG
}
func (logger Logger) IsErrorEnabled() bool {
	return LogLevel >= LOG_ERROR
}
func (logger Logger) IsInfoEnabled() bool {
	return LogLevel >= LOG_INFO
}
func (logger Logger) IsWarnEnabled() bool {
	return LogLevel >= LOG_WARN
}
func (logger Logger) IsSqlEnabled() bool {
	return LogLevel >= LOG_SQL
}
func (logger Logger) Debug(msg string) {
	if logger.IsDebugEnabled() {
		logger.println(logger.formatHead("DEBUG") + msg)
	}
}
func (logger Logger) DebugWithErr(msg string, err error) {
	if logger.IsDebugEnabled() {
		if e, ok := err.(*DmError); ok {
			logger.println(logger.formatHead("DEBUG") + msg + util.LINE_SEPARATOR + e.FormatStack())
		} else {
			logger.println(logger.formatHead("DEBUG") + msg + util.LINE_SEPARATOR + err.Error())
		}
	}
}
func (logger Logger) Info(msg string) {
	if logger.IsInfoEnabled() {
		logger.println(logger.formatHead("INFO ") + msg)
	}
}
func (logger Logger) Sql(msg string) {
	if logger.IsSqlEnabled() {
		logger.println(logger.formatHead("SQL  ") + msg)
	}
}
func (logger Logger) Warn(msg string) {
	if logger.IsWarnEnabled() {
		logger.println(logger.formatHead("WARN ") + msg)
	}
}
func (logger Logger) ErrorWithErr(msg string, err error) {
	if e, ok := err.(*DmError); ok {
		logger.println(logger.formatHead("ERROR") + msg + util.LINE_SEPARATOR + e.FormatStack())
	} else {
		logger.println(logger.formatHead("ERROR") + msg + util.LINE_SEPARATOR + err.Error())
	}
}

//TODO: 获取goroutine objId
func (logger Logger) formatHead(head string) string {
	// return "[" + head + " - " + StringUtil.formatTime() + "] tid:" + Thread.currentThread().getId();
	return "[" + head + " - " + util.StringUtil.FormatTime() + "]"
}
func (logger Logger) println(msg string) {
	goMap["log"].(*logWriter).WriteLine(msg)
}

/*************************************************************************************************/
func formatSource(source interface{}) string {
	if source == nil {
		return ""
	}
	var str string
	switch src := source.(type) {
	case string:
		str += src
	case *DmDriver:
		str += formatDriver(src)
	case *DmConnector:
		str += formatContor(src)
	case *DmConnection:
		str += formatConn(src)
	case *DmStatement:
		str += formatConn(src.dmConn) + ", "
		str += formatStmt(src)
	case *DmResult:
		str += formatConn(src.dmStmt.dmConn) + ", "
		str += formatStmt(src.dmStmt) + ", "
		str += formatRs(src)
	case *DmRows:
		str += formatConn(src.CurrentRows.dmStmt.dmConn) + ", "
		str += formatStmt(src.CurrentRows.dmStmt) + ", "
		str += formatRows(src)
	default:
		str += reflect.TypeOf(src).String() + "@" + reflect.ValueOf(src).Addr().String()
	}
	return str
}

func formatDriver(driver *DmDriver) string {
	if driver != nil && driver.logInfo != nil {
		return "driver-" + strconv.FormatInt(driver.getID(), 10)
	}
	return "driver-nil"
}

func formatContor(contor *DmConnector) string {
	if contor != nil && contor.logInfo != nil {
		return "contor-" + strconv.FormatInt(contor.getID(), 10)
	}
	return "contor-nil"
}

func formatConn(conn *DmConnection) string {
	if conn != nil && conn.logInfo != nil {
		return "conn-" + strconv.FormatInt(conn.getID(), 10)
	}
	return "conn-nil"
}

func formatStmt(stmt *DmStatement) string {
	if stmt != nil && stmt.logInfo != nil {
		return "stmt-" + strconv.FormatInt(stmt.getID(), 10)
	}
	return "stmt-nil"
}

func formatRs(result *DmResult) string {
	if result != nil && result.logInfo != nil {
		return "rs-" + strconv.FormatInt(result.getID(), 10)
	}
	return "rs-nil"
}

func formatRows(rows *DmRows) string {
	if rows != nil && rows.logInfo != nil {
		return "rows-" + strconv.FormatInt(rows.getID(), 10)
	}
	return "rows-nil"
}

func formatTrace(source string, method string, returnValue interface{}, params ...interface{}) string {
	var str string
	if source != "" {
		str += "{ " + source + " } "
	}
	str += method + "("
	var paramStartIndex = 0
	if params != nil && len(params) > paramStartIndex {
		for i := paramStartIndex; i < len(params); i++ {
			if i != paramStartIndex {
				str += ", "
			}
			if params[i] != nil {
				str += reflect.TypeOf(params[i]).String()
			} else {
				str += "nil"
			}
		}
	}
	str += ")"
	if returnValue != nil {
		str += ": " + formatReturn(returnValue)
	}
	str += ";  "
	if params != nil && len(params) > paramStartIndex {
		str += "[PARAMS]: "
		for i := paramStartIndex; i < len(params); i++ {
			if i != 0 {
				str += ", "
			}
			if s, ok := params[i].(string); ok {
				str += "\"" + s + "\""
			} else {
				str += reflect.ValueOf(params[i]).String()
			}
		}
		str += "; "
	}
	return str
}

func formatReturn(returnObj interface{}) string {
	var str string
	switch o := returnObj.(type) {
	case *DmConnection:
		str = formatConn(o)
		sessID := o.SessionID
		if sessID != -1 {
			str += ", sessionID-0x" + strconv.FormatInt(sessID, 16)
		}
	case *DmStatement:
		str = formatStmt(o)
		str += ", handle-" + strconv.Itoa(int(o.id))
	case *DmResult:
		str = formatRs(o)
	case *DmRows:
		str = formatRows(o)
	case string:
		str = `"` + o + `"`
	case nullData:
		str = "nil"
	default:
		str = "unknown"
	}
	return str
}

func formatUsedTime(duration time.Duration) string {
	return "[USED TIME]: " + duration.String()
}

/************************************************************************************************************/

type nullData struct{}

var null = nullData{}

type LogRecord struct {
	source      string
	method      string
	params      []interface{}
	returnValue interface{}
	e           error
	usedTime    time.Duration
	sql         string
}

func (record *LogRecord) Set(source interface{}, method string, params ...interface{}) {
	record.source = formatSource(source)
	record.method = method
	record.params = params
}

func (record *LogRecord) SetReturnValue(retValue interface{}) {
	if retValue == nil {
		record.returnValue = null
	} else {
		record.returnValue = retValue
	}
}

func (record *LogRecord) GetReturnValue() interface{} {
	return record.returnValue
}

func (record *LogRecord) SetSql(sql string) {
	record.sql = sql
}

func (record *LogRecord) GetSql() string {
	return record.sql
}

func (record *LogRecord) SetUsedTime(usedTime time.Duration) {
	record.usedTime = usedTime
}

func (record *LogRecord) GetUsedTime() time.Duration {
	return record.usedTime
}

func (record *LogRecord) SetError(err error) {
	record.e = err
}

func (record *LogRecord) GetError() error {
	return record.e
}

func (record *LogRecord) ToString() string {
	var str string
	str += formatTrace(record.source, record.method, record.returnValue, record.params...)
	if record.usedTime > 0 {
		str += formatUsedTime(record.usedTime)
	}
	return str
}
