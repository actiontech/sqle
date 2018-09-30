package inspector

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pingcap/tidb/ast"
	"sqle/storage"
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

type Rule struct {
	DfConfig *storage.InspectConfig
	CheckFn  func(db *gorm.DB, node ast.StmtNode, config string) (string, error)
}

func NewRule(code, configType int, valiable string, stmtType, level int, disable bool,
	fn func(db *gorm.DB, node ast.StmtNode, config string) (string, error), desc string) *Rule {
	return &Rule{
		DfConfig: &storage.InspectConfig{
			Code:       code,
			ConfigType: configType,
			StmtType:   stmtType,
			Variable:   valiable,
			Desc:       desc,
			Level:      level,
			Disable:    disable,
		},
		CheckFn: fn,
	}
}

func (s *Rule) Check(config *storage.InspectConfig, db *gorm.DB, node ast.StmtNode) (errMsgs, warnMsgs string, err error) {
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

var Rules []*Rule

func init() {
	Rules = []*Rule{
		NewRule(SELECT_STMT_TABLE_MUST_EXIST, 0, "", 0, 0, false, SelectStmtTableMustExist, ""),
	}
}

func SelectStmtTableMustExist(db *gorm.DB, node ast.StmtNode, variable string) (string, error) {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return "", nil
	}
	fmt.Println("table: ", selectStmt.From)
	_ = selectStmt
	return "", nil
}
