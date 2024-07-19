package logs

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
	uuid "github.com/satori/go.uuid"
)

const LogPrefixAmazonRds string = "%t:%r:%u@%d:[%p]:"
const LogPrefixAzure string = "%t-%c-"
const LogPrefixCustom1 string = "%m [%p][%v] : [%l-1] %q[app=%a] "
const LogPrefixCustom2 string = "%t [%p-%l] %q%u@%d "
const LogPrefixCustom3 string = "%m [%p] %q[user=%u,db=%d,app=%a] "
const LogPrefixCustom4 string = "%m [%p] %q[user=%u,db=%d,app=%a,host=%h] "
const LogPrefixCustom5 string = "%t [%p]: [%l-1] user=%u,db=%d - PG-%e "
const LogPrefixCustom6 string = "%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h "
const LogPrefixCustom7 string = "%t [%p]: [%l-1] [trx_id=%x] user=%u,db=%d "
const LogPrefixCustom8 string = "[%p]: [%l-1] db=%d,user=%u "
const LogPrefixCustom9 string = "%m %r %u %a [%c] [%p] "
const LogPrefixCustom10 string = "%m [%p]: [%l-1] db=%d,user=%u "
const LogPrefixCustom11 string = "pid=%p,user=%u,db=%d,app=%a,client=%h "
const LogPrefixCustom12 string = "user=%u,db=%d,app=%a,client=%h "
const LogPrefixCustom13 string = "%p-%s-%c-%l-%h-%u-%d-%m "
const LogPrefixCustom14 string = "%m [%p][%b][%v][%x] %q[user=%u,db=%d,app=%a] "
const LogPrefixCustom15 string = "%m [%p] %q%u@%d "
const LogPrefixCustom16 string = "%t [%p] %q%u@%d %h "
const LogPrefixSimple string = "%m [%p] "
const LogPrefixHeroku1 string = " sql_error_code = %e "
const LogPrefixHeroku2 string = ` sql_error_code = %e time_ms = "%m" pid="%p" proc_start_time="%s" session_id="%c" vtid="%v" tid="%x" log_line="%l" %qdatabase="%d" connection_source="%r" user="%u" application_name="%a" `

// Used only to recognize the Heroku hobby tier log_line_prefix to give a warning (logs are not supported
// on hobby tier) and avoid errors during prefix check; logs with this prefix are never actually received
const LogPrefixHerokuHobbyTier string = " database = %d connection_source = %r sql_error_code = %e "
const LogPrefixEmpty string = ""

var RecommendedPrefixIdx = 4

var SupportedPrefixes = []string{
	LogPrefixAmazonRds, LogPrefixAzure, LogPrefixCustom1, LogPrefixCustom2,
	LogPrefixCustom3, LogPrefixCustom4, LogPrefixCustom5, LogPrefixCustom6,
	LogPrefixCustom7, LogPrefixCustom8, LogPrefixCustom9, LogPrefixCustom10,
	LogPrefixCustom11, LogPrefixCustom12, LogPrefixCustom13, LogPrefixCustom14,
	LogPrefixCustom15, LogPrefixCustom16,
	LogPrefixSimple, LogPrefixHeroku1, LogPrefixHeroku2, LogPrefixEmpty,
}

// Every one of these regexps should produce exactly one matching group
var TimeRegexp = `(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? [\-+]?\w+)` // %t or %m (or %s)
var HostAndPortRegexp = `(.+(?:\(\d+\))?)?`                                  // %r
var PidRegexp = `(\d+)`                                                      // %p
var UserRegexp = `(\S*)`                                                     // %u
var DbRegexp = `(\S*)`                                                       // %d
var AppBeforeSpaceRegexp = `(\S*)`                                           // %a
var AppBeforeCommaRegexp = `([^,]*)`                                         // %a
var AppBeforeQuoteRegexp = `([^"]*)`                                         // %a
var AppInsideBracketsRegexp = `(\[unknown\]|[^,\]]*)`                        // %a
var HostRegexp = `(\S*)`                                                     // %h
var VirtualTxRegexp = `(\d+/\d+)?`                                           // %v
var LogLineCounterRegexp = `(\d+)`                                           // %l
var SqlstateRegexp = `(\w{5})`                                               // %e
var TransactionIdRegexp = `(\d+)`                                            // %x
var SessionIdRegexp = `(\w+\.\w+)`                                           // %c
var BackendTypeRegexp = `([\w ]+)`                                           // %b
// Missing:
// - %n (unix timestamp)
// - %i (command tag)

