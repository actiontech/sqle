package errors

import "fmt"

type ErrorCode int

const (
	StatusOK ErrorCode = 0

	TaskRunning    ErrorCode = 1001
	TaskActionDone ErrorCode = 1002

	HttpRequestFormatError ErrorCode = 2001

	ErrAccessDeniedError      ErrorCode = 3001
	EnterpriseEditionFeatures ErrorCode = 3002

	LoginAuthFail     ErrorCode = 4001
	UserDisabled      ErrorCode = 4005
	TaskNotExist      ErrorCode = 4006
	TaskActionInvalid ErrorCode = 4009
	DataExist         ErrorCode = 4010
	DataNotExist      ErrorCode = 4011
	DataConflict      ErrorCode = 4012
	DataInvalid       ErrorCode = 4013
	DataParseFail     ErrorCode = 4014
	UserNotPermission ErrorCode = 4015

	ConnectStorageError        ErrorCode = 5001
	ConnectRemoteDatabaseError ErrorCode = 5002
	ReadUploadFileError        ErrorCode = 5003
	ParseMyBatisXMLFileError   ErrorCode = 5006
	WriteDataToTheFileError    ErrorCode = 5007

	DriverNotExist ErrorCode = 5001
	LoadDriverFail ErrorCode = 5008

	FeatureNotImplemented ErrorCode = 7001

	SQLAnalysisSQLIsNotSupported ErrorCode = 8001

	SQLAnalysisCommunityNotSupported = 8002

	CustomRuleEditionNotSupported = 8003

	SQLOptimizationCommunityNotSupported = 8004

	SQLVersionNotAllTasksExecutedSuccess = 8005

	// 需要隐藏所有错误细节或不确定时使用
	GenericError ErrorCode = 9999
)

var (
	ErrSQLTypeConflict = New(-1, fmt.Errorf("cannot submit both DDL and DML statements simultaneously"))
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

func NewUserNotPermissionError(op string) error {
	return New(UserNotPermission, fmt.Errorf("the current user does not have permission for %v to perform this operation", op))
}

func NewAuditPlanNotExistErr() error {
	return New(DataNotExist, fmt.Errorf("audit plan is not exist"))
}

func NewInstanceAuditPlanNotExistErr() error {
	return New(DataNotExist, fmt.Errorf("instance audit plan is not exist"))
}

func NewNotSupportGetAuditPlanAnalysisDataErr() error {
	return New(EnterpriseEditionFeatures, fmt.Errorf("get audit plan analysis data is enterprise version function"))
}

func NewOnlySupportForEnterpriseVersion() error {
	return New(EnterpriseEditionFeatures, fmt.Errorf("this api or function is only supported in enterprise version"))
}

func NewNotSupportGetTaskAnalysisDataErr() error {
	return New(EnterpriseEditionFeatures, fmt.Errorf("get task analysis data is enterprise version function"))
}

func NewTaskNoExistOrNoAccessErr() error {
	return New(DataNotExist, fmt.Errorf("task is not exist or you can't access it"))
}

func NewInstanceNoExistErr() error {
	return New(DataNotExist, fmt.Errorf("instance is not exist"))
}

func NewNotSupportGetSqlFileOrderMethodErr() error {
	return New(EnterpriseEditionFeatures, fmt.Errorf("get sql file order method is enterprise version function"))
}

func NewNotSupportBackupErr() error {
	return New(EnterpriseEditionFeatures, fmt.Errorf("backup and recovery is enterprise version function"))
}

func NewAuditPlanExecuteExtractErr(err error, instanceID string, auditPlan string) error {
	return New(DataInvalid, fmt.Errorf("execute extract sql failed, instance[%s], audit plan[%s], err: %v", instanceID, auditPlan, err))
}
