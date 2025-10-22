package ai

import (
	"fmt"
	"math"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	tidbTypes "github.com/pingcap/tidb/types"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00112 = "SQLE00112"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00112,
			Desc:       plocale.Rule00112Desc,
			Annotation: plocale.Rule00112Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00112Message,
		Func:    RuleSQLE00112,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== 规则总体逻辑 ====

规则目标：检查 WHERE/ON/USING 子句中条件字段与值的数据类型是否一致，避免隐式类型转换导致的性能问题。

【解析器的局限性】
当前使用的 SQL 解析器在判断常量类型时存在局限性：
1. 数值常量类型判断不准确
   - 例如：常量 100 会被解析器判断为 TypeLong（BIGINT）
   - 实际上：100 可以是 TINYINT/SMALLINT/INT/BIGINT 等多种类型
   - 问题：如果严格按照解析器类型判断，会导致 WHERE id = 100（id为INT列）被误报

2. 字符串常量类型判断
   - 短字符串可能被判断为 VARCHAR
   - 长字符串可能被判断为 TEXT
   - 实际上：字符串常量可以与 CHAR/VARCHAR/TEXT 类型的列兼容比较

3. BLOB 常量判断
   - 解析器可能根据长度判断为不同的 BLOB 子类型
   - 实际上：BLOB 各子类型之间应该兼容

因此，本规则采用"大类匹配"和"范围检查"的策略：
- 不完全依赖解析器给出的精确类型
- 而是判断常量的值是否在列类型的合理范围内
- 只要 MySQL 转换常量值而非列，就认为兼容

核心判断原则：
1. 【列与列比较】- 严格匹配原则
   - 两列的数据类型必须完全一致
   - 任何类型不一致都会报错
   - 示例：INT 列与 VARCHAR 列比较 → 报错
   - 适用场景：
     a) WHERE/ON 子句中的列与列比较：WHERE t1.col1 = t2.col2
     b) USING 子句：JOIN ... USING (column_name) - 检查两表中同名列的类型是否一致

2. 【列与值比较】- 宽松兼容原则
   - 判断依据：值的大类与列的大类是否一致，且 MySQL 会转换值而非列
   - 如果 MySQL 在比较时转换常量值（而非列），则不影响列的索引使用，不报错
   - 如果 MySQL 需要转换列，则会导致索引失效，需要报错
   - 适用场景：WHERE/ON 子句中的列与常量比较：WHERE column = 'value'

   兼容场景（不报错）：
   a) 整数类型之间：TINYINT/SMALLINT/INT/BIGINT 常量与列比较，只要数值在列范围内即兼容
   b) 字符串类型之间：VARCHAR/CHAR/TEXT 常量与列比较，MySQL 会转换值，兼容
   c) BLOB 类型之间：BLOB/TINYBLOB/MEDIUMBLOB/LONGBLOB 互相兼容
   d) 字符串与 BLOB：VARCHAR/CHAR 常量与 BLOB 列比较，MySQL 会转换值，兼容

   不兼容场景（报错）：
   a) 跨大类比较：整数与字符串、整数与日期时间等
   b) 字符串与日期时间：字符串常量与 DATE/DATETIME 列比较（强制使用显式转换函数）
   c) DECIMAL 与整数：小数类型与整数类型不同

3. 【USING 子句特殊处理】
   - USING (column_name) 语法用于 JOIN，表示用同名列进行连接
   - 检查逻辑：获取 USING 中指定列在所有关联表中的类型，如果存在类型不一致则报错
   - 示例：
     SELECT * FROM t1 JOIN t2 USING (id)
     → 检查 t1.id 和 t2.id 的类型是否一致
     → 如果 t1.id 是 INT，t2.id 是 VARCHAR，则报错

==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00112): "在 MySQL 中，禁止WHERE子句中条件字段与值的数据类型不一致."
您应遵循以下逻辑：
1. 对所有 DML 语句，解析 SQL 语句，获取所有的 WHERE 和 ON 条件字段。
   1. 若条件两侧均为列字段：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查两列字段类型是否一致，不一致则报告违规。

   2. 若左侧为列字段，右侧为常量：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查常量类型是否与列字段类型一致，不一致则报告违规。

   3. 若左侧为常量，右侧为列字段：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查常量类型是否与列字段类型一致，不一致则报告违规。

   4. 若条件两侧均为常量：
      - 检查两常量类型是否一致，不一致则报告违规。

