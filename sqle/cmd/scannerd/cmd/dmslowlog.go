package cmd

import (
	"context"
	"fmt"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/dmslowlog"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	slowLog                string
	syncDmSlowLogStartTime string
	skipErrorDmSlowQuery   bool
	skipErrorDmSlowLogFile bool
	skipDmSlowLogAudit     bool

	slowLogFileCmd = &cobra.Command{
		Use:   "dm-slow-log",
		Short: "Parse dm slow log file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &dmslowlog.Params{
				SQLDir:                 slowLog,
				APName:                 rootCmdFlags.auditPlanName,
				SyncDmSlowLogStartTime: syncDmSlowLogStartTime,
				SkipErrorQuery:         skipErrorDmSlowQuery,
				SkipErrorSqlFile:       skipErrorDmSlowLogFile,
				SkipAudit:              skipDmSlowLogAudit,
			}
			log := logrus.WithField("scanner", "dmSlowLogFile")
			client := scanner.NewSQLEClient(time.Second*time.Duration(rootCmdFlags.timeout), rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
			scanner, err := dmslowlog.New(param, log, client)
			if err != nil {
				fmt.Println(color.RedString(err.Error()))
				os.Exit(1)
			}

			err = supervisor.Start(context.TODO(), scanner, 30, 1024)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}
)

func init() {
	slowLogFileCmd.Flags().StringVarP(&slowLog, "slow-log", "D", "", "dm slow log file directory")
	slowLogFileCmd.Flags().StringVarP(&syncDmSlowLogStartTime, "sync-Dm-Slow-Log-StartTime", "B", "", "sync dm slowLog start time")
	slowLogFileCmd.Flags().BoolVarP(&skipErrorDmSlowQuery, "skip-error-dm-slow-log-query", "S", false, "skip the statement that the scanner failed to parse from within the log file")
	slowLogFileCmd.Flags().BoolVarP(&skipErrorDmSlowLogFile, "skip-error-dm-slow-log-file", "X", false, "skip the dm slow log file that failed to parse")
	slowLogFileCmd.Flags().BoolVarP(&skipDmSlowLogAudit, "skip-dm-slow-log-file-audit", "K", false, "only upload dm slow log file to sqle, not audit")
	_ = slowLogFileCmd.MarkFlagRequired("slow-log")
	rootCmd.AddCommand(slowLogFileCmd)
}
