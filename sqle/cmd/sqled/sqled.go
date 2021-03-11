package main

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api"
	"actiontech.cloud/universe/sqle/v4/sqle/api/server"
	"actiontech.cloud/universe/sqle/v4/sqle/inspector"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"actiontech.cloud/universe/sqle/v4/sqle/sqlserverClient"
	"actiontech.cloud/universe/sqle/v4/sqle/utils"
	"actiontech.cloud/universe/ucommon/v4/ubootstrap"
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"syscall"
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
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sqle",
		Short: "SQLe",
		Long:  "SQLe\n\nVersion:\n  " + version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); nil != err {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 10000, "http server port")
	rootCmd.PersistentFlags().StringVarP(&mysqlUser, "mysql-user", "", "sqle", "mysql user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPass, "mysql-password", "", "sqle", "mysql password")
	rootCmd.PersistentFlags().StringVarP(&mysqlHost, "mysql-host", "", "localhost", "mysql host")
	rootCmd.PersistentFlags().StringVarP(&mysqlPort, "mysql-port", "", "3306", "mysql port")
	rootCmd.PersistentFlags().StringVarP(&mysqlSchema, "mysql-schema", "", "sqle", "mysql schema")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&pidFile, "pidfile", "", "", "pid file path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug mode, print more log")
	rootCmd.PersistentFlags().BoolVarP(&autoMigrateTable, "auto-migrate-table", "", false, "auto migrate table if table model has changed")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, _ []string) error {

	mysqlPass, err := utils.DecodeString(mysqlPass)
	if err != nil {
		return fmt.Errorf("decode mysql password to string error : %v", err)
	}

	option := &SqleTaskOptions{
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
	}

	if option.ConfigPath != "" {
		conf := model.Config{}
		b, err := ioutil.ReadFile(option.ConfigPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed error :%v", option.ConfigPath, err)
		}
		err = yaml.Unmarshal(b, &conf)
		if err != nil {
			return fmt.Errorf("unmarshal config file error %v", err)
		}

		option.MysqlUser = conf.Server.DBCnf.MysqlCnf.User
		option.MysqlPass = conf.Server.DBCnf.MysqlCnf.Password
		option.MysqlHost = conf.Server.DBCnf.MysqlCnf.Host
		option.MysqlPort = conf.Server.DBCnf.MysqlCnf.Port
		option.MysqlSchema = conf.Server.DBCnf.MysqlCnf.Schema
		option.Port = conf.Server.SqleCnf.SqleServerPort
		option.AutoMigrateTable = conf.Server.SqleCnf.AutoMigrateTable
		option.Debug = conf.Server.SqleCnf.DebugLog
		option.LogPath = conf.Server.SqleCnf.LogPath
		option.SqlServerParserServerHost = conf.Server.DBCnf.SqlServerCnf.Host
		option.SqlServerParserServerPort = conf.Server.DBCnf.SqlServerCnf.Port
	}

	// init logger
	log.InitLogger(option.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")

	if pidFile != "" {
		f, err := os.Create(pidFile)
		if err != nil {
			return err
		}
		fmt.Fprintf(f, "%d\n", os.Getpid())
		f.Close()
		defer func() {
			os.Remove(pidFile)
		}()
	}

	err = inspector.LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template")
	if err != nil {
		return fmt.Errorf("load './scripts/pt-online-schema-change.template/' failed: %v", err)
	}

	s, err := model.NewStorage(option.MysqlUser, option.MysqlPass, option.MysqlHost, option.MysqlPort, option.MysqlSchema, option.Debug)
	if err != nil {
		return fmt.Errorf("get new storage failed: %v", err)
	}
	model.InitStorage(s)
	_ = sqlserverClient.InitClient(option.SqlServerParserServerHost, option.SqlServerParserServerPort)

	if option.AutoMigrateTable {
		if err := s.AutoMigrate(); err != nil {
			return fmt.Errorf("auto migrate table failed: %v", err)
		}
		if err := s.CreateRulesIfNotExist(inspector.InitRules); err != nil {
			return fmt.Errorf("create rules failed while auto migrating table: %v", err)
		}
		if err := s.CreateDefaultTemplate(inspector.DefaultTemplateRules); err != nil {
			return fmt.Errorf("create default template failed while auto migrating table: %v", err)
		}
		if err := s.CreateAdminUser(); err != nil {
			return fmt.Errorf("create default admin user failed while auto migrating table: %v", err)
		}
	}

	exitChan := make(chan struct{}, 0)
	server.InitSqled(exitChan)
	go api.StartApi(option.Port, exitChan, option.LogPath)

	net := gracenet.Net{}

	killChan := ubootstrap.ListenKillSignal()
	select {
	case <-exitChan:
		log.Logger().Infoln("sqled server will exit")
	case sig := <-killChan:
		switch sig {
		case syscall.SIGUSR2:
			if pid, err := net.StartProcess(); nil != err {
				log.Logger().Infoln("Graceful restarted by signal SIGUSR2, but failed: %v", err)
				return err
			} else {
				log.Logger().Infoln("Graceful restarted, new pid is %v", pid)
			}
			log.Logger().Infoln("old sqled exit")
		default:
			log.Logger().Infoln("Exit by signal %v", sig)
		}
	}

	log.Logger().Info("stop sqled server")
	return nil
}
