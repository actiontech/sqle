package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	rotate "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"math/rand"
	"os"
	"strings"
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

func InitLogger(filePath string) {
	std.SetOutput(NewRotateFile(filePath, "/sqled.log", 1024 /*1GB*/))
}

func ExitLogger() {
	w := std.Out
	std.SetOutput(os.Stderr)
	if wc, ok := w.(io.Closer); ok {
		wc.Close()
	}
}

func NewRotateFile(filePath, fileName string, maxSize int) *rotate.Logger {
	return &rotate.Logger{
		Filename: strings.TrimRight(filePath, "/") + fileName,
		MaxSize:  maxSize,
	}
}

func genRandomThreadId() string {
	seq := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l := len(seq)
	a := rand.Intn(l * l * l)
	return fmt.Sprintf("%c%c%c", seq[a%l], seq[(a/l)%l], seq[(a/l/l)%l])
}