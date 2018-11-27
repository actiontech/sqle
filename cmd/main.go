package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sqle/api"
	"sqle/api/server"
	"sqle/inspector"
	"sqle/model"
	"sqle/utils"
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sqle",
		Short: "SQLe",
		Long:  "SQLe\n\nVersion:\n  " + version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); nil != err {
				//os.ErrExit(err)
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 5799, "http server port")
	rootCmd.PersistentFlags().StringVarP(&mysqlUser, "mysql-user", "", "", "mysql user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPass, "mysql-password", "", "", "mysql password")
	rootCmd.PersistentFlags().StringVarP(&mysqlHost, "mysql-host", "", "localhost", "mysql host")
	rootCmd.PersistentFlags().StringVarP(&mysqlPort, "mysql-port", "", "3306", "mysql port")
	rootCmd.PersistentFlags().StringVarP(&mysqlSchema, "mysql-schema", "", "sqle", "mysql schema")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&pidFile, "pidfile", "", "", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug mode, print more log")
	rootCmd.PersistentFlags().BoolVarP(&autoMigrateTable, "auto-migrate-table", "", false, "auto migrate table if table model has changed")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, _ []string) error {

	// if conf path is exist, load option from conf
	if configPath != "" {
		conf, err := utils.LoadIniConf(configPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed", configPath)
		}
		mysqlUser = conf.GetString("server", "mysql_user", "")
		mysqlPass = conf.GetString("server", "mysql_password", "")
		mysqlHost = conf.GetString("server", "mysql_host", "")
		mysqlPort = conf.GetString("server", "mysql_port", "")
		mysqlSchema = conf.GetString("server", "mysql_schema", "")
		port = conf.GetInt("server", "port", 12160)
		autoMigrateTable = conf.GetBool("server", "auto_migrate_table", false)
		debug = conf.GetBool("server", "debug", false)
	}

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

	s, err := model.NewStorage(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlSchema, debug)
	if err != nil {
		return err
	}
	model.InitStorage(s)

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
		if err := s.CreateConfigsIfNotExist(inspector.GetAllConfig()); err != nil {
			return err
		}
	}

	exitChan := make(chan struct{}, 0)
	server.InitSqled(exitChan)
	go api.StartApi(port, exitChan)

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
	return nil
}
