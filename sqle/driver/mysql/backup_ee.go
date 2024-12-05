//go:build enterprise
// +build enterprise

package mysql

import (
	"context"

	"fmt"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"golang.org/x/text/language"
)

const (
	BackupStrategyNone        string = "none"         // 不备份(不支持备份、无需备份、选择不备份)
	BackupStrategyReverseSql  string = "reverse_sql"  // 备份为反向SQL
	BackupStrategyOriginalRow string = "original_row" // 备份为原始行
	BackupStrategyManually    string = "manual"       // 标记为人工备份
)

func (i *MysqlDriverImpl) Backup(ctx context.Context, backupStrategy string, sql string) (BackupSql []string, ExecuteInfo string, err error) {
	if i.IsOfflineAudit() {
		return nil, "暂不支持不连库备份", nil
	}
	if i.HasInvalidSql {
		return nil, "SQL无法正常解析，无法进行备份", nil
	}
	nodes, err := i.ParseSql(sql)
	if err != nil {
		return nil, "", err
	}
	if len(nodes) == 0 {
		return []string{}, fmt.Sprintf("in plugin, when backup ParseSql for sql %v extract 0 ast node of sql", sql), nil
	}
	var info i18nPkg.I18nStr
	var rollbackSqls []string
	switch backupStrategy {
	case BackupStrategyReverseSql:
		rollbackSqls, info, err = i.GenerateRollbackSqls(nodes[0])
		if err != nil {
			i.Logger().Errorf("in plugin when backup GenerateRollbackSqls for sql %v failed: %v", sql, err)
			return nil, err.Error(), err
		}
	case BackupStrategyOriginalRow:
		rollbackSqls, info, err = i.GetOriginalRow(nodes[0])
		if err != nil {
			i.Logger().Errorf("in plugin when backup GetOriginalRow for sql %v failed: %v", sql, err)
			return nil, err.Error(), err
		}
	case BackupStrategyManually, BackupStrategyNone:
	default:
		return []string{}, fmt.Sprintf("不支持的备份类型: %v,未执行备份", backupStrategy), nil
	}
	i.Ctx.UpdateContext(nodes[0])
	ExecuteInfo = info.GetStrInLang(language.Chinese)

	if ExecuteInfo == "" {
		if len(rollbackSqls) == 0 {
			ExecuteInfo = "无影响范围或不支持回滚，无备份回滚语句"
		} else {
			ExecuteInfo = "备份成功"
		}
	}
	return rollbackSqls, ExecuteInfo, nil
}

func (i *MysqlDriverImpl) RecommendBackupStrategy(ctx context.Context, sql string) (*driver.RecommendBackupStrategyRes, error) {
	var BackupStrategy string
	var BackupStrategyTip string
	var TablesRefer []string
	var SchemasRefer []string
	nodes, err := i.ParseSql(sql)
	if err != nil {
		i.Logger().Errorf("in plugin when RecommendBackupStrategy ParseSql %v failed: %v", sql, err)
		return nil, err
	}
	switch nodes[0].(type) {
	case *ast.DeleteStmt:
		BackupStrategy = BackupStrategyOriginalRow
		BackupStrategyTip = "删除行操作，建议使用行备份，完整保存被删除行的原始数据，便于精确恢复"
	case *ast.UpdateStmt:
		BackupStrategy = BackupStrategyOriginalRow
		BackupStrategyTip = "更新行操作，推荐使用行备份，保存更新前的完整行数据，也可以使用反向SQL仅备份受影响的列"
	case *ast.InsertStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "插入行操作，推荐备份为反向SQL，通过DELETE语句实现快速回滚"

	case *ast.CreateTableStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "建表操作，推荐备份为反向SQL，回滚时仅需删除创建的表，注意：回滚前应进行人工备份"
	case *ast.DropTableStmt:
		BackupStrategy = BackupStrategyManually
		BackupStrategyTip = "删表操作，强烈建议手工全量备份，避免不可逆的数据丢失风险"
	case *ast.AlterTableStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "表结构变更操作，推荐备份为反向SQL，回滚时使用对应的反向DDL语句，同时建议保存原始表结构定义"

	case *ast.CreateIndexStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "创建索引操作，推荐备份为反向SQL，回滚时仅需删除创建的索引"

	case *ast.DropIndexStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "删除索引操作，推荐备份为反向SQL，回滚时重新创建被删除的索引"

	case *ast.CreateDatabaseStmt:
		BackupStrategy = BackupStrategyReverseSql
		BackupStrategyTip = "创建数据库操作，推荐备份为反向SQL，回滚时仅需删除创建的数据库"

	case *ast.DropDatabaseStmt:
		BackupStrategy = BackupStrategyManually
		BackupStrategyTip = "删除数据库操作，推荐进行手工备份，涉及所有表和数据，需要确保完全备份所有数据库对象和数据"
	case *ast.SelectStmt:
		BackupStrategy = BackupStrategyNone
		BackupStrategyTip = "SELECT语句，未对数据库进行变更，无需备份"
	default:
		BackupStrategy = BackupStrategyManually
		BackupStrategyTip = "暂不支持备份该SQL，请手工备份"
	}
	tableSourceExtractor := &util.TableSourceExtractor{TableSources: make(map[string]*ast.TableSource)}
	nodes[0].Accept(tableSourceExtractor)
	for _, tableSource := range tableSourceExtractor.TableSources {
		if tableName, ok := tableSource.Source.(*ast.TableName); ok {
			TablesRefer = append(TablesRefer, tableName.Name.O)
			if tableName.Schema.String() != "" {
				SchemasRefer = append(SchemasRefer, tableName.Schema.String())
			} else {
				SchemasRefer = append(SchemasRefer, i.Ctx.CurrentSchema())
			}
		}
	}
	return &driver.RecommendBackupStrategyRes{
		BackupStrategyTip: BackupStrategyTip,
		BackupStrategy:    BackupStrategy,
		TablesRefer:       TablesRefer,
		SchemasRefer:      SchemasRefer,
	}, nil
}
