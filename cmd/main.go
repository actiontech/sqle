package main

import (
	"fmt"
	"io"
	"os"
	"sqle/api"
	"sqle/api/server"
	"sqle/inspector"
	"sqle/log"
	"sqle/model"
	"sqle/sqlserverClient"
	"sqle/utils"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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

	var createConfigFileCmd = &cobra.Command{
		Use:   "load",
		Short: "create config file using the filled in parameters",
		Long:  "create config file using the filled in parameters",
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLogger(logPath)
			defer log.ExitLogger()
			log.Logger().Info("create config file using the filled in parameters")

			//fileName := fmt.Sprintf("/etc/%v", configFile)
			f, err := os.Create(configPath)
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
		},
	}
	rootCmd.AddCommand(createConfigFileCmd)

	rootCmd.Execute()
}

func run(cmd *cobra.Command, _ []string) error {

	// if conf path is exist, load option from conf
	if configPath != "" {
		conf, err := utils.LoadIniConf(configPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed", configPath)
		}
		mysqlUser = conf.GetString("server", "mysql_user", "sqle")
		mysqlPass = conf.GetString("server", "mysql_password", "sqle")
		mysqlHost = conf.GetString("server", "mysql_host", "localhost")
		mysqlPort = conf.GetString("server", "mysql_port", "3306")
		mysqlSchema = conf.GetString("server", "mysql_schema", "")
		port = conf.GetInt("server", "port", 10000)
		autoMigrateTable = conf.GetBool("server", "auto_migrate_table", false)
		debug = conf.GetBool("server", "debug", false)
		logPath = conf.GetString("server", "log_path", "./logs")
		sqlServerParserServerHost = conf.GetString("ms_parser_server", "host", "localhost")
		sqlServerParserServerPort = conf.GetString("ms_parser_server", "port", "10001")
	}

	// init logger
	log.InitLogger(logPath)
	defer log.ExitLogger()

	log.Logger().Info("starting sqled server")

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

	err := inspector.LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template")
	if err != nil {
		return err
	}

	s, err := model.NewStorage(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlSchema, debug)
	if err != nil {
		return err
	}
	model.InitStorage(s)
	_ = sqlserverClient.InitClient(sqlServerParserServerHost, sqlServerParserServerPort)

	if autoMigrateTable {
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
	go api.StartApi(port, exitChan, logPath)

	select {
	case <-exitChan:
		//log.UserInfo(stage, "Beego exit unexpectly")
		//case sig := <-killChan:
		//
		//case syscall.SIGUSR2:
		//doesn't support graceful shutdown because beego uses its own graceful-way
		//
		//os.HaltIfShutdown(stage)
		//log.UserInfo(stage, "Exit by signal %v", sig)
	}
	log.Logger().Info("stop sqled server")
	return nil
}
