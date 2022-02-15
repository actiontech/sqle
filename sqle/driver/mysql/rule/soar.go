package rule

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	parserDriver "github.com/pingcap/tidb/types/parser_driver"
)

var SoarRuleHandlers = []RuleHandler{
	{
		Rule: driver.Rule{ //select a as id, id , b as user  from mysql.user;
			Name:     "ALI003",
			Desc:     "别名不要与表或列的名字相同",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "这些别名(%v)与列名或表名相同",
		Func:    ALI003,
	},
	//{
	//	Rule: driver.Rule{ //用不了
	//		Name:     "ALI002",
	//		Desc:     "不建议给列通配符'*'设置别名",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "不建议给列通配符'*'设置别名",
	//	Func:    ALI002,
	//},
	{

		Rule: driver.Rule{ //ALTER TABLE logtest CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;
			Name:     "ALT001",
			Desc:     "修改表的默认字符集不会改表各个字段的字符集",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
		},
		Message: "修改表的默认字符集不会改表各个字段的字符集",
		Func:    ALT001,
	}, {
		Rule: driver.Rule{ //ALTER TABLE tbl DROP COLUMN col;
			Name:     "ALT003",
			Desc:     "删除列为高危操作",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message: "删除列为高危操作",
		Func:    ALT003,
	}, {
		Rule: driver.Rule{ //ALTER TABLE tbl DROP PRIMARY KEY;
			Name:     "ALT004",
			Desc:     "删除主键为高危操作",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message: "删除主键为高危操作",
		Func:    ALT004,
	}, {
		Rule: driver.Rule{ //ALTER TABLE tbl DROP FOREIGN KEY a;
			Name:     "ALT004_1",
			Desc:     "提示删除外键为高危操作",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message: "删除外键为高危操作",
		Func:    ALT004_1,
	},
	//{
	//	Rule: driver.Rule{ //select * from user where id like "%a";
	//		Name:     "ARG001",
	//		Desc:     "不建议使用前项通配符查找",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "不建议使用前项通配符查找",
	//	Func:    ARG001,
	//},
	{
		Rule: driver.Rule{ //select * from user where id like "a";
			Name:     "ARG002",
			Desc:     "不建议使用没有通配符的 LIKE 查询",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议使用没有通配符的 LIKE 查询",
		Func:    ARG002,
	}, {
		Rule: driver.Rule{ //SELECT * FROM tb WHERE col IN (NULL);
			Name:     "ARG004",
			Desc:     "IN (NULL)/NOT IN (NULL) 永远非真",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message: "IN (NULL)/NOT IN (NULL) 永远非真",
		Func:    ARG004,
	}, {
		Rule: driver.Rule{ //select * from user where id in (a);
			Name:     "ARG005",
			Desc:     "尽量不要使用IN",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "尽量不要使用IN",
		Func:    ARG005,
	},
	//{
	//	Rule: driver.Rule{ //select id from t where num is null
	//		Name:     "ARG006",
	//		Desc:     "应尽量避免在 WHERE 子句中对字段进行 NULL 值判断",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "应尽量避免在 WHERE 子句中对字段进行 NULL 值判断",
	//	Func:    ARG006,
	//}, {
	//	Rule: driver.Rule{ //select * from user where id like "%a";
	//		Name:     "ARG007",
	//		Desc:     "应避免使用模式匹配",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "应避免使用模式匹配",
	//	Func:    ARG007,
	//},
	{
		Rule: driver.Rule{ //select * from user where id = ' 1';
			Name:     "ARG009",
			Desc:     "引号中的字符串开头或结尾包含空格",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message: "引号中的字符串开头或结尾包含空格",
		Func:    ARG009,
	}, {
		Rule: driver.Rule{ //CREATE TABLE tb (a varchar(10) default '“');
			Name:     "ARG013",
			Desc:     "DDL 语句中使用了中文全角引号",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message: "DDL 语句中使用了中文全角引号",
		Func:    ARG013,
	}, {
		Rule: driver.Rule{ //select name from tbl where id < 1000 order by rand(1)
			Name:     "CLA002",
			Desc:     "不建议使用 ORDER BY RAND()",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议使用 ORDER BY RAND()",
		Func:    CLA002,
	}, {
		Rule: driver.Rule{ //select col1,col2 from tbl group by 1
			Name:     "CLA004",
			Desc:     "不建议对常量进行 GROUP BY",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议对常量进行 GROUP BY",
		Func:    CLA004,
	}, {
		Rule: driver.Rule{ //select c1,c2,c3 from t1 where c1='foo' order by c2 desc, c3 asc
			Name:     "CLA007",
			Desc:     "ORDER BY 语句对多个不同条件使用不同方向的排序无法使用索引",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "ORDER BY 语句对多个不同条件使用不同方向的排序无法使用索引",
		Func:    CLA007,
	}, {
		Rule: driver.Rule{ //select col1,col2 from tbl group by 1
			Name:     "CLA008",
			Desc:     "请为 GROUP BY 显示添加 ORDER BY 条件",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "请为 GROUP BY 显示添加 ORDER BY 条件",
		Func:    CLA008,
	}, {
		Rule: driver.Rule{ //select description from film where title ='ACADEMY DINOSAUR' order by length-language_id;
			Name:     "CLA009",
			Desc:     "不建议ORDER BY 的条件为表达式",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议ORDER BY 的条件为表达式",
		Func:    CLA009,
	}, {
		Rule: driver.Rule{ //select description from film where title ='ACADEMY DINOSAUR' order by length-language_id;
			Name:     "CLA012",
			Desc:     "建议将过长的SQL分解成几个简单的SQL",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "64",
					Desc:  "SQL最大长度",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "建议将过长的SQL分解成几个简单的SQL",
		Func:    CLA012,
	}, {
		Rule: driver.Rule{ //SELECT s.c_id,count(s.c_id) FROM s where c = test GROUP BY s.c_id HAVING s.c_id <> '1660' AND s.c_id <> '2' order by s.c_id
			Name:     "CLA013",
			Desc:     "不建议使用 HAVING 子句",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议使用 HAVING 子句",
		Func:    CLA013,
	}, {
		Rule: driver.Rule{ //delete from tbl
			Name:     "CLA014",
			Desc:     "删除全表时建议使用 TRUNCATE 替代 DELETE",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "删除全表时建议使用 TRUNCATE 替代 DELETE",
		Func:    CLA014,
	}, {
		Rule: driver.Rule{ //update mysql.func set name ="asdf";
			Name:     "CLA016",
			Desc:     "不要 UPDATE 主键",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message: "不要 UPDATE 主键",
		Func:    CLA016,
	}, {
		Rule: driver.Rule{ //create table t(c1 int,c2 int,c3 int,c4 int,c5 int,c6 int);
			Name:     "COL006",
			Desc:     "表中包含有太多的列",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "最大列数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "表中包含有太多的列",
		Func:    COL006,
	}, {
		Rule: driver.Rule{ //CREATE TABLE `tb2` ( `id` int(11) DEFAULT NULL, `col` char(10) CHARACTER SET utf8 DEFAULT NULL)
			Name:     "COL014",
			Desc:     "建议列与表使用同一个字符集",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
		},
		Message: "建议列与表使用同一个字符集",
		Func:    COL014,
	}, {
		Rule: driver.Rule{ //CREATE TABLE tab (a INT(1));
			Name:     "COL016",
			Desc:     "整型定义建议采用 INT(10) 或 BIGINT(20)",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message: "整型定义建议采用 INT(10) 或 BIGINT(20)",
		Func:    COL016,
	}, {
		Rule: driver.Rule{ //CREATE TABLE tab (a varchar(3500));
			Name:     "COL017",
			Desc:     "VARCHAR 定义长度过长",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "1024",
					Desc:  "VARCHAR最大长度",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "VARCHAR 定义长度过长",
		Func:    COL017,
	}, {
		Rule: driver.Rule{ //select id from t where substring(name,1,3)='abc'
			Name:     "DIS003",
			Desc:     "应避免在 WHERE 条件中使用函数或其他运算符",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "应避免在 WHERE 条件中使用函数或其他运算符",
		Func:    DIS003,
	}, {
		Rule: driver.Rule{ //SELECT SYSDATE();
			Name:     "FUN004",
			Desc:     "不建议使用 SYSDATE() 函数",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "不建议使用 SYSDATE() 函数",
		Func:    FUN004,
	}, {
		Rule: driver.Rule{ //SELECT SUM(COL) FROM tbl;
			Name:     "FUN006",
			Desc:     "使用 SUM(COL) 时需注意 NPE 问题",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "使用 SUM(COL) 时需注意 NPE 问题",
		Func:    FUN006,
	}, {
		Rule: driver.Rule{ //CREATE TABLE tbl ( a int, b int, c int, PRIMARY KEY(`a`,`b`,`c`));
			Name:     "KEY006",
			Desc:     "检测主键中的列过多",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "2",
					Desc:  "主键应当不超过多少列",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "主键中的列过多",
		Func:    KEY006,
	}, {
		Rule: driver.Rule{ //select col1,col2 from tbl where name=xx limit 10
			Name:     "RES002",
			Desc:     "未使用 ORDER BY 的 LIMIT 查询",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "未使用 ORDER BY 的 LIMIT 查询",
		Func:    RES002,
	},
	// {
	//	Rule: driver.Rule{ //UPDATE film SET length = 120 WHERE title = 'abc' LIMIT 1;
	//		Name:     "RES003",
	//		Desc:     "UPDATE/DELETE 操作不应使用LIMIT",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "UPDATE/DELETE 操作不应使用LIMIT",
	//	Func:    RES003,
	//},
	//{
	//	Rule: driver.Rule{ //UPDATE film SET length = 120 WHERE title = 'abc' ORDER BY title
	//		Name:     "RES004",
	//		Desc:     "UPDATE/DELETE 操作不要指定 ORDER BY",
	//		Level:    driver.RuleLevelWarn,
	//		Category: RuleTypeDMLConvention,
	//	},
	//	Message: "UPDATE/DELETE 操作不要指定 ORDER BY",
	//	Func:    RES004,
	//},
	{
		Rule: driver.Rule{ //TRUNCATE TABLE tbl_name
			Name:     "SEC001",
			Desc:     "请谨慎使用TRUNCATE操作",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "请谨慎使用TRUNCATE操作",
		Func:    SEC001,
	}, {
		Rule: driver.Rule{ //SELECT BENCHMARK(10, RAND())
			Name:     "SEC004",
			Desc:     "发现常见 SQL 注入函数",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message: "发现常见 SQL 注入函数",
		Func:    SEC004,
	}, {
		Rule: driver.Rule{ //CREATE TABLE tbl (a int) AUTO_INCREMENT = 10;
			Name:     "TBL004",
			Desc:     "表的初始AUTO_INCREMENT值不为0",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
		},
		Message: "表的初始AUTO_INCREMENT值不为0",
		Func:    TBL004,
	},
}

func init() {
	for _, rh := range SoarRuleHandlers {
		RuleHandlers = append(RuleHandlers, rh)
		RuleHandlerMap[rh.Rule.Name] = rh
		InitRules = append(InitRules, rh.Rule)
	}
}

func ALI002(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		for _, field := range stmt.Fields.Fields {
			if field.Text() == "*" && field.AsName.String() != "" {
				addResult(res, rule, rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

func ALI003(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		repeats := []string{}
		fields := map[string]struct{}{}
		if stmt.From != nil {
			if source, ok := stmt.From.TableRefs.Left.(*ast.TableSource); ok {
				if tableName, ok := source.Source.(*ast.TableName); ok {
					fields[tableName.Name.L] = struct{}{}
				}

			}
		}
		for _, field := range stmt.Fields.Fields {
			if selectColumn, ok := field.Expr.(*ast.ColumnNameExpr); ok && selectColumn.Name.Name.L != "" {
				fields[selectColumn.Name.Name.L] = struct{}{}
			}
		}
		for _, field := range stmt.Fields.Fields {
			if _, ok := fields[field.AsName.L]; ok {
				repeats = append(repeats, field.AsName.String())
			}
		}
		if len(repeats) > 0 {
			addResult(res, rule, rule.Name, strings.Join(repeats, ","))
		}
		return nil
	default:
		return nil
	}
}

func ALT001(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, option := range spec.Options {
				if option.Tp == ast.TableOptionCharset {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func ALT003(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropColumn {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func ALT004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.DropIndexStmt:
		if strings.ToLower(stmt.IndexName) == "primary" {
			addResult(res, rule, rule.Name)
		}
		return nil
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropPrimaryKey {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func ALT004_1(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropForeignKey {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func ARG001(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternLikeExpr:
				switch pattern := x.Pattern.(type) {
				case *parserDriver.ValueExpr:
					datum := pattern.Datum.GetString()
					if strings.HasPrefix(datum, "%") {
						trigger = true
						return true
					}
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG002(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternLikeExpr:
				switch pattern := x.Pattern.(type) {
				case *parserDriver.ValueExpr:
					datum := pattern.Datum.GetString()
					if !strings.HasPrefix(datum, "%") && !strings.HasSuffix(datum, "%") {
						trigger = true
						return true
					}
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternInExpr:
				for _, exprNode := range x.List {
					switch pattern := exprNode.(type) {
					case *parserDriver.ValueExpr:
						if pattern.Datum.GetString() == "" {
							trigger = true
							return true
						}

					}
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG005(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.PatternInExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG006(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.IsNullExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG007(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.PatternRegexpExpr, *ast.PatternLikeExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG009(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *parserDriver.ValueExpr:
				datum := pattern.Datum.GetString()
				if strings.HasPrefix(datum, " ") || strings.HasSuffix(datum, " ") {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func ARG013(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case ast.DDLNode:
		if strings.Index(node.Text(), "“") != -1 {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func CLA002(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			if expr, ok := orderBy.Items[0].Expr.(*ast.FuncCallExpr); ok && expr.FnName.L == "rand" {
				addResult(res, rule, rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

func CLA004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		groupBy := stmt.GroupBy
		if groupBy != nil {
			if _, ok := groupBy.Items[0].Expr.(*ast.PositionExpr); ok {
				addResult(res, rule, rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

func CLA007(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			isDesc := false
			for i, item := range orderBy.Items {
				if i == 0 {
					isDesc = item.Desc
				}
				if item.Desc != isDesc {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func CLA008(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.GroupBy != nil && stmt.OrderBy == nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func CLA009(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			for _, item := range orderBy.Items {
				if _, ok := item.Expr.(*ast.BinaryOperationExpr); ok {
					addResult(res, rule, rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func CLA012(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if len(node.Text()) > rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
		addResult(res, rule, rule.Name)
	}
	return nil
}

func CLA013(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Having != nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func CLA014(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.DeleteStmt:
		if stmt.Where == nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func CLA016(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		createTable, exist, err := ctx.GetCreateTableStmt(stmt.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName))
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		primary := map[string]struct{}{}
		for _, col := range createTable.Constraints {
			if col.Tp == ast.ConstraintPrimaryKey {
				for _, key := range col.Keys {
					primary[key.Column.Name.L] = struct{}{}
				}
				break
			}
		}
		for _, assignment := range stmt.List {
			if _, ok := primary[assignment.Column.Name.L]; ok {
				addResult(res, rule, rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func COL006(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if len(stmt.Cols) > rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func COL014(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp.Charset != "" {
				addResult(res, rule, rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func COL016(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if (col.Tp.Tp == mysql.TypeLong && col.Tp.Flen != 10) || (col.Tp.Tp == mysql.TypeLonglong && col.Tp.Flen != 20) {
				addResult(res, rule, rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func COL017(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeVarchar && col.Tp.Flen > rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
				addResult(res, rule, rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func DIS003(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.FuncCallExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func FUN004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.FuncCallExpr); ok && fu.FnName.L == "sysdate" {
				addResult(res, rule, rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.FuncCallExpr:
				if pattern.FnName.L == "sysdate" {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func FUN006(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.AggregateFuncExpr); ok && strings.ToLower(fu.F) == "sum" {
				addResult(res, rule, rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.AggregateFuncExpr:
				if strings.ToLower(pattern.F) == "sum" {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func KEY006(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey && len(constraint.Keys) > rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
				addResult(res, rule, rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func RES002(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.OrderBy == nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func RES003(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func RES004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			addResult(res, rule, rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func SEC001(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case *ast.TruncateTableStmt:
		addResult(res, rule, rule.Name)
		return nil
	default:
		return nil
	}
}

func SEC004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	funcs := []string{"sleep", "benchmark", "get_lock", "release_lock"}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.FuncCallExpr); ok && inSlice(funcs, fu.FnName.L) {
				addResult(res, rule, rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.FuncCallExpr:
				if inSlice(funcs, pattern.FnName.L) {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}

func inSlice(ss []string, s string) bool {
	for _, s2 := range ss {
		if s2 == s {
			return true
		}
	}
	return false
}

func TBL004(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, option := range stmt.Options {
			if option.Tp == ast.TableOptionAutoIncrement && option.UintValue != 0 {
				addResult(res, rule, rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

//
//func temp(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
//	switch stmt := node.(type) {
//	case *ast.SelectStmt:
//
//		return nil
//	default:
//		return nil
//	}
//}