2. 对所有 DML 语句，解析 SQL 语句，获取 USING 条件字段。
   1. 使用辅助函数GetCreateTableStmt获取USING字段在关联表中的数据类型。
   2. 若 USING 字段在两表中的数据类型不一致，则报告违规。

3. 对所有 DML 语句，解析 SQL 语句中的 WHERE 子句，考虑其在 SELECT、UPDATE、DELETE 语句中的位置。
   - 对于 SELECT 语句中的 WHERE 子句：用于筛选符合条件的记录。
   - 对于 UPDATE 语句中的 WHERE 子句：用于指定更新哪些记录。
   - 对于 DELETE 语句中的 WHERE 子句：用于指定删除哪些记录。
==== Prompt end ====
*/

// ==== Rule code start ====

// getValueFromExpr 从ValueExpr中提取实际的数值
func getValueFromExpr(expr *parserdriver.ValueExpr) (int64, bool) {
	switch expr.Datum.Kind() {
	case tidbTypes.KindInt64:
		return expr.GetInt64(), true
	case tidbTypes.KindUint64:
		val := expr.GetUint64()
		if val <= math.MaxInt64 {
			return int64(val), true
		}
	}
	return 0, false
}

// isIntegerValueInRange 检查整数值是否在目标类型范围内
func isIntegerValueInRange(value int64, targetType byte) bool {
	switch targetType {
	case mysql.TypeTiny:
		return value >= -128 && value <= 127
	case mysql.TypeShort:
		return value >= -32768 && value <= 32767
	case mysql.TypeInt24:
		return value >= -8388608 && value <= 8388607
	case mysql.TypeLong:
		return value >= -2147483648 && value <= 2147483647
	case mysql.TypeLonglong:
		return true
	}
	return false
}

// isIntegerValueInUnsignedRange 检查整数值是否在UNSIGNED类型范围内
func isIntegerValueInUnsignedRange(value int64, targetType byte) bool {
	switch targetType {
	case mysql.TypeTiny:
		return value >= 0 && value <= 255
	case mysql.TypeShort:
		return value >= 0 && value <= 65535
	case mysql.TypeInt24:
		return value >= 0 && value <= 16777215
	case mysql.TypeLong:
		return value >= 0 && value <= 4294967295
	case mysql.TypeLonglong:
		return value >= 0
	}
	return false
}

// isIntegerType 检查是否为整数类型
func isIntegerType(tp byte) bool {
	switch tp {
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24,
		mysql.TypeLong, mysql.TypeLonglong:
		return true
	}
	return false
}

// isStringType 检查是否为字符串类型
func isStringType(tp byte) bool {
	switch tp {
	case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString:
		return true
	}
	return false
}

// isBlobType 检查是否为BLOB类型
func isBlobType(tp byte) bool {
	switch tp {
	case mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
		return true
	}
	return false
}

// isDecimalType 检查是否为DECIMAL类型
func isDecimalType(tp byte) bool {
	switch tp {
	case mysql.TypeNewDecimal, mysql.TypeDecimal:
		return true
	}
	return false
}

// isDateTimeType 检查是否为日期时间类型
func isDateTimeType(tp byte) bool {
	switch tp {
	case mysql.TypeDate, mysql.TypeDatetime, mysql.TypeTimestamp:
		return true
	}
	return false
}

// checkIntegerCompatibility 检查整数类型的兼容性
func checkIntegerCompatibility(leftType, rightType byte, leftExpr, rightExpr ast.ExprNode, leftIsUnsigned, rightIsUnsigned bool) bool {
	leftValue, leftIsValue := leftExpr.(*parserdriver.ValueExpr)
	rightValue, rightIsValue := rightExpr.(*parserdriver.ValueExpr)

	// 常量与列的比较：检查常量数值是否在列类型范围内
	if leftIsValue && !rightIsValue {
		if val, ok := getValueFromExpr(leftValue); ok {
			if rightIsUnsigned {
				return isIntegerValueInUnsignedRange(val, rightType)
			}
			return isIntegerValueInRange(val, rightType)
		}
	}

	if rightIsValue && !leftIsValue {
		if val, ok := getValueFromExpr(rightValue); ok {
			if leftIsUnsigned {
				return isIntegerValueInUnsignedRange(val, leftType)
			}
			return isIntegerValueInRange(val, leftType)
		}
	}

	// 常量与常量：如果都是整数类型则兼容
	if leftIsValue && rightIsValue {
		return true
	}

	return false
}

