//go:build enterprise
// +build enterprise

package command

func init() {
	tbaseLog.addFather(&rootCmd)
	tbaseLog.addStringFlag(FlagDirectory, FlagDirectorySort, EmptyDefaultValue, "log file absolute path")
	tbaseLog.addStringFlag(FlagFileFormat, FlagFileFormatSort, "postgresql-*.csv", "log file name format")
	tbaseLog.addRequiredFlag(FlagDirectory)
}

func init() {
	tidbAuditLog.addFather(&rootCmd)
	tidbAuditLog.addStringFlag(FlagFile, FlagFileSort, EmptyDefaultValue, "audit log file path")
	tidbAuditLog.addRequiredFlag(FlagFile)
}
