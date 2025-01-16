package log

import (
	"fmt"
	"context"
	"github.com/sirupsen/logrus"
	rotate "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	gormLog "gorm.io/gorm/logger"
)

var std *logrus.Logger

func Logger() *logrus.Logger {
	return std
}

func NewEntry() *logrus.Entry {
	return std.WithFields(logrus.Fields{
		"thread_id": genRandomThreadId(),
	})
}

func init() {
	std = logrus.New()
}

func InitLogger(filePath string, maxSize, maxBackupNum int) {
	std.SetOutput(NewRotateFile(filePath, "/sqled.log", maxSize /*MB*/, maxBackupNum))
}

func ExitLogger() {
	w := std.Out
	std.SetOutput(os.Stderr)
	if wc, ok := w.(io.Closer); ok {
		wc.Close()
	}
}

func NewRotateFile(filePath, fileName string, maxSize, maxBackupNum int) *rotate.Logger {
	return &rotate.Logger{
		Filename:   strings.TrimRight(filePath, "/") + fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackupNum,
	}
}

func genRandomThreadId() string {
	seq := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l := len(seq)
	a := rand.Intn(l * l * l)
	return fmt.Sprintf("%c%c%c", seq[a%l], seq[(a/l)%l], seq[(a/l/l)%l])
}


type gormLogWrapper struct {
	logger   *logrus.Entry
	logLevel gormLog.LogLevel
}

func NewGormLogWrapper(level gormLog.LogLevel) *gormLogWrapper {
	h := &gormLogWrapper{
		logger:   Logger().WithField("type", "sql"),
		logLevel: level,
	}
	return h
}

func (h *gormLogWrapper) LogMode(level gormLog.LogLevel) gormLog.Interface {
	h.logLevel = level
	return h
}

func (h *gormLogWrapper) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if h.logLevel != gormLog.Info {
		return
	}
	sql, rowsAffected := fc()
	h.logger.Debugf(fmt.Sprintf("trace: sql: %v; rowsAffected: %v; err: %v", sql, rowsAffected, err))
}

func (h *gormLogWrapper) Error(ctx context.Context, format string, a ...interface{}) {
	if h.logLevel >= gormLog.Error {
		h.logger.Error(fmt.Sprintf(format, a...))
	}
}

func (h *gormLogWrapper) Warn(ctx context.Context, format string, a ...interface{}) {
	if h.logLevel >= gormLog.Warn {
		h.logger.Warn(fmt.Sprintf(format, a...))
	}
}

func (h *gormLogWrapper) Info(ctx context.Context, format string, a ...interface{}) {
	if h.logLevel >= gormLog.Info {
		h.logger.Info(fmt.Sprintf(format, a...))
	}
}
