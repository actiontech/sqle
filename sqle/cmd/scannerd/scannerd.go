package main

import (
	"fmt"
	"os"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/config"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/supervisor"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// arguments
var (
	host          string
	port          string
	dir           string
	typ           string
	auditPlanName string
	token         string

	logFilePath string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "SQLE Scanner",
		Short: "SQLE Scanner",
		Long:  "SQLE Scanner\nVersion:\n  " + "version",
	}
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "sqle host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "P", "10000", "sqle port")
	rootCmd.PersistentFlags().StringVarP(&auditPlanName, "name", "N", "", "audit plan name")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "A", "", "sqle token")
	rootCmd.MarkPersistentFlagRequired("host")
	rootCmd.MarkPersistentFlagRequired("port")
	rootCmd.MarkPersistentFlagRequired("name")
	rootCmd.MarkPersistentFlagRequired("token")

	mybatisCmd := &cobra.Command{
		Use:   "mybatis",
		Short: "Parse MyBatis XML file",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args, config.ScannerTypeSlowQuery)
			if err != nil {
				fmt.Println(color.RedString("Error: %v", err))
			}
		},
	}
	mybatisCmd.Flags().StringVarP(&dir, "dir", "D", "", "xml directory")
	mybatisCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(mybatisCmd)

	slowlogCmd := &cobra.Command{
		Use:   "slowquery",
		Short: "Parse slow query",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args, config.ScannerTypeSlowQuery)
			if err != nil {
				fmt.Println(color.RedString("Error: %v", err))
			}
		},
	}
	slowlogCmd.Flags().StringVarP(&logFilePath, "log-file", "", "", "log file absolute path")
	slowlogCmd.MarkFlagRequired("log-file")
	rootCmd.AddCommand(slowlogCmd)

	code := 0

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		code = 1
	}

	color.Unset()
	if code != 0 {
		os.Exit(code)
	}
}

func run(_ *cobra.Command, _ []string, typ config.ScannerType) error {
	cfg := &config.Config{
		Host:          host,
		Port:          port,
		Dir:           dir,
		Typ:           typ,
		AuditPlanName: auditPlanName,
		Token:         token,
		LogFilePath:   logFilePath,
	}
	return supervisor.Start(cfg)
}
