package onlineddl

import (
	"github.com/openark/golib/log"
	"github.com/sirupsen/logrus"
)

type logAdaptor struct {
	inner *logrus.Entry
}

func newLogAdaptor(l *logrus.Entry) *logAdaptor {
	return &logAdaptor{
		inner: l,
	}
}

func (l *logAdaptor) Debug(args ...interface{}) {
	l.inner.Debug(args...)
}

func (l *logAdaptor) Debugf(format string, args ...interface{}) {
	l.inner.Debugf(format, args...)
}

func (l *logAdaptor) Info(args ...interface{}) {
	l.inner.Info(args...)
}

func (l *logAdaptor) Infof(format string, args ...interface{}) {
	l.inner.Infof(format, args...)
}

func (l *logAdaptor) Warning(args ...interface{}) error {
	l.inner.Warning(args...)
	return nil
}

func (l *logAdaptor) Warningf(format string, args ...interface{}) error {
	l.inner.Warningf(format, args...)
	return nil
}

func (l *logAdaptor) Error(args ...interface{}) error {
	l.inner.Error(args...)
	return nil
}

func (l *logAdaptor) Errorf(format string, args ...interface{}) error {
	l.inner.Errorf(format, args...)
	return nil
}

func (l *logAdaptor) Errore(err error) error {
	l.inner.Errorln(err)
	return nil
}

func (l *logAdaptor) Fatal(args ...interface{}) error {
	l.inner.Fatal(args...)
	return nil
}

func (l *logAdaptor) Fatalf(format string, args ...interface{}) error {
	l.inner.Fatalf(format, args...)
	return nil
}

func (l *logAdaptor) Fatale(err error) error {
	l.inner.Fatalln(err)
	return nil
}

func (l *logAdaptor) SetLevel(level log.LogLevel) {
	switch level {
	case log.DEBUG:
		l.inner.Logger.SetLevel(logrus.DebugLevel)
	case log.INFO, log.NOTICE:
		l.inner.Logger.SetLevel(logrus.InfoLevel)
	case log.WARNING:
		l.inner.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR, log.CRITICAL:
		l.inner.Logger.SetLevel(logrus.ErrorLevel)
	case log.FATAL:
		l.inner.Logger.SetLevel(logrus.FatalLevel)
	}
}

func (l *logAdaptor) SetPrintStackTrace(printStackTraceFlag bool) {
}
