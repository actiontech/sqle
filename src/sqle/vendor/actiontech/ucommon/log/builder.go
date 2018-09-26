package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// FilterPasswordFunc is the type of filterPassword. The default function used to filter password is Builder's
// filterPassword, user can implement this function and replace the default filter with Builder's SetFilter.
type FilterPasswordFunc func(string) string

type Builder struct {
	lines                   []string
	prefix                  Prefix
	logger                  Logger
	dilutes                 dilutes //protected by logger.Mutex
	diluteKey               string
	diluteDurationSeconds   int
	diluteCheckpointSeconds int

	filterPassword FilterPasswordFunc
}

func Write(f Prefix) *Builder {
	b := Builder{}
	b.lines = make([]string, 4)
	b.prefix = f
	b.logger = instance
	b.dilutes = dilutesInstance
	b.filterPassword = defaultFilterPassword
	return &b
}

func WriteDilute(f Prefix, diluteKey string, diluteDurationSeconds int, diluteCheckpointSeconds int) *Builder {
	b := Write(f)
	b.diluteKey = diluteKey
	b.diluteDurationSeconds = diluteDurationSeconds
	b.diluteCheckpointSeconds = diluteCheckpointSeconds
	return b
}

func (b *Builder) SetFilter(f FilterPasswordFunc) *Builder {
	b.filterPassword = f
	return b
}

// password=, Password=, PASSWORD=
// password"=, Password"=, PASSWORD"=
// password:, Password:, PASSWORD:
// password":, Password":, PASSWORD":
// password(, Password(, PASSWORD(
// password"(, Password"(, PASSWORD"(
// identified by, IDENTIFIED BY, -P
//
// Behind those samples is real password. Password may be surrounded by "".
// If it is surrounded by "", filter content should be "******".
var sample = []string{
	"password=", "Password=", "PASSWORD=", "password\"=", "Password\"=", "PASSWORD\"=",
	"password:", "Password:", "PASSWORD:", "password\":", "Password\":", "PASSWORD\":",
	"password(", "Password(", "PASSWORD(", "password\"(", "Password\"(", "PASSWORD\"(",
	"identified by ", "IDENTIFIED BY ", "-P ",
}

func defaultFilterPassword(line string) string {
	for _, key := range sample {
		if idx := strings.Index(line, key); -1 != idx {
			s, e := idx+len(key), idx+len(key)
			if s >= len(line) {
				line += "******"
				return line
			}

			prevPart := line[:idx]
			nextPart := ""
			hidePass := ""
			if line[s] == '"' {
				for i := s + 1; i < len(line) && line[i] != '"'; {
					if line[i] == '\\' && i+1 < len(line) && line[i+1] == '"' {
						i += 2
					} else {
						i++
					}
					e = i + 1
				}
				hidePass = `"******"`
				nextPart = line[e:]
			} else {
				for i := s; i < len(line); i, e = i+1, i+1 {
					if strings.HasSuffix(key, "(") && line[i] == ')' {
						break
					}
					if line[i] == ' ' || line[i] == ',' {
						break
					}
				}
				hidePass = "******"
				nextPart = line[e:]
			}
			line = defaultFilterPassword(prevPart) + line[idx:idx+len(key)] + hidePass + defaultFilterPassword(nextPart)
			break
		}
	}
	return line
}

func (b *Builder) buildLine(level int, msg string, args ...interface{}) *Builder {
	line := msg
	if len(args) > 0 {
		line = fmt.Sprintf(msg, args...)
	}
	if level >= len(b.lines) {
		return b
	}
	b.lines[level] = b.lines[level] + line
	return b
}

func (b *Builder) UserInfo(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(user) {
		return b.buildLine(user, "[INFO] "+msg, args...)
	}
	return b
}

func (b *Builder) UserWarn(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(user) {
		return b.buildLine(user, "[WARN] "+msg, args...)
	}
	return b
}

func (b *Builder) UserError(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(user) {
		return b.buildLine(user, "[ERROR] "+msg, args...)
	}
	return b
}

func (b *Builder) Brief(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(brief) {
		return b.buildLine(brief, msg, args...)
	}
	return b
}

func (b *Builder) Detail(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(detail) {
		return b.buildLine(detail, msg, args...)
	}
	return b
}

func (b *Builder) Key(msg string, args ...interface{}) *Builder {
	if shouldPrintLog(key) {
		return b.buildLine(key, msg, args...)
	}
	return b
}

//Done is a warpper that transfer call stack depth for logging filename and line
func (b *Builder) Done() { b.done(3) }

func (b *Builder) done(calldepth int) {
	ts := time.Now().Format(LogTimeStamp)

	mutex := b.logger.getLock()
	mutex.Lock()
	defer mutex.Unlock()

	shouldWriteLog := false
	lastSeconds := 0
	{
		tsNow := time.Now()
		if "" != b.diluteKey {
			dilute := dilutesInstance[b.diluteKey]
			if nil == dilute || dilute.lastTimestamp.Add(time.Duration(b.diluteDurationSeconds)*time.Second).Before(tsNow) {
				dilutesInstance[b.diluteKey] = &diluteRecord{
					firstTimestamp:      tsNow,
					lastTimestamp:       tsNow,
					checkpointTimestamp: tsNow,
				}
				shouldWriteLog = true
			} else {
				dilutesInstance[b.diluteKey].lastTimestamp = tsNow
				if dilute.checkpointTimestamp.Add(time.Duration(b.diluteCheckpointSeconds) * time.Second).Before(tsNow) {
					dilutesInstance[b.diluteKey].checkpointTimestamp = tsNow
					lastSeconds = int(tsNow.Sub(dilute.firstTimestamp).Seconds())
					shouldWriteLog = true
				}
			}

			//so-called "recycle", hope it's no use
			if len(dilutesInstance) > DILUTE_LIMITS {
				dilutesInstance = map[string]*diluteRecord{}
			}
		} else {
			shouldWriteLog = true
		}
	}

	if !shouldWriteLog {
		return
	}

	rawContent := ""
	for i, levelContent := range b.lines {
		if "" != rawContent && "" != levelContent {
			rawContent = rawContent + ","
		}
		rawContent = rawContent + levelContent

		if "" == rawContent {
			continue
		}

		//filter password
		content := b.filterPassword(rawContent)

		//add timestamp & prefix
		//set user.log histroy
		if user == i {
			content = "[" + ts + "] " + content

		} else if _, ok := b.logger.(*fileLogger); !ok || detail == i {
			_, file, line, ok := runtime.Caller(calldepth)
			if !ok {
				file = "???"
				line = 0
			}
			content = "[" + ts + "] " + b.prefix.ToPrefix() + " " +
				filepath.Base(file) + ":" + strconv.Itoa(line) + ": " + content

		} else {
			content = "[" + ts + "] " + b.prefix.ToPrefix() + " " + content
		}

		if lastSeconds > 0 {
			content = fmt.Sprintf("%s. Last for %vs", content, lastSeconds)
		}

		b.logger.printLog(i, content)
	}
}
