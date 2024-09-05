package command

import (
	"fmt"

	"strconv"
)

const (
	FlagDirectory        string = "dir"
	FlagDirectorySort    string = "D"
	FlagFile             string = "file"
	FlagFileSort         string = "f"
	FlagInstanceName     string = "instance-name"
	FlagInstanceNameSort string = "I"
	FlagDbType           string = "db-type"
	FlagDbTypeSort       string = "B"
	FlagSchemaName       string = "schema-name"
	FlagSchemaNameSort   string = "C"
	// empty
	EmptyDefaultValue string = ""
	EmptyFlagSort     string = ""
	// root
	FlagHost        string = "host"
	FlagHostSort    string = "H"
	FlagPort        string = "port"
	FlagPortSort    string = "P"
	FlagAuditPlanID string = "audit_plan_id"
	FlagToken       string = "token"
	FlagTokenSort   string = "A"
	FlagTimeout     string = "timeout"
	FlagTimeoutSort string = "T"
	FlagProject     string = "project"
	FlagProjectSort string = "J"
	// mybatis
	FlagSkipErrorQuery     string = "skip-error-query"
	FlagSkipErrorQuerySort string = "S"
	FlagSkipErrorXml       string = "skip-error-xml"
	FlagSkipErrorXmlSort   string = "X"
	// sqlfile
	FlagSkipErrorSqlFile     string = "skip-error-sql-file"
	FlagSkipErrorSqlFileSort string = "S"
	// slow log
	FlagLogFile           string = "log-file"
	FlagIncludeUserList   string = "include-user-list"
	FlagExcludeUserList   string = "exclude-user-list"
	FlagIncludeSchemaList string = "include-schema-list"
	FlagExcludeSchemaList string = "exclude-schema-list"
	// tbase
	FlagFileFormat     string = "format"
	FlagFileFormatSort string = "F"
)

func newScannerCmd(scannerType string) scannerCmd {
	return scannerCmd{
		ScannerType:  scannerType,
		StringFlagFn: make(map[string]func(variable *string) (p *string, name string, shorthand string, value string, usage string)),
		BoolFlagFn:   make(map[string]func(variable *bool) (p *bool, name string, shorthand string, value bool, usage string)),
		IntFlagFn:    make(map[string]func(variable *int) (p *int, name string, shorthand string, value int, usage string)),
	}
}

type scannerCmd struct {
	ScannerType   string
	FatherCmds    []*scannerCmd
	StringFlagFn  map[string]func(variable *string) (p *string, name string, shorthand string, value string, usage string)
	BoolFlagFn    map[string]func(variable *bool) (p *bool, name string, shorthand string, value bool, usage string)
	IntFlagFn     map[string]func(variable *int) (p *int, name string, shorthand string, value int, usage string)
	RequiredFlags []string
}

func GetScannerdCmd(scannerType string) (*scannerCmd, error) {
	switch scannerType {
	case TypeRootScannerd:
		return &rootCmd, nil
	case TypeMySQLMybatis:
		return &myBatis, nil
	case TypeMySQLSlowLog:
		return &slowLog, nil
	case TypeTiDBAuditLog:
		return &tidbAuditLog, nil
	case TypeSQLFile:
		return &sqlFile, nil
	case TypeTBaseSlowLog:
		return &tbaseLog, nil
	default:
		return nil, fmt.Errorf("unsupport scannerd type %s", scannerType)
	}
}

func (newCmd *scannerCmd) addFather(cmd *scannerCmd) {
	newCmd.FatherCmds = append(newCmd.FatherCmds, cmd)
}

func (cmd *scannerCmd) addStringFlag(name string, shorthand string, value string, usage string) {
	cmd.StringFlagFn[name] = func(variable *string) (*string, string, string, string, string) {
		return variable, name, shorthand, value, usage
	}
}

func (cmd *scannerCmd) addIntFlag(name string, shorthand string, value int, usage string) {
	cmd.IntFlagFn[name] = func(variable *int) (*int, string, string, int, string) {
		return variable, name, shorthand, value, usage
	}
}

func (cmd *scannerCmd) addBoolFlag(name string, shorthand string, value bool, usage string) {
	cmd.BoolFlagFn[name] = func(variable *bool) (*bool, string, string, bool, string) {
		return variable, name, shorthand, value, usage
	}
}

func (cmd *scannerCmd) addRequiredFlag(name string) {
	cmd.RequiredFlags = append(cmd.RequiredFlags, name)
}

func (cmd scannerCmd) Type() string {
	return cmd.ScannerType
}

// path can be relative path or absolute path. params is flagName:flagValue map, bool type input true or false string.
func (cmd scannerCmd) GenCommand(path string, params map[string] /* flag name */ string /* flag value */) (string, error) {
	// check required flag exist
	for _, father := range cmd.FatherCmds {
		for _, requiredFlag := range father.RequiredFlags {
			if value, exist := params[requiredFlag]; !exist || value == "" {
				return "", fmt.Errorf("required flag: %s value: %s", requiredFlag, value)
			}
		}
	}
	for _, requiredFlag := range cmd.RequiredFlags {
		if value, exist := params[requiredFlag]; !exist || value == "" {
			return "", fmt.Errorf("required flag: %s value: %s", requiredFlag, value)
		}
	}
	var command string = fmt.Sprintf("%s %s", path, cmd.Type())
	var addParamTpl string = "%s --%s %s"
	// check is flag valid and add flag
	for flagName, flagValue := range params {
		var err error
		var exist bool
		for _, father := range cmd.FatherCmds {
			exist, err = father.checkFlag(flagName, flagValue)
			if err != nil {
				return "", fmt.Errorf("when checking flag: %s,error %w", flagName, err)
			}
			if exist {
				break
			}
		}

		if !exist {
			exist, err = cmd.checkFlag(flagName, flagValue)
			if err != nil {
				return "", fmt.Errorf("when checking flag: %s,error %w", flagName, err)
			}
		}
		if exist {
			if flagValue == "" {
				continue
			}
			command = fmt.Sprintf(addParamTpl, command, flagName, flagValue)
			continue
		}
		return "", fmt.Errorf("unsupport flag %s", flagName)
	}
	return command, nil
}

func (cmd scannerCmd) checkFlag(flagName string, flagValue string) (exist bool, err error) {
	if _, exist = cmd.StringFlagFn[flagName]; exist {
		return true, nil
	}
	if _, exist = cmd.BoolFlagFn[flagName]; exist {
		if flagValue != "false" && flagValue != "true" {
			return true, fmt.Errorf("flage %s is bool type, should input false or true", flagName)
		}
	}
	if _, exist = cmd.IntFlagFn[flagName]; exist {
		_, err = strconv.Atoi(flagValue)
		return true, err
	}
	return false, nil
}
