package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/mybatis"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dir            string
	skipErrorQuery bool
	skipErrorXml   bool
	dbTypeXml      string
	instNameXml    string
	schemaNameXml  string

	mybatisCmd = &cobra.Command{
		Use:   scannerCmd.TypeMySQLMybatis,
		Short: "Parse MyBatis XML file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &mybatis.Params{
				XMLDir:         dir,
				SkipErrorQuery: skipErrorQuery,
				SkipErrorXml:   skipErrorXml,
				DbType:         dbTypeXml,
				InstName:       instNameXml,
				SchemaName:     schemaNameXml,
			}
			log := logrus.WithField("scanner", "mybatis")
			client := scanner.NewSQLEClient(time.Second*time.Duration(rootCmdFlags.timeout), rootCmdFlags.host, rootCmdFlags.port).WithToken(rootCmdFlags.token).WithProject(rootCmdFlags.project)
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
	mybatis, err := scannerCmd.GetScannerdCmd(scannerCmd.TypeMySQLMybatis)
	if err != nil {
		panic(err)
	}
	mybatisCmd.Flags().StringVarP(mybatis.StringFlagFn[scannerCmd.FlagDirectory](&dir))
	mybatisCmd.Flags().BoolVarP(mybatis.BoolFlagFn[scannerCmd.FlagSkipErrorQuery](&skipErrorQuery))
	mybatisCmd.Flags().BoolVarP(mybatis.BoolFlagFn[scannerCmd.FlagSkipErrorXml](&skipErrorXml))
	mybatisCmd.Flags().StringVarP(mybatis.StringFlagFn[scannerCmd.FlagDbType](&dbTypeXml))
	mybatisCmd.Flags().StringVarP(mybatis.StringFlagFn[scannerCmd.FlagInstanceName](&instNameXml))
	mybatisCmd.Flags().StringVarP(mybatis.StringFlagFn[scannerCmd.FlagSchemaName](&schemaNameXml))

	for _, requiredFlag := range mybatis.RequiredFlags {
		_ = mybatisCmd.MarkFlagRequired(requiredFlag)
	}

	rootCmd.AddCommand(mybatisCmd)
}