// checkStringCompatibility 检查字符串类型的兼容性
// VARCHAR、CHAR、TEXT 类型在常量与列比较时互相兼容
func checkStringCompatibility(leftType, rightType byte, leftExpr, rightExpr ast.ExprNode) bool {
	_, leftIsValue := leftExpr.(*parserdriver.ValueExpr)
	_, rightIsValue := rightExpr.(*parserdriver.ValueExpr)

	// 常量与列的比较：字符串类型互相兼容
	if (leftIsValue && !rightIsValue) || (!leftIsValue && rightIsValue) {
		return true
	}

	// 常量与常量：兼容
	if leftIsValue && rightIsValue {
		return true
	}

	return false
}

// checkBlobStringCompatibility 检查BLOB与字符串类型的兼容性
// VARCHAR/CHAR/TEXT 常量可以与 BLOB 列比较
func checkBlobStringCompatibility(leftType, rightType byte, leftExpr, rightExpr ast.ExprNode) bool {
	_, leftIsValue := leftExpr.(*parserdriver.ValueExpr)
	_, rightIsValue := rightExpr.(*parserdriver.ValueExpr)

	// 只允许常量与列的比较
	if (leftIsValue && !rightIsValue) || (!leftIsValue && rightIsValue) {
		return true
	}

	return false
}

// isTypeCompatible 检查两个类型是否兼容
//
// 规则判断总体逻辑：
// 1. 【列与列比较】- 严格匹配原则
//   - 两列数据类型必须完全一致，任何类型不一致都返回 false（报错）
//   - 原因：列与列比较时，MySQL可能转换其中一列，导致索引失效
//
// 2. 【列与值比较】- 宽松兼容原则
//
//   - 判断依据：值的大类与列的大类是否一致，且 MySQL 会转换值而非列
//
//   - 如果 MySQL 转换常量值，不影响列索引使用，返回 true（不报错）
//
//   - 如果需要转换列，会导致索引失效，返回 false（报错）
//
//     兼容场景示例：
//
//   - INT 常量 100 与 BIGINT 列比较 → true（MySQL转换值）
//
//   - VARCHAR 常量 'abc' 与 CHAR 列比较 → true（MySQL转换值）
//
//   - VARCHAR 常量与 BLOB 列比较 → true（MySQL转换值）
//
//     不兼容场景示例：
//
//   - INT 列与 VARCHAR 列比较 → false（列与列，严格匹配）
//
//   - INT 常量与 VARCHAR 列比较 → false（跨大类，需要转换列）
//
//   - 字符串常量与 DATE 列比较 → false（需要显式转换函数）
func isTypeCompatible(leftType, rightType byte, leftExpr, rightExpr ast.ExprNode, leftIsUnsigned, rightIsUnsigned bool) bool {
	// 类型完全相同，直接返回兼容
	if leftType == rightType {
		return true
	}

	// 判断是否为列表达式
	_, leftIsColumn := leftExpr.(*ast.ColumnNameExpr)
	_, rightIsColumn := rightExpr.(*ast.ColumnNameExpr)

	// 【关键】列与列比较，必须严格匹配
	// 任何类型不一致都不兼容，因为可能导致其中一列被转换，造成索引失效
	if leftIsColumn && rightIsColumn {
		return false
	}

	// 以下都是常量与列的比较场景
	// 原则：只要 MySQL 转换常量值（而非列），就认为兼容

	// 场景1: 整数类型之间的兼容性
	if isIntegerType(leftType) && isIntegerType(rightType) {
		return checkIntegerCompatibility(leftType, rightType, leftExpr, rightExpr, leftIsUnsigned, rightIsUnsigned)
	}

	// 场景2: 字符串类型之间的兼容性（VARCHAR/CHAR/TEXT）
	if isStringType(leftType) && isStringType(rightType) {
		return checkStringCompatibility(leftType, rightType, leftExpr, rightExpr)
	}

	// 场景3: BLOB类型之间的兼容性
	if isBlobType(leftType) && isBlobType(rightType) {
		return true
	}

	// 场景4: 字符串与BLOB的兼容性
	if (isStringType(leftType) && isBlobType(rightType)) || (isBlobType(leftType) && isStringType(rightType)) {
		return checkBlobStringCompatibility(leftType, rightType, leftExpr, rightExpr)
	}

	// 场景5: DECIMAL类型之间的兼容性
	if isDecimalType(leftType) && isDecimalType(rightType) {
		return true
	}

	// 场景6: 明确不兼容的场景
	// 字符串与日期时间类型不兼容（强制使用显式转换）
	if (isStringType(leftType) && isDateTimeType(rightType)) || (isDateTimeType(leftType) && isStringType(rightType)) {
		return false
	}

	// 场景7: DECIMAL与整数类型不兼容（小数与整数类型不同）
	if (isIntegerType(leftType) && isDecimalType(rightType)) || (isDecimalType(leftType) && isIntegerType(rightType)) {
		return false
	}

	// 其他情况：不兼容
	return false
}

