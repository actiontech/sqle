package logs

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/pganalyze/collector/logs/querysample"
	"github.com/pganalyze/collector/logs/util"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
)

type match struct {
	prefixes      []string
	regexp        *regexp.Regexp
	secrets       []state.LogSecretKind
	remainderKind state.LogSecretKind
}

type analyzeGroup struct {
	classification pganalyze_collector.LogLineInformation_LogClassification
	primary        match
	detail         match
	hint           match
}

var utcTimestampRegexp = `(\d+-\d+-\d+ \d+:\d+:\d+(?:\.\d+)?(?:[\d:+-]+| \w+))`

var autoExplain = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STATEMENT_AUTO_EXPLAIN,
	primary: match{
		prefixes:      []string{"duration: "},
		regexp:        regexp.MustCompile(`^duration: ([\d\.]+) ms\s+ plan:\s+`),
		secrets:       []state.LogSecretKind{0},
		remainderKind: state.StatementTextLogSecret,
	},
}
var duration = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STATEMENT_DURATION,
	primary: match{
		prefixes:      []string{"duration: "},
		regexp:        regexp.MustCompile(`^duration: ([\d\.]+) ms(?:  (?:statement|(parse|bind|execute|execute fetch from) ([^:]+)(?:/([^:]+))?):\s+)?`),
		secrets:       []state.LogSecretKind{0, 0, 0, 0},
		remainderKind: state.StatementTextLogSecret,
	},
	detail: match{
		regexp:  regexp.MustCompile(`(?:parameters: |, )\$\d+ = (?:(NULL)|'([^']*)')`),
		secrets: []state.LogSecretKind{state.StatementParameterLogSecret, state.StatementParameterLogSecret},
	},
}
var autovacuumCancel = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_CANCEL,
	primary: match{
		prefixes: []string{"canceling autovacuum task"},
	},
}
var skippingAnalyze = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SKIPPING_ANALYZE_LOCK_NOT_AVAILABLE,
	primary: match{
		prefixes: []string{"skipping analyze of"},
		regexp:   regexp.MustCompile(`^skipping analyze of "([^"]+)" --- lock not available`),
		secrets:  []state.LogSecretKind{0},
	},
}
var skippingVacuum = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SKIPPING_VACUUM_LOCK_NOT_AVAILABLE,
	primary: match{
		prefixes: []string{"skipping vacuum of"},
		regexp:   regexp.MustCompile(`^skipping vacuum of "([^"]+)" --- lock not available`),
		secrets:  []state.LogSecretKind{0},
	},
}
var autoVacuum = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_COMPLETED,
	primary: match{
		prefixes: []string{"automatic vacuum of table", "automatic aggressive vacuum of table", "automatic aggressive vacuum to prevent wraparound of table"},
		regexp: regexp.MustCompile(`^automatic (aggressive )?vacuum (to prevent wraparound )?of table "(.+?)": index scans: (\d+),?\s*` +
			`(?:elapsed time: \d+ \w+, index vacuum time: \d+ \w+,)?\s*` + // Google AlloyDB for PostgreSQL
			`pages: (\d+) removed, (\d+) remain, (?:(\d+) skipped due to pins, (\d+) skipped frozen|(\d+) scanned \(([\d.]+)% of total\)),?\s*` +
			`(?:\d+ skipped using mintxid)?,?\s*` + // Google AlloyDB for PostgreSQL
			`tuples: (\d+) removed, (\d+) remain, (\d+) are dead but not yet removable(?:, oldest xmin: (\d+))?,?\s*` +
			`(?:tuples missed: (\d+) dead from (\d+) pages not removed due to cleanup lock contention)?,?\s*` + // Postgres 15+
			`(?:removable cutoff: (\d+), which was (\d+) XIDs old when operation ended)?,?\s*` + // Postgres 15+
			`(?:new relfrozenxid: (\d+), which is (\d+) XIDs ahead of previous value)?,?\s*` + // Postgres 15+
			`(?:new relminmxid: (\d+), which is (\d+) MXIDs ahead of previous value)?,?\s*` + // Postgres 15+
			`(?:index scan (not needed|needed|bypassed|bypassed by failsafe): (\d+) pages from table \(([\d.]+)% of total\) (?:have|had) (\d+) dead item identifiers(?: removed)?)?,?\s*` + // Postgres 14+
			`((?:index ".+?": pages: \d+ in total, \d+ newly deleted, \d+ currently deleted, \d+ reusable,?\s*)*)?` + // Postgres 14+
			`(?:I/O timings: read: ([\d.]+) ms, write: ([\d.]+) ms)?,?\s*` + // Postgres 14+
			`(?:avg read rate: ([\d.]+) MB/s, avg write rate: ([\d.]+) MB/s)?,?\s*` + // Postgres 14+
			`buffer usage: (\d+) hits, (\d+) misses, (\d+) dirtied,?\s*` +
			`(?:avg read rate: ([\d.]+) MB/s, avg write rate: ([\d.]+) MB/s)?,?\s*` + // Postgres 13 and older
			`(?:WAL usage: (\d+) records, (\d+) full page images, (\d+) bytes)?,?\s*` + // Postgres 14+
			`system usage: CPU(?:(?: ([\d.]+)s/([\d.]+)u sec elapsed ([\d.]+) sec)|(?:: user: ([\d.]+) s, system: ([\d.]+) s, elapsed: ([\d.]+) s))`),
		secrets: []state.LogSecretKind{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}
var autoAnalyze = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_AUTOANALYZE_COMPLETED,
	primary: match{
		prefixes: []string{"automatic analyze of table"},
		regexp: regexp.MustCompile(`^automatic analyze of table "(.+?)"\s*` +
			`(?:I/O timings: read: ([\d.]+) ms, write: ([\d.]+) ms)?\s*` + // Postgres 14+
			`(?:avg read rate: ([\d.]+) MB/s, avg write rate: ([\d.]+) MB/s)?\s*` + // Postgres 14+
			`(?:buffer usage: (\d+) hits, (\d+) misses, (\d+) dirtied)?\s*` + // Postgres 14+
			`system usage: CPU(?:(?: ([\d.]+)s/([\d.]+)u sec elapsed ([\d.]+) sec)|(?:: user: ([\d.]+) s, system: ([\d.]+) s, elapsed: ([\d.]+) s))`),
		secrets: []state.LogSecretKind{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}
var checkpointStarting = analyzeGroup{
	primary: match{
		prefixes: []string{"checkpoint", "restartpoint"},
		regexp:   regexp.MustCompile(`^(checkpoint|restartpoint) starting: ([a-z- ]+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var checkpointComplete = analyzeGroup{
	primary: match{
		regexp: regexp.MustCompile(`^(checkpoint|restartpoint) complete: wrote (\d+) buffers \(([\d\.]+)%\); ` +
			`(\d+) (?:transaction log|WAL) file\(s\) added, (\d+) removed, (\d+) recycled; ` +
			`write=([\d\.]+) s, sync=([\d\.]+) s, total=([\d\.]+) s; ` +
			`sync files=(\d+), longest=([\d\.]+) s, average=([\d\.]+) s` +
			`; distance=(\d+) kB, estimate=(\d+) kB`),
		secrets: []state.LogSecretKind{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}
var checkpointsTooFrequent = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CHECKPOINT_TOO_FREQUENT,
	primary: match{
		prefixes: []string{"checkpoints"},
		regexp:   regexp.MustCompile(`^checkpoints are occurring too frequently \((\d+) seconds? apart\)`),
		secrets:  []state.LogSecretKind{0},
	},
	hint: match{
		prefixes: []string{"Consider increasing the configuration parameter \"max_wal_size\"."},
	},
}
var restartpointAt = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_RESTARTPOINT_AT,
	primary: match{
		prefixes: []string{"recovery restart point at"},
		regexp:   regexp.MustCompile(`^recovery restart point at (\w+)/(\w+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
	detail: match{
		prefixes: []string{"last completed transaction was at log time "},
		regexp:   regexp.MustCompile(`^last completed transaction was at log time (\d+-\d+-\d+ \d+:\d+:\d+\.\d+[\d:+-]+)`),
		secrets:  []state.LogSecretKind{0},
	},
}
var connectionReceived = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_RECEIVED,
	primary: match{
		prefixes: []string{"connection received: "},
		regexp:   regexp.MustCompile(`^connection received: host=([^ ]+)( port=\w+)?`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var connectionAuthorized = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_AUTHORIZED,
	primary: match{
		prefixes: []string{"connection authorized: "},
		regexp:   regexp.MustCompile(`^connection authorized: user=\w+( database=\w+)?( application_name=.+)?( SSL enabled \(protocol=([\w.]+), cipher=[\w-]+, compression=\w+\))?`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
}
var connectionRejected = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
	primary: match{
		prefixes: []string{"pg_hba.conf rejects connection ", "no pg_hba.conf entry for"},
		regexp:   regexp.MustCompile(`^(?:(?:pg_hba.conf rejects connection|no pg_hba.conf entry) for host "[^"]+", user "[^"]+", database "[^"]+"(, SSL on|, SSL off)?|password authentication failed for user "[^"]+")`),
		secrets:  []state.LogSecretKind{0},
	},
}
var authenticationFailed = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
	primary: match{
		prefixes: []string{"password authentication failed for user", "Ident authentication failed for user", "could not connect to Ident server"},
		regexp:   regexp.MustCompile(`^(?:(?:Ident|password) authentication failed for user "([^"]+)"|could not connect to Ident server at address "([^"]+)", port \d+: ([\w ]+))`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^(?:(?:Role|User|Password does not match for user|Password of user) "([^"]+)" ?(?:does not have a valid SCRAM secret|does not exist|has no password assigned|has an expired password|has a password that cannot be used with MD5 authentication|is in unrecognized format)?\.\s+)?Connection matched pg_hba.conf line \d+: "([^"]+)"`),
		secrets: []state.LogSecretKind{0, state.OpsLogSecret},
	},
}
var roleNotAllowedLogin = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
	primary: match{
		prefixes: []string{"role"},
		regexp:   regexp.MustCompile(`^role "([^"]+)" is not permitted to log in`),
		secrets:  []state.LogSecretKind{0},
	},
}
var databaseNotAcceptingConnections = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
	primary: match{
		prefixes: []string{"database"},
		regexp:   regexp.MustCompile(`^database "([^"]+)" is not currently accepting connections`),
		secrets:  []state.LogSecretKind{0},
	},
}
var disconnection = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_DISCONNECTED,
	primary: match{
		prefixes: []string{"disconnection: "},
		regexp:   regexp.MustCompile(`^disconnection: session time: (\d+):(\d+):([\d\.]+) user=\w+ database=\w+ host=[^ ]+( port=\w+)?`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
}
var connectionClientFailedToConnect = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_CLIENT_FAILED_TO_CONNECT,
	primary: match{
		prefixes: []string{"incomplete startup packet"},
	},
}
var connectionLostOpenTx = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST_OPEN_TX,
	primary: match{
		prefixes: []string{"unexpected EOF on client connection with an open transaction"},
	},
}
var connectionLost = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
	primary: match{
		prefixes: []string{
			"unexpected EOF on client connection",
			"connection to client lost",
			"terminating connection because protocol synchronization was lost",
		},
	},
}
var connectionLostSocketError = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
	primary: match{
		prefixes: []string{"could not receive data from client", "could not send data to client"},
		regexp:   regexp.MustCompile(`^could not (?:receive data from|send data to) client: [\w ]+`),
		secrets:  []state.LogSecretKind{0},
	},
}
var connectionTerminated = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CONNECTION_TERMINATED,
	primary: match{
		prefixes: []string{"terminating connection due to administrator command"},
	},
}
var outOfConnections = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_OUT_OF_CONNECTIONS,
	primary: match{
		prefixes: []string{
			"remaining connection slots are reserved for non-replication superuser connections",
			"sorry, too many clients already",
		},
	},
}
var tooManyConnectionsRole = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_TOO_MANY_CONNECTIONS_ROLE,
	primary: match{
		prefixes: []string{"too many connections for role"},
		regexp:   regexp.MustCompile(`^too many connections for role "([^"]+)"`),
		secrets:  []state.LogSecretKind{0},
	},
}
var tooManyConnectionsDatabase = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_TOO_MANY_CONNECTIONS_DATABASE,
	primary: match{
		prefixes: []string{"too many connections for database"},
		regexp:   regexp.MustCompile(`^too many connections for database "([^"]+)"`),
		secrets:  []state.LogSecretKind{0},
	},
}
var couldNotAcceptSSLConnection = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COULD_NOT_ACCEPT_SSL_CONNECTION,
	primary: match{
		prefixes: []string{"could not accept SSL connection: "},
		regexp:   regexp.MustCompile(`^could not accept SSL connection: [\w ]+`),
		secrets:  []state.LogSecretKind{0},
	},
}
var protocolErrorUnsupportedVersion = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_PROTOCOL_ERROR_UNSUPPORTED_VERSION,
	primary: match{
		prefixes: []string{"unsupported frontend protocol"},
		regexp:   regexp.MustCompile(`^unsupported frontend protocol \d+\.\d+: server supports \d+\.\d+ to \d+\.\d+`),
		secrets:  []state.LogSecretKind{},
	},
}
var protocolErrorIncompleteMessage = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_PROTOCOL_ERROR_INCOMPLETE_MESSAGE,
	primary: match{
		prefixes: []string{"incomplete message from client"},
	},
}
var walInvalidRecordLength = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_INVALID_RECORD_LENGTH,
	primary: match{
		prefixes: []string{"invalid record length at "},
		regexp:   regexp.MustCompile(`^invalid record length at (\w+)/(\w+)(?:: wanted \d+, got \d+)?`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var walRedo = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_REDO,
	primary: match{
		prefixes: []string{"redo starts at", "redo done at", "redo is not required"},
		regexp:   regexp.MustCompile(`^redo (?:(?:starts|done) at (\w+)/(\w+)|is not required)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var walRedoLastTx = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_REDO,
	primary: match{
		prefixes: []string{"last completed transaction was at log time "},
		regexp:   regexp.MustCompile(`^last completed transaction was at log time (\d+-\d+-\d+ \d+:\d+:\d+\.\d+[\d:+-]+)`),
		secrets:  []state.LogSecretKind{0},
	},
}
var archiveCommandFailed = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_ARCHIVE_COMMAND_FAILED,
	primary: match{
		prefixes: []string{"archive command"},
		regexp:   regexp.MustCompile(`^archive command (?:failed with exit code (\d+)|was terminated by signal (\d+)(: [\w ]+)?)`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^The failed archive command was: (.+)`),
		secrets: []state.LogSecretKind{state.OpsLogSecret},
	},
}
var archiverProcessExited = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_ARCHIVE_COMMAND_FAILED,
	primary: match{
		prefixes: []string{"archiver process"},
		regexp:   regexp.MustCompile(`^archiver process \(PID \d+\) exited with exit code \d+`),
		secrets:  []state.LogSecretKind{},
	},
}
var walBaseBackupComplete = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_WAL_BASE_BACKUP_COMPLETE,
	primary: match{
		prefixes: []string{"pg_stop_backup complete, all required WAL segments have been archived"},
	},
}
var lockAcquired = analyzeGroup{
	primary: match{
		prefixes: []string{"process"},
		regexp:   regexp.MustCompile(`^process \d+ acquired (\w+Lock) on (\w+)(?: [\(\)\d,]+)?( of \w+ \d+)* after ([\d\.]+) ms`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
}
var lockWait = analyzeGroup{
	primary: match{
		prefixes: []string{"process"},
		regexp:   regexp.MustCompile(`^process \d+ (still waiting|avoided deadlock|detected deadlock while waiting) for (\w+) on (\w+) (?:.+?) after ([\d\.]+) ms`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Process(?:es)? holding the lock: ([\d, ]+). Wait queue: ([\d, ]+)\.?`),
		secrets: []state.LogSecretKind{0, 0},
	},
}
var deadlock = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_DETECTED,
	primary: match{
		prefixes: []string{"deadlock detected"},
		regexp:   regexp.MustCompile(`^deadlock detected`),
		secrets:  []state.LogSecretKind{},
	},
	detail: match{
		regexp:  regexp.MustCompile(`(?m)^Process (\d+)(?: waits for \w+ on transaction \d+; blocked by process \d+.\s+|: (.+))`),
		secrets: []state.LogSecretKind{0, state.StatementTextLogSecret},
	},
	hint: match{
		prefixes: []string{"See server log for query details."},
	},
}
var lockTimeout = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_LOCK_TIMEOUT,
	primary: match{
		prefixes: []string{"canceling statement due to lock timeout"},
	},
}
var wraparoundWarning = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_WARNING,
	primary: match{
		prefixes: []string{"database"},
		regexp:   regexp.MustCompile(`^database (with OID (\d+)|"(.+?)") must be vacuumed within (\d+) transactions`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
	hint: match{
		prefixes: []string{"To avoid a database shutdown, execute a full-database VACUUM in"},
		regexp:   regexp.MustCompile(`^To avoid a database shutdown, execute a full-database VACUUM in "(.+)".\s+You might also need to commit or roll back old prepared transactions.`),
		secrets:  []state.LogSecretKind{0},
	},
}
var wraparoundError = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_ERROR,
	primary: match{
		prefixes: []string{"database is not accepting commands to avoid wraparound data loss in database"},
		regexp:   regexp.MustCompile(`^database is not accepting commands to avoid wraparound data loss in database (with OID (\d+)|"(.+?)")`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	hint: match{
		prefixes: []string{"Stop the postmaster and use a standalone backend to vacuum that database. You might also need to commit or roll back old prepared transactions."},
	},
}
var autovacuumLauncherStarted = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_LAUNCHER_STARTED,
	primary: match{
		prefixes: []string{"autovacuum launcher started"},
	},
}
var autovacuumLauncherShuttingDown = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_LAUNCHER_SHUTTING_DOWN,
	primary: match{
		prefixes: []string{"autovacuum launcher shutting down", "terminating autovacuum process due to administrator command"},
	},
}
var serverCrashed = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
	primary: match{
		prefixes: []string{"server process"},
		regexp:   regexp.MustCompile(`^server process \(PID (\d+)\) was terminated by signal (6|11)(: [\w ]+)?`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Failed process was running: (.*)`),
		secrets: []state.LogSecretKind{state.StatementTextLogSecret},
	},
}
var serverCrashedOtherProcesses = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
	primary: match{
		prefixes: []string{
			"terminating any other active server processes",
			"terminating connection because of crash of another server process",
			"all server processes terminated; reinitializing",
		},
	},
	detail: match{
		prefixes: []string{"The postmaster has commanded this server process to roll back the current transaction and exit, because another server process exited abnormally and possibly corrupted shared memory."},
	},
	hint: match{
		prefixes: []string{"In a moment you should be able to reconnect to the database and repeat your command."},
	},
}
var serverOutOfMemory = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_OUT_OF_MEMORY,
	primary: match{
		prefixes: []string{"out of memory"},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Failed on request of size (\d+)\.`),
		secrets: []state.LogSecretKind{0},
	},
}
var serverOutOfMemoryCrash = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_OUT_OF_MEMORY,
	primary: match{
		prefixes: []string{"out of memory", "server process"},
		regexp:   regexp.MustCompile(`^server process \(PID (\d+)\) was terminated by signal (9)(: [\w ]+)?`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
}
var serverStart = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_START,
	primary: match{
		prefixes: []string{
			"database system is ready to accept connections",
			"database system is ready to accept read only connections",
			"MultiXact member wraparound protections are now enabled",
			"entering standby mode",
			"redirecting log output to logging collector process",
			"ending log output to stderr",
		},
	},
}
var serverStartShutdownAt = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_START,
	primary: match{
		prefixes: []string{
			"database system was shut down at ",
			"database system was shut down in recovery at ",
		},
		regexp:  regexp.MustCompile(`^database system was shut down(?: in recovery)? at ` + utcTimestampRegexp),
		secrets: []state.LogSecretKind{0},
	},
}
var serverStartRecovering = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_START_RECOVERING,
	primary: match{
		prefixes: []string{
			"database system was interrupted; last known up at ",
			"database system shutdown was interrupted; last known up at ",
			"database system was interrupted while in recovery at ",
			"database system was not properly shut down; automatic recovery in progress",
		},
		regexp:  regexp.MustCompile(`^(?:database system was not properly shut down; automatic recovery in progress|(?:database system was interrupted; last known up at|database system shutdown was interrupted; last known up at|database system was interrupted while in recovery at(?: log time)?) ` + utcTimestampRegexp + `)`),
		secrets: []state.LogSecretKind{0},
	},
	hint: match{
		prefixes: []string{
			"This probably means that some data is corrupted and you will have to use the last backup for recovery.",
			"If this has occurred more than once some data might be corrupted and you might need to choose an earlier recovery target.",
		},
	},
}
var temporaryFile = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_TEMP_FILE_CREATED,
	primary: match{
		prefixes: []string{"temporary file: path "},
		regexp:   regexp.MustCompile(`^temporary file: path "(.+?)", size (\d+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var serverMiscCouldNotOpenUsermap = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
	primary: match{
		prefixes: []string{"could not open usermap file"},
		regexp:   regexp.MustCompile(`^could not open usermap file "(.+)": (.+)`),
		secrets:  []state.LogSecretKind{state.OpsLogSecret, state.OpsLogSecret},
	},
}
var serverMiscCouldNotLinkFile = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
	primary: match{
		prefixes: []string{"could not link file"},
		regexp:   regexp.MustCompile(`^could not link file "(.+)" to "(.+)": (.+)`),
		secrets:  []state.LogSecretKind{state.OpsLogSecret, state.OpsLogSecret, state.OpsLogSecret},
	},
}
var serverMiscUnexpectedAddr = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
	primary: match{
		prefixes: []string{"unexpected pageaddr"},
		regexp:   regexp.MustCompile(`^unexpected pageaddr \w+/\w+ in log segment \w+, offset \d+`),
		secrets:  []state.LogSecretKind{},
	},
}
var serverReload = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
	primary: match{
		prefixes: []string{"received SIGHUP, reloading configuration files"},
	},
}
var serverShutdown = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
	primary: match{
		prefixes: []string{
			"received fast shutdown request",
			"received smart shutdown request",
			"aborting any active transactions",
			"shutting down",
			"the database system is shutting down",
			"database system is shut down",
		},
	},
}
var pageVerificationFailed = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_INVALID_CHECKSUM,
	primary: match{
		prefixes: []string{"page verification failed"},
		regexp:   regexp.MustCompile(`^page verification failed, calculated checksum (\d+) but expected (\d+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var invalidChecksum = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_INVALID_CHECKSUM,
	primary: match{
		prefixes: []string{"invalid page in block"},
		regexp:   regexp.MustCompile(`^invalid page in block (\d+) of relation (\w+/\d+/\d+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var parameterCannotBeChanged = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
	primary: match{
		prefixes: []string{"parameter"},
		regexp:   regexp.MustCompile(`^parameter "([^"]+)" (changed to "([^"]+)"|cannot be changed without restarting the server)`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
}
var configFileContainsErrors = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
	primary: match{
		prefixes: []string{"configuration file"},
		regexp:   regexp.MustCompile(`^configuration file "([^"]+)" contains errors; unaffected changes were applied`),
		secrets:  []state.LogSecretKind{0},
	},
}
var workerProcessExited = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_PROCESS_EXITED,
	primary: match{
		prefixes: []string{"worker process: "},
		regexp:   regexp.MustCompile(`^worker process: (.+?) \(PID (\d+)\) (?:exited with exit code (\d+)|was terminated by signal (\d+))`),
		secrets:  []state.LogSecretKind{0, 0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Failed process was running: (.*)`),
		secrets: []state.LogSecretKind{state.StatementTextLogSecret},
	},
}
var statsCollectorTimeout = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SERVER_STATS_COLLECTOR_TIMEOUT,
	primary: match{
		prefixes: []string{
			"using stale statistics instead of current ones because stats collector is not responding",
			"pgstat wait timeout",
		},
	},
}
var standbyRestoredLogFile = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_RESTORED_WAL_FROM_ARCHIVE,
	primary: match{
		prefixes: []string{"restored log file"},
		regexp:   regexp.MustCompile(`^restored log file "([^"]+)" from archive`),
		secrets:  []state.LogSecretKind{0},
	},
}
var standbyStreamingStarted = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_STARTED_STREAMING,
	primary: match{
		prefixes: []string{"started streaming WAL", "restarted WAL streaming"},
		regexp:   regexp.MustCompile(`^(?:started streaming WAL from primary|restarted WAL streaming) at (\w+)/(\w+) on timeline (\d+)`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
}
var standbyStreamingInterrupted = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_STREAMING_INTERRUPTED,
	primary: match{
		prefixes: []string{"could not receive data from WAL stream"},
		regexp:   regexp.MustCompile(`^could not receive data from WAL stream: ([\w: ]+)`),
		secrets:  []state.LogSecretKind{state.OpsLogSecret},
	},
}
var standbyStoppedStreaming = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_STOPPED_STREAMING,
	primary: match{
		prefixes: []string{"terminating walreceiver process due to administrator command"},
	},
}
var standbyConsistentRecoveryState = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_CONSISTENT_RECOVERY_STATE,
	primary: match{
		prefixes: []string{"consistent recovery state reached at"},
		regexp:   regexp.MustCompile(`^consistent recovery state reached at (\w+)/(\w+)`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var standbyStatementCanceled = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_STATEMENT_CANCELED,
	primary: match{
		prefixes: []string{"canceling statement due to conflict with recovery"},
	},
	detail: match{
		prefixes: []string{"User query might have needed to see row versions that must be removed."},
	},
}
var regexpInvalidTimeline = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STANDBY_INVALID_TIMELINE,
	primary: match{
		prefixes: []string{"according to history file, WAL location"},
		regexp:   regexp.MustCompile(`^according to history file, WAL location .+? belongs to timeline \d+, but previous recovered WAL file came from timeline \d+`),
		secrets:  []state.LogSecretKind{},
	},
}
var uniqueConstraintViolation = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_UNIQUE_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"duplicate key value violates unique constraint"},
		regexp:   regexp.MustCompile(`^duplicate key value violates unique constraint "(.+)"`),
		secrets:  []state.LogSecretKind{0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Key \((.+)\)=\((.+)\) already exists.`),
		secrets: []state.LogSecretKind{0, state.TableDataLogSecret},
	},
}
var foreignKeyConstraintViolation1 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_FOREIGN_KEY_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"insert or update on table"},
		regexp:   regexp.MustCompile(`^insert or update on table "(.+?)" violates foreign key constraint "(.+?)"`),
		secrets:  []state.LogSecretKind{0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Key \((.+)\)=\((.+)\) is not present in table "(.+)".`),
		secrets: []state.LogSecretKind{0, state.TableDataLogSecret, 0},
	},
}
var foreignKeyConstraintViolation2 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_FOREIGN_KEY_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"update or delete on table"},
		regexp:   regexp.MustCompile(`^update or delete on table "(.+?)" violates foreign key constraint "(.+?)" on table "(.+?)"`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Key \((.+)\)=\((.+)\) is still referenced from table "(.+)".`),
		secrets: []state.LogSecretKind{0, state.TableDataLogSecret, 0},
	},
}
var nullConstraintViolation = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_NOT_NULL_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"null value in column"},
		regexp:   regexp.MustCompile(`^null value in column "(.+?)" violates not-null constraint`),
		secrets:  []state.LogSecretKind{0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Failing row contains \((.+)\).`),
		secrets: []state.LogSecretKind{state.TableDataLogSecret},
	},
}
var checkConstraintViolation1 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"new row for relation"},
		regexp:   regexp.MustCompile(`^new row for relation "(.+?)" violates check constraint "(.+?)"`),
		secrets:  []state.LogSecretKind{0, 0},
		// FIXME: Store constraint name and relation name
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Failing row contains \((.+)\).`),
		secrets: []state.LogSecretKind{state.TableDataLogSecret},
	},
}
var checkConstraintViolation2 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"check constraint"},
		regexp:   regexp.MustCompile(`^check constraint "(.+?)" is violated by some row`),
		secrets:  []state.LogSecretKind{0},
		// FIXME: Store constraint name
	},
}
var checkConstraintViolation3 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column "(.+?)" of table "(.+?)" contains values that violate the new constraint`),
		secrets:  []state.LogSecretKind{0, 0},
		// FIXME: Store relation name
	},
}
var checkConstraintViolation4 = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"value for domain"},
		regexp:   regexp.MustCompile(`^value for domain (.+?) violates check constraint "(.+?)"`),
		secrets:  []state.LogSecretKind{0, 0},
		// FIXME: Store constraint name
	},
}
var exclusionConstraintViolation = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_EXCLUSION_CONSTRAINT_VIOLATION,
	primary: match{
		prefixes: []string{"conflicting key value violates exclusion constraint"},
		regexp:   regexp.MustCompile(`^conflicting key value violates exclusion constraint "(.+?)"`),
		secrets:  []state.LogSecretKind{0},
		// FIXME: Store constraint name
	},
	detail: match{
		regexp:  regexp.MustCompile(`^Key \([^)]+\)=\((.+)\) conflicts with existing key \([^)]+\)=\((.+)\).`),
		secrets: []state.LogSecretKind{state.TableDataLogSecret, state.TableDataLogSecret},
	},
}
var syntaxError = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SYNTAX_ERROR,
	primary: match{
		prefixes: []string{"syntax error at"},
		regexp:   regexp.MustCompile(`^syntax error at (?:end of input|or near "(.+)")(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{state.ParsingErrorLogSecret},
	},
}
var columnMissingFromGroupBy = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COLUMN_MISSING_FROM_GROUP_BY,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column "([^"]+)" must appear in the GROUP BY clause or be used in an aggregate function(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
	},
}
var columnDoesNotExist = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COLUMN_DOES_NOT_EXIST,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column "([^"]+)" does not exist(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
	},
}
var columnDoesNotExistOnTable = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COLUMN_DOES_NOT_EXIST,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column "([^"]+)" of relation "([^"]+)" does not exist(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var columnReferenceAmbiguous = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COLUMN_REFERENCE_AMBIGUOUS,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column reference "([^"]+)" is ambiguous(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
	},
}
var relationDoesNotExist = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_RELATION_DOES_NOT_EXIST,
	primary: match{
		prefixes: []string{"relation"},
		regexp:   regexp.MustCompile(`^relation "([^"]+)" does not exist(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
	},
}
var functionDoesNotExist = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_FUNCTION_DOES_NOT_EXIST,
	primary: match{
		prefixes: []string{"function"},
		regexp:   regexp.MustCompile(`^function ([^"]+) does not exist(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
	},
	hint: match{
		prefixes: []string{"No function matches the given name and argument types. You might need to add explicit type casts."},
	},
}
var invalidInputSyntax = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INVALID_INPUT_SYNTAX,
	primary: match{
		prefixes: []string{"invalid input syntax for"},
		regexp:   regexp.MustCompile(`^invalid input syntax for [\w ]+(?:: "([^"]+)")?(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{state.TableDataLogSecret},
	},
	detail: match{
		prefixes: []string{"Escape sequence \""},
		regexp:   regexp.MustCompile(`^Escape sequence "(.+)" is invalid\.`),
		secrets:  []state.LogSecretKind{state.TableDataLogSecret},
	},
}
var valueTooLongForType = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_VALUE_TOO_LONG_FOR_TYPE,
	primary: match{
		prefixes: []string{"value too long for type"},
		regexp:   regexp.MustCompile(`^value too long for type ([\w ()]+)`),
		secrets:  []state.LogSecretKind{0},
	},
}
var invalidValue = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INVALID_VALUE,
	primary: match{
		prefixes: []string{"invalid value"},
		regexp:   regexp.MustCompile(`^invalid value "([^"]+)" for "([^"]+)"`),
		secrets:  []state.LogSecretKind{state.TableDataLogSecret, state.TableDataLogSecret},
	},
	detail: match{
		prefixes: []string{"Value must be an integer."},
	},
}
var malformedArrayLiteral = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_MALFORMED_ARRAY_LITERAL,
	primary: match{
		prefixes: []string{"malformed array literal"},
		regexp:   regexp.MustCompile(`^malformed array literal: "(.+)"(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{state.TableDataLogSecret},
	},
	detail: match{
		prefixes: []string{
			"Array value must start with \"{\" or dimension information.",
			"Unexpected array element.",
		},
	},
}
var subqueryMissingAlias = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_SUBQUERY_MISSING_ALIAS,
	primary: match{
		prefixes: []string{"subquery in FROM must have an alias"},
		regexp:   regexp.MustCompile(`^subquery in FROM must have an alias(?: at character \d+)?`),
	},
	hint: match{
		prefixes: []string{"For example, FROM (SELECT ...) [AS] foo."},
	},
}
var insertTargetColumnMismatch = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INSERT_TARGET_COLUMN_MISMATCH,
	primary: match{
		prefixes: []string{"INSERT has more expressions than target columns"},
		regexp:   regexp.MustCompile(`^INSERT has more expressions than target columns(?: at character \d+)?`),
	},
}
var anyAllRequiresArray = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_ANY_ALL_REQUIRES_ARRAY,
	primary: match{
		prefixes: []string{"op ANY/ALL (array) requires array on right side"},
		regexp:   regexp.MustCompile(`^op ANY/ALL \(array\) requires array on right side(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{},
	},
}
var operatorDoesNotExist = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_OPERATOR_DOES_NOT_EXIST,
	primary: match{
		prefixes: []string{"operator does not exist: "},
		regexp:   regexp.MustCompile(`^operator does not exist: (\w+) ([` + regexp.QuoteMeta("+*/<>=~!@#%^&|`?-") + `]+) (\w+)(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0, 0, 0},
	},
	hint: match{
		prefixes: []string{"No operator matches the given name and argument type(s). You might need to add explicit type casts."},
	},
}
var permissionDenied = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_PERMISSION_DENIED,
	primary: match{
		prefixes: []string{"permission denied"},
		regexp:   regexp.MustCompile(`^permission denied for (?:column|relation|table|sequence|database|function|operator|type|language|large object|schema|operator class|operator family|collation|conversion|tablespace|text search dictionary|text search configuration|foreign-data wrapper|foreign server|event trigger|extension) ([\w_-]+)(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{0},
		// FIXME: Store relation name when this is "permission denied for relation [relation name]"
	},
}
var transactionIsAborted = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_TRANSACTION_IS_ABORTED,
	primary: match{
		prefixes: []string{"current transaction is aborted, commands ignored until end of transaction block"},
	},
}
var onConflictNoConstraintMatch = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_ON_CONFLICT_NO_CONSTRAINT_MATCH,
	primary: match{
		prefixes: []string{"there is no unique or exclusion constraint matching the ON CONFLICT specification"},
	},
}
var onConflictRowAffectedTwice = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_ON_CONFLICT_ROW_AFFECTED_TWICE,
	primary: match{
		prefixes: []string{"ON CONFLICT DO UPDATE command cannot affect row a second time"},
	},
	hint: match{
		prefixes: []string{"Ensure that no rows proposed for insertion within the same command have duplicate constrained values."},
	},
}
var columnCannotBeCast = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COLUMN_CANNOT_BE_CAST,
	primary: match{
		prefixes: []string{"column"},
		regexp:   regexp.MustCompile(`^column "([^"]+)" cannot be cast to type "([^"]+)"`),
		secrets:  []state.LogSecretKind{0, 0},
	},
}
var divisionByZero = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_DIVISION_BY_ZERO,
	primary: match{
		prefixes: []string{"division by zero"},
	},
}
var cannotDrop = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_CANNOT_DROP,
	primary: match{
		prefixes: []string{"cannot drop"},
		regexp:   regexp.MustCompile(`^cannot drop ([^"]+) because other objects depend on it`),
		secrets:  []state.LogSecretKind{0},
	},
	detail: match{
		regexp:  regexp.MustCompile(`^\w+ (.+) depends on \w+ (.+)`),
		secrets: []state.LogSecretKind{0, 0},
	},
	hint: match{
		prefixes: []string{"Use DROP ... CASCADE to drop the dependent objects too."},
	},
}
var integerOutOfRange = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INTEGER_OUT_OF_RANGE,
	primary: match{
		prefixes: []string{"integer out of range"},
	},
}
var invalidRegexp = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INVALID_REGEXP,
	primary: match{
		prefixes: []string{"invalid regular expression: "},
		// Error messages from Postgres' src/include/regex/regerrs.h
		regexp: regexp.MustCompile(`^invalid regular expression: (?:` +
			`no errors detected` +
			`|failed to match` +
			`|invalid regexp \(reg version 0.8\)` +
			`|invalid collating element` +
			`|invalid character class` +
			`|invalid escape \\ sequence` +
			`|invalid backreference number` +
			`|brackets \[\] not balanced` +
			`|parentheses \(\) not balanced` +
			`|braces \{\} not balanced` +
			`|invalid repetition count\(s\)` +
			`|invalid character range` +
			`|out of memory` +
			`|quantifier operand invalid` +
			`|"cannot happen" -- you found a bug` +
			`|invalid argument to regex function` +
			`|character widths of regex and string differ` +
			`|invalid embedded option` +
			`|regular expression is too complex` +
			`|too many colors` +
			`|operation cancelled` +
			`)`),
	},
}
var paramMissing = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_PARAM_MISSING,
	primary: match{
		prefixes: []string{"there is no parameter $"},
		regexp:   regexp.MustCompile(`^there is no parameter \$\d+(?: at character \d+)?`),
	},
}
var noSuchSavepoint = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_NO_SUCH_SAVEPOINT,
	primary: match{
		prefixes: []string{"no such savepoint"},
	},
}
var unterminatedQuotedString = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_UNTERMINATED_QUOTED_STRING,
	primary: match{
		prefixes: []string{"unterminated quoted string"},
		regexp:   regexp.MustCompile(`^unterminated quoted string(?: at or near "(.+?)")?(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{state.ParsingErrorLogSecret},
	},
}
var unterminatedQuotedIdentifier = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_UNTERMINATED_QUOTED_IDENTIFIER,
	primary: match{
		prefixes: []string{"unterminated quoted identifier"},
		regexp:   regexp.MustCompile(`^unterminated quoted identifier(?: at or near "(.+?)")?(?: at character \d+)?`),
		secrets:  []state.LogSecretKind{state.ParsingErrorLogSecret},
	},
}
var invalidByteSequence = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INVALID_BYTE_SEQUENCE,
	primary: match{
		prefixes: []string{"invalid byte sequence for encoding"},
		regexp:   regexp.MustCompile(`^invalid byte sequence for encoding "([^"]+)": (.*)`),
		secrets:  []state.LogSecretKind{0, state.TableDataLogSecret},
	},
}
var couldNotSerializeRepeatableRead = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COULD_NOT_SERIALIZE_REPEATABLE_READ,
	primary: match{
		prefixes: []string{"could not serialize access due to concurrent update"},
	},
}
var couldNotSerializeSerializable = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_COULD_NOT_SERIALIZE_SERIALIZABLE,
	primary: match{
		prefixes: []string{"could not serialize access due to read/write dependencies among transactions"},
	},
	detail: match{
		prefixes: []string{
			"Reason code: Canceled on identification as a pivot, during conflict out checking.",
			"Reason code: Canceled on identification as a pivot, during conflict in checking.",
			"Reason code: Canceled on identification as a pivot, during write.",
			"Reason code: Canceled on identification as a pivot, during commit attempt.",
			"Reason code: Canceled on identification as a pivot, with conflict out to old committed transaction %u.",
			"Reason code: Canceled on commit attempt with conflict in from prepared pivot.",
			"Reason code: Canceled on conflict out to pivot %u, during read.",
			"Reason code: Canceled on conflict out to old pivot %u.",
			"Reason code: Canceled on conflict out to old pivot.",
		},
		regexp: regexp.MustCompile(`^Reason code: Canceled on (?:` +
			`identification as a pivot, during conflict out checking` +
			`|identification as a pivot, during conflict in checking` +
			`|identification as a pivot, during write` +
			`|identification as a pivot, during commit attempt` +
			`|identification as a pivot, with conflict out to old committed transaction \d+` +
			`|commit attempt with conflict in from prepared pivot` +
			`|conflict out to pivot \d+, during read` +
			`|conflict out to old pivot \d+` +
			`|conflict out to old pivot` +
			`)\.`),
	},
	hint: match{
		prefixes: []string{"The transaction might succeed if retried."},
	},
}
var inconsistentRangeBounds = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_INCONSISTENT_RANGE_BOUNDS,
	primary: match{
		prefixes: []string{"range lower bound must be less than or equal to range upper bound"},
	},
}
var statementLog = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STATEMENT_LOG,
	primary: match{
		prefixes: []string{"statement: ", "execute "},
		regexp:   regexp.MustCompile(`^(?:statement|(?:execute|execute fetch from) (?:[^:]+)(?:/(?:[^:]+))?): (.*)`),
		secrets:  []state.LogSecretKind{state.StatementTextLogSecret},
	},
	detail: match{
		regexp:  regexp.MustCompile(`(?:parameters: |, )\$\d+ = '([^']*)'|^prepare: (.+)`),
		secrets: []state.LogSecretKind{state.StatementParameterLogSecret, state.StatementTextLogSecret},
	},
}
var statementCanceledTimeout = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STATEMENT_CANCELED_TIMEOUT,
	primary: match{
		prefixes: []string{"canceling statement due to statement timeout"},
	},
}
var statementCanceledUser = analyzeGroup{
	classification: pganalyze_collector.LogLineInformation_STATEMENT_CANCELED_USER,
	primary: match{
		prefixes: []string{"canceling statement due to user request"},
	},
}
var pgaCollectorIdentify = analyzeGroup{
	primary: match{
		prefixes: []string{"pganalyze-collector-identify: "},
		regexp:   regexp.MustCompile(`^pganalyze-collector-identify: (.*)`),
		secrets:  []state.LogSecretKind{state.UnidentifiedLogSecret},
	},
}

