package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/slowquery"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logFilePath    string
	includeUsers   string
	excludeUsers   string
	includeSchemas string
	excludeSchemas string

	slowlogCmd = &cobra.Command{
		Use:   scannerCmd.TypeMySQLSlowLog,
		Short: "Parse slow query",
		Run: func(cmd *cobra.Command, args []string) {
			param := &slowquery.Params{
				LogFilePath:    logFilePath,
				AuditPlanID:    rootCmdFlags.auditPlanID,
				IncludeUsers:   includeUsers,
				ExcludeUsers:   excludeUsers,
				IncludeSchemas: includeSchemas,
				ExcludeSchemas: excludeSchemas,
			}
			log := logrus.WithField("scanner", "slowquery")
			client := scanner.NewSQLEClient(time.Second*time.Duration(rootCmdFlags.timeout), rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
			scanner, err := slowquery.New(param, log, client)
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
	slowlog, err := scannerCmd.GetScannerdCmd(scannerCmd.TypeMySQLSlowLog)
	if err != nil {
		panic(err)
	}
	slowlogCmd.Flags().StringVarP(slowlog.StringFlagFn[scannerCmd.FlagLogFile](&logFilePath))
	slowlogCmd.Flags().StringVarP(slowlog.StringFlagFn[scannerCmd.FlagIncludeUserList](&includeUsers))
	slowlogCmd.Flags().StringVarP(slowlog.StringFlagFn[scannerCmd.FlagExcludeUserList](&excludeUsers))
	slowlogCmd.Flags().StringVarP(slowlog.StringFlagFn[scannerCmd.FlagIncludeSchemaList](&includeSchemas))
	slowlogCmd.Flags().StringVarP(slowlog.StringFlagFn[scannerCmd.FlagExcludeSchemaList](&excludeSchemas))

	for _, requiredFlag := range slowlog.RequiredFlags {
		_ = slowlogCmd.MarkFlagRequired(requiredFlag)
	}

	rootCmd.AddCommand(slowlogCmd)
}
