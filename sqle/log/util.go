package log

import "github.com/sirupsen/logrus"

//AutoPrintError  will automatically print the error log according to whether the execution reports an error
//printing effect: '${msg} | error: %e'
func AutoPrintError(logger *logrus.Entry, f func() error, msg string) {
	if err := f(); err != nil {
		logger.Error(msg, "| error:", err)
	}
}