// CONTEXT patterns (errcontext calls in Postgres)
var contextCancelAutovacuum = match{ // do_autovacuum in autovacuum.c
	regexp:  regexp.MustCompile(`automatic (?:vacuum|analyze) of table "(.+?)"`),
	secrets: []state.LogSecretKind{0},
}
var contextCancelAutovacuumDetail = match{ // vacuum_error_callback in vacuumlazy.c
	regexp:  regexp.MustCompile(`while (?:scanning|vacuuming|vacuuming index \"[^"]+\"|cleaning up index \"[^"]+\"|truncating) (?:block \d+ )?(?:of )?relation "(.+?)"(?: to \d+ blocks)?`),
	secrets: []state.LogSecretKind{0},
}
var otherContextPatterns = []match{
	{
		prefixes: []string{"COPY"},
		regexp:   regexp.MustCompile(`^COPY \w+, line \d+(?:, column \w+)?`),
	},
	{
		prefixes: []string{"PL/pgSQL function"},
		regexp:   regexp.MustCompile(`PL/pgSQL function (?:[^(]+\([^)]+\)|inline_code_block)(.*)`),
		secrets:  []state.LogSecretKind{0},
	},
	{
		prefixes: []string{"while updating tuple"},
		regexp:   regexp.MustCompile(`while updating tuple \(\d+,\d+\) in relation "([^"]+)"`),
		secrets:  []state.LogSecretKind{0},
	},
	{
		prefixes: []string{"while inserting index tuple"},
		regexp:   regexp.MustCompile(`while inserting index tuple \(\d+,\d+\) in relation "([^"]+)"`),
		secrets:  []state.LogSecretKind{0},
	},
	{
		prefixes: []string{"JSON data, line "},
		regexp:   regexp.MustCompile(`^JSON data, line (\d+): (.+)`),
		secrets:  []state.LogSecretKind{0, state.TableDataLogSecret},
	},
}

