package onlineddl

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"

	"github.com/github/gh-ost/go/base"
	"github.com/github/gh-ost/go/logic"
	"github.com/go-ini/ini"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Executor struct {
	l  base.Logger
	mc *base.MigrationContext
}

func NewExecutor(logger *logrus.Entry, inst *driver.DSN, schema string, query string) (*Executor, error) {
	logger = logger.WithFields(logrus.Fields{
		"onlineddl": "gh-ost",
		"host":      inst.Host,
		"port":      inst.Port,
		"alter":     query,
	})

	la := newLogAdaptor(logger)

	mc := base.NewMigrationContext()

	// get migration context fields from SQLE
	{
		mc.Log = la
		mc.InspectorConnectionConfig.Key.Hostname = inst.Host
		port, _ := strconv.ParseInt(inst.Port, 10, 64)
		mc.InspectorConnectionConfig.Key.Port = int(port)
		mc.CliUser = inst.User
		mc.CliPassword = inst.Password
		mc.UseTLS = false
		schemaAtQuery, table, alterOpts, err := parseAlterTableOptions(query)
		if err != nil {
			return nil, errors.Wrap(err, "parse alter")
		}
		if schemaAtQuery == "" {
			mc.DatabaseName = schema
		} else {
			mc.DatabaseName = schemaAtQuery
		}
		mc.OriginalTableName = table
		mc.AlterStatementOptions = alterOpts
	}

	// get migration context fields from config file
	{
		cfg := newDefaultConfig()

		if _, err := os.Stat(cfgPath); err == nil {
			f, err := ini.Load(cfgPath)
			if err != nil {
				return nil, errors.Wrap(err, "load config for gh-ost")
			}

			if err := f.Section("DEFAULT").MapTo(cfg); err != nil {
				return nil, errors.Wrap(err, "map config to struct")
			}
		}

		if err := cfg.apply(mc); err != nil {
			return nil, errors.Wrap(err, "apply config to migration context")
		}
	}

	if err := checkMigrationContext(mc); err != nil {
		return nil, errors.Wrap(err, "check migration context")
	}

	return &Executor{
		l:  la,
		mc: mc,
	}, nil
}

func (e *Executor) Execute(ctx context.Context, dryRun bool) error {
	if dryRun {
		e.mc.Noop = true
	}

	m := logic.NewMigrator(e.mc)
	err := m.Migrate()
	if err != nil {
		return errors.Wrapf(err, "migrate table, dry-run(%v)", dryRun)
	}

	return nil
}

const cfgPath = "./etc/gh-ost.ini"

