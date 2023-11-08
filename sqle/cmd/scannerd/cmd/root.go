package cmd

import (
	"fmt"
	"context"

	pkgScanner "github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/spf13/cobra"
)

var (
	rootCmdFlags struct {
		host          string
		port          string
		token         string
		project       string
		auditPlanName string
		timeout       int
	}

	rootCmd = &cobra.Command{
		Use:     "SQLE Scanner",
		Short:   "SQLE Scanner",
		Version: "SQLE version", // cobra设置--version的固定写法
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.host, "host", "H", "127.0.0.1", "sqle host")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.port, "port", "P", "10000", "sqle port")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.auditPlanName, "name", "N", "", "audit plan name")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.token, "token", "A", "", "sqle token")
	rootCmd.PersistentFlags().IntVarP(&rootCmdFlags.timeout, "timeout", "T", pkgScanner.DefaultTimeoutNum, "request sqle timeout in seconds")
	rootCmd.PersistentFlags().StringVarP(&rootCmdFlags.project, "project", "J", "default", "project name")
	_ = rootCmd.MarkPersistentFlagRequired("name")
	_ = rootCmd.MarkPersistentFlagRequired("token")
}

func Execute(ctx context.Context) error {
	rootCmd.SetVersionTemplate(fmt.Sprintln(ctx.Value(VersionKey)))
	return rootCmd.Execute()
}
