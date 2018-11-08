package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sqle"
	"sqle/api"
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sqle",
		Short: "Universe Database Platform",
		Long:  "Universe Database Platform\n\nVersion:\n  " + version,
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
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "f", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&pidFile, "pidfile", "", "", "config file path")
	//rootCmd.SetHelpTemplate(ucobra.HELP_TEMPLATE)

	var docsCmd = &cobra.Command{
		Use:   "docs",
		Short: "a swagger server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := docs(cmd, args); err != nil {
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(docsCmd)
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

	s, err := model.NewMysql(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlSchema)
	if err != nil {
		return err
	}
	model.InitStorage(s)

	exitChan := make(chan struct{}, 0)
	sqle.InitSqled(exitChan)
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

func docs(cmd *cobra.Command, _ []string) error {
	exitChan := make(chan struct{}, 0)
	go api.StartDocs(port, exitChan)
	fmt.Printf("open browser: http://localhost:%d/swagger/index.html\n", port)
	select {
	case <-exitChan:
	}
	return nil
}
