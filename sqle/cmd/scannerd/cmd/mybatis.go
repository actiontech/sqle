package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/mybatis"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dir            string
	skipErrorQuery bool
	skipAudit      bool

	mybatisCmd = &cobra.Command{
		Use:   "mybatis",
		Short: "Parse MyBatis XML file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &mybatis.Params{
				XMLDir:         dir,
				APName:         rootCmdFlags.auditPlanName,
				SkipErrorQuery: skipErrorQuery,
				SkipAudit:      skipAudit,
			}
			log := logrus.WithField("scanner", "mybatis")
			client := scanner.NewSQLEClient(scanner.DefaultTimeout, rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token)
			scanner, err := mybatis.New(param, log, client)
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
	mybatisCmd.Flags().StringVarP(&dir, "dir", "D", "", "xml directory")
	mybatisCmd.Flags().BoolVarP(&skipErrorQuery, "skip-error-query", "S", false, "skip the statement that the scanner failed to parse from within the xml file")
	mybatisCmd.Flags().BoolVarP(&skipAudit, "skip-audit", "K", false, "only upload sql to sqle, not audit")
	_ = mybatisCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(mybatisCmd)
}
