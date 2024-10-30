package command

const (
	TypeMySQLMybatis       = "mysql_mybatis"
	TypeMySQLSlowLog       = "mysql_slow_log"
	TypeTDSQLInnodbSlowLog = "tdsql_for_innodb_slow_log"
	TypeSQLFile            = "sql_file"
	TypeTBaseSlowLog       = "TBase_slow_log"
	TypeTiDBAuditLog       = "tidb_audit_log"
	TypeRootScannerd       = "root"
)

var (
	rootCmd      scannerCmd = newScannerCmd(TypeRootScannerd)
	myBatis      scannerCmd = newScannerCmd(TypeMySQLMybatis)
	slowLog      scannerCmd = newScannerCmd(TypeMySQLSlowLog)
	sqlFile      scannerCmd = newScannerCmd(TypeSQLFile)
	tbaseLog     scannerCmd = newScannerCmd(TypeTBaseSlowLog)
	tidbAuditLog scannerCmd = newScannerCmd(TypeTiDBAuditLog)
)

func init() {
	rootCmd.addStringFlag(FlagHost, FlagHostSort, "127.0.0.1", "sqle host")
	rootCmd.addStringFlag(FlagPort, FlagPortSort, "10000", "sqle port")
	rootCmd.addStringFlag(FlagAuditPlanID, EmptyFlagSort, "", "audit plan id")
	rootCmd.addStringFlag(FlagToken, FlagTokenSort, "", "sqle token")
	rootCmd.addIntFlag(FlagTimeout, FlagTimeoutSort, 10, "request sqle timeout in seconds")
	rootCmd.addStringFlag(FlagProject, FlagProjectSort, "default", "project name")
	rootCmd.addRequiredFlag(FlagToken)
}

func init() {
	myBatis.addFather(&rootCmd)
	myBatis.addStringFlag(FlagDirectory, FlagDirectorySort, EmptyDefaultValue, "xml directory")
	myBatis.addStringFlag(FlagDbType, FlagDbTypeSort, EmptyDefaultValue, "database type")
	myBatis.addStringFlag(FlagInstanceName, FlagInstanceNameSort, EmptyDefaultValue, "instance name")
	myBatis.addStringFlag(FlagSchemaName, FlagSchemaNameSort, EmptyDefaultValue, "schema name")
	myBatis.addBoolFlag(FlagSkipErrorQuery, FlagSkipErrorQuerySort, false,
		"skip the statement that the scanner failed to parse from within the xml file")
	myBatis.addBoolFlag(FlagSkipErrorXml, FlagSkipErrorXmlSort, false, "skip the xml file that failed to parse")
	myBatis.addRequiredFlag(FlagDirectory)
}

func init() {
	slowLog.addFather(&rootCmd)
	slowLog.addStringFlag(FlagLogFile, EmptyFlagSort, EmptyDefaultValue, "log file absolute path")
	slowLog.addStringFlag(FlagIncludeUserList, EmptyFlagSort, EmptyDefaultValue, "include mysql user list, split by \",\"")
	slowLog.addStringFlag(FlagExcludeUserList, EmptyFlagSort, EmptyDefaultValue, "exclude mysql user list, split by \",\"")
	slowLog.addStringFlag(FlagIncludeSchemaList, EmptyFlagSort, EmptyDefaultValue, "include mysql schema list, split by \",\"")
	slowLog.addStringFlag(FlagExcludeSchemaList, EmptyFlagSort, EmptyDefaultValue, "exclude mysql schema list, split by \",\"")
	slowLog.addRequiredFlag(FlagLogFile)
}

func init() {
	sqlFile.addFather(&rootCmd)
	sqlFile.addStringFlag(FlagDirectory, FlagDirectorySort, EmptyDefaultValue, "sql file directory")
	sqlFile.addBoolFlag(FlagSkipErrorSqlFile, FlagSkipErrorSqlFileSort, false, "skip the sql file that failed to parse")
	sqlFile.addStringFlag(FlagDbType, FlagDbTypeSort, EmptyDefaultValue, "database type")
	sqlFile.addStringFlag(FlagInstanceName, FlagInstanceNameSort, EmptyDefaultValue, "instance name")
	sqlFile.addStringFlag(FlagSchemaName, FlagSchemaNameSort, EmptyDefaultValue, "schema name")
	sqlFile.addRequiredFlag(FlagDirectory)
}
