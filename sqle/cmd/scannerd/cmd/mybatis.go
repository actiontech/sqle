package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/mybatis"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	pkgAP "github.com/actiontech/sqle/sqle/server/auditplan"
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
		Use:   pkgAP.TypeMySQLMybatis,
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
	mybatisCmd.Flags().StringVarP(&dir, "dir", "D", "", "xml directory")
	mybatisCmd.Flags().BoolVarP(&skipErrorQuery, "skip-error-query", "S", false, "skip the statement that the scanner failed to parse from within the xml file")
	mybatisCmd.Flags().BoolVarP(&skipErrorXml, "skip-error-xml", "X", false, "skip the xml file that failed to parse")
	mybatisCmd.Flags().StringVarP(&dbTypeXml, "db-type", "B", "", "database type")
	mybatisCmd.Flags().StringVarP(&instNameXml, "instance-name", "I", "", "instance name")
	mybatisCmd.Flags().StringVarP(&schemaNameXml, "schema-name", "C", "", "schema name")
	_ = mybatisCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(mybatisCmd)
}
