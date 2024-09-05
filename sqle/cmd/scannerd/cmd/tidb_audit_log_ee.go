//go:build enterprise
// +build enterprise

package cmd

import (
	"context"
	"fmt"
	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/tidb_audit_log"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	tidbAuditLogPath string

	tidbAuditLogCmd = &cobra.Command{
		Use:   scannerCmd.TypeTiDBAuditLog,
		Short: "Parse TiDB audit log file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &tidb_audit_log.Params{
				AuditLogPath: tidbAuditLogPath,
				AuditPlanID:  rootCmdFlags.auditPlanID,
			}
			log := logrus.WithField("scanner", "tidb-audit-log")
			client := scanner.NewSQLEClient(scanner.DefaultTimeout, rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
			scanner, err := tidb_audit_log.New(param, log, client)
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
	tidbAuditLog, err := scannerCmd.GetScannerdCmd(scannerCmd.TypeTiDBAuditLog)
	if err != nil {
		panic(err)
	}
	tidbAuditLogCmd.Flags().StringVarP(tidbAuditLog.StringFlagFn[scannerCmd.FlagFile](&tidbAuditLogPath))

	for _, requiredFlag := range tidbAuditLog.RequiredFlags {
		_ = tidbAuditLogCmd.MarkFlagRequired(requiredFlag)
	}
	rootCmd.AddCommand(tidbAuditLogCmd)
}
