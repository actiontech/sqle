package cmd

import (
	"context"
	"fmt"

	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"

	"github.com/spf13/cobra"
)

var (
	rootCmdFlags struct {
		host        string
		port        string
		token       string
		project     string
		auditPlanID string
		timeout     int
	}

	rootCmd = &cobra.Command{
		Use:     "SQLE Scanner",
		Short:   "SQLE Scanner",
		Version: "SQLE version", // cobra设置--version的固定写法
	}
)

func init() {
	root, err := scannerCmd.GetScannerdCmd(scannerCmd.TypeRootScannerd)
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVarP(root.StringFlagFn[scannerCmd.FlagHost](&rootCmdFlags.host))
	rootCmd.PersistentFlags().StringVarP(root.StringFlagFn[scannerCmd.FlagPort](&rootCmdFlags.port))
	rootCmd.PersistentFlags().StringVarP(root.StringFlagFn[scannerCmd.FlagAuditPlanID](&rootCmdFlags.auditPlanID))
	rootCmd.PersistentFlags().StringVarP(root.StringFlagFn[scannerCmd.FlagToken](&rootCmdFlags.token))
	rootCmd.PersistentFlags().IntVarP(root.IntFlagFn[scannerCmd.FlagTimeout](&rootCmdFlags.timeout))
	rootCmd.PersistentFlags().StringVarP(root.StringFlagFn[scannerCmd.FlagProject](&rootCmdFlags.project))

	for _, requiredFlag := range root.RequiredFlags {
		_ = rootCmd.MarkPersistentFlagRequired(requiredFlag)
	}
}

func Execute(ctx context.Context) error {
	rootCmd.SetVersionTemplate(fmt.Sprintln(ctx.Value(VersionKey)))
	return rootCmd.Execute()
}
