package inspector

//import (
//	"fmt"
//	"github.com/pingcap/tidb/ast"
//	"sqle/executor"
//	"sqle/storage"
//	"strings"
//)
//
//func Inspect(config map[int]*storage.InspectConfig, task *storage.Task) ([]*storage.Sql, error) {
//	return nil, nil
//	sqls := []*storage.Sql{}
//	stmts, err := parseSql(task.Inst.DbType, task.ReqSql)
//	if err != nil {
//		return nil, err
//	}
//	conn, err := executor.OpenDbWithTask(task)
//	if err != err {
//		return nil, err
//	}
//	defer conn.Close()
//
//	for _, stmt := range stmts {
//		errMsgs := []string{}
//		warnMsgs := []string{}
//		for _, rule := range Rules {
//			errMsg, warnMsg, err := rule.check(config[rule.DfConfig.Code], conn, stmt)
//			if err != err {
//				return nil, err
//			}
//			errMsgs = append(errMsgs, errMsg)
//			warnMsgs = append(warnMsgs, warnMsg)
//		}
//		sql := &storage.Sql{}
//		//sql.CommitSql = stmt.Text()
//		//sql.InspectError = strings.Join(errMsgs, "\n")
//		//sql.InspectWarn = strings.Join(warnMsgs, "\n")
//		sqls = append(sqls, sql)
//	}
//	return sqls, nil
//}
//
//var Rules []*Rule
//
//type Rule struct {
//	DfConfig *storage.InspectConfig
//	CheckFn  func(conn *executor.Conn, node ast.StmtNode, config string) (string, error)
//}
//
//func NewRule(code int, fn func(conn *executor.Conn, node ast.StmtNode, config string) (string, error)) *Rule {
//	return &Rule{
//		DfConfig: storage.DfConfigMap[code],
//		CheckFn:  fn,
//	}
//}
//
//func (s *Rule) check(config *storage.InspectConfig, db *executor.Conn, node ast.StmtNode) (errMsgs, warnMsgs string, err error) {
//	var currentConfig *storage.InspectConfig
//	if config != nil {
//		currentConfig = config
//	} else {
//		currentConfig = s.DfConfig
//	}
//
//	if currentConfig.Disable {
//		return errMsgs, warnMsgs, nil
//	}
//
//	msg, err := s.CheckFn(db, node, currentConfig.Variable)
//	if err != nil {
//		return errMsgs, warnMsgs, err
//	}
//
//	switch currentConfig.Level {
//	case storage.RULE_LEVEL_ERROR:
//		errMsgs = msg
//	case storage.RULE_LEVEL_WARN:
//		warnMsgs = msg
//	}
//	return errMsgs, warnMsgs, nil
//}
//
//func init() {
//	Rules = []*Rule{
//		NewRule(storage.SELECT_STMT_TABLE_MUST_EXIST, selectStmtTableMustExist),
//	}
//}
//
//func selectStmtTableMustExist(conn *executor.Conn, node ast.StmtNode, variable string) (string, error) {
//	selectStmt, ok := node.(*ast.SelectStmt)
//	if !ok {
//		return "", nil
//	}
//	tablerefs := selectStmt.From.TableRefs
//	tables := getTables(tablerefs)
//
//	tablesName := map[string]struct{}{}
//	for _, t := range tables {
//		tablesName[getTableName(t)] = struct{}{}
//	}
//	msgs := []string{}
//	for name, _ := range tablesName {
//		exist := conn.HasTable(name)
//		if conn.Error != nil {
//			return "", nil
//		}
//		if !exist {
//			msgs = append(msgs, name)
//		}
//	}
//	if len(msgs) > 0 {
//		return fmt.Sprintf("%v is not exist", strings.Join(msgs, ", ")), nil
//	}
//	return "", nil
//}