var LevelAndContentRegexp = `(\w+):\s+(.*\n?)$`
var LogPrefixAmazonRdsRegexp = regexp.MustCompile(`(?s)^` + TimeRegexp + `:` + HostAndPortRegexp + `:` + UserRegexp + `@` + DbRegexp + `:\[` + PidRegexp + `\]:` + LevelAndContentRegexp)
var LogPrefixAzureRegexp = regexp.MustCompile(`(?s)^` + TimeRegexp + `-` + SessionIdRegexp + `-` + LevelAndContentRegexp)
var LogPrefixCustom1Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]\[` + VirtualTxRegexp + `\] : \[` + LogLineCounterRegexp + `-1\] (?:\[app=` + AppInsideBracketsRegexp + `\] )?` + LevelAndContentRegexp)
var LogPrefixCustom2Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `-` + LogLineCounterRegexp + `\] ` + `(?:` + UserRegexp + `@` + DbRegexp + ` )?` + LevelAndContentRegexp)
var LogPrefixCustom3Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\] (?:\[user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppInsideBracketsRegexp + `\] )?` + LevelAndContentRegexp)
var LogPrefixCustom4Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\] (?:\[user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppBeforeCommaRegexp + `,host=` + HostRegexp + `\] )?` + LevelAndContentRegexp)
var LogPrefixCustom5Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]: \[` + LogLineCounterRegexp + `-1\] user=` + UserRegexp + `,db=` + DbRegexp + ` - PG-` + SqlstateRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom6Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]: \[` + LogLineCounterRegexp + `-1\] user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppBeforeCommaRegexp + `,client=` + HostRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom7Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]: \[` + LogLineCounterRegexp + `-1\] \[trx_id=` + TransactionIdRegexp + `\] user=` + UserRegexp + `,db=` + DbRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom8Regexp = regexp.MustCompile(`(?s)^\[` + PidRegexp + `\]: \[` + LogLineCounterRegexp + `-1\] db=` + DbRegexp + `,user=` + UserRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom9Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` ` + HostAndPortRegexp + ` ` + UserRegexp + ` ` + AppBeforeSpaceRegexp + ` \[` + SessionIdRegexp + `\] \[` + PidRegexp + `\] ` + LevelAndContentRegexp)
var LogPrefixCustom10Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]: \[` + LogLineCounterRegexp + `-1\] db=` + DbRegexp + `,user=` + UserRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom11Regexp = regexp.MustCompile(`(?s)^pid=` + PidRegexp + `,user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppBeforeCommaRegexp + `,client=` + HostRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom12Regexp = regexp.MustCompile(`(?s)^user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppBeforeCommaRegexp + `,client=` + HostRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom13Regexp = regexp.MustCompile(`(?s)^` + PidRegexp + `-` + TimeRegexp + `-` + SessionIdRegexp + `-` + LogLineCounterRegexp + `-` + HostRegexp + `-` + UserRegexp + `-` + DbRegexp + `-` + TimeRegexp + ` ` + LevelAndContentRegexp)
var LogPrefixCustom14Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\]\[` + BackendTypeRegexp + `\]\[` + VirtualTxRegexp + `\]\[` + TransactionIdRegexp + `\] (?:\[user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppInsideBracketsRegexp + `\] )?` + LevelAndContentRegexp)
var LogPrefixCustom15Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\] ` + `(?:` + UserRegexp + `@` + DbRegexp + ` )?` + LevelAndContentRegexp)
var LogPrefixCustom16Regexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\] ` + `(?:` + UserRegexp + `@` + DbRegexp + ` ` + HostRegexp + ` )?` + LevelAndContentRegexp)
var LogPrefixSimpleRegexp = regexp.MustCompile(`(?s)^` + TimeRegexp + ` \[` + PidRegexp + `\] ` + LevelAndContentRegexp)
var LogPrefixNoTimestampUserDatabaseAppRegexp = regexp.MustCompile(`(?s)^\[user=` + UserRegexp + `,db=` + DbRegexp + `,app=` + AppInsideBracketsRegexp + `\] ` + LevelAndContentRegexp)
var LogPrefixHeroku1Regexp = regexp.MustCompile(`^ sql_error_code = ` + SqlstateRegexp + " " + LevelAndContentRegexp)
var LogPrefixHeroku2Regexp = regexp.MustCompile(`^ sql_error_code = ` + SqlstateRegexp + ` time_ms = "` + TimeRegexp + `" pid="` + PidRegexp + `" proc_start_time="` + TimeRegexp + `" session_id="` + SessionIdRegexp + `" vtid="` + VirtualTxRegexp + `" tid="` + TransactionIdRegexp + `" log_line="` + LogLineCounterRegexp + `" (?:database="` + DbRegexp + `" connection_source="` + HostAndPortRegexp + `" user="` + UserRegexp + `" application_name="` + AppBeforeQuoteRegexp + `" )?` + LevelAndContentRegexp)