var autoVacuumIndexRegexp = regexp.MustCompile(`index "(.+?)": pages: (\d+) in total, (\d+) newly deleted, (\d+) currently deleted, (\d+) reusable,?\s*`)
var parallelWorkerProcessTextRegexp = regexp.MustCompile(`^parallel worker for PID (\d+)`)

func AnalyzeLogLines(logLinesIn []state.LogLine) (logLinesOut []state.LogLine, samples []state.PostgresQuerySample) {
	// Split log lines by backend to ensure we have the right context
	backendLogLines := make(map[int32][]state.LogLine)

	for _, logLine := range logLinesIn {
		backendLogLines[logLine.BackendPid] = append(backendLogLines[logLine.BackendPid], logLine)
	}

	for _, logLines := range backendLogLines {
		backendLogLinesOut, backendSamples := AnalyzeBackendLogLines(logLines)
		for _, logLine := range backendLogLinesOut {
			logLinesOut = append(logLinesOut, logLine)
		}
		for _, sample := range backendSamples {
			samples = append(samples, sample)
		}
	}

	return
}

func classifyAndSetDetails(logLine state.LogLine, statementLine state.LogLine, detailLine state.LogLine, contextLine state.LogLine, hintLine state.LogLine, samples []state.PostgresQuerySample) (state.LogLine, state.LogLine, state.LogLine, state.LogLine, state.LogLine, []state.PostgresQuerySample) {
	var parts []string

	// Generic handlers
	groupX := []analyzeGroup{
		connectionRejected,
		authenticationFailed,
		databaseNotAcceptingConnections,
		roleNotAllowedLogin,
		connectionClientFailedToConnect,
		connectionLostOpenTx,
		connectionLost,
		connectionLostSocketError,
		connectionTerminated,
		outOfConnections,
		tooManyConnectionsRole,
		tooManyConnectionsDatabase,
		couldNotAcceptSSLConnection,
		protocolErrorUnsupportedVersion,
		protocolErrorIncompleteMessage,
		restartpointAt,
		walInvalidRecordLength,
		walRedo,
		archiverProcessExited,
		walBaseBackupComplete,
		lockTimeout,
		statementCanceledUser,
		statementCanceledTimeout,
		serverCrashedOtherProcesses,
		serverOutOfMemory,
		serverMiscCouldNotOpenUsermap,
		serverMiscCouldNotLinkFile,
		serverMiscUnexpectedAddr,
		serverReload,
		serverShutdown,
		serverStart,
		serverStartShutdownAt,
		serverStartRecovering,
		pageVerificationFailed,
		autovacuumLauncherStarted,
		autovacuumLauncherShuttingDown,
		statsCollectorTimeout,
		standbyRestoredLogFile,
		standbyStreamingStarted,
		standbyStreamingInterrupted,
		standbyStoppedStreaming,
		standbyConsistentRecoveryState,
		standbyStatementCanceled,
		regexpInvalidTimeline,
		checkConstraintViolation1,
		checkConstraintViolation2,
		checkConstraintViolation3,
		checkConstraintViolation4,
		exclusionConstraintViolation,
		syntaxError,
		columnMissingFromGroupBy,
		columnDoesNotExist,
		columnDoesNotExistOnTable,
		columnReferenceAmbiguous,
		relationDoesNotExist,
		functionDoesNotExist,
		invalidInputSyntax,
		valueTooLongForType,
		invalidValue,
		malformedArrayLiteral,
		subqueryMissingAlias,
		insertTargetColumnMismatch,
		anyAllRequiresArray,
		operatorDoesNotExist,
		permissionDenied,
		transactionIsAborted,
		onConflictNoConstraintMatch,
		onConflictRowAffectedTwice,
		columnCannotBeCast,
		divisionByZero,
		cannotDrop,
		integerOutOfRange,
		invalidRegexp,
		paramMissing,
		noSuchSavepoint,
		unterminatedQuotedString,
		unterminatedQuotedIdentifier,
		invalidByteSequence,
		couldNotSerializeRepeatableRead,
		couldNotSerializeSerializable,
		inconsistentRangeBounds,
	}
	for _, m := range groupX {
		if matchesPrefix(logLine, m.primary.prefixes) {
			logLine, parts = matchLogLine(logLine, m.primary)
			if parts != nil {
				logLine.Classification = m.classification
				detailLine, _ = matchLogLine(detailLine, m.detail)
				hintLine, _ = matchLogLine(hintLine, m.hint)
				contextLine = matchOtherContextLogLine(contextLine)
				return logLine, statementLine, detailLine, contextLine, hintLine, samples
			}
		}
	}

	// Connects/Disconnects
	if matchesPrefix(logLine, connectionReceived.primary.prefixes) {
		logLine.Classification = connectionReceived.classification
		logLine, parts = matchLogLine(logLine, connectionReceived.primary)
		if len(parts) == 3 {
			logLine.Details = map[string]interface{}{"host": parts[1]}
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}
	if matchesPrefix(logLine, connectionAuthorized.primary.prefixes) {
		logLine.Classification = connectionAuthorized.classification
		logLine, parts = matchLogLine(logLine, connectionAuthorized.primary)
		if len(parts) == 5 && parts[4] != "" {
			logLine.Details = map[string]interface{}{"ssl_protocol": parts[4]}
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}
	if matchesPrefix(logLine, disconnection.primary.prefixes) {
		logLine.Classification = disconnection.classification
		logLine, parts = matchLogLine(logLine, disconnection.primary)
		if len(parts) == 5 {
			timeSecs, _ := strconv.ParseFloat(parts[3], 64)
			timeMinutes, _ := strconv.ParseFloat(parts[2], 64)
			timeHours, _ := strconv.ParseFloat(parts[1], 64)
			logLine.Details = map[string]interface{}{"session_time_secs": timeSecs + timeMinutes*60 + timeHours*3600}
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}

	// Checkpointer
	if matchesPrefix(logLine, checkpointStarting.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, checkpointStarting.primary)
		if len(parts) == 3 {
			if parts[1] == "checkpoint" {
				logLine.Classification = pganalyze_collector.LogLineInformation_CHECKPOINT_STARTING
			} else if parts[1] == "restartpoint" {
				logLine.Classification = pganalyze_collector.LogLineInformation_RESTARTPOINT_STARTING
			}

			logLine.Details = map[string]interface{}{"reason": parts[2]}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}

		logLine, parts = matchLogLine(logLine, checkpointComplete.primary)
		if len(parts) == 15 {
			if parts[1] == "checkpoint" {
				logLine.Classification = pganalyze_collector.LogLineInformation_CHECKPOINT_COMPLETE
			} else if parts[1] == "restartpoint" {
				logLine.Classification = pganalyze_collector.LogLineInformation_RESTARTPOINT_COMPLETE
			}

			bufsWritten, _ := strconv.ParseInt(parts[2], 10, 64)
			bufsWrittenPct, _ := strconv.ParseFloat(parts[3], 64)
			segsAdded, _ := strconv.ParseInt(parts[4], 10, 64)
			segsRemoved, _ := strconv.ParseInt(parts[5], 10, 64)
			segsRecycled, _ := strconv.ParseInt(parts[6], 10, 64)
			writeSecs, _ := strconv.ParseFloat(parts[7], 64)
			syncSecs, _ := strconv.ParseFloat(parts[8], 64)
			totalSecs, _ := strconv.ParseFloat(parts[9], 64)
			syncRels, _ := strconv.ParseInt(parts[10], 10, 64)
			longestSecs, _ := strconv.ParseFloat(parts[11], 64)
			averageSecs, _ := strconv.ParseFloat(parts[12], 64)
			distanceKb, _ := strconv.ParseInt(parts[13], 10, 64)
			estimateKb, _ := strconv.ParseInt(parts[14], 10, 64)
			logLine.Details = map[string]interface{}{
				"bufs_written": bufsWritten, "segs_added": segsAdded,
				"segs_removed": segsRemoved, "segs_recycled": segsRecycled,
				"sync_rels":        syncRels,
				"bufs_written_pct": bufsWrittenPct, "write_secs": writeSecs,
				"sync_secs": syncSecs, "total_secs": totalSecs,
				"longest_secs": longestSecs, "average_secs": averageSecs,
				"distance_kb": distanceKb,
				"estimate_kb": estimateKb,
			}

			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, checkpointsTooFrequent.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, checkpointsTooFrequent.primary)
		if len(parts) == 2 {
			logLine.Classification = checkpointsTooFrequent.classification
			elapsedSecs, _ := strconv.ParseFloat(parts[1], 64)
			logLine.Details = map[string]interface{}{
				"elapsed_secs": elapsedSecs,
			}
			hintLine, _ = matchLogLine(hintLine, checkpointsTooFrequent.hint)
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}

	// WAL/Archiving
	if matchesPrefix(logLine, walRedoLastTx.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, walRedoLastTx.primary)
		if len(parts) == 2 {
			logLine.Classification = walRedoLastTx.classification
			logLine.Details = map[string]interface{}{
				"last_transaction": parts[1],
			}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, archiveCommandFailed.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, archiveCommandFailed.primary)
		if len(parts) == 4 {
			logLine.Classification = archiveCommandFailed.classification
			logLine.Details = map[string]interface{}{}
			if parts[1] != "" {
				exitCode, _ := strconv.ParseInt(parts[1], 10, 32)
				logLine.Details["exit_code"] = exitCode
			}
			if parts[2] != "" {
				signal, _ := strconv.ParseInt(parts[2], 10, 32)
				logLine.Details["signal"] = signal
			}
			detailLine, _ = matchLogLine(detailLine, archiveCommandFailed.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}

	// Lock waits
	if matchesPrefix(logLine, lockAcquired.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, lockAcquired.primary)
		if len(parts) == 5 {
			logLine.Classification = pganalyze_collector.LogLineInformation_LOCK_ACQUIRED
			afterMs, _ := strconv.ParseFloat(parts[4], 64)
			logLine.Details = map[string]interface{}{
				"lock_mode": parts[1],
				"lock_type": parts[2],
				"after_ms":  afterMs,
			}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, lockWait.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, lockWait.primary)
		if len(parts) == 5 {
			if parts[1] == "still waiting" {
				logLine.Classification = pganalyze_collector.LogLineInformation_LOCK_WAITING
			} else if parts[1] == "avoided deadlock" {
				logLine.Classification = pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_AVOIDED
			} else if parts[1] == "detected deadlock while waiting" {
				logLine.Classification = pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_DETECTED
			}
			lockType := parts[3]
			// Match lock types to names in pg_locks.locktype
			if lockType == "extension" {
				lockType = "extend"
			} else if lockType == "transaction" {
				lockType = "transactionid"
			} else if lockType == "virtual" {
				lockType = "virtualxid"
			}
			afterMs, _ := strconv.ParseFloat(parts[4], 64)
			logLine.Details = map[string]interface{}{"lock_mode": parts[2], "lock_type": lockType, "after_ms": afterMs}
			if detailLine.Content != "" {
				detailLine, parts = matchLogLine(detailLine, lockWait.detail)
				if len(parts) == 3 {
					logLine.RelatedPids = []int32{}
					lockHolders := []int32{}
					for _, s := range strings.Split(parts[1], ", ") {
						i, _ := strconv.ParseInt(s, 10, 32)
						lockHolders = append(lockHolders, int32(i))
						logLine.RelatedPids = append(logLine.RelatedPids, int32(i))
					}
					lockWaiters := []int32{}
					for _, s := range strings.Split(parts[2], ", ") {
						i, _ := strconv.ParseInt(s, 10, 32)
						lockWaiters = append(lockWaiters, int32(i))
						logLine.RelatedPids = append(logLine.RelatedPids, int32(i))
					}
					logLine.Details["lock_holders"] = lockHolders
					logLine.Details["lock_waiters"] = lockWaiters
				}
			}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, deadlock.primary.prefixes) {
		logLine, _ = matchLogLine(logLine, deadlock.primary)
		logLine.Classification = deadlock.classification
		logLine.RelatedPids = []int32{}
		var allParts [][]string
		detailLine, allParts = matchLogLineAll(detailLine, deadlock.detail)
		for _, parts = range allParts {
			pid, _ := strconv.ParseInt(parts[1], 10, 32)
			logLine.RelatedPids = append(logLine.RelatedPids, int32(pid))
		}
		hintLine, _ = matchLogLineAll(hintLine, deadlock.hint)
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}
	// Statement duration (log_min_duration_statement output) and auto_explain
	if matchesPrefix(logLine, autoExplain.primary.prefixes) {
		// auto_explain needs to come before statement duration since its a subset of that regexp
		logLine, parts = matchLogLine(logLine, autoExplain.primary)
		if len(parts) == 2 {
			logLine.Classification = autoExplain.classification

			explainText := strings.TrimSpace(logLine.Content[len(parts[0]):len(logLine.Content)])
			sample, err := querysample.TransformAutoExplainToQuerySample(logLine, explainText, parts[1])
			if err != nil {
				logLine.Details = map[string]interface{}{"query_sample_error": fmt.Sprintf("%s", err)}
			} else {
				samples = append(samples, sample)
				logLine.Query = sample.Query
			}

			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, duration.primary.prefixes) {
		logLine.Classification = duration.classification
		logLine, parts = matchLogLine(logLine, duration.primary)
		if util.WasTruncated(logLine.Content) {
			logLine.Details = map[string]interface{}{"truncated": true}
		} else if len(parts) == 5 {
			var parameterParts [][]string
			if strings.HasPrefix(detailLine.Content, "parameters: ") {
				detailLine, parameterParts = matchLogLineAll(detailLine, duration.detail)
			}
			queryText := logLine.Content[len(parts[0]):len(logLine.Content)]

			sample, ok := querysample.TransformLogMinDurationStatementToQuerySample(logLine, queryText, parts[1], parts[2], parameterParts)
			if ok {
				samples = append(samples, sample)
				logLine.Query = sample.Query
			}
		}

		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}

	// Statement log (log_statement output)
	if matchesPrefix(logLine, statementLog.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, statementLog.primary)
		if len(parts) == 2 {
			logLine.Classification = statementLog.classification
			logLine.Query = strings.TrimSpace(parts[1])
			detailLine, _ = matchLogLineAll(detailLine, statementLog.detail)
		}
	}

	// Autovacuum
	if matchesPrefix(logLine, autovacuumCancel.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, autovacuumCancel.primary)
		if parts != nil {
			logLine.Classification = autovacuumCancel.classification
			contextLine, parts = matchLogLine(contextLine, contextCancelAutovacuum)
			if len(parts) == 2 {
				subParts := strings.SplitN(parts[1], ".", 3)
				logLine.Database = subParts[0]
				if len(subParts) >= 2 {
					logLine.SchemaName = subParts[1]
				}
				if len(subParts) >= 3 {
					logLine.RelationName = subParts[2]
				}
			} else {
				contextLine, parts = matchLogLine(contextLine, contextCancelAutovacuumDetail)
				if len(parts) == 2 {
					subParts := strings.SplitN(parts[1], ".", 2)
					if len(subParts) == 2 {
						logLine.SchemaName = subParts[0]
						logLine.RelationName = subParts[1]
					}
				}
			}
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, skippingVacuum.primary.prefixes) {
		logLine.Classification = skippingVacuum.classification
		logLine, parts = matchLogLine(logLine, skippingVacuum.primary)
		if len(parts) == 2 {
			// Unfortunately Postgres doesn't log a schema here, so we need to store this
			// outside of the usual relation reference (which requires a schema)
			logLine.Details = map[string]interface{}{"relation_name": parts[1]}
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}
	if matchesPrefix(logLine, skippingAnalyze.primary.prefixes) {
		logLine.Classification = skippingAnalyze.classification
		logLine, parts = matchLogLine(logLine, skippingAnalyze.primary)
		if len(parts) == 2 {
			// Unfortunately Postgres doesn't log a schema here, so we need to store this
			// outside of the usual relation reference (which requires a schema)
			logLine.Details = map[string]interface{}{"relation_name": parts[1]}
		}
		contextLine = matchOtherContextLogLine(contextLine)
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}
	if matchesPrefix(logLine, wraparoundWarning.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, wraparoundWarning.primary)
		if len(parts) == 5 {
			logLine.Classification = wraparoundWarning.classification
			remainingXids, _ := strconv.ParseInt(parts[4], 10, 64)
			logLine.Details = map[string]interface{}{"remaining_xids": remainingXids}
			if parts[2] != "" {
				databaseOid, _ := strconv.ParseInt(parts[2], 10, 64)
				logLine.Details["database_oid"] = databaseOid
			}
			if parts[3] != "" {
				logLine.Details["database_name"] = parts[3]
			}
			hintLine, _ = matchLogLine(hintLine, wraparoundWarning.hint)
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, wraparoundError.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, wraparoundError.primary)
		if len(parts) == 4 {
			logLine.Classification = wraparoundError.classification
			if parts[2] != "" {
				databaseOid, _ := strconv.ParseInt(parts[2], 10, 64)
				logLine.Details = map[string]interface{}{"database_oid": databaseOid}
			}
			if parts[3] != "" {
				logLine.Details = map[string]interface{}{"database_name": parts[3]}
			}
			hintLine, _ = matchLogLine(hintLine, wraparoundError.hint)
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, autoVacuum.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, autoVacuum.primary)
		if len(parts) == 46 {
			var readRatePart, writeRatePart, kernelPart, userPart, elapsedPart string

			aggressiveVacuum := parts[1] == "aggressive "

			// Part 2 (anti-wraparound) is only present on Postgres 12+, and dealt with at the end

			logLine.Classification = autoVacuum.classification
			subParts := strings.SplitN(parts[3], ".", 3)
			logLine.Database = subParts[0]
			if len(subParts) >= 2 {
				logLine.SchemaName = subParts[1]
			}
			if len(subParts) >= 3 {
				logLine.RelationName = subParts[2]
			}

			numIndexScans, _ := strconv.ParseInt(parts[4], 10, 64)
			pagesRemoved, _ := strconv.ParseInt(parts[5], 10, 64)
			relPages, _ := strconv.ParseInt(parts[6], 10, 64)

			// Parts 7 to 10 (scanned, pinskipped and frozenskipped pages) are version dependent and dealt with later

			tuplesDeleted, _ := strconv.ParseInt(parts[11], 10, 64)
			newRelTuples, _ := strconv.ParseInt(parts[12], 10, 64)

			// In Postgres 15+ the internal name changed from "new_dead_tuples" to "recently_dead_tuples",
			// to distinguish it from the added "missed_dead_tuples" - we're keeping the old name for compatibility
			newDeadTuples, _ := strconv.ParseInt(parts[13], 10, 64)

			// Parts 14 to 22 (missed dead tuples, OldestXmin, frozenxid and minmxid advancement) are version dependent and dealt with later
			// Parts 23 to 29 (Index scan info, I/O read/write timings) are only present on Postgres 14+ and dealt with later

			if parts[30] != "" { // Postgres 14+, with I/O information before buffers
				readRatePart = parts[30]
				writeRatePart = parts[31]
			} else { // Postgres 13 and older
				readRatePart = parts[35]
				writeRatePart = parts[36]
			}
			readRateMb, _ := strconv.ParseFloat(readRatePart, 64)
			writeRateMb, _ := strconv.ParseFloat(writeRatePart, 64)

			vacuumPageHit, _ := strconv.ParseInt(parts[32], 10, 64)
			vacuumPageMiss, _ := strconv.ParseInt(parts[33], 10, 64)
			vacuumPageDirty, _ := strconv.ParseInt(parts[34], 10, 64)

			// Parts 37 to 39 (WAL Usage) are only present on Postgres 13+ and dealt with later

			if parts[40] != "" {
				kernelPart = parts[40]
				userPart = parts[41]
				elapsedPart = parts[42]
			} else {
				userPart = parts[43]
				kernelPart = parts[44]
				elapsedPart = parts[45]
			}
			rusageKernelMode, _ := strconv.ParseFloat(kernelPart, 64)
			rusageUserMode, _ := strconv.ParseFloat(userPart, 64)
			rusageElapsed, _ := strconv.ParseFloat(elapsedPart, 64)

			logLine.Details = map[string]interface{}{
				"aggressive":      aggressiveVacuum,
				"num_index_scans": numIndexScans, "pages_removed": pagesRemoved,
				"rel_pages": relPages, "tuples_deleted": tuplesDeleted,
				"new_rel_tuples": newRelTuples, "new_dead_tuples": newDeadTuples,
				"vacuum_page_hit": vacuumPageHit, "vacuum_page_miss": vacuumPageMiss,
				"vacuum_page_dirty": vacuumPageDirty, "read_rate_mb": readRateMb,
				"write_rate_mb": writeRateMb, "rusage_kernel": rusageKernelMode,
				"rusage_user": rusageUserMode, "elapsed_secs": rusageElapsed,
			}
			// List anti-wraparound status either if the message indicates that it is, or if
			// our Postgres version is new enough (13+) as determined by the presence of WAL
			// record information (parts[37])
			//
			// Note that Postgres 12 is the odd one out, because it already had anti-wraparound
			// status displayed, but we have no way to distinguish it from versions that didn't
			// have it - there, only include the case when the vacuum indeed is a anti-wraparound
			// vacuum.
			if parts[2] != "" || parts[37] != "" {
				antiWraparound := parts[2] == "to prevent wraparound "
				logLine.Details["anti_wraparound"] = antiWraparound
			}
			if parts[9] != "" { // Postgres 15+, with scanned pages, but no pinskipped/frozenskipped counter
				scannedPages, _ := strconv.ParseInt(parts[9], 10, 64)
				scannedPagesPercent, _ := strconv.ParseFloat(parts[10], 64)
				logLine.Details["scanned_pages"] = scannedPages
				logLine.Details["scanned_pages_percent"] = scannedPagesPercent
			} else { // Postgres 14 and older
				pinskippedPages, _ := strconv.ParseInt(parts[7], 10, 64)
				frozenskippedPages, _ := strconv.ParseInt(parts[8], 10, 64)
				logLine.Details["pinskipped_pages"] = pinskippedPages
				logLine.Details["frozenskipped_pages"] = frozenskippedPages
			}
			if parts[14] != "" { // Postgres 10 to 14
				oldestXmin, _ := strconv.ParseInt(parts[14], 10, 64)
				logLine.Details["oldest_xmin"] = oldestXmin
			} else if parts[17] != "" { // Postgres 15+
				oldestXmin, _ := strconv.ParseInt(parts[17], 10, 64)
				oldestXminAge, _ := strconv.ParseInt(parts[18], 10, 64)
				logLine.Details["oldest_xmin"] = oldestXmin
				logLine.Details["oldest_xmin_age"] = oldestXminAge
			}
			if parts[15] != "" { // Postgres 15+, if dead tuples were skipped due to cleanup lock contention
				missedDeadTuples, _ := strconv.ParseInt(parts[15], 10, 64)
				missedDeadPages, _ := strconv.ParseInt(parts[16], 10, 64)
				logLine.Details["missed_dead_tuples"] = missedDeadTuples
				logLine.Details["missed_dead_pages"] = missedDeadPages
			}
			if parts[19] != "" { // Postgres 15+, if frozenxid was updated
				newRelfrozenXid, _ := strconv.ParseInt(parts[19], 10, 64)
				newRelfrozenXidDiff, _ := strconv.ParseInt(parts[20], 10, 64)
				logLine.Details["new_relfrozenxid"] = newRelfrozenXid
				logLine.Details["new_relfrozenxid_diff"] = newRelfrozenXidDiff
			}
			if parts[21] != "" { // Postgres 15+, if minmxid was updated
				newRelminMxid, _ := strconv.ParseInt(parts[21], 10, 64)
				newRelminMxidDiff, _ := strconv.ParseInt(parts[22], 10, 64)
				logLine.Details["new_relminmxid"] = newRelminMxid
				logLine.Details["new_relminmxid_diff"] = newRelminMxidDiff
			}
			if parts[23] != "" {
				lpdeadItemPages, _ := strconv.ParseInt(parts[24], 10, 64)
				lpdeadItemPagePercent, _ := strconv.ParseFloat(parts[25], 64)
				lpdeadItems, _ := strconv.ParseInt(parts[26], 10, 64)
				logLine.Details["lpdead_index_scan"] = parts[23] // not needed / needed / bypassed / bypassed by failsafe
				logLine.Details["lpdead_item_pages"] = lpdeadItemPages
				logLine.Details["lpdead_item_page_percent"] = lpdeadItemPagePercent
				logLine.Details["lpdead_items"] = lpdeadItems
			}
			if parts[27] != "" {
				indexParts := autoVacuumIndexRegexp.FindAllStringSubmatch(parts[27], -1)
				index_vacuums := make(map[string]interface{})
				for _, p := range indexParts {
					numPages, _ := strconv.ParseInt(p[2], 10, 64)
					pagesNewlyDeleted, _ := strconv.ParseInt(p[3], 10, 64)
					pagesDeleted, _ := strconv.ParseInt(p[4], 10, 64)
					pagesFree, _ := strconv.ParseInt(p[5], 10, 64)
					index_vacuums[p[1]] = map[string]interface{}{
						"num_pages":           numPages,
						"pages_newly_deleted": pagesNewlyDeleted,
						"pages_deleted":       pagesDeleted,
						"pages_free":          pagesFree,
					}
				}
				logLine.Details["index_vacuums"] = index_vacuums
			}
			if parts[28] != "" {
				blkReadTime, _ := strconv.ParseFloat(parts[28], 64)
				blkWriteTime, _ := strconv.ParseFloat(parts[29], 64)
				logLine.Details["blk_read_time"] = blkReadTime
				logLine.Details["blk_write_time"] = blkWriteTime
			}
			if parts[37] != "" {
				walRecords, _ := strconv.ParseInt(parts[37], 10, 64)
				walFpi, _ := strconv.ParseInt(parts[38], 10, 64)
				walBytes, _ := strconv.ParseInt(parts[39], 10, 64)
				logLine.Details["wal_records"] = walRecords
				logLine.Details["wal_fpi"] = walFpi
				logLine.Details["wal_bytes"] = walBytes
			}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, autoAnalyze.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, autoAnalyze.primary)
		if len(parts) == 15 {
			var kernelPart, userPart, elapsedPart string
			logLine.Classification = autoAnalyze.classification
			subParts := strings.SplitN(parts[1], ".", 3)
			logLine.Database = subParts[0]
			if len(subParts) >= 2 {
				logLine.SchemaName = subParts[1]
			}
			if len(subParts) >= 3 {
				logLine.RelationName = subParts[2]
			}

			// Parts 2 to 8 (I/O and buffers information) are only present on Postgres 14+ and dealt with at the end

			if parts[9] != "" {
				kernelPart = parts[9]
				userPart = parts[10]
				elapsedPart = parts[11]
			} else {
				userPart = parts[12]
				kernelPart = parts[13]
				elapsedPart = parts[14]
			}
			rusageKernelMode, _ := strconv.ParseFloat(kernelPart, 64)
			rusageUserMode, _ := strconv.ParseFloat(userPart, 64)
			rusageElapsed, _ := strconv.ParseFloat(elapsedPart, 64)
			logLine.Details = map[string]interface{}{
				"rusage_kernel": rusageKernelMode, "rusage_user": rusageUserMode,
				"elapsed_secs": rusageElapsed,
			}
			if parts[2] != "" {
				blkReadTime, _ := strconv.ParseFloat(parts[2], 64)
				blkWriteTime, _ := strconv.ParseFloat(parts[3], 64)
				readRateMb, _ := strconv.ParseFloat(parts[4], 64)
				writeRateMb, _ := strconv.ParseFloat(parts[5], 64)
				analyzePageHit, _ := strconv.ParseInt(parts[6], 10, 64)
				analyzePageMiss, _ := strconv.ParseInt(parts[7], 10, 64)
				analyzePageDirty, _ := strconv.ParseInt(parts[8], 10, 64)
				logLine.Details["blk_read_time"] = blkReadTime
				logLine.Details["blk_write_time"] = blkWriteTime
				logLine.Details["read_rate_mb"] = readRateMb
				logLine.Details["write_rate_mb"] = writeRateMb
				logLine.Details["analyze_page_hit"] = analyzePageHit
				logLine.Details["analyze_page_miss"] = analyzePageMiss
				logLine.Details["analyze_page_dirty"] = analyzePageDirty
			}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}

	// Server events
	if matchesPrefix(logLine, serverCrashed.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, serverCrashed.primary)
		if len(parts) == 4 {
			logLine.Classification = serverCrashed.classification
			detailLine, _ = matchLogLine(detailLine, serverCrashed.detail)
			processPid, _ := strconv.ParseInt(parts[1], 10, 32)
			signal, _ := strconv.ParseInt(parts[2], 10, 32)
			logLine.Details = map[string]interface{}{
				"process_type": "server process",
				"process_pid":  processPid,
				"signal":       signal,
			}
			logLine.RelatedPids = []int32{int32(processPid)}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, serverOutOfMemoryCrash.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, serverOutOfMemoryCrash.primary)
		if len(parts) == 4 {
			logLine.Classification = serverOutOfMemoryCrash.classification
			processPid, _ := strconv.ParseInt(parts[1], 10, 32)
			signal, _ := strconv.ParseInt(parts[2], 10, 32)
			logLine.Details = map[string]interface{}{
				"process_type": "server process",
				"process_pid":  processPid,
				"signal":       signal,
			}
			logLine.RelatedPids = []int32{int32(processPid)}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, invalidChecksum.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, invalidChecksum.primary)
		if len(parts) == 3 {
			logLine.Classification = invalidChecksum.classification
			blockNo, _ := strconv.ParseInt(parts[1], 10, 64)
			logLine.Details = map[string]interface{}{"block": blockNo, "file": parts[2]}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, temporaryFile.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, temporaryFile.primary)
		if len(parts) == 3 {
			logLine.Classification = temporaryFile.classification
			size, _ := strconv.ParseInt(parts[2], 10, 64)
			logLine.Details = map[string]interface{}{"size": size, "file": parts[1]}
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, parameterCannotBeChanged.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, parameterCannotBeChanged.primary)
		if len(parts) == 4 {
			logLine.Classification = parameterCannotBeChanged.classification
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, configFileContainsErrors.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, configFileContainsErrors.primary)
		if len(parts) == 2 {
			logLine.Classification = configFileContainsErrors.classification
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, workerProcessExited.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, workerProcessExited.primary)
		if len(parts) == 5 {
			logLine.Classification = workerProcessExited.classification
			processPid, _ := strconv.ParseInt(parts[2], 10, 32)
			logLine.RelatedPids = []int32{int32(processPid)}
			logLine.Details = map[string]interface{}{
				"process_type": parts[1],
				"process_pid":  processPid,
			}

			if parts[3] != "" {
				exitCode, _ := strconv.ParseInt(parts[3], 10, 32)
				logLine.Details["exit_code"] = exitCode
			}
			if parts[4] != "" {
				signal, _ := strconv.ParseInt(parts[4], 10, 32)
				logLine.Details["signal"] = signal
			}

			if strings.HasPrefix(parts[1], "parallel worker for PID") {
				textParts := parallelWorkerProcessTextRegexp.FindStringSubmatch(parts[1])
				if len(textParts) == 2 {
					parentPid, _ := strconv.ParseInt(textParts[1], 10, 32)
					logLine.Details["process_type"] = "parallel worker"
					logLine.Details["parent_pid"] = parentPid
					logLine.RelatedPids = append(logLine.RelatedPids, int32(parentPid))
				}
			}
			detailLine, _ = matchLogLine(detailLine, workerProcessExited.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}

	// Constraint violations
	if matchesPrefix(logLine, uniqueConstraintViolation.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, uniqueConstraintViolation.primary)
		if len(parts) == 2 {
			logLine.Classification = uniqueConstraintViolation.classification
			detailLine, _ = matchLogLine(detailLine, uniqueConstraintViolation.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			// FIXME: Store constraint name
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, foreignKeyConstraintViolation1.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, foreignKeyConstraintViolation1.primary)
		if len(parts) == 3 {
			logLine.Classification = foreignKeyConstraintViolation1.classification
			detailLine, _ = matchLogLine(detailLine, foreignKeyConstraintViolation1.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			// FIXME: Store constraint name and relation name
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, foreignKeyConstraintViolation2.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, foreignKeyConstraintViolation2.primary)
		if len(parts) == 4 {
			logLine.Classification = foreignKeyConstraintViolation2.classification
			detailLine, _ = matchLogLine(detailLine, foreignKeyConstraintViolation2.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			// FIXME: Store constraint name and both relation names
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}
	if matchesPrefix(logLine, nullConstraintViolation.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, nullConstraintViolation.primary)
		if len(parts) == 2 {
			logLine.Classification = nullConstraintViolation.classification
			detailLine, _ = matchLogLine(detailLine, nullConstraintViolation.detail)
			contextLine = matchOtherContextLogLine(contextLine)
			return logLine, statementLine, detailLine, contextLine, hintLine, samples
		}
	}

	// pganalyze-collector-identify
	if matchesPrefix(logLine, pgaCollectorIdentify.primary.prefixes) {
		logLine, parts = matchLogLine(logLine, pgaCollectorIdentify.primary)
		if len(parts) == 2 {
			logLine.Classification = pganalyze_collector.LogLineInformation_PGA_COLLECTOR_IDENTIFY
			logLine.Details = map[string]interface{}{
				"config_section": strings.TrimSpace(parts[1]),
			}
			contextLine = matchOtherContextLogLine(contextLine)
		}
		return logLine, statementLine, detailLine, contextLine, hintLine, samples
	}

	return logLine, statementLine, detailLine, contextLine, hintLine, samples
}

func matchLogLineCommon(logLine state.LogLine, m match, matchAll bool) (state.LogLine, [][]string) {
	if logLine.Content == "" {
		return logLine, nil
	}

	if m.regexp == nil {
		for _, prefix := range m.prefixes {
			if strings.HasPrefix(logLine.Content, prefix) {
				logLine.ReviewedForSecrets = true
				if len(prefix) < len(logLine.Content) {
					markerByteStart := len(prefix)
					markerByteEnd := len(logLine.Content)

					// Avoid including trailing new lines in the secret markers
					if logLine.Content[len(logLine.Content)-1] == '\n' {
						markerByteEnd--
					}
					if markerByteEnd-markerByteStart > 0 {
						logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
							ByteStart: markerByteStart,
							ByteEnd:   markerByteEnd,
							Kind:      state.UnidentifiedLogSecret,
						})
					}
				}
				return logLine, [][]string{[]string{prefix}}
			}
		}
		return logLine, nil
	}

	var p [][]int
	if matchAll {
		p = m.regexp.FindAllStringSubmatchIndex(logLine.Content, -1)
	} else {
		pp := m.regexp.FindStringSubmatchIndex(logLine.Content)
		if pp != nil {
			p = [][]int{pp}
		}
	}
	if p == nil {
		return logLine, nil
	}

	logLine.ReviewedForSecrets = true
	parts := make([][]string, len(p))
	for x := 0; x < len(p); x++ {
		parts[x] = make([]string, 1+m.regexp.NumSubexp())
		for i := range parts[x] {
			if 2*i < len(p[x]) && p[x][2*i] >= 0 {
				parts[x][i] = logLine.Content[p[x][2*i]:p[x][2*i+1]]
			}
		}

		if x == 0 && p[x][0] > 0 {
			logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
				ByteStart: 0,
				ByteEnd:   p[x][0],
				Kind:      state.UnidentifiedLogSecret,
			})
		}
		for idx := 0; idx < m.regexp.NumSubexp(); idx++ {
			start := p[x][2*(idx+1)]
			end := p[x][2*(idx+1)+1]
			if start < 0 {
				continue
			}
			if idx >= len(m.secrets) {
				logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
					ByteStart: start,
					ByteEnd:   end,
					Kind:      state.UnidentifiedLogSecret})
				continue
			}

			kind := m.secrets[idx]
			if kind != 0 {
				logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
					ByteStart: start,
					ByteEnd:   end,
					Kind:      kind})
			}
		}
		if x > 0 && p[x-1][1] < p[x][0]-1 {
			logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
				ByteStart: p[x-1][1],
				ByteEnd:   p[x][0],
				Kind:      state.UnidentifiedLogSecret,
			})
		}
	}
	if p[len(p)-1][1] < len(logLine.Content)-1 {
		var kind state.LogSecretKind
		if m.remainderKind != 0 {
			kind = m.remainderKind
		} else {
			kind = state.UnidentifiedLogSecret
		}
		markerByteStart := p[len(p)-1][1]
		markerByteEnd := len(logLine.Content)

		// Avoid including trailing new lines in the secret markers
		if logLine.Content[len(logLine.Content)-1] == '\n' {
			markerByteEnd--
		}

		if markerByteEnd-markerByteStart > 0 {
			logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
				ByteStart: markerByteStart,
				ByteEnd:   markerByteEnd,
				Kind:      kind,
			})
		}
	}
	return logLine, parts
}

