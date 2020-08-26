package main

import (
	"fmt"
	"io/ioutil"
	"syscall"

	"actiontech.cloud/universe/sqle/v3/sqle/utils"

	"github.com/facebookgo/grace/gracenet"
	yaml "gopkg.in/yaml.v2"

	"actiontech.cloud/universe/sqle/v3/sqle/api"
	"actiontech.cloud/universe/sqle/v3/sqle/api/server"
	"actiontech.cloud/universe/sqle/v3/sqle/inspector"
	"actiontech.cloud/universe/sqle/v3/sqle/log"
	"actiontech.cloud/universe/sqle/v3/sqle/model"
	"actiontech.cloud/universe/sqle/v3/sqle/sqlserverClient"
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
		mysqlPass, err := utils.DecodeString(mysqlPass)
		if err != nil {
			log.Logger().Errorf("decode mysql password to string error :%v", err)
			return

		}
		conf := model.Config{}
		conf.Server.DBCnf.MysqlCnf.Port = mysqlPort
		conf.Server.DBCnf.MysqlCnf.Host = mysqlHost
		conf.Server.DBCnf.MysqlCnf.Schema = mysqlSchema
		conf.Server.DBCnf.MysqlCnf.Password = mysqlPass
		conf.Server.DBCnf.MysqlCnf.User = mysqlUser
		conf.Server.SqleCnf.DebugLog = debug
		conf.Server.SqleCnf.LogPath = logPath
		conf.Server.SqleCnf.SqleServerPort = port
		conf.Server.SqleCnf.AutoMigrateTable = autoMigrateTable
		data, err := yaml.Marshal(conf)
		if err != nil {

			log.Logger().Errorf("marshal sqle config error %v", err)
			return

		}
		err = ioutil.WriteFile(configPath, data, 0666)
		if err != nil {
			log.Logger().Errorf("%v write sqle config file error %v", configPath, err)
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
		StartComponentSipService(m.opts.CompType, m.opts.CompGroupId, m.opts.CompId, m.opts.ServerId)
}

func (m *SqleOnDmpManager) Run() {
	m.initSqleOnDmpManagerTask()
	m.r.Run()
}

// for sqle running independently
type SqleManager struct {
	r    *component.Runner
	opts *component.RunnerOptions
}

func NewSqleManager(opts *component.RunnerOptions, blockedTask component.ComponentLifeCycleTask) *SqleManager {
	m := &SqleManager{
		r:    component.NewRunner(),
		opts: opts,
	}
	m.r.AddBlockedTask(blockedTask)
	return m
}

func (m *SqleManager) initSqleManagerTask() {
	m.r.InitLoggerWithHouseKeep(m.opts.LogFileLimit, m.opts.LogTotalLimit, m.opts.RunUser, m.opts.EnableDetailLog).
		InitComponentInfo(m.opts.RunUser, m.opts.GrpcPort, m.opts.PgrpcPort, m.opts.Caps,
			m.opts.CompType, m.opts.CompId, m.opts.CompGroupId, m.opts.Version, m.opts.ServerId).
		CheckPrivileges(m.opts.RunUser, m.opts.RunUserBackupGround, m.opts.Caps).
		PersistFlags(m.opts.Flags, m.opts.ExceptPersistFlags).
		InitAndCheckResourceLimit(m.opts.NoFile, m.opts.NProc).
		StartUagentWatchOrGuardService(m.opts.PIDFile)
}

func (m *SqleManager) Run() {
	m.initSqleManagerTask()
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
		conf := model.Config{}
		b, err := ioutil.ReadFile(t.opts.ConfigPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed error :%v", t.opts.ConfigPath, err)
		}
		err = yaml.Unmarshal(b, &conf)
		if err != nil {
			return fmt.Errorf("unmarshal config file error %v", err)
		}

		t.opts.MysqlUser = conf.Server.DBCnf.MysqlCnf.User
		t.opts.MysqlPass = conf.Server.DBCnf.MysqlCnf.Password
		t.opts.MysqlHost = conf.Server.DBCnf.MysqlCnf.Host
		t.opts.MysqlPort = conf.Server.DBCnf.MysqlCnf.Port
		t.opts.MysqlSchema = conf.Server.DBCnf.MysqlCnf.Schema
		t.opts.Port = conf.Server.SqleCnf.SqleServerPort
		t.opts.AutoMigrateTable = conf.Server.SqleCnf.AutoMigrateTable
		t.opts.Debug = conf.Server.SqleCnf.DebugLog
		t.opts.LogPath = conf.Server.SqleCnf.LogPath
		t.opts.SqlServerParserServerHost = conf.Server.DBCnf.SqlServerCnf.Host
		t.opts.SqlServerParserServerPort = conf.Server.DBCnf.SqlServerCnf.Port
	}

	// init logger
	log.InitLogger(t.opts.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server, runOnDmp=", t.opts.RunOnDmp)

	err := inspector.LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template")
	if err != nil {
		return fmt.Errorf("load './scripts/pt-online-schema-change.template/' failed: %v", err)
	}

	s, err := model.NewStorage(t.opts.MysqlUser, t.opts.MysqlPass, t.opts.MysqlHost, t.opts.MysqlPort, t.opts.MysqlSchema, t.opts.Debug)
	if err != nil {
		return fmt.Errorf("get new storage failed: %v", err)
	}
	model.InitStorage(s)
	_ = sqlserverClient.InitClient(t.opts.SqlServerParserServerHost, t.opts.SqlServerParserServerPort)

	if t.opts.AutoMigrateTable {
		if err := s.AutoMigrate(); err != nil {
			return fmt.Errorf("auto migrate table failed: %v", err)
		}
		if err := s.CreateRulesIfNotExist(inspector.DefaultRules); err != nil {
			return fmt.Errorf("create rules failed while auto migrating table: %v", err)
		}
		if err := s.CreateDefaultTemplate(inspector.DefaultRules); err != nil {
			return fmt.Errorf("create default template failed while auto migrating table: %v", err)
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