var SyslogSequenceAndSplitRegexp = `(\[[\d-]+\])?`

var RsyslogLevelAndContentRegexp = `(?:(\w+):\s+)?(.*\n?)$`
var RsyslogTimeRegexp = `(\w+\s+\d+ \d{2}:\d{2}:\d{2})`
var RsyslogHostnameRegxp = `(\S+)`
var RsyslogProcessNameRegexp = `(\w+)`
var RsyslogRegexp = regexp.MustCompile(`^` + RsyslogTimeRegexp + ` ` + RsyslogHostnameRegxp + ` ` + RsyslogProcessNameRegexp + `\[` + PidRegexp + `\]: ` + SyslogSequenceAndSplitRegexp + ` ` + RsyslogLevelAndContentRegexp)

func IsSupportedPrefix(prefix string) bool {
	for _, supportedPrefix := range SupportedPrefixes {
		if supportedPrefix == prefix {
			return true
		}
	}
	return false
}

func ParseLogLineWithPrefix(prefix string, line string, tz *time.Location) (logLine state.LogLine, ok bool) {
	var timePart, userPart, dbPart, appPart, pidPart, logLineNumberPart, levelPart, contentPart string

	rsyslog := false

	// Only read the first 1000 characters of a log line to parse the log_line_prefix
	//
	// This reduces the overhead for very long loglines, because we don't pass in the
	// whole line to the regexp engine (twice).
	lineExtra := ""
	if len(line) > 1000 {
		lineExtra = line[1000:]
		line = line[0:1000]
	}

	if prefix == "" {
		if LogPrefixAmazonRdsRegexp.MatchString(line) {
			prefix = LogPrefixAmazonRds
		} else if LogPrefixAzureRegexp.MatchString(line) {
			prefix = LogPrefixAzure
		} else if LogPrefixCustom1Regexp.MatchString(line) {
			prefix = LogPrefixCustom1
		} else if LogPrefixCustom2Regexp.MatchString(line) {
			prefix = LogPrefixCustom2
		} else if LogPrefixCustom4Regexp.MatchString(line) { // 4 is more specific than 3, so needs to go first
			prefix = LogPrefixCustom4
		} else if LogPrefixCustom3Regexp.MatchString(line) {
			prefix = LogPrefixCustom3
		} else if LogPrefixCustom5Regexp.MatchString(line) {
			prefix = LogPrefixCustom5
		} else if LogPrefixCustom6Regexp.MatchString(line) {
			prefix = LogPrefixCustom6
		} else if LogPrefixCustom7Regexp.MatchString(line) {
			prefix = LogPrefixCustom7
		} else if LogPrefixCustom8Regexp.MatchString(line) {
			prefix = LogPrefixCustom8
		} else if LogPrefixCustom9Regexp.MatchString(line) {
			prefix = LogPrefixCustom9
		} else if LogPrefixCustom10Regexp.MatchString(line) {
			prefix = LogPrefixCustom10
		} else if LogPrefixCustom11Regexp.MatchString(line) {
			prefix = LogPrefixCustom11
		} else if LogPrefixCustom12Regexp.MatchString(line) {
			prefix = LogPrefixCustom12
		} else if LogPrefixCustom13Regexp.MatchString(line) {
			prefix = LogPrefixCustom13
		} else if LogPrefixCustom14Regexp.MatchString(line) {
			prefix = LogPrefixCustom14
		} else if LogPrefixCustom15Regexp.MatchString(line) {
			prefix = LogPrefixCustom15
		} else if LogPrefixCustom16Regexp.MatchString(line) {
			prefix = LogPrefixCustom16
		} else if LogPrefixSimpleRegexp.MatchString(line) {
			prefix = LogPrefixSimple
		} else if LogPrefixHeroku2Regexp.MatchString(line) {
			prefix = LogPrefixHeroku2
		} else if LogPrefixHeroku1Regexp.MatchString(line) {
			// LogPrefixHeroku1 is a subset of 2, so it must be matched second
			prefix = LogPrefixHeroku1
		} else if RsyslogRegexp.MatchString(line) {
			rsyslog = true
		}
	}

	if rsyslog {
		parts := RsyslogRegexp.FindStringSubmatch(line)
		if len(parts) == 0 {
			return
		}

		timePart = fmt.Sprintf("%d %s", time.Now().Year(), parts[1])
		// ignore syslog hostname
		// ignore syslog process name
		pidPart = parts[4]
		// ignore syslog postgres sequence and split number
		levelPart = parts[6]
		contentPart = strings.Replace(parts[7], "#011", "\t", -1)

		parts = LogPrefixNoTimestampUserDatabaseAppRegexp.FindStringSubmatch(contentPart)
		if len(parts) == 6 {
			userPart = parts[1]
			dbPart = parts[2]
			appPart = parts[3]
			levelPart = parts[4]
			contentPart = parts[5]
		}
	} else {
		switch prefix {
		case LogPrefixAmazonRds: // "%t:%r:%u@%d:[%p]:"
			parts := LogPrefixAmazonRdsRegexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}

			timePart = parts[1]
			// skip %r (ip+port)
			userPart = parts[3]
			dbPart = parts[4]
			pidPart = parts[5]
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixAzure: // "%t-%c-"
			parts := LogPrefixAzureRegexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}

			timePart = parts[1]
			// skip %c (session id)
			levelPart = parts[3]
			contentPart = parts[4]
		case LogPrefixCustom1: // "%m [%p][%v] : [%l-1] %q[app=%a] "
			parts := LogPrefixCustom1Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			// skip %v (virtual TX)
			logLineNumberPart = parts[4]
			appPart = parts[5]
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixCustom2: // "%t [%p-1] %q%u@%d "
			parts := LogPrefixCustom2Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			logLineNumberPart = parts[3]
			userPart = parts[4]
			dbPart = parts[5]
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixCustom3: // "%m [%p] %q[user=%u,db=%d,app=%a] ""
			parts := LogPrefixCustom3Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			userPart = parts[3]
			dbPart = parts[4]
			appPart = parts[5]
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixCustom4: // "%m [%p] %q[user=%u,db=%d,app=%a,host=%h] "
			parts := LogPrefixCustom4Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			userPart = parts[3]
			dbPart = parts[4]
			appPart = parts[5]
			// skip %h (host)
			levelPart = parts[7]
			contentPart = parts[8]
		case LogPrefixCustom5: // "%t [%p]: [%l-1] user=%u,db=%d - PG-%e "
			parts := LogPrefixCustom5Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			logLineNumberPart = parts[3]
			userPart = parts[4]
			dbPart = parts[5]
			// skip %e (SQLSTATE)
			levelPart = parts[7]
			contentPart = parts[8]
		case LogPrefixCustom6: // "%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h "
			parts := LogPrefixCustom6Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			logLineNumberPart = parts[3]
			userPart = parts[4]
			dbPart = parts[5]
			// skip %a (application name)
			// skip %h (host)
			levelPart = parts[8]
			contentPart = parts[9]
		case LogPrefixCustom7: // "%t [%p]: [%l-1] [trx_id=%x] user=%u,db=%d "
			parts := LogPrefixCustom7Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			logLineNumberPart = parts[3]
			// skip %x (transaction id)
			userPart = parts[5]
			dbPart = parts[6]
			levelPart = parts[7]
			contentPart = parts[8]
		case LogPrefixCustom8: // "[%p]: [%l-1] db=%d,user=%u "
			parts := LogPrefixCustom8Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			pidPart = parts[1]
			logLineNumberPart = parts[2]
			dbPart = parts[3]
			userPart = parts[4]
			levelPart = parts[5]
			contentPart = parts[6]
		case LogPrefixCustom9: // "%m %r %u %a [%c] [%p] "
			parts := LogPrefixCustom9Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			// skip %r (ip+port)
			userPart = parts[3]
			appPart = parts[4]
			// skip %c (session id)
			pidPart = parts[6]
			levelPart = parts[7]
			contentPart = parts[8]
		case LogPrefixCustom10: // "%t [%p]: [%l-1] db=%d,user=%u "
			parts := LogPrefixCustom10Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			logLineNumberPart = parts[3]
			dbPart = parts[4]
			userPart = parts[5]
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixCustom11: // "pid=%p,user=%u,db=%d,app=%a,client=%h "
			parts := LogPrefixCustom11Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			pidPart = parts[1]
			userPart = parts[2]
			dbPart = parts[3]
			// skip %a (application name)
			// skip %h (host)
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixCustom12: // "user=%u,db=%d,app=%a,client=%h "
			parts := LogPrefixCustom12Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			userPart = parts[1]
			dbPart = parts[2]
			// skip %a (application name)
			// skip %h (host)
			levelPart = parts[5]
			contentPart = parts[6]
		case LogPrefixCustom13: // "%p-%s-%c-%l-%h-%u-%d-%m "
			parts := LogPrefixCustom13Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			pidPart = parts[1]
			// skip %s
			// skip %c
			logLineNumberPart = parts[4]
			// skip %h (host)
			userPart = parts[6]
			dbPart = parts[7]
			timePart = parts[8]
			levelPart = parts[9]
			contentPart = parts[10]
		case LogPrefixCustom14: // "%m [%p][%b][%v][%x] %q[user=%u,db=%d,app=%a] "
			parts := LogPrefixCustom14Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}

			timePart = parts[1]
			pidPart = parts[2]
			// skip %b
			// skip %v
			// skip %x
			userPart = parts[6]
			dbPart = parts[7]
			appPart = parts[8]
			levelPart = parts[9]
			contentPart = parts[10]
		case LogPrefixCustom15: // "%m [%p] %q%u@%d "
			parts := LogPrefixCustom15Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			userPart = parts[3]
			dbPart = parts[4]
			levelPart = parts[5]
			contentPart = parts[6]
		case LogPrefixCustom16: // "%t [%p] %q%u@%d %h "
			parts := LogPrefixCustom16Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			userPart = parts[3]
			dbPart = parts[4]
			// skip %h (host)
			levelPart = parts[6]
			contentPart = parts[7]
		case LogPrefixSimple: // "%t [%p] "
			parts := LogPrefixSimpleRegexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			timePart = parts[1]
			pidPart = parts[2]
			levelPart = parts[3]
			contentPart = parts[4]
		case LogPrefixHeroku1:
			parts := LogPrefixHeroku1Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			// skip %e
			levelPart = parts[2]
			contentPart = parts[3]
		case LogPrefixHeroku2:
			parts := LogPrefixHeroku2Regexp.FindStringSubmatch(line)
			if len(parts) == 0 {
				return
			}
			// skip %e
			timePart = parts[2]
			pidPart = parts[3]
			// skip %s
			// skip %c
			// skip %v
			// skip %x
			logLineNumberPart = parts[8]
			dbPart = parts[9]
			// skip %r
			userPart = parts[11]
			appPart = parts[12]
			levelPart = parts[13]
			contentPart = parts[14]
		default:
			// Some callers use the content of unparsed lines to stitch multi-line logs together
			logLine.Content = line + lineExtra
			return
		}
	}

	if timePart != "" {
		occurredAt := getOccurredAt(timePart, tz, rsyslog)
		if occurredAt.IsZero() {
			return
		}
		logLine.OccurredAt = occurredAt
	}

	if userPart != "[unknown]" {
		logLine.Username = userPart
	}
	if dbPart != "[unknown]" {
		logLine.Database = dbPart
	}
	if appPart != "[unknown]" {
		logLine.Application = appPart
	}
	if logLineNumberPart != "" {
		logLineNumber, _ := strconv.ParseInt(logLineNumberPart, 10, 32)
		logLine.LogLineNumber = int32(logLineNumber)
	}

	backendPid, _ := strconv.ParseInt(pidPart, 10, 32)
	logLine.BackendPid = int32(backendPid)
	logLine.Content = contentPart + lineExtra

	// This is actually a continuation of a previous line
	if levelPart == "" {
		return
	}

	logLine.LogLevel = pganalyze_collector.LogLineInformation_LogLevel(pganalyze_collector.LogLineInformation_LogLevel_value[levelPart])
	ok = true

	return
}