func RuleSQLE00112(input *rulepkg.RuleHandlerInput) error {
	// 当前解析器的弊端：对于值类型，没法准确判断其类型，例如100会被判断成TypeLong实际上应该是INT
	// 因此这里需要手动判断类型是否兼容
	var defaultTable string
	var alias []*util.TableAliasInfo
	getTableName := func(col *ast.ColumnNameExpr) string {
		if col.Name.Table.L != "" {
			for _, a := range alias {
				if a.TableAliasName == col.Name.Table.String() {
					return a.TableName
				}
			}
			return col.Name.Table.O
		}
		return defaultTable
	}

	// 内部辅助函数：获取表达式的类型和UNSIGNED信息
	getExprTypeInfo := func(expr ast.ExprNode) (byte, bool, error) {
		switch node := expr.(type) {
		case *ast.ColumnNameExpr:
			// 获取列的表名
			tableName := getTableName(node)

			// 获取CREATE TABLE语句
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, &ast.TableName{Name: model.NewCIStr(tableName)})
			if err != nil {
				return 0, false, err
			}
			// 获取列定义
			for _, colDef := range createTableStmt.Cols {
				if strings.EqualFold(util.GetColumnName(colDef), node.Name.Name.O) {
					isUnsigned := mysql.HasUnsignedFlag(colDef.Tp.Flag)
					return colDef.Tp.Tp, isUnsigned, nil
				}
			}
			return 0, false, fmt.Errorf("列未找到: %s", node.Name.Name.O)
		case *ast.FuncCallExpr:
			names := util.GetFuncName(node)
			if len(names) == 1 {
				switch strings.ToUpper(names[0]) {
				case "CURRENT_DATE":
					return mysql.TypeDate, false, nil
				case "CURRENT_TIME", "NOW":
					return mysql.TypeDatetime, false, nil
				case "CURRENT_TIMESTAMP":
					return mysql.TypeTimestamp, false, nil
				case "CURDATE":
					return mysql.TypeDate, false, nil
				}
			}
			return 0, false, fmt.Errorf("不支持的函数: %s", strings.Join(names, "."))
		case *parserdriver.ValueExpr:
			// 直接使用解析器确定的精确类型，而不是通过 Datum.Kind() 手动映射
			return node.Type.Tp, false, nil
		default:
			return 0, false, fmt.Errorf("不支持的表达式类型: %T", expr)
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UnionStmt:
		// 获取DELETE、UPDATE语句的表
		alias = util.GetTableAliasInfoFromJoin(util.GetFirstJoinNodeFromStmt(stmt))

		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// 获取默认表名
			if t := util.GetDefaultTable(selectStmt); t != nil {
				defaultTable = t.Name.O
			}

			// 获取表别名信息
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				// 获取FROM子句中的所有表及其别名
				alias = append(alias, util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)...)
			}

			// 获取所有WHERE和ON条件表达式
			whereExprs := util.GetWhereExprFromDMLStmt(selectStmt)
			onExprs := []ast.ExprNode{}
			usingExprs := [][]*ast.ColumnName{}
			for _, join := range util.GetAllJoinsFromNode(selectStmt) {
				if join != nil && join.On != nil {
					onExprs = append(onExprs, join.On.Expr)
				}
				if join != nil && join.Using != nil {
					usingExprs = append(usingExprs, join.Using)
				}
			}

			// Combine WHERE and ON expressions
			allExprs := append(whereExprs, onExprs...)

			// 遍历所有条件表达式
			for _, expr := range allExprs {
				// 遍历表达式中的所有二元操作
				util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
					binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
					if !ok {
						return false
					}

					// 仅处理等于操作符
					if binExpr.Op != opcode.EQ {
						return false
					}

					// 获取左侧和右侧的类型和UNSIGNED信息
					leftType, leftIsUnsigned, err := getExprTypeInfo(binExpr.L)
					if err != nil {
						// 记录错误但不阻止其他检查
						log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
						return false
					}

					rightType, rightIsUnsigned, err := getExprTypeInfo(binExpr.R)
					if err != nil {
						log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
						return false
					}

					// 比较类型是否兼容
					if !isTypeCompatible(leftType, rightType, binExpr.L, binExpr.R, leftIsUnsigned, rightIsUnsigned) {
						// 报告违规
						rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
					}
					return false
				}, expr)
			}
			for _, using := range usingExprs {
				for _, column := range using {
					colTyp := make(map[byte]struct{})
					for _, aliasInfo := range alias {
						tableName := &ast.TableName{Schema: model.NewCIStr(aliasInfo.SchemaName), Name: model.NewCIStr(aliasInfo.TableName)}
						createTableStmt, err := util.GetCreateTableStmt(input.Ctx, tableName)
						if err != nil {
							log.NewEntry().Errorf("获取表 %s 的CREATE TABLE语句失败: %v", tableName.Name.L, err)
							continue
						}

						for _, colDef := range createTableStmt.Cols {
							if strings.EqualFold(util.GetColumnName(colDef), column.Name.L) {
								colTyp[colDef.Tp.Tp] = struct{}{}
							}
						}
					}

					if len(colTyp) > 1 {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
					}
				}
			}
		}

	}

	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		// 获取默认表名
		if t := util.GetTableNames(stmt); len(t) != 0 {
			defaultTable = t[0].Name.O
		}

		// 获取表别名信息
		if stmt.TableRefs != nil && stmt.TableRefs.TableRefs != nil {
			// 获取FROM子句中的所有表及其别名
			alias = util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		}

		// 获取所有WHERE和ON条件表达式
		whereExprs := util.GetWhereExprFromDMLStmt(stmt)
		onExprs := []ast.ExprNode{}
		joinExpr := util.GetFirstJoinNodeFromStmt(stmt)
		if joinExpr != nil && joinExpr.On != nil {
			onExprs = append(onExprs, joinExpr.On.Expr)
		}
		// Combine WHERE and ON expressions
		allExprs := append(whereExprs, onExprs...)

		// 遍历所有条件表达式
		for _, expr := range allExprs {
			// 遍历表达式中的所有二元操作
			util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
				binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
				if !ok {
					return false
				}

				// 仅处理等于操作符
				if binExpr.Op != opcode.EQ {
					return false
				}

				// 获取左侧和右侧的类型和UNSIGNED信息
				leftType, leftIsUnsigned, err := getExprTypeInfo(binExpr.L)
				if err != nil {
					// 记录错误但不阻止其他检查
					log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
					return false
				}

				rightType, rightIsUnsigned, err := getExprTypeInfo(binExpr.R)
				if err != nil {
					log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
					return false
				}

				// 比较类型是否兼容
				if !isTypeCompatible(leftType, rightType, binExpr.L, binExpr.R, leftIsUnsigned, rightIsUnsigned) {
					// 报告违规
					rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
				}
				return false
			}, expr)
		}
	case *ast.DeleteStmt:
		// 获取默认表名
		if t := util.GetTableNames(stmt); len(t) != 0 {
			defaultTable = t[0].Name.O
		}

		// 获取表别名信息
		if stmt.TableRefs != nil && stmt.TableRefs.TableRefs != nil {
			// 获取FROM子句中的所有表及其别名
			alias = util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		}

		// 获取所有WHERE和ON条件表达式
		whereExprs := util.GetWhereExprFromDMLStmt(stmt)
		onExprs := []ast.ExprNode{}
		joinExpr := util.GetFirstJoinNodeFromStmt(stmt)
		if joinExpr != nil && joinExpr.On != nil {
			onExprs = append(onExprs, joinExpr.On.Expr)
		}
		// Combine WHERE and ON expressions
		allExprs := append(whereExprs, onExprs...)
		// 遍历所有条件表达式
		for _, expr := range allExprs {
			// 遍历表达式中的所有二元操作
			util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
				binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
				if !ok {
					return false
				}

				// 仅处理等于操作符
				if binExpr.Op != opcode.EQ {
					return false
				}

				// 获取左侧和右侧的类型和UNSIGNED信息
				leftType, leftIsUnsigned, err := getExprTypeInfo(binExpr.L)
				if err != nil {
					// 记录错误但不阻止其他检查
					log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
					return false
				}

				rightType, rightIsUnsigned, err := getExprTypeInfo(binExpr.R)
				if err != nil {
					log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
					return false
				}

				// 比较类型是否兼容
				if !isTypeCompatible(leftType, rightType, binExpr.L, binExpr.R, leftIsUnsigned, rightIsUnsigned) {
					// 报告违规
					rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
				}
				return false
			}, expr)
		}
	}

	return nil
}

// ==== Rule code end ====
