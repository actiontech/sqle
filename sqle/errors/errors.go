package errors

import "fmt"

type ErrorCode int

const (
	STATUS_OK                   ErrorCode = 0
	CONNECT_STORAGE_ERROR       ErrorCode = 5001
	CONNECT_REMOTE_DB_ERROR     ErrorCode = 5002
	READ_UPLOAD_FILE_ERROR      ErrorCode = 5003
	CONNECT_SQLSERVER_RPC_ERROR ErrorCode = 5004
	PARSER_MYCAT_CONFIG_ERROR   ErrorCode = 5005

	INSTANCE_EXIST          ErrorCode = 4001
	RULE_TEMPLATE_EXIST     ErrorCode = 4002
	INSTANCE_NOT_EXIST      ErrorCode = 4003
	RULE_TEMPLATE_NOT_EXIST ErrorCode = 4004
	RULE_TEMPLATE_IS_USED   ErrorCode = 4010
	RULE_NOT_EXIST          ErrorCode = 4005
	TASK_NOT_EXIST          ErrorCode = 4006
	TASK_ACTION_INVALID     ErrorCode = 4009

	DATA_EXIST     ErrorCode = 4010
	DATA_NOT_EXIST ErrorCode = 4011

	TASK_RUNNING     ErrorCode = 1001
	TASK_ACTION_DONE ErrorCode = 1002
)

var (
	SQL_STMT_CONFLICT_ERROR           = New(-1, fmt.Errorf("不能同时提交 DDL 和 DML 语句"))
	SQL_STMT_PROCEUDRE_FUNCTION_ERROR = New(-1, fmt.Errorf("包含存储过程或者函数的任务不能包含其他DDL、DML语句"))
)

type CodeError struct {
	code ErrorCode
	err  error
}

func (e *CodeError) Error() string {
	if e.err == nil {
		return "ok"
	}
	return e.err.Error()
}

func (e *CodeError) Code() int {
	if e.err == nil {
		return int(STATUS_OK)
	}
	return int(e.code)
}

func New(code ErrorCode, err error) error {
	if err == nil {
		return nil
	}
	return &CodeError{
		code: code,
		err:  err,
	}
}