// config refer to https://github.com/github/gh-ost/blob/master/go/cmd/gh-ost/main.go
type config struct {
	//user             string  `ini:"user"`
	//password         string  `ini:"password"`
	//database 		   string `ini:"database"`
	//table            string `ini:"table"`
	//alter            string `ini:"alter"`
	MySQLTimeout float64 `ini:"mysql_timeout"`

	AssumeMasterHost string `ini:"assume_master_host"`
	MasterUser       string `ini:"master_user"`
	MasterPassword   string `ini:"master_password"`

	//ssl              bool   `ini:"ssl"`
	//sslCA            string `ini:"ssl_ca"`
	//sslCert          string `ini:"ssl_cert"`
	//sslKey           string `ini:"ssl_key"`
	//sslAllowInsecure bool   `ini:"ssl_allow_insecure"`

	ExactRowcount          bool `ini:"exact_rowcount"`
	ConcurrentRowcount     bool `ini:"concurrent_rowcount"`
	AllowOnMaster          bool `ini:"allow_on_master"`
	AllowMasterMaster      bool `ini:"allow_master_master"`
	AllowNullableUniqueKey bool `ini:"allow_nullable_unique_key"`
	ApproveRenamedColumns  bool `ini:"approve_renamed_columns"`
	SkipRenamedColumns     bool `ini:"skip_renamed_columns"`
	Tungsten               bool `ini:"tungsten"`
	DiscardForeignKeys     bool `ini:"discard_foreign_keys"`
	SkipForeignKeyChecks   bool `ini:"skip_foreign_key_checks"`
	SkipStrictMode         bool `ini:"skip_strict_mode"`
	AliyunRDS              bool `ini:"aliyun_rds"`
	GCP                    bool `ini:"gcp"`
	Azure                  bool `ini:"azure"`

	TestOnReplica                bool `ini:"test_on_replica"`
	TestOnReplicaSkipReplicaStop bool `ini:"test_on_replica_skip_replica_stop"`
	MigrateOnReplica             bool `ini:"migrate_on_replica"`
	OKToDropTable                bool `ini:"ok_to_drop_table"`
	InitiallyDropOldTable        bool `ini:"initially_drop_old_table"`
	InitiallyDropGhostTable      bool `ini:"initially_drop_ghost_table"`
	TimestampOldTable            bool `ini:"timestamp_old_table"`

	CutOver string `ini:"cut_over"`

	ForceNamedCutOver bool `ini:"force_named_cut_over"`
	ForceNamedPanic   bool `ini:"force_named_panic"`

	SwitchToRBR                   bool    `ini:"switch_to_rbr"`
	AssumeRBR                     bool    `ini:"assume_rbr"`
	CutOverExponentialBackoff     bool    `ini:"cut_over_exponential_backoff"`
	ExponentialBackoffMaxInterval int64   `ini:"exponential_backoff_max_interval"`
	ChunkSize                     int64   `ini:"chunk_size"`
	DMLBatchSize                  int64   `ini:"dml_batch_size"`
	DefaultRetries                int64   `ini:"default_retries"`
	CutOverLockTimeoutSeconds     int64   `ini:"cut_over_lock_timeout_seconds"`
	NiceRatio                     float64 `ini:"nice_ratio"`

	MaxLagMillis int64 `ini:"max_lag_millis"`
	//replicationLagQuery        string `ini:"replication_lag_query"`
	ThrottleControlReplicas    string `ini:"throttle_control_replicas"`
	ThrottleQuery              string `ini:"throttle_query"`
	ThrottleHTTP               string `ini:"throttle_http"`
	IgnoreHTTPErrors           bool   `ini:"ignore_http_errors"`
	HeartbeatIntervalMillis    int64  `ini:"heartbeat_interval_millis"`
	ThrottleFlagFile           string `ini:"throttle_flag_file"`
	ThrottleAdditionalFlagFile string `ini:"throttle_additional_flag_file"`
	PostponeCutOverFlagFile    string `ini:"postpone_cut_over_flag_file"`
	PanicFlagFile              string `ini:"panic_flag_file"`

	InitiallyDropSocketFile bool   `ini:"initially_drop_socket_file"`
	ServeSocketFile         string `ini:"serve_socket_file"`
	ServeTCPPort            int64  `ini:"serve_tcp_port"`

	//hooksPath      string `ini:"hooks_path"`
	//hooksHint      string `ini:"hooks_hint"`
	//hooksHintOwner string `ini:"hooks_hint_owner"`
	//hooksHintToken string `ini:"hooks_hint_token"`

	ReplicaServerID uint `ini:"replica_server_id"`

	MaxLoad                      string `ini:"max_load"`
	CriticalLoad                 string `ini:"critical_load"`
	CriticalLoadIntervalMillis   int64  `ini:"critical_load_interval_millis"`
	CriticalLoadHibernateSeconds int64  `ini:"critical_load_hibernate_seconds"`
	//quiet                        bool   `ini:"quiet"`
	//verbose                      bool   `ini:"verbose"`
	//debug                        bool   `ini:"debug"`
	//stack                        bool   `ini:"stack"`
	//help                         bool   `ini:"help"`
	//version                      bool   `ini:"version"`
	//checkFlag       bool   `ini:"check_flag"`
	ForceTableNames string `ini:"force_table_names"`
}

