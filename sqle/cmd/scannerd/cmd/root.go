package cmd

import (
	pkgScanner "github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/spf13/cobra"
)

var (
	rootCmdFlags struct {
		host          string
		port          string
		token         string
		auditPlanName string
		timeout       int
	}

	rootCmd = &cobra.Command{
		Use:   "SQLE Scanner",
		Short: "SQLE Scanner",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.host, "host", "H", "127.0.0.1", "sqle host")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.port, "port", "P", "10000", "sqle port")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.auditPlanName, "name", "N", "", "audit plan name")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.token, "token", "A", "", "sqle token")
	rootCmd.PersistentFlags().IntVarP(&rootCmdFlags.timeout, "timeout", "T", pkgScanner.DefaultTimeoutNum, "request sqle timeout in seconds")
	_ = rootCmd.MarkPersistentFlagRequired("name")
	_ = rootCmd.MarkPersistentFlagRequired("token")
}

func Execute() error {
	return rootCmd.Execute()
}
