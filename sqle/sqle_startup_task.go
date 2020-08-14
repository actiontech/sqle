package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"io"
	os_ "os"
	"strconv"
	"strings"
	"syscall"

	"actiontech.cloud/universe/sqle/v3/sqle/api"
	"actiontech.cloud/universe/sqle/v3/sqle/api/server"
	"actiontech.cloud/universe/sqle/v3/sqle/inspector"
	"actiontech.cloud/universe/sqle/v3/sqle/log"
	"actiontech.cloud/universe/sqle/v3/sqle/model"
	"actiontech.cloud/universe/sqle/v3/sqle/sqlserverClient"
	"actiontech.cloud/universe/sqle/v3/sqle/utils"
	ucommonLog "actiontech.cloud/universe/ucommon/v3/log"
	"actiontech.cloud/universe/ucommon/v3/os"
	"actiontech.cloud/universe/ucommon/v3/ubootstrap"
	"actiontech.cloud/universe/ucore-common/v3/component"
)

func createConfigFileCmd() *component.Cmd {
	cmd := component.NewCmd(
		"load",
		"create config file using the filled in parameters",
		"create config file using the filled in parameters",
	)

	cmd.RegisterRun(func() {
		log.InitLogger(logPath)
		defer log.ExitLogger()
		log.Logger().Info("create config file using the filled in parameters")

		//fileName := fmt.Sprintf("/etc/%v", configFile)
		f, err := os_.Create(configPath)
		if err != nil {
			log.Logger().Errorf("open %v file error :%v", configPath, err)
			return
		}
		fileContent := `
[server]
port={{SERVER_PORT}}
mysql_host={{MYSQL_HOST}}
mysql_port={{MYSQL_PORT}}
mysql_user={{MYSQL_USER}}
mysql_password={{MYSQL_PASSWORD}}
mysql_schema={{MYSQL_SCHEMA}}
log_path=./logs
#
auto_migrate_table={{AUTO_MIGRATE_TABLE}}
debug={{DEBUG}}

# SQLServer parser server config
[ms_parser_server]
host=
port=
`

		fileContent = strings.Replace(fileContent, "{{SERVER_PORT}}", strconv.Itoa(port), -1)
		fileContent = strings.Replace(fileContent, "{{MYSQL_HOST}}", mysqlHost, -1)
		fileContent = strings.Replace(fileContent, "{{MYSQL_PORT}}", mysqlPort, -1)
		fileContent = strings.Replace(fileContent, "{{MYSQL_USER}}", mysqlUser, -1)
		fileContent = strings.Replace(fileContent, "{{MYSQL_PASSWORD}}", mysqlPass, -1)
		fileContent = strings.Replace(fileContent, "{{MYSQL_SCHEMA}}", mysqlSchema, -1)
		fileContent = strings.Replace(fileContent, "{{AUTO_MIGRATE_TABLE}}", strconv.FormatBool(autoMigrateTable), -1)
		fileContent = strings.Replace(fileContent, "{{DEBUG}}", strconv.FormatBool(debug), -1)
		_, err = io.WriteString(f, fileContent)
		if nil != err {
			log.Logger().Errorf("write config file error :%v", err)
			return
		}
	})

	return cmd
}

//TODO create struct for startup tasks that not run on DMP

// for sqle started on DMP
type SqleOnDmpManager struct {
	r    *component.Runner
	opts *component.RunnerOptions
}

func NewSqleOnDmpManager(opts *component.RunnerOptions, blockedTask component.ComponentLifeCycleTask) *SqleOnDmpManager {
	m := &SqleOnDmpManager{
		r:    component.NewRunner(),
		opts: opts,
	}
	m.r.AddBlockedTask(blockedTask)
	return m
}

func (m *SqleOnDmpManager) initSqleOnDmpManagerTask() {
	m.r.InitLoggerWithHouseKeep(m.opts.LogFileLimit, m.opts.LogTotalLimit, m.opts.RunUser, m.opts.EnableDetailLog).
		InitComponentInfo(m.opts.RunUser, m.opts.GrpcPort, m.opts.PgrpcPort, m.opts.Caps,
			m.opts.CompType, m.opts.CompId, m.opts.CompGroupId, m.opts.Version, m.opts.ServerId).
		CheckPrivileges(m.opts.RunUser, m.opts.RunUserBackupGround, m.opts.Caps).
		PersistFlags(m.opts.Flags, m.opts.ExceptPersistFlags).
		InitAndCheckResourceLimit(m.opts.NoFile, m.opts.NProc).
		InitNetworkConfig(m.opts.CompId, m.opts.EnableGrpcSecurityMode, m.opts.EnableLogGrpcMessage).
		InitAndWatchUcore(m.opts.CompId, m.opts.UcoreIps, m.opts.UcorePort, m.opts.UcoreChangeHandle).
		WatchLogConfig(m.opts.CompId).
		WatchGrpcConfig().
		StartVersionUpdater(m.opts.CompId, m.opts.Version).
		StartUagentWatchOrGuardService(m.opts.PIDFile).
		StartUcoreGuardService(m.opts.CompId, m.opts.UcoreHeartbeatPeriod).
		StartDiagnosisService().
		InitPlatform().
		SetupQaDebugger()

	m.r.StartComponentSipService(m.opts.CompType, m.opts.CompGroupId, m.opts.CompId, m.opts.ServerId)
}

