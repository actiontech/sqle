//go:build enterprise
// +build enterprise

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/tbase_audit_log"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tbaseLogFolder string
	fileFormat     string

	tbaseLogCmd = &cobra.Command{
		Use:   "tbase-audit-log",
		Short: "Parse tbase pg_log",
		Run: func(cmd *cobra.Command, args []string) {
			param := &tbase_audit_log.Params{
				LogFolder:      tbaseLogFolder,
				AuditPlanID:    rootCmdFlags.auditPlanID,
				FileNameFormat: fileFormat,
			}
			log := logrus.WithField("scanner", "tbase-audit-log")
			client := scanner.NewSQLEClient(time.Second*time.Duration(rootCmdFlags.timeout), rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
			scanner, err := tbase_audit_log.New(param, log, client)
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
	tbaseLogCmd.Flags().StringVarP(&tbaseLogFolder, "dir", "D", "", "log file absolute path")
	tbaseLogCmd.Flags().StringVarP(&fileFormat, "format", "F", "postgresql-*.csv", "log file name format")
	_ = tbaseLogCmd.MarkFlagRequired("log-folder")
	rootCmd.AddCommand(tbaseLogCmd)
}
