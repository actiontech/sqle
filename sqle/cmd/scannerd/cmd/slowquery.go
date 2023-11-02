package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

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
		Use:   "slowquery",
		Short: "Parse slow query",
		Run: func(cmd *cobra.Command, args []string) {
			param := &slowquery.Params{
				LogFilePath:    logFilePath,
				APName:         rootCmdFlags.auditPlanName,
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
	slowlogCmd.Flags().StringVarP(&logFilePath, "log-file", "", "", "log file absolute path")
	slowlogCmd.Flags().StringVarP(&includeUsers, "include-user-list", "", "", "include mysql user list, split by \",\"")
	slowlogCmd.Flags().StringVarP(&excludeUsers, "exclude-user-list", "", "", "exclude mysql user list, split by \",\"")
	slowlogCmd.Flags().StringVarP(&includeSchemas, "include-schema-list", "", "", "include mysql schema list, split by \",\"")
	slowlogCmd.Flags().StringVarP(&excludeSchemas, "exclude-schema-list", "", "", "exclude mysql schema list, split by \",\"")
	_ = slowlogCmd.MarkFlagRequired("log-file")
	rootCmd.AddCommand(slowlogCmd)
}
