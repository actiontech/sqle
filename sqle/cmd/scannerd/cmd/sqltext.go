package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/sqltext"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sql    string
	sqldir string
	audit  string
	synctype  string

	sqltextCmd = &cobra.Command{
		Use:   "sqltext",
		Short: "Parse sql text file",
		Run: func(cmd *cobra.Command, args []string) {

			param := &sqltext.Params{
				SQL:    sql,
				SQLDir: sqldir,
				APName: rootCmdFlags.auditPlanName,
				AUDIT:  audit,
				SYNCTYPE:synctype,
			}
			log := logrus.WithField("scanner", "sqltext")
			client := scanner.NewSQLEClient(time.Second, rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token)
			scanner, err := sqltext.New(param, log, client)
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
	sqltextCmd.Flags().StringVarP(&sql, "sql", "S", "", "sql query")
	sqltextCmd.Flags().StringVarP(&sqldir, "dir", "D", "", "sql directory")
	sqltextCmd.Flags().StringVarP(&audit, "audit", "U", "true", "trigger audit immediately")
	sqltextCmd.Flags().StringVarP(&synctype, "synctype", "T", "1", "sync type 1:full,2:partial")
	//sqltextCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(sqltextCmd)
}
