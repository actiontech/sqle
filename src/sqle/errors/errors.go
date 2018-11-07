package errors

type ErrorCode int

const (
	CONNECT_STORAGE_ERROR   ErrorCode = -1
	INSTANCE_EXIST          ErrorCode = 4001
	RULE_TEMPLATE_NOT_EXIST ErrorCode = 4002
	INSTANCE_NOT_EXIST      ErrorCode = 4003
	RULE_TEMPLATE_EXIST     ErrorCode = 4004
	RULE_NOT_EXIST          ErrorCode = 4005
	TASK_NOT_EXIST          ErrorCode = 4006
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
		return 0
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
