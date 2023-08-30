//go:build enterprise
// +build enterprise

package scanners

type Logger interface {
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Tracef(format string, v ...interface{})
}
