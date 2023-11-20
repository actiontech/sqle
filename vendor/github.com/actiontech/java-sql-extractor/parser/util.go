package parser

import (
	"strings"

	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"
)

var jdbcExecSqlFuncs = []string{"executeQuery", "addBatch", "queryForList", "executeUpdate", "execute", "prepareStatement"}

func getStrValueFromExpression(expr *Expression) []string {
	var results []string

	if expr.RuleIndex == javaAntlr.JavaParserRULE_identifier {
		values := getVariableValueFromTree(expr.Content, expr.Node)
		if len(values) == 0 {
			return results
		}
		results = values
		if expr.Symbol == PLUS {
			var tmpSlice []string
			nextStrs := getStrValueFromExpression(expr.Next)
			for _, str := range nextStrs {
				for _, result := range results {
					tmpSlice = append(tmpSlice, result+str)
				}
			}
			if len(tmpSlice) > 0 {
				results = tmpSlice
			}
		}
	} else if expr.RuleIndex == javaAntlr.JavaParserRULE_literal {
		result := strings.Trim(expr.Content, "\"")
		if result == "" || result == "null" {
			return results
		}

		if expr.Symbol == PLUS {
			nextStrs := getStrValueFromExpression(expr.Next)
			if len(nextStrs) == 0 {
				return []string{result}
			}
			for _, str := range nextStrs {
				results = append(results, result+str)
			}
		} else {
			results = append(results, result)
		}
	}
	return results
}

func getValuesFromVariables(variableName string, variables []*Variable) []string {
	var results []string
	for _, vari := range variables {
		// 获取字符串类型的变量,不是字符串类型添加空字符串
		if vari.Name == variableName {
			if vari.Type == "String" {
				results = getStrValueFromExpression(vari.Value)
			} else {
				results = append(results, "")
			}
		}
	}
	return results
}

func (f *FuncBlock) getValueFromCallExpr(argumentIndex int) []string {
	var values []string
	for _, calledExpr := range f.CalledExprs {
		// 判断参数是字符串还是变量，变量只寻找string类型可以使用getVariableValueFromTree
		for ind, arg := range calledExpr.Arguments {
			if ind != argumentIndex {
				continue
			}
			if arg.RuleIndex == javaAntlr.JavaParserRULE_identifier {
				recursionSql := getVariableValueFromTree(arg.Content, arg.Node)
				if len(recursionSql) > 0 {
					values = append(values, recursionSql...)
				}
			} else if arg.RuleIndex == javaAntlr.JavaParserRULE_literal {
				values = append(values, strings.Trim(arg.Content, "\""))
				if arg.Symbol == PLUS {
					var tmpSlice []string
					nextStrs := getStrValueFromExpression(arg.Next)
					for _, str := range nextStrs {
						for _, result := range values {
							tmpSlice = append(tmpSlice, result+str)
						}
					}
					if len(tmpSlice) > 0 {
						values = tmpSlice
					}
				}
			}
		}
	}
	return values
}

func (f *FuncBlock) getValueFromParams(variableName string) []string {
	var result []string
	for ind, param := range f.Params {
		if param.Content == variableName {
			result = append(result, f.getValueFromCallExpr(ind)...)
			break
		}
	}
	return result
}

func getVarFromFunc(variableName string, funcBlock *FuncBlock) []string {
	var sqls []string
	sqls = getValuesFromVariables(variableName, funcBlock.Variables)
	if len(sqls) == 0 {
		sqls = funcBlock.getValueFromParams(variableName)
	}
	return sqls
}

func getVarFromClass(variableName string, classBlock *ClassBlock) []string {
	return getValuesFromVariables(variableName, classBlock.Variables)
}

func getVariableValueFromTree(variableName string, currentNode interface{}) []string {
	var variableValues []string
	switch node := currentNode.(type) {
	case *ClassBlock:
		variableValues = getVarFromClass(variableName, node)
	case *CodeBlock:
		variableValues = getVarFromFunc(variableName, node.Parent)
		if len(variableValues) == 0 {
			if classBlock, success := node.Parent.Parent.(*ClassBlock); success {
				variableValues = getVarFromClass(variableName, classBlock)
			}
		}
	}
	return variableValues
}

func GetSqlsFromVisitor(ctx *JavaVisitor) []string {
	sqls := []string{}
	for _, expression := range ctx.ExecSqlExpressions {
		if len(expression.Arguments) == 0 {
			continue
		}
		arg := expression.Arguments[0]
		// 参数为变量
		if arg.RuleIndex == javaAntlr.JavaParserRULE_identifier {
			sqls = append(sqls, getVariableValueFromTree(arg.Content, expression.Node)...)
			if arg.Symbol == PLUS {
				tmpSlice := []string{}
				nextSqls := getVariableValueFromTree(arg.Next.Content, expression.Node)
				for _, str := range nextSqls {
					for _, result := range sqls {
						tmpSlice = append(tmpSlice, result+str)
					}
				}
				if len(tmpSlice) > 0 {
					sqls = tmpSlice
				}
			}
			// 参数为字符串
		} else if arg.RuleIndex == javaAntlr.JavaParserRULE_literal {
			// anltr为了区分字符串和其他变量，会为字符串的值左右添加双引号，获取sql时需要去除左右的双引号
			sql := strings.Trim(arg.Content, "\"")
			if sql == "null" || sql == "" {
				continue
			}
			sqls = append(sqls, sql)
			if arg.Symbol == PLUS {
				tmpSlice := []string{}
				nextSqls := getVariableValueFromTree(arg.Next.Content, expression.Node)
				for _, str := range nextSqls {
					for _, result := range sqls {
						tmpSlice = append(tmpSlice, result+str)
					}
				}
				if len(tmpSlice) > 0 {
					sqls = tmpSlice
				}
			}
		}
	}

	return sqls
}

func judgeIsJdbcFunc(funcName string) bool {
	for _, jdbcExecFunc := range jdbcExecSqlFuncs {
		if jdbcExecFunc == funcName {
			return true
		}
	}
	return false
}
