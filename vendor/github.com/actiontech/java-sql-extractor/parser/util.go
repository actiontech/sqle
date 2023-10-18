package parser

import (
	"strings"

	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"
)

func getValueFromVariables(variableName string, variables []*Variable) (string, bool) {
	var result string
	var ok bool
	for _, vari := range variables {
		if vari.Name == variableName {
			result = vari.Value
			ok = true
		}
	}
	return result, ok
}

func getVarFromFunc(variableName string, funcBlock *FuncBlock) (string, bool) {
	return getValueFromVariables(variableName, funcBlock.Variables)
}

func getVarFromClass(variableName string, classBlock *ClassBlock) (string, bool) {
	return getValueFromVariables(variableName, classBlock.Variables)
}

func getVariableValueFromTree(variableName string, currentNode interface{}) (string, bool) {
	var variableValue string
	var ok bool
	switch node := currentNode.(type) {
	case *ClassBlock:
		variableValue, ok = getVarFromClass(variableName, node)
	case *CodeBlock:
		variableValue, ok = getVarFromFunc(variableName, node.Parent)
		if !ok {
			if classBlock, success := node.Parent.Parent.(*ClassBlock); success {
				variableValue, ok = getVarFromClass(variableName, classBlock)
			}
		}
	}
	return variableValue, ok
}

func GetSqlsFromVisitor(ctx *JavaVisitor) []string {
	sqls := []string{}
	for _, expression := range ctx.ExecSqlExpressions {
		for _, arg := range expression.Arguments {
			// 参数为变量
			if arg.Content.RuleIndex == javaAntlr.JavaParserRULE_identifier {
				sql, ok := getVariableValueFromTree(arg.Content.Value, expression.Node)
				if !ok {
					continue
				}
				// anltr为了区分字符串和其他变量，会为字符串的值左右添加双引号，获取sql时需要去除左右的双引号
				sqls = append(sqls, strings.Trim(sql, "\""))
			// 参数为字符串
			} else if arg.Content.RuleIndex == javaAntlr.JavaParserRULE_literal {
				sqls = append(sqls, strings.Trim(arg.Content.Value, "\""))
			}

		}
	}

	return sqls
}
