package main

import (
	"fmt"
	"os"

	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/config"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/scanners/supervisor"

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

	databaseHost    string
	databasePort    string
	databaseUser    string
	databasePass    string
	slowQuerySecond int
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
		Short: "Parse MyBatis xml file",
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
	slowlogCmd.Flags().StringVarP(&databaseHost, "dbHost", "", "127.0.0.1", "database host")
	slowlogCmd.Flags().StringVarP(&databasePort, "dbPort", "", "3306", "database port")
	slowlogCmd.Flags().StringVarP(&databaseUser, "dbUser", "", "", "database user")
	// todo: need prompt, try https://github.com/manifoldco/promptui
	slowlogCmd.Flags().StringVarP(&databasePass, "dbPass", "", "", "database password")
	// default value keep consistent with MySQL long-query-time
	slowlogCmd.Flags().IntVarP(&slowQuerySecond, "slow-query-second", "", 10, "slow query second")
	slowlogCmd.MarkFlagRequired("dbHost")
	slowlogCmd.MarkFlagRequired("dbPort")
	slowlogCmd.MarkFlagRequired("dbUser")
	slowlogCmd.MarkFlagRequired("dbPass")
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
		Host:            host,
		Port:            port,
		Dir:             dir,
		Typ:             typ,
		AuditPlanName:   auditPlanName,
		Token:           token,
		ScannerDBHost:   databaseHost,
		ScannerDBPort:   databasePort,
		ScannerDBUser:   databaseUser,
		ScannerDBPass:   databasePass,
		SlowQuerySecond: slowQuerySecond,
	}
	return supervisor.Start(cfg)
}
