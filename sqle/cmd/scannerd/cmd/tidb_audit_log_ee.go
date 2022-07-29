//go:build enterprise
// +build enterprise

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/tidb_audit_log"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tidbAuditLogPath string

	tidbAuditLogCmd = &cobra.Command{
		Use:   "tidb-audit-log",
		Short: "Parse TiDB audit log file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &tidb_audit_log.Params{
				AuditLogPath: tidbAuditLogPath,
				APName:       rootCmdFlags.auditPlanName,
			}
			log := logrus.WithField("scanner", "tidb-audit-log")
			client := scanner.NewSQLEClient(scanner.DefaultTimeout, rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token)
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
	tidbAuditLogCmd.Flags().StringVarP(&tidbAuditLogPath, "file", "f", "", "audit log file path")
	_ = tidbAuditLogCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(tidbAuditLogCmd)
}
