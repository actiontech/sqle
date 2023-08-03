package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	sqlFile "github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/sql_file"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	skipErrorSqlFile bool

	sqlFileCmd = &cobra.Command{
		Use:   "sqlFile",
		Short: "Parse sql file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &sqlFile.Params{
				SQLDir:           dir,
				APName:           rootCmdFlags.auditPlanName,
				SkipErrorQuery:   skipErrorQuery,
				SkipErrorSqlFile: skipErrorSqlFile,
				SkipAudit:        skipAudit,
			}
			log := logrus.WithField("scanner", "sqlFile")
			client := scanner.NewSQLEClient(time.Second*time.Duration(rootCmdFlags.timeout), rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
			scanner, err := sqlFile.New(param, log, client)
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
	sqlFileCmd.Flags().StringVarP(&dir, "dir", "D", "", "sql file directory")
	sqlFileCmd.Flags().BoolVarP(&skipErrorSqlFile, "skip-error-sql-file", "S", false, "skip the sql file that failed to parse")
	sqlFileCmd.Flags().BoolVarP(&skipAudit, "skip-sql-file-audit", "K", false, "only upload sql to sqle, not audit")
	_ = sqlFileCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(sqlFileCmd)
}
