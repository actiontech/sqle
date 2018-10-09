package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/executor"
	"sqle/storage"
	"strings"
)

// inspector rule code
const (
	SELECT_STMT_TABLE_MUST_EXIST = iota
)

// inspector rule level
const (
	RULE_LEVEL_ERROR = iota
	RULE_LEVEL_WARN
)

var DfConfigMap = []*storage.InspectConfig{
	&storage.InspectConfig{
		Code:       SELECT_STMT_TABLE_MUST_EXIST,
		ConfigType: 0,
		Variable:   "",
		StmtType:   0,
		Level:      RULE_LEVEL_WARN,
		Disable:    false,
	},
}

var Rules []*Rule

func init() {
	Rules = []*Rule{
		NewRule(SELECT_STMT_TABLE_MUST_EXIST, selectStmtTableMustExist),
	}
}

type Rule struct {
	DfConfig *storage.InspectConfig
	CheckFn  func(conn *executor.Conn, node ast.StmtNode, config string) (string, error)
}

func NewRule(code int, fn func(conn *executor.Conn, node ast.StmtNode, config string) (string, error)) *Rule {
	return &Rule{
		DfConfig: DfConfigMap[code],
		CheckFn:  fn,
	}
}

func (s *Rule) Check(config *storage.InspectConfig, db *executor.Conn, node ast.StmtNode) (errMsgs, warnMsgs string, err error) {
	var currentConfig *storage.InspectConfig
	if config != nil {
		currentConfig = config
	} else {
		currentConfig = s.DfConfig
	}

	if currentConfig.Disable {
		return errMsgs, warnMsgs, nil
	}

	msg, err := s.CheckFn(db, node, currentConfig.Variable)
	if err != nil {
		return errMsgs, warnMsgs, err
	}

	switch currentConfig.Level {
	case RULE_LEVEL_ERROR:
		errMsgs = msg
	case RULE_LEVEL_WARN:
		warnMsgs = msg
	}
	return errMsgs, warnMsgs, nil
}

func selectStmtTableMustExist(conn *executor.Conn, node ast.StmtNode, variable string) (string, error) {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return "", nil
	}
	tablerefs := selectStmt.From.TableRefs
	tables := getTables(tablerefs)

	tablesName := map[string]struct{}{}
	for _, t := range tables {
		tablesName[getTableName(t)] = struct{}{}
	}
	msgs := []string{}
	for name, _ := range tablesName {
		exist := conn.HasTable(name)
		if conn.Error != nil {
			return "", nil
		}
		if !exist {
			msgs = append(msgs, name)
		}
	}
	if len(msgs) > 0 {
		return fmt.Sprintf("%v is not exist", strings.Join(msgs, ", ")), nil
	}
	return "", nil
}
