package errors

import "fmt"

type ErrorCode int

const (
	StatusOK                   ErrorCode = 0
	ConnectStorageError        ErrorCode = 5001
	ConnectRemoteDatabaseError ErrorCode = 5002
	ReadUploadFileError        ErrorCode = 5003
	ParseMyBatisXMLFileError   ErrorCode = 5006

	TaskNotExist      ErrorCode = 4006
	TaskActionInvalid ErrorCode = 4009

	TaskRunning    ErrorCode = 1001
	TaskActionDone ErrorCode = 1002

	LoginAuthFail ErrorCode = 4001
	DataExist     ErrorCode = 4010
	DataNotExist  ErrorCode = 4011
	DataConflict  ErrorCode = 4012
	DataInvalid   ErrorCode = 4013

	DriverNotExist ErrorCode = 5001

	FeatureNotImplemented ErrorCode = 7001
)

var (
	ErrSQLTypeConflict = New(-1, fmt.Errorf("不能同时提交 DDL 和 DML 语句"))
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
		return int(StatusOK)
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

func NewNotImplemented(feature string) *CodeError {
	return &CodeError{code: FeatureNotImplemented, err: fmt.Errorf("Not available feature: %v, it is only supported for enterprise edition", feature)}
}
