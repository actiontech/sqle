package mysql

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"
	"github.com/pingcap/parser/ast"
)

var ptTemplate = `pt-online-schema-change D={{.Schema}},t={{.Table}} --alter='{{.Alter}}' --host={{.Host}} --user={{.User}} --port={{.Port}} --ask-pass --print --execute`
var ptTemplateMutex sync.Mutex

func LoadPtTemplateFromFile(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	ptTemplateMutex.Lock()
	ptTemplate = string(b)
	ptTemplateMutex.Unlock()
	return nil
}

// generateOSCCommandLine generate pt-online-schema-change command-line statement;
// see https://www.percona.com/doc/percona-toolkit/LATEST/pt-online-schema-change.html.
func (i *MysqlDriverImpl) generateOSCCommandLine(node ast.Node) (i18nPkg.I18nStr, error) {
	if i.cnf.DDLOSCMinSize < 0 {
		return nil, nil
	}

	stmt, ok := node.(*ast.AlterTableStmt)
	if !ok {
		return nil, nil
	}
	tableSize, err := i.Ctx.GetTableSize(stmt.Table)
	if err != nil {
		return nil, err
	}

	if int64(tableSize) < i.cnf.DDLOSCMinSize {
		return nil, err
	}

	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if !exist || err != nil {
		return nil, err
	}

	// In almost all cases a PRIMARY KEY or UNIQUE INDEX needs to be present in the table.
	// This is necessary because the tool creates a DELETE trigger to keep the new table
	// updated while the process is running.
	if !util.HasPrimaryKey(createTableStmt) && !util.HasUniqIndex(createTableStmt) {
		return plocale.ShouldLocalizeAll(plocale.PTOSCNoUniqueIndexOrPrimaryKey), nil
	}

	// The RENAME clause cannot be used to rename the table.
	if len(util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameTable)) > 0 {
		return plocale.ShouldLocalizeAll(plocale.PTOSCAvoidRenameTable), nil
	}

	// If you add a column without a default value and make it NOT NULL, the tool will fail,
	// as it will not try to guess a default value for you; You must specify the default.
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		for _, col := range spec.NewColumns {
			if util.HasOneInOptions(col.Options, ast.ColumnOptionNotNull) {
				if !util.HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					return plocale.ShouldLocalizeAll(plocale.PTOSCAvoidNoDefaultValueOnNotNullColumn), nil
				}
			}
		}
	}

	// Avoid pt-online-schema-change to run if the specified statement for --alter is trying to add an unique index.
	// Since pt-online-schema-change uses INSERT IGNORE to copy rows to the new table, if the row being written
	// produces a duplicate key, it will fail silently and data will be lost.
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintUniq:
			return plocale.ShouldLocalizeAll(plocale.PTOSCAvoidUniqueIndex), nil
		}
	}

	// generate pt-online-change-schema command line
	changes := []string{}
	for _, spec := range stmt.Specs {
		/*
			DROP FOREIGN KEY constraint_name requires specifying _constraint_name rather than the real constraint_name.
			Due to a limitation in MySQL, pt-online-schema-change adds a leading underscore to foreign key constraint
			names when creating the new table.For example, to drop this constraint:
			CONSTRAINT `fk_foo` FOREIGN KEY (`foo_id`) REFERENCES `bar` (`foo_id`)
			You must specify --alter "DROP FOREIGN KEY _fk_foo".
		*/
		if spec.Tp == ast.AlterTableDropPrimaryKey {
			spec.Name = fmt.Sprintf("_%s", spec.Name)
		}
		change := util.AlterTableSpecFormat(spec)
		if change != "" {
			changes = append(changes, change)
		}
	}

	if len(changes) <= 0 {
		return nil, nil
	}

	ptTemplateMutex.Lock()
	text := ptTemplate
	ptTemplateMutex.Unlock()
	tp, err := template.New("tp").Parse(text)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBufferString("[osc]")
	err = tp.Execute(buff, map[string]interface{}{
		"Alter":  strings.Join(changes, ","),
		"Host":   i.inst.Host,
		"Port":   i.inst.Port,
		"User":   i.inst.User,
		"Schema": i.Ctx.GetSchemaName(stmt.Table),
		"Table":  stmt.Table.Name.String(),
	})
	return i18nPkg.ConvertStr2I18nAsDefaultLang(buff.String()), err
}
