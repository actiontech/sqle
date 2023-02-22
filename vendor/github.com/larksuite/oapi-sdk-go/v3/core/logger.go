/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkcore

import (
	"context"
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelWarn  LogLevel = 3
	LogLevelError LogLevel = 4
)

type loggerProxy struct {
	LogLevel LogLevel
	Logger   Logger
}

func newLoggerProxy(logLevel LogLevel, logger Logger) *loggerProxy {
	return &loggerProxy{
		LogLevel: logLevel,
		Logger:   logger,
	}
}

func (p *loggerProxy) Debug(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelDebug {
		p.Logger.Debug(ctx, args...)
	}
}

func (p *loggerProxy) Info(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelInfo {
		p.Logger.Info(ctx, args...)
	}
}

func (p *loggerProxy) Warn(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelWarn {
		p.Logger.Warn(ctx, args...)
	}
}

func (p *loggerProxy) Error(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelError {
		p.Logger.Error(ctx, args...)
	}
}

type Logger interface {
	Debug(context.Context, ...interface{})
	Info(context.Context, ...interface{})
	Warn(context.Context, ...interface{})
	Error(context.Context, ...interface{})
}

func NewLogger(config *Config) {
	if config.Logger != nil {
		logLevel := LogLevelInfo
		if config.LogLevel != 0 {
			logLevel = config.LogLevel
		}
		logger := newLoggerProxy(logLevel, config.Logger)
		config.Logger = logger
	} else {
		logLevel := LogLevelInfo
		if config.LogLevel != 0 {
			logLevel = config.LogLevel
		}
		logger := newLoggerProxy(logLevel, defaultLogger{
			logger: log.New(os.Stdout, "", log.LstdFlags),
		})
		config.Logger = logger
	}
}

func NewEventLogger() Logger {
	logger := defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
	return logger
}

type defaultLogger struct {
	logger *log.Logger
}

func (l defaultLogger) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Debug] %v", args)
}

func (l defaultLogger) Info(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Info] %v", args)
}

func (l defaultLogger) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Warn] %v", args)
}

func (l defaultLogger) Error(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Error] %v", args)
}
