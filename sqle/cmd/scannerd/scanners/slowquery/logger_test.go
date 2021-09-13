package slowquery

import "testing"

var _ Logger = (*testLogger)(nil)

type testLogger struct {
	t testing.T
}

func (tl *testLogger) Warnf(format string, v ...interface{}) {
	tl.t.Helper()
	tl.t.Logf("WARN : "+format, v...)
}

func (tl *testLogger) Infof(format string, v ...interface{}) {
	tl.t.Helper()
	tl.t.Logf("INFO : "+format, v...)
}

func (tl *testLogger) Debugf(format string, v ...interface{}) {
	tl.t.Helper()
	tl.t.Logf("DEBUG: "+format, v...)
}

func (tl *testLogger) Tracef(format string, v ...interface{}) {
	tl.t.Helper()
	tl.t.Logf("TRACE: "+format, v...)
}