func getOccurredAt(timePart string, tz *time.Location, rsyslog bool) time.Time {
	if tz != nil && !rsyslog {
		lastSpaceIdx := strings.LastIndex(timePart, " ")
		if lastSpaceIdx == -1 {
			return time.Time{}
		}
		timePartNoTz := timePart[0:lastSpaceIdx]
		result, err := time.ParseInLocation("2006-01-02 15:04:05", timePartNoTz, tz)
		if err != nil {
			return time.Time{}
		}

		return result
	}

	// Assume Postgres time format unless overriden by the prefix (e.g. syslog)
	var timeFormat, timeFormatAlt string
	if rsyslog {
		timeFormat = "2006 Jan  2 15:04:05"
		timeFormatAlt = ""
	} else {
		timeFormat = "2006-01-02 15:04:05 -0700"
		timeFormatAlt = "2006-01-02 15:04:05 MST"
	}

	ts, err := time.Parse(timeFormat, timePart)
	if err != nil {
		if timeFormatAlt != "" {
			// Ensure we have the correct format remembered for ParseInLocation call that may happen later
			timeFormat = timeFormatAlt
			ts, err = time.Parse(timeFormat, timePart)
		}
		if err != nil {
			return time.Time{}
		}
	}

	// Handle non-UTC timezones in systems that have log_timezone set to a different
	// timezone value than their system timezone. This is necessary because Go otherwise
	// only reads the timezone name but does not set the timezone offset, see
	// https://pkg.go.dev/time#Parse
	zone, offset := ts.Zone()
	if offset == 0 && zone != "UTC" && zone != "" {
		var zoneLocation *time.Location
		zoneNum, err := strconv.Atoi(zone)
		if err == nil {
			zoneLocation = time.FixedZone(zone, zoneNum*3600)
		} else {
			zoneLocation, err = time.LoadLocation(zone)
			if err != nil {
				// We don't know which timezone this is (and a timezone name is present), so we can't process this log line
				return time.Time{}
			}
		}
		ts, err = time.ParseInLocation(timeFormat, timePart, zoneLocation)
		if err != nil {
			// Technically this should not occur (as we should have already failed previously in time.Parse)
			return time.Time{}
		}
	}
	return ts
}

