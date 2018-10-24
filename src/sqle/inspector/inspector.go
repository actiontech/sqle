package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type Inspector struct {
	Config        map[string]*model.Rule
	Db            model.Instance
	CurrentSchema string
	Sql           string
	sqlStmt       string
	dbConn        *executor.Conn
	isConnected   bool
}

type Result struct {
	Level   string
	Message string
}

type Results []Result

// level find highest level in result
func (rs Results) level() string {
	level := model.RULE_LEVEL_NOTICE
	for _, result := range rs {
		if model.RuleLevelMap[level] < model.RuleLevelMap[result.Level] {
			level = result.Level
		}
	}
	return level
}

func (rs Results) message() string {
	messages := make([]string, len(rs))
	for n, result := range rs {
		messages[n] = fmt.Sprintf("[%s]%s", result.Level, result.Message)
	}
	return strings.Join(messages, "\n")
}

func NewInspector(config map[string]*model.Rule, db model.Instance, Schema, sql string) *Inspector {
	return &Inspector{
		Config:        config,
		Db:            db,
		CurrentSchema: Schema,
		Sql:           sql,
	}
}

func (i *Inspector) getDbConn() (*executor.Conn, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	fmt.Println("get conn")
	conn, err := executor.NewConn(i.Db.DbType, i.Db.User, i.Db.Password, i.Db.Host, i.Db.Port, i.CurrentSchema)
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

func (i *Inspector) closeDbConn() {
	if i.isConnected {
		i.dbConn.Close()
		i.isConnected = false
	}
}

func (i *Inspector) Inspect() ([]*model.CommitSql, error) {
	defer i.closeDbConn()

	stmts, err := parseSql(i.Db.DbType, i.Sql)
	if err != nil {
		return nil, err
	}
	commitSqls := make([]*model.CommitSql, len(stmts))
	for n, stmt := range stmts {
		var results Results
		var err error

		switch s := stmt.(type) {
		case *ast.SelectStmt:
			results, err = i.inspectSelectStmt(s)
		default:
		}
		if err != nil {
			return nil, err
		}
		commitSqls[n] = &model.CommitSql{
			Number:        uint(n),
			Sql:           stmt.Text(),
			InspectResult: results.message(),
			InspectLevel:  results.level(),
		}
	}
	return commitSqls, nil
}

func (i *Inspector) inspectSelectStmt(stmt *ast.SelectStmt) (Results, error) {
	results := Results{}

	// check table must exist
	tablerefs := stmt.From.TableRefs
	tables := getTables(tablerefs)
	tablesName := map[string]struct{}{}
	for _, t := range tables {
		tablesName[getTableName(t)] = struct{}{}
	}
	conn, err := i.getDbConn()
	if err != nil {
		return results, err
	}
	notExistTables := []string{}
	for name, _ := range tablesName {
		exist := conn.HasTable(name)
		if conn.Error != nil {
			return results, conn.Error
		}
		if !exist {
			notExistTables = append(notExistTables, name)
		}
	}
	if len(notExistTables) > 0 {
		results = append(results, Result{model.RULE_LEVEL_ERROR,
			fmt.Sprintf("table %s is not exist", strings.Join(notExistTables, ", "))})
	}

	return results, nil
}

func (i *Inspector) inspectAlterTableStmt(stmt *ast.AlterTableSpec) {

}
