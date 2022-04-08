package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

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

	mybatisCmd = &cobra.Command{
		Use:   "mybatis",
		Short: "Parse MyBatis XML file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &mybatis.Params{
				XMLDir:         dir,
				APName:         rootCmdFlags.auditPlanName,
				SkipErrorQuery: skipErrorQuery,
			}
			log := logrus.WithField("scanner", "mybatis")
			client := scanner.NewSQLEClient(time.Second, rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token)
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
	mybatisCmd.Flags().BoolVarP(&skipErrorQuery, "skip-unqualified-sql", "S", false, "skip unqualified sql")
	_ = mybatisCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(mybatisCmd)
}