type LineReader interface {
	ReadString(delim byte) (string, error)
}

func ParseAndAnalyzeBuffer(logStream LineReader, linesNewerThan time.Time, server *state.Server) ([]state.LogLine, []state.PostgresQuerySample) {
	var logLines []state.LogLine
	var currentByteStart int64 = 0
	var tz = server.GetLogTimezone()

	for {
		line, err := logStream.ReadString('\n')
		byteStart := currentByteStart
		currentByteStart += int64(len(line))

		// This is intentionally after updating currentByteStart, since we consume the
		// data in the file even if an error is returned
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Log Read ERROR: %s", err)
			}
			break
		}

		logLine, ok := ParseLogLineWithPrefix("", line, tz)
		if !ok {
			// Assume that a parsing error in a follow-on line means that we actually
			// got additional data for the previous line
			if len(logLines) > 0 && logLine.Content != "" {
				logLines[len(logLines)-1].Content += logLine.Content
				logLines[len(logLines)-1].ByteEnd += int64(len(logLine.Content))
			}
			continue
		}

		// Ignore loglines which are outside our time window
		if logLine.OccurredAt.Before(linesNewerThan) {
			continue
		}

		// Ignore loglines that are ignored server-wide (e.g. because they are
		// log_statement=all/log_duration=on lines). Note this intentionally
		// runs after multi-line log lines have been stitched together.
		if server.IgnoreLogLine(logLine.Content) {
			continue
		}

		logLine.ByteStart = byteStart
		logLine.ByteContentStart = byteStart + int64(len(line)-len(logLine.Content))
		logLine.ByteEnd = byteStart + int64(len(line))

		// Generate unique ID that can be used to reference this line
		logLine.UUID = uuid.NewV4()

		logLines = append(logLines, logLine)
	}

	newLogLines, newSamples := AnalyzeLogLines(logLines)
	return newLogLines, newSamples
}
