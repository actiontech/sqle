package main

import (
	ucobra "actiontech/ucommon/cobra"
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"actiontech/ucommon/ubootstrap"
	"github.com/spf13/cobra"
	"sqle"
	"sqle/storage"
	"sqle/web"
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
				os.ErrExit(err)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlUser, "mysql_user", "", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPass, "mysql_password", "", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlHost, "mysql_host", "", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPort, "mysql_port", "", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().StringVarP(&mysqlSchema, "mysql_schema", "", "actiontech-universe", "run uagent by which user")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 5799, "http server port")
	rootCmd.SetHelpTemplate(ucobra.HELP_TEMPLATE)
	rootCmd.Execute()
}

func run(cmd *cobra.Command, _ []string) error {
	stage := log.NewStage().Enter("run")
	go ubootstrap.DumpLoop()

	if err := ubootstrap.ChangeRunUser(user, true); nil != err {
		return err
	}

	if err := ubootstrap.StartPid(PID_FILE); nil != err {
		return err
	}
	defer ubootstrap.StopPid(PID_FILE)

	killChan := ubootstrap.ListenKillSignal()

	log.InitFileLoggerWithHouseKeep(100, 1024, user, true)

	db, err := storage.NewMysql(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlSchema)
	if err != nil {
		return err
	}

	sqle.InitSqled(stage, db)

	beegoExitChan := make(chan struct{}, 0)
	go web.StartBeego(port, beegoExitChan)
	go sqle.GetSqled().TaskLoop(beegoExitChan)

	select {
	case <-beegoExitChan:
		log.UserInfo(stage, "Beego exit unexpectly")
	case sig := <-killChan:
		//case syscall.SIGUSR2:
		//doesn't support graceful shutdown because beego uses its own graceful-way

		os.HaltIfShutdown(stage)
		log.UserInfo(stage, "Exit by signal %v", sig)
	}

	return nil
}
