package network

type OracleError struct {
	ErrCode int
	ErrMsg  string
}

func (err *OracleError) Error() string {
	return err.ErrMsg
}

func (err *OracleError) translate() {
	switch err.ErrCode {
	case 12564:
		err.ErrMsg = "ORA-12564: TNS connection refused"
		return
	case 12514:
		err.ErrMsg = "ORA-12514: TNS:listener does nto currently know of service requested in connect descriptor"
		return
	default:
		err.ErrMsg = ""
		return
	}
}
