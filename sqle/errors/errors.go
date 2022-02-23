package errors

import "fmt"

type ErrorCode int

const (
	StatusOK ErrorCode = 0

	TaskRunning    ErrorCode = 1001
	TaskActionDone ErrorCode = 1002

	HttpRequestFormatError ErrorCode = 2001

	ErrAccessDeniedError ErrorCode = 3001

	LoginAuthFail     ErrorCode = 4001
	UserDisabled      ErrorCode = 4005
	TaskNotExist      ErrorCode = 4006
	TaskActionInvalid ErrorCode = 4009
	DataExist         ErrorCode = 4010
	DataNotExist      ErrorCode = 4011
	DataConflict      ErrorCode = 4012
	DataInvalid       ErrorCode = 4013
	DataParseFail     ErrorCode = 4014

	ConnectStorageError        ErrorCode = 5001
	ConnectRemoteDatabaseError ErrorCode = 5002
	ReadUploadFileError        ErrorCode = 5003
	ParseMyBatisXMLFileError   ErrorCode = 5006
	WriteDataToTheFileError    ErrorCode = 5007

	DriverNotExist ErrorCode = 5001
	LoadDriverFail ErrorCode = 5008

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

func NewNotImplementedError(format string, a ...interface{}) error {
	return New(FeatureNotImplemented, fmt.Errorf(format, a...))
}

func NewDataInvalidErr(format string, a ...interface{}) error {
	return New(DataInvalid, fmt.Errorf(format, a...))
}

func NewUserDisabledErr(format string, a ...interface{}) error {
	return New(UserDisabled, fmt.Errorf(format, a...))
}

func NewDataNotExistErr(format string, a ...interface{}) error {
	return New(DataNotExist, fmt.Errorf(format, a...))
}

func HttpRequestFormatErrWrapper(err error) error {
	return New(HttpRequestFormatError, err)
}

func ConnectStorageErrWrapper(err error) error {
	if err == nil {
		return nil
	}
	return New(ConnectStorageError, err)
}

func NewAccessDeniedErr(format string, a ...interface{}) error {
	return New(ErrAccessDeniedError, fmt.Errorf(format, a...))
}
