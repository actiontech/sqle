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
	slowLogFile string

	slowLogFileCmd = &cobra.Command{
		Use:   "dm-slow-log",
		Short: "Parse dm slow log file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &dmslowlog.Params{
				SlowLogFile: slowLogFile,
				APName:      rootCmdFlags.auditPlanName,
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
	slowLogFileCmd.Flags().StringVarP(&slowLogFile, "slow-log-file", "D", "", "slow log file")
	_ = slowLogFileCmd.MarkFlagRequired("slow-log-file")
	rootCmd.AddCommand(slowLogFileCmd)
}
