package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/pflag"

	_ "github.com/pingcap/tidb/types/parser_driver"

	ucobra "actiontech.cloud/universe/ucommon/v3/cobra"
	"actiontech.cloud/universe/ucommon/v3/os"
	"actiontech.cloud/universe/ucore-common/v3/component"
	ucoreModel "actiontech.cloud/universe/ucore-common/v3/model"
)

var version string
var port int
var user string
var mysqlUser string
var mysqlPass string
var mysqlHost string
var mysqlPort string
var mysqlSchema string
var configPath string
var pidFile string
var debug bool
var autoMigrateTable bool
var logPath = "./logs"
var sqlServerParserServerHost = "127.0.0.1"
var sqlServerParserServerPort = "10001"

var defaultUser = "actiontech-universe"
var ucorePort int
var noFile uint64
var nProc uint64
var securityEnabled bool
var ucoreIps string
var compId string
var serverId string
var caps string
var runOnDmpStr string

const (
	PID_FILE    = "sqled.pid"
	CONFIG_FILE = "./etc/sqled.cnf"
)

func main() {
	rootCmd := component.NewRootCmd(ucoreModel.ComponentTypeSqle, version)
	rootCmd.AddUserFlag(&user, defaultUser)
	rootCmd.AddUlimitNoFileFlag(&noFile)
	rootCmd.AddUlimitNprocFlag(&nProc)
	rootCmd.AddPortFlag(&port, 5801)

	runOnDmp, err := strconv.ParseBool(runOnDmpStr)
	if nil != err {
		fmt.Printf("parse runOnDmpStr failed, runOnDmpStr=%v, err=%v", runOnDmpStr, err)
		runOnDmp = false
	}

	if runOnDmp {
		rootCmd.AddSecurityModeFlag(&securityEnabled)
		rootCmd.AddUcorePortFlag(&ucorePort)
		rootCmd.AddServerIdFlag(&serverId)
		rootCmd.AddCompIdFlag(&compId)
		rootCmd.AddUcoreIpsFlag(&ucoreIps)
	}

	mysqlUserFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "mysql-user",
			Shorthand: "",
			Usage:     "mysql user",
		},
		PString:      &mysqlUser,
		DefaultValue: "sqle",
	}
	rootCmd.AddStringCmdFlag(mysqlUserFlag)

	mysqlPassFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "mysql-password",
			Shorthand: "",
			Usage:     "mysql password",
		},
		PString:      &mysqlPass,
		DefaultValue: "sqle",
	}
	rootCmd.AddStringCmdFlag(mysqlPassFlag)

	mysqlHostFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "mysql-host",
			Shorthand: "",
			Usage:     "mysql host",
		},
		PString:      &mysqlHost,
		DefaultValue: "localhost",
	}
	rootCmd.AddStringCmdFlag(mysqlHostFlag)

	mysqlPortFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "mysql-port",
			Shorthand: "",
			Usage:     "mysql port",
		},
		PString:      &mysqlPort,
		DefaultValue: "3306",
	}
	rootCmd.AddStringCmdFlag(mysqlPortFlag)

	mysqlSchemaFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "mysql-schema",
			Shorthand: "",
			Usage:     "mysql schema",
		},
		PString:      &mysqlSchema,
		DefaultValue: "sqle",
	}
	rootCmd.AddStringCmdFlag(mysqlSchemaFlag)

	configPathFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "config",
			Shorthand: "",
			Usage:     "config file path",
			Persisted: true,
		},
		PString:      &configPath,
		DefaultValue: CONFIG_FILE,
	}
	rootCmd.AddStringCmdFlag(configPathFlag)

	pidFileFlag := &component.StringCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "pidfile",
			Shorthand: "",
			Usage:     "pid file path",
		},
		PString:      &pidFile,
		DefaultValue: PID_FILE,
	}
	rootCmd.AddStringCmdFlag(pidFileFlag)

	debugFlag := &component.BoolCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "debug",
			Shorthand: "",
			Usage:     "debug mode, print more log",
		},
		PBool:        &debug,
		DefaultValue: false,
	}
	rootCmd.AddBoolCmdFlag(debugFlag)
	autoMigrateTableFlag := &component.BoolCmdFlag{
		BaseCmdFlag: component.BaseCmdFlag{
			Name:      "auto-migrate-table",
			Shorthand: "",
			Usage:     "auto migrate table if table model has changed",
		},
		PBool:        &autoMigrateTable,
		DefaultValue: false,
	}
	rootCmd.AddBoolCmdFlag(autoMigrateTableFlag)

	rootCmd.RegisterRun(func() {
		flags, excepts := rootCmd.PersistFlags()
		if err := run(runOnDmp, flags, excepts); nil != err {
			os.ErrExit(err)
		}
	})

	rootCmd.AddCmd(createConfigFileCmd())

	rootCmd.SetHelpTemplate(ucobra.HELP_TEMPLATE)
	rootCmd.AddCmd(component.PPROFCmd(ucoreModel.ComponentTypeSqle, PID_FILE))

	rootCmd.Execute()
}

func run(runOnDmp bool, flags *pflag.FlagSet, excepts []string) error {
	task := NewSqleTask(&SqleTaskOptions{
		ConfigPath:                configPath,
		MysqlUser:                 mysqlUser,
		MysqlPass:                 mysqlPass,
		MysqlHost:                 mysqlHost,
		MysqlPort:                 mysqlPort,
		MysqlSchema:               mysqlSchema,
		Port:                      port,
		AutoMigrateTable:          autoMigrateTable,
		Debug:                     debug,
		LogPath:                   logPath,
		SqlServerParserServerHost: sqlServerParserServerHost,
		SqlServerParserServerPort: sqlServerParserServerPort,
		RunOnDmp:                  runOnDmp,
	})
	runnerOpts := &component.RunnerOptions{
		ComponentRuntimeConfig: component.ComponentRuntimeConfig{
			RunUser:             user,
			RunUserBackupGround: true,
			CompType:            ucoreModel.ComponentTypeUmc,
			CompId:              compId,
			CompGroupId:         ucoreModel.ComponentTypeUmc,
			Version:             version,
			Caps:                caps,
			ServerId:            serverId,
			PIDFile:             PID_FILE,
			Flags:               flags,
			ExceptPersistFlags:  excepts,
			LogFileLimit:        100,
			LogTotalLimit:       1024,
			EnableDetailLog:     true,
			NoFile:              noFile,
			NProc:               nProc,
		},
		ComponentUcoreConfig: component.ComponentUcoreConfig{
			UcoreIps:             ucoreIps,
			UcorePort:            ucorePort,
			UcoreHeartbeatPeriod: 5 * time.Second,
		},
		ComponentGrpcConfig: component.ComponentGrpcConfig{
			EnableGrpcSecurityMode: securityEnabled,
			GrpcPort:               port,
		},
	}

	if runOnDmp {
		NewSqleOnDmpManager(runnerOpts, task).Run()
	} else {
		//TODO startup without DMP
	}
	return nil
}