func newDefaultConfig() *config {
	cfg := &config{
		MySQLTimeout:                  0.0,
		AssumeMasterHost:              "",
		MasterUser:                    "",
		MasterPassword:                "",
		ExactRowcount:                 false,
		ConcurrentRowcount:            true,
		AllowOnMaster:                 false,
		AllowMasterMaster:             false,
		AllowNullableUniqueKey:        false,
		ApproveRenamedColumns:         false,
		SkipRenamedColumns:            false,
		Tungsten:                      false,
		DiscardForeignKeys:            false,
		SkipForeignKeyChecks:          false,
		SkipStrictMode:                false,
		AliyunRDS:                     false,
		GCP:                           false,
		Azure:                         false,
		TestOnReplica:                 false,
		TestOnReplicaSkipReplicaStop:  false,
		MigrateOnReplica:              false,
		OKToDropTable:                 false,
		InitiallyDropOldTable:         false,
		InitiallyDropGhostTable:       false,
		TimestampOldTable:             false,
		CutOver:                       "",
		ForceNamedCutOver:             false,
		ForceNamedPanic:               false,
		SwitchToRBR:                   false,
		AssumeRBR:                     false,
		CutOverExponentialBackoff:     false,
		ExponentialBackoffMaxInterval: 64,
		ChunkSize:                     1000,
		DMLBatchSize:                  10,
		DefaultRetries:                120,
		CutOverLockTimeoutSeconds:     3,
		NiceRatio:                     0,
		MaxLagMillis:                  1500,
		ThrottleControlReplicas:       "",
		ThrottleQuery:                 "",
		ThrottleHTTP:                  "",
		IgnoreHTTPErrors:              false,
		HeartbeatIntervalMillis:       100,
		ThrottleFlagFile:              "",
		ThrottleAdditionalFlagFile:    "/tmp/gh-ost.throttle",
		PostponeCutOverFlagFile:       "",
		PanicFlagFile:                 "",
		InitiallyDropSocketFile:       false,
		ServeSocketFile:               "",
		ServeTCPPort:                  0,
		ReplicaServerID:               uint(rand.Uint32()),
		MaxLoad:                       "Threads_running=80,Threads_connected=1000",
		CriticalLoad:                  "Threads_running=80,Threads_connected=1000",
		CriticalLoadIntervalMillis:    0,
		CriticalLoadHibernateSeconds:  0,
		ForceTableNames:               "",
	}
	return cfg
}

func (cfg *config) apply(mc *base.MigrationContext) error {
	mc.ReplicaServerId = cfg.ReplicaServerID
	mc.AssumeMasterHostname = cfg.AssumeMasterHost
	mc.CliMasterUser = cfg.MasterUser
	mc.CliMasterPassword = cfg.MasterPassword
	mc.CountTableRows = cfg.ExactRowcount
	mc.ConcurrentCountTableRows = cfg.ConcurrentRowcount
	mc.AllowedRunningOnMaster = cfg.AllowOnMaster
	mc.AllowedMasterMaster = cfg.AllowMasterMaster
	mc.NullableUniqueKeyAllowed = cfg.AllowNullableUniqueKey
	mc.ApproveRenamedColumns = cfg.ApproveRenamedColumns
	mc.SkipRenamedColumns = cfg.SkipRenamedColumns // todo
	mc.IsTungsten = cfg.Tungsten
	mc.DiscardForeignKeys = cfg.DiscardForeignKeys
	mc.SkipForeignKeyChecks = cfg.SkipForeignKeyChecks
	mc.AliyunRDS = cfg.AliyunRDS
	mc.GoogleCloudPlatform = cfg.GCP
	mc.AzureMySQL = cfg.Azure
	mc.TestOnReplica = cfg.TestOnReplica
	mc.TestOnReplicaSkipReplicaStop = cfg.TestOnReplicaSkipReplicaStop
	mc.MigrateOnReplica = cfg.MigrateOnReplica
	mc.OkToDropTable = cfg.OKToDropTable
	mc.InitiallyDropOldTable = cfg.InitiallyDropOldTable
	mc.InitiallyDropGhostTable = cfg.InitiallyDropGhostTable
	mc.TimestampOldTable = cfg.TimestampOldTable
	mc.ForceTmpTableName = cfg.ForceTableNames
	mc.CriticalLoadIntervalMilliseconds = cfg.CriticalLoadIntervalMillis
	mc.CriticalLoadHibernateSeconds = cfg.CriticalLoadHibernateSeconds
	mc.ThrottleFlagFile = cfg.ThrottleFlagFile
	mc.ThrottleAdditionalFlagFile = cfg.ThrottleAdditionalFlagFile
	mc.PostponeCutOverFlagFile = cfg.PostponeCutOverFlagFile
	mc.IgnoreHTTPErrors = cfg.IgnoreHTTPErrors
	mc.DropServeSocket = cfg.InitiallyDropSocketFile
	mc.ServeTCPPort = cfg.ServeTCPPort

	mc.ServeSocketFile = cfg.ServeSocketFile
	if mc.ServeSocketFile == "" {
		mc.ServeSocketFile = fmt.Sprintf("/tmp/gh-ost.%s.%s.sock", mc.DatabaseName, mc.OriginalTableName)
	}

	switch cfg.CutOver {
	case "atomic", "default", "":
		mc.CutOverType = base.CutOverAtomic
	case "two-step":
		mc.CutOverType = base.CutOverTwoStep
	default:
		return fmt.Errorf("Unknown cut-over: %s", cfg.CutOver)
	}

	mc.ForceNamedCutOverCommand = cfg.ForceNamedCutOver
	mc.ForceNamedPanicCommand = cfg.ForceNamedPanic
	mc.PanicFlagFile = cfg.PanicFlagFile
	mc.SwitchToRowBinlogFormat = cfg.SwitchToRBR
	mc.AssumeRBR = cfg.AssumeRBR
	mc.CutOverExponentialBackoff = cfg.CutOverExponentialBackoff

	if err := mc.ReadThrottleControlReplicaKeys(cfg.ThrottleControlReplicas); err != nil {
		return errors.Wrap(err, "read throttle_control_replicas")
	}
	if err := mc.ReadMaxLoad(cfg.MaxLoad); err != nil {
		return errors.Wrap(err, "read max_load")
	}
	if err := mc.ReadCriticalLoad(cfg.CriticalLoad); err != nil {
		return errors.Wrap(err, "read critical_load")
	}
	if err := mc.SetCutOverLockTimeoutSeconds(cfg.CutOverLockTimeoutSeconds); err != nil {
		return errors.Wrap(err, "set cut_over_lock_timeout_seconds")
	}
	if err := mc.SetExponentialBackoffMaxInterval(cfg.ExponentialBackoffMaxInterval); err != nil {
		return errors.Wrap(err, "set exponential_backoff_max_interval")
	}

	mc.SetChunkSize(cfg.ChunkSize)
	mc.SetHeartbeatIntervalMilliseconds(cfg.HeartbeatIntervalMillis)
	mc.SetNiceRatio(cfg.NiceRatio)
	mc.SetDMLBatchSize(cfg.DMLBatchSize)
	mc.SetMaxLagMillisecondsThrottleThreshold(cfg.MaxLagMillis)
	mc.SetThrottleQuery(cfg.ThrottleQuery)
	mc.SetThrottleHTTP(cfg.ThrottleHTTP)
	mc.SetDefaultNumRetries(cfg.DefaultRetries)
	mc.ApplyCredentials()

	return nil
}

