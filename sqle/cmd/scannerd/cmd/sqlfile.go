package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	sqlFile "github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/sql_file"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	skipErrorSqlFile  bool
	dbTypeSqlFile     string
	instNameSqlFile   string
	schemaNameSqlFile string

	sqlFileCmd = &cobra.Command{
		Use:   scannerCmd.TypeSQLFile,
		Short: "Parse sql file",
		Run: func(cmd *cobra.Command, args []string) {
			param := &sqlFile.Params{
				SQLDir:           dir,
				SkipErrorQuery:   skipErrorQuery,
				SkipErrorSqlFile: skipErrorSqlFile,
				DbType:           dbTypeSqlFile,
				InstName:         instNameSqlFile,
				SchemaName:       schemaNameSqlFile,
				ShowFileContent:  ShowFileContent,
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
	sqlfile, err := scannerCmd.GetScannerdCmd(scannerCmd.TypeSQLFile)
	if err != nil {
		panic(err)
	}
	sqlFileCmd.Flags().StringVarP(sqlfile.StringFlagFn[scannerCmd.FlagDirectory](&dir))
	sqlFileCmd.Flags().BoolVarP(sqlfile.BoolFlagFn[scannerCmd.FlagSkipErrorSqlFile](&skipErrorSqlFile))
	sqlFileCmd.Flags().StringVarP(sqlfile.StringFlagFn[scannerCmd.FlagDbType](&dbTypeSqlFile))
	sqlFileCmd.Flags().StringVarP(sqlfile.StringFlagFn[scannerCmd.FlagInstanceName](&instNameSqlFile))
	sqlFileCmd.Flags().StringVarP(sqlfile.StringFlagFn[scannerCmd.FlagSchemaName](&schemaNameSqlFile))
	sqlFileCmd.Flags().BoolVarP(sqlfile.BoolFlagFn[scannerCmd.FlagShowFileContent](&ShowFileContent))

	for _, requiredFlag := range sqlfile.RequiredFlags {
		_ = sqlFileCmd.MarkFlagRequired(requiredFlag)
	}

	rootCmd.AddCommand(sqlFileCmd)
}
