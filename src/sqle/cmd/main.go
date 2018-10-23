package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sqle"
	"sqle/api"
	"sqle/storage"
)

const (
	PID_FILE = "sqle.pid"
)

var version string
var port int
var user string
var mysqlUser string
var mysqlPass string
var mysqlHost string
var mysqlPort string
var mysqlSchema string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sqle",
		Short: "Universe Database Platform",
		Long:  "Universe Database Platform\n\nVersion:\n  " + version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); nil != err {
				//os.ErrExit(err)
				os.Exit(1)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlUser, "mysql_user", "", "", "mysql user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPass, "mysql_password", "", "", "mysql password")
	rootCmd.PersistentFlags().StringVarP(&mysqlHost, "mysql_host", "", "localhost", "mysql host")
	rootCmd.PersistentFlags().StringVarP(&mysqlPort, "mysql_port", "", "3306", "mysql port")
	rootCmd.PersistentFlags().StringVarP(&mysqlSchema, "mysql_schema", "", "sqle", "mysql schema")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 5799, "http server port")
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
	//stage := log.NewStage().Enter("run")
	//go ubootstrap.DumpLoop()

	//if err := ubootstrap.ChangeRunUser(user, true); nil != err {
	//	return err
	//}

	//if err := ubootstrap.StartPid(PID_FILE); nil != err {
	//	return err
	//}
	//defer ubootstrap.StopPid(PID_FILE)

	//killChan := ubootstrap.ListenKillSignal()
	//
	//log.InitFileLoggerWithHouseKeep(100, 1024, user, true)

	s, err := storage.NewMysql(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlSchema)
	if err != nil {
		return err
	}
	storage.InitStorage(s)

	exitChan := make(chan struct{}, 0)
	sqle.InitSqled(exitChan)
	go api.StartApi(port, exitChan)

	select {
	case <-exitChan:
		//log.UserInfo(stage, "Beego exit unexpectly")
		//case sig := <-killChan:

		//case syscall.SIGUSR2:
		//doesn't support graceful shutdown because beego uses its own graceful-way

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
