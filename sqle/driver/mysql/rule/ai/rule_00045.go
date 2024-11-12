package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00045 = "SQLE00045"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00045,
			Desc:       "避免在分页查询中使用过大偏移量",
			Annotation: "在数据库中，分页查询通常使用 LIMIT 和 OFFSET 语句进行。当数据量较大时，使用大的偏移量（OFFSET）进行分页查询可能会导致性能下降，因为数据库需要跳过大量的行来获得所需的结果集。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "10000",
					Desc:  "最大偏移量",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "避免在分页查询中使用过大偏移量, 最大偏移量:%v",
		AllowOffline: true,
		Func:         RuleSQLE00045,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00045): "在 MySQL 中，避免在分页查询中使用过大偏移量.最大偏移量:10000"
您应遵循以下逻辑：
1. 对于 “SELECT ... ” 语句
  1. 定义一个集合
  2. 把语句中 LIMIT 子句中的偏移量存入集合
  3. 从集合中读取偏移量数值，与规则变量max_offset_size 对比，如果比它大，则报告违反规则。
2. 对于 “INSERT ... SELECT ” 语句，执行与上述同样检查。
3. 对于 “... UNION ALL ... ” 语句，执行与上述同样检查。
4. 对于 "UNION ... " 语句, 对于其中的所有SELECT子句进行与 “SELECT ... ” 语句相同的检查。
5. 对于嵌套在其他DML语句（如UPDATE、DELETE）中的SELECT子句，执行与 “SELECT ... ” 语句相同的检查。
6. 对于WITH子句中的所有SELECT子句，执行与 “SELECT ... ” 语句相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00045(input *rulepkg.RuleHandlerInput) error {
	// 获取规则参数 max_offset_size
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	// 将 max_offset_size 转换为整数
	maxOffsetSize, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("parameter 'max_offset_size' must be an integer, got '%s'", param.Value)
	}

	// 根据输入节点的类型进行不同的处理
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		// 处理 SELECT 语句
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			offset := util.GetLimitOffsetValue(selectStmt)
			if offset > int64(maxOffsetSize) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00045, maxOffsetSize)
				return nil
			}
		}

	case *ast.InsertStmt:
		// 处理 INSERT ... SELECT 语句
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			offset := util.GetLimitOffsetValue(selectStmt)
			if offset > int64(maxOffsetSize) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00045, maxOffsetSize)
				return nil
			}
		}

	case *ast.UnionStmt:
		// 处理 UNION 和 UNION ALL 语句
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			offset := util.GetLimitOffsetValue(selectStmt)
			if offset > int64(maxOffsetSize) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00045, maxOffsetSize)
				return nil
			}
		}
		offset := util.GetLimitOffsetValueByUnionStmt(stmt)
		if offset > int64(maxOffsetSize) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00045, maxOffsetSize)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
