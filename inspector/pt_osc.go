package inspector

import (
	"bytes"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"io/ioutil"
	"os"
	"sqle/log"
	"sqle/model"
	"strings"
	"sync"
	"text/template"
)

var ptTemplate = `pt-online-schema-change D={{.Schema}},t={{.Table}} \
--alter="{{.Alter}}" \
--host={{.Host}} \
--user={{.User}} \
--port={{.Port}} \
--ask-pass \
--print \
--execute`

var ptTemplateMutex sync.Mutex

func LoadPtTemplateFromFile(fileName string) error {
	log.Logger().Info("loading pt-online-schema-change template")
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger().Infof("file %s not found, skip", fileName)
			return nil
		}
		return err
	}
	ptTemplateMutex.Lock()
	ptTemplate = string(b)
	ptTemplateMutex.Unlock()
	log.Logger().Info("loaded pt-online-schema-change template")
	return nil
}

const (
	OSC_NO_UNIQUE_INDEX_AND_PRIMARY_KEY      = "至少要包含主键或者唯一键索引才能使用 pt-online-schema-change"
	OSC_AVOID_ADD_UNIQUE_INDEX               = "添加唯一键使用 pt-online-schema-change，可能会导致数据丢失，在数据迁移到新表时使用了insert ignore"
	OSC_AVOID_RENAME_TABLE                   = "pt-online-schema-change 不支持使用rename table 来重命名表"
	OSC_AVOID_ADD_NOT_NULL_NO_DEFAULT_COLUMN = "非空字段必须设置默认值，不然 pt-online-schema-change 会执行失败"
)

// see https://www.percona.com/doc/percona-toolkit/LATEST/pt-online-schema-change.html
func (i *Inspect) generateOSCCommandLine(node ast.StmtNode) (string, error) {
	// just support mysql
	if i.Task.Instance.DbType != model.DB_TYPE_MYSQL {
		return "", nil
	}
	stmt, ok := node.(*ast.AlterTableStmt)
	if !ok {
		return "", nil
	}
	tableName := i.getTableName(stmt.Table)
	tableSize, err := i.getTableSize(tableName)
	if err != nil {
		return "", err
	}

	if int64(tableSize) < GetConfigInt(CONFIG_DDL_OSC_SIZE_LIMIT) {
		return "", err
	}

	createTableStmt, exist, err := i.getCreateTableStmt(tableName)
	if !exist || err != nil {
		return "", err
	}

	// In almost all cases a PRIMARY KEY or UNIQUE INDEX needs to be present in the table.
	// This is necessary because the tool creates a DELETE trigger to keep the new table
	// updated while the process is running.
	if !hasPrimaryKey(createTableStmt) && !hasUniqIndex(createTableStmt) {
		return OSC_NO_UNIQUE_INDEX_AND_PRIMARY_KEY, nil
	}

	// The RENAME clause cannot be used to rename the table.
	if len(getAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameTable)) > 0 {
		return OSC_AVOID_RENAME_TABLE, nil
	}

	// If you add a column without a default value and make it NOT NULL, the tool will fail,
	// as it will not try to guess a default value for you; You must specify the default.
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		for _, col := range spec.NewColumns {
			if HasOneInOptions(col.Options, ast.ColumnOptionNotNull) {
				if !HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					return OSC_AVOID_ADD_NOT_NULL_NO_DEFAULT_COLUMN, nil
				}
			}
		}
	}

	// Avoid pt-online-schema-change to run if the specified statement for --alter is trying to add an unique index.
	// Since pt-online-schema-change uses INSERT IGNORE to copy rows to the new table, if the row being written
	// produces a duplicate key, it will fail silently and data will be lost.
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintUniq:
			return OSC_AVOID_ADD_UNIQUE_INDEX, nil
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
		change := alterTableSpecFormat(spec)
		if change != "" {
			changes = append(changes, change)
		}
	}

	if len(changes) <= 0 {
		return "", nil
	}

	ptTemplateMutex.Lock()
	text := ptTemplate
	ptTemplateMutex.Unlock()
	tp, err := template.New("tp").Parse(text)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	err = tp.Execute(buff, map[string]interface{}{
		"Alter":  strings.Join(changes, ","),
		"Host":   i.Task.Instance.Host,
		"Port":   i.Task.Instance.Port,
		"User":   i.Task.Instance.User,
		"Schema": i.getSchemaName(stmt.Table),
		"Table":  stmt.Table.Name.String(),
	})
	return buff.String(), err
}