func (m *SqleOnDmpManager) Run() {
	m.initSqleOnDmpManagerTask()
	m.r.Run()
}

type SqleTaskOptions struct {
	ConfigPath                string
	MysqlUser                 string
	MysqlPass                 string
	MysqlHost                 string
	MysqlPort                 string
	MysqlSchema               string
	Port                      int
	AutoMigrateTable          bool
	Debug                     bool
	LogPath                   string
	SqlServerParserServerHost string
	SqlServerParserServerPort string
	RunOnDmp                  bool
}

type SqleTask struct {
	*component.NopComponentLifeCycleTask
	opts *SqleTaskOptions
}

func NewSqleTask(options *SqleTaskOptions) *SqleTask {
	return &SqleTask{
		NopComponentLifeCycleTask: &component.NopComponentLifeCycleTask{},
		opts:                      options,
	}
}

func (t *SqleTask) Initialize(stage *ucommonLog.Stage) error {
	stage.Enter("SqleTask.Initialize")
	defer stage.Exit()

	// if conf path is exist, load option from conf
	if t.opts.ConfigPath != "" {
		conf, err := utils.LoadIniConf(t.opts.ConfigPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed", t.opts.ConfigPath)
		}
		t.opts.MysqlUser = conf.GetString("server", "mysql_user", "sqle")
		t.opts.MysqlPass = conf.GetString("server", "mysql_password", "sqle")
		t.opts.MysqlHost = conf.GetString("server", "mysql_host", "localhost")
		t.opts.MysqlPort = conf.GetString("server", "mysql_port", "3306")
		t.opts.MysqlSchema = conf.GetString("server", "mysql_schema", "")
		t.opts.Port = conf.GetInt("server", "port", 10000)
		t.opts.AutoMigrateTable = conf.GetBool("server", "auto_migrate_table", false)
		t.opts.Debug = conf.GetBool("server", "debug", false)
		t.opts.LogPath = conf.GetString("server", "log_path", "./logs")
		t.opts.SqlServerParserServerHost = conf.GetString("ms_parser_server", "host", "localhost")
		t.opts.SqlServerParserServerPort = conf.GetString("ms_parser_server", "port", "10001")
	}

	// init logger
	log.InitLogger(t.opts.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server, runOnDmp=", t.opts.RunOnDmp)

	err := inspector.LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template")
	if err != nil {
		return err
	}

	s, err := model.NewStorage(t.opts.MysqlUser, t.opts.MysqlPass, t.opts.MysqlHost, t.opts.MysqlPort, t.opts.MysqlSchema, t.opts.Debug)
	if err != nil {
		return err
	}
	model.InitStorage(s)
	_ = sqlserverClient.InitClient(t.opts.SqlServerParserServerHost, t.opts.SqlServerParserServerPort)

	if t.opts.AutoMigrateTable {
		if err := s.AutoMigrate(); err != nil {
			return err
		}
		if err := s.CreateRulesIfNotExist(inspector.DefaultRules); err != nil {
			return err
		}
		if err := s.CreateDefaultTemplate(inspector.DefaultRules); err != nil {
			return err
		}
	}

	exitChan := make(chan struct{}, 0)
	server.InitSqled(exitChan)
	go api.StartApi(t.opts.Port, exitChan, t.opts.LogPath)

	net := gracenet.Net{}

	killChan := ubootstrap.ListenKillSignal()
	select {
	case <-exitChan:
		ucommonLog.Key(stage, "sqled server will exit")
	case sig := <-killChan:
		switch sig {
		case syscall.SIGUSR2:
			ubootstrap.StopPid(PID_FILE)
			if pid, err := net.StartProcess(); nil != err {
				ucommonLog.UserError(stage, "Graceful restarted by signal SIGUSR2, but failed: %v", err)
				return err
			} else {
				ucommonLog.UserInfo(stage, "Graceful restarted, new pid is %v", pid)
			}
			ucommonLog.Key(stage, "old sqled exit")
		default:
			os.HaltIfShutdown(stage)
			ucommonLog.Key(stage, "Exit by signal %v", sig)
		}
	}

	log.Logger().Info("stop sqled server")
	return nil
}

func (t *SqleTask) Finalized(stage *ucommonLog.Stage) error {
	stage.Enter("SqleTask.Finalize")
	defer stage.Exit()

	return nil
}
