package main

import (
	"fmt"
	"os"

	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/config"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/scanners/mybatis"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// arguments
var (
	host          string
	port          string
	dir           string
	typ           string
	auditPlanName string
	token         string
)

func main() {
	rootCmd = &cobra.Command{
		Use:   "SQLE Scanner",
		Short: "SQLE Scanner",
		Long:  "SQLE Scanner\nVersion:\n  " + "version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "sqle host")
	rootCmd.Flags().StringVarP(&port, "port", "P", "10000", "sqle port")
	rootCmd.Flags().StringVarP(&dir, "dir", "D", "", "xml directory")
	rootCmd.Flags().StringVarP(&auditPlanName, "name", "N", "", "audit plan name")
	rootCmd.Flags().StringVarP(&typ, "typ", "T", "", "scanner type")
	rootCmd.Flags().StringVarP(&token, "token", "A", "", "sqle token")

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

func run(_ *cobra.Command, _ []string) error {
	cfg := &config.Config{
		Host:          host,
		Port:          port,
		Dir:           dir,
		Typ:           typ,
		AuditPlanName: auditPlanName,
		Token:         token,
	}
	switch cfg.Typ {
	case "mybatis":
		return mybatis.MybatisScanner(cfg)
	default:
		return nil
	}
}