func matchesPrefix(logLine state.LogLine, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(logLine.Content, prefix) {
			return true
		}
	}
	return false
}

func matchLogLine(logLine state.LogLine, m match) (state.LogLine, []string) {
	logLine, parts := matchLogLineCommon(logLine, m, false)
	if parts != nil {
		return logLine, parts[0]
	}
	return logLine, nil
}

func matchLogLineAll(logLine state.LogLine, m match) (state.LogLine, [][]string) {
	return matchLogLineCommon(logLine, m, true)
}

func matchOtherContextLogLine(logLine state.LogLine) state.LogLine {
	if logLine.Content == "" {
		return logLine
	}
	for _, match := range otherContextPatterns {
		logLine, parts := matchLogLine(logLine, match)
		if parts != nil {
			return logLine
		}
	}
	return logLine
}

func markLineAsSecret(logLine state.LogLine, markerKind state.LogSecretKind) state.LogLine {
	logLine.ReviewedForSecrets = true
	logLine.SecretMarkers = append(logLine.SecretMarkers, state.LogSecretMarker{
		ByteStart: 0,
		ByteEnd:   len(logLine.Content),
		Kind:      markerKind})
	return logLine
}

func AnalyzeBackendLogLines(logLines []state.LogLine) (logLinesOut []state.LogLine, samples []state.PostgresQuerySample) {
	additionalLines := 0

	for idx, logLine := range logLines {
		if additionalLines > 0 {
			logLinesOut = append(logLinesOut, logLine)
			additionalLines--
			continue
		}

		// Look up to 3 lines in the future to find context for this line
		var statementLine state.LogLine
		var statementLineIdx int
		var detailLine state.LogLine
		var detailLineIdx int
		var contextLine state.LogLine
		var contextLineIdx int
		var hintLine state.LogLine
		var hintLineIdx int

		lowerBound := int(math.Min(float64(len(logLines)), float64(idx+1)))
		upperBound := int(math.Min(float64(len(logLines)), float64(idx+5)))
		for idx, futureLine := range logLines[lowerBound:upperBound] {
			if futureLine.LogLevel == pganalyze_collector.LogLineInformation_STATEMENT || futureLine.LogLevel == pganalyze_collector.LogLineInformation_DETAIL ||
				futureLine.LogLevel == pganalyze_collector.LogLineInformation_HINT || futureLine.LogLevel == pganalyze_collector.LogLineInformation_CONTEXT ||
				futureLine.LogLevel == pganalyze_collector.LogLineInformation_QUERY {
				if futureLine.LogLevel == pganalyze_collector.LogLineInformation_STATEMENT && !util.WasTruncated(futureLine.Content) {
					logLine.Query = futureLine.Content
					statementLine = futureLine
					statementLine.ParentUUID = logLine.UUID
					statementLineIdx = lowerBound + idx
					// Ensure STATEMENT line is consistently marked as statement text log secret
					statementLine = markLineAsSecret(statementLine, state.StatementTextLogSecret)
				} else if futureLine.LogLevel == pganalyze_collector.LogLineInformation_DETAIL {
					detailLine = futureLine
					detailLine.ParentUUID = logLine.UUID
					detailLineIdx = lowerBound + idx
				} else if futureLine.LogLevel == pganalyze_collector.LogLineInformation_CONTEXT {
					contextLine = futureLine
					contextLine.ParentUUID = logLine.UUID
					contextLineIdx = lowerBound + idx
				} else if futureLine.LogLevel == pganalyze_collector.LogLineInformation_HINT {
					hintLine = futureLine
					hintLine.ParentUUID = logLine.UUID
					hintLineIdx = lowerBound + idx
				} else if futureLine.LogLevel == pganalyze_collector.LogLineInformation_QUERY {
					logLines[lowerBound+idx].ParentUUID = logLine.UUID
					// Ensure QUERY line is consistently marked as statement text log secret
					logLines[lowerBound+idx] = markLineAsSecret(logLines[lowerBound+idx], state.StatementTextLogSecret)
				} else {
					logLines[lowerBound+idx].ParentUUID = logLine.UUID
				}
				additionalLines++
			} else {
				break
			}
		}

		logLine, statementLine, detailLine, contextLine, hintLine, samples = classifyAndSetDetails(logLine, statementLine, detailLine, contextLine, hintLine, samples)

		if statementLineIdx != 0 {
			logLines[statementLineIdx] = statementLine
		}
		if detailLineIdx != 0 {
			logLines[detailLineIdx] = detailLine
		}
		if contextLineIdx != 0 {
			logLines[contextLineIdx] = contextLine
		}
		if hintLineIdx != 0 {
			logLines[hintLineIdx] = hintLine
		}

		logLinesOut = append(logLinesOut, logLine)
	}

	// Ensure no other part of the system accidentally sends log line contents, as
	// they should be considered opaque from here on
	for idx := range logLinesOut {
		logLinesOut[idx].Content = ""
	}

	return
}