func checkMigrationContext(mc *base.MigrationContext) error {
	if mc.CliMasterUser != "" && mc.AssumeMasterHostname == "" {
		return errors.New("--master-user requires --assume-master-host")
	}
	if mc.CliMasterPassword != "" && mc.AssumeMasterHostname == "" {
		return errors.New("--master-password requires --assume-master-host")
	}
	if mc.AllowedRunningOnMaster && mc.TestOnReplica {
		return errors.New("--allow-on-master and --test-on-replica are mutually exclusive")
	}
	if mc.AllowedRunningOnMaster && mc.MigrateOnReplica {
		return errors.New("--allow-on-master and --migrate-on-replica are mutually exclusive")
	}
	if mc.MigrateOnReplica && mc.TestOnReplica {
		return errors.New("--migrate-on-replica and --test-on-replica are mutually exclusive")
	}
	if mc.SwitchToRowBinlogFormat && mc.AssumeRBR {
		return errors.New("--switch-to-rbr and --assume-rbr are mutually exclusive")
	}
	if mc.TestOnReplicaSkipReplicaStop {
		if !mc.TestOnReplica {
			return errors.New("--test-on-replica-skip-replica-stop requires --test-on-replica to be enabled")
		}
		mc.Log.Warning("--test-on-replica-skip-replica-stop enabled. We will not stop replication before cut-over. Ensure you have a plugin that does this.")
	}
	return nil
}

func parseAlterTableOptions(alter string) (schema, table, alterOpts string, err error) {
	p := parser.New()
	nodes, _, err := p.PerfectParse(alter, "", "")
	alterStmt, ok := nodes[0].(*ast.AlterTableStmt)
	if !ok {
		return "", "", "", fmt.Errorf("want alter stmt, but got %v", alter)
	}

	var alterOptPart []string
	var builder strings.Builder
	for _, spec := range alterStmt.Specs {
		builder.Reset()
		restoreCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, &builder)
		if err = spec.Restore(restoreCtx); err != nil {
			return "", "", "", errors.Wrap(err, "")
		}
		alterOptPart = append(alterOptPart, builder.String())
	}

	schema = alterStmt.Table.Schema.String()
	table = alterStmt.Table.Name.String()
	alterOpts = strings.Join(alterOptPart, ",")

	return
}
