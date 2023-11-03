package parser

import (
	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"

	"github.com/antlr4-go/antlr/v4"
)

const ASSIGN = "="
const STATIC = "static"

var jdbcExecSqlFuncs = []string{"executeQuery", "addBatch", "queryForList", "executeUpdate", "execute", "prepareStatement"}

func (v *JavaVisitor) VisitCompilationUnit(ctx *javaAntlr.CompilationUnitContext) interface{} {
	fileBlock := new(FileBlock)
	v.Tree = fileBlock
	for _, child := range ctx.GetChildren() {
		node, ok := child.(*javaAntlr.TypeDeclarationContext)
		if !ok {
			continue
		}
		v.VisitTypeDeclaration(node)
	}
	return nil
}

func (v *JavaVisitor) VisitTypeDeclaration(ctx *javaAntlr.TypeDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ClassDeclarationContext:
			v.Tree.ClassBlocks = append(v.Tree.ClassBlocks, v.VisitClassDeclaration(node).(*ClassBlock))
		case *javaAntlr.EnumDeclarationContext:
		}
	}
	return nil
}

func (v *JavaVisitor) VisitClassDeclaration(ctx *javaAntlr.ClassDeclarationContext) interface{} {
	classBlock := new(ClassBlock)
	v.currentClass = classBlock

	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.IdentifierContext:
			classBlock.Name = v.VisitIdentifier(node).(string)
		case *javaAntlr.ClassBodyContext:
			v.currentNode = classBlock
			v.VisitClassBody(node)
		}
	}
	return classBlock
}

func (v *JavaVisitor) VisitIdentifier(ctx *javaAntlr.IdentifierContext) interface{} {
	return ctx.GetText()
}

func (v *JavaVisitor) VisitClassBody(ctx *javaAntlr.ClassBodyContext) interface{} {
	for _, child := range ctx.GetChildren() {
		declaration, ok := child.(*javaAntlr.ClassBodyDeclarationContext)
		if !ok {
			continue
		}
		v.VisitClassBodyDeclaration(declaration)
	}

	return nil
}

func (v *JavaVisitor) VisitClassBodyDeclaration(ctx *javaAntlr.ClassBodyDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.MemberDeclarationContext:
			v.VisitMemberDeclaration(node)
		case *javaAntlr.ModifierContext:
			if node.GetText() == STATIC {
				v.isStatic = true
			}
		}
	}
	v.isStatic = false
	return nil
}

func (v *JavaVisitor) VisitMemberDeclaration(ctx *javaAntlr.MemberDeclarationContext) interface{} {
	classBlock := v.currentNode.(*ClassBlock)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.FieldDeclarationContext:
			v.VisitFieldDeclaration(node)
		case *javaAntlr.MethodDeclarationContext:
			classBlock.FuncBlocks = append(classBlock.FuncBlocks, v.VisitMethodDeclaration(node).(*FuncBlock))
		}
	}

	return classBlock
}

func (v *JavaVisitor) VisitFieldDeclaration(ctx *javaAntlr.FieldDeclarationContext) interface{} {
	children := ctx.GetChildren()
	// 根据java定义变量规范，第一个child为变量类型，第二个child为变量名以及值
	firstChild := children[0]
	secondChild := children[1]
	fieldType := ""
	variables := []*Variable{}
	if node, ok := firstChild.(*javaAntlr.TypeTypeContext); ok {
		fieldType = v.VisitTypeType(node).(string)
	}
	if node, ok := secondChild.(*javaAntlr.VariableDeclaratorsContext); ok {
		variables = v.VisitVariableDeclarators(node).([]*Variable)
	}
	for _, vari := range variables {
		vari.Type = fieldType
	}

	return variables
}

func (v *JavaVisitor) VisitTypeType(ctx *javaAntlr.TypeTypeContext) interface{} {
	return ctx.GetText()
}

func (v *JavaVisitor) VisitVariableDeclarators(ctx *javaAntlr.VariableDeclaratorsContext) interface{} {
	variables := []*Variable{}
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.VariableDeclaratorContext:
			variable, ok := v.VisitVariableDeclarator(node).(*Variable)
			if ok {
				variables = append(variables, variable)
			}
		}
	}
	return variables
}

func (v *JavaVisitor) VisitVariableDeclarator(ctx *javaAntlr.VariableDeclaratorContext) interface{} {
	variable := new(Variable)
	expr := new(Expression)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.VariableDeclaratorIdContext:
			variable.Name = v.VisitVariableDeclaratorId(node).(string)
		case *javaAntlr.VariableInitializerContext:
			// 变量的值以表达式形式返回，兼容非基础类型的变量
			// ResultSet resultSet = statement.executeQuery(sqlQuery);
			expr = v.VisitVariableInitializer(node).(*Expression)
		}
	}
	// 根据表达式的content判断变量的值是否为基础数据类型
	if expr.Content == nil {
		return nil
	}
	variable.Value = expr.Content
	switch v.currentNode.(type) {
	case *ClassBlock:
		if v.isStatic {
			v.currentClass.Variables = append(v.currentClass.Variables, variable)
		}
	case *CodeBlock:
		v.currentFunc.Variables = append(v.currentFunc.Variables, variable)
	}
	return variable
}

func (v *JavaVisitor) VisitVariableDeclaratorId(ctx *javaAntlr.VariableDeclaratorIdContext) interface{} {
	return ctx.GetText()
}

func (v *JavaVisitor) VisitVariableInitializer(ctx *javaAntlr.VariableInitializerContext) interface{} {
	expr := &Expression{}
	// 部分变量定义为表达式，如: ResultSet resultSet = statement.executeQuery(sqlQuery);
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ExpressionContext:
			expr = v.VisitExpression(node).(*Expression)
		}
	}
	return expr
}

func (v *JavaVisitor) VisitMethodDeclaration(ctx *javaAntlr.MethodDeclarationContext) interface{} {
	funcBlock := new(FuncBlock)
	funcBlock.Parent = v.currentNode
	v.currentFunc = funcBlock

	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.MethodBodyContext:
			v.currentNode = funcBlock
			funcBlock.CodeBlock = v.VisitMethodBody(node).(*CodeBlock)
		case *javaAntlr.IdentifierContext:
			funcBlock.Name = v.VisitIdentifier(node).(string)
		}
	}
	return funcBlock
}

func (v *JavaVisitor) VisitMethodBody(ctx *javaAntlr.MethodBodyContext) interface{} {
	block := new(CodeBlock)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.BlockContext:
			v.currentNode = block
			block = v.VisitBlock(node).(*CodeBlock)
		}
	}
	block.Parent = v.currentFunc
	return block
}

func (v *JavaVisitor) VisitBlock(ctx *javaAntlr.BlockContext) interface{} {
	block := new(CodeBlock)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.BlockStatementContext:
			v.currentNode = block
			v.VisitBlockStatement(node)
		}
	}
	return block
}

func (v *JavaVisitor) VisitBlockStatement(ctx *javaAntlr.BlockStatementContext) interface{} {
	block := v.currentNode.(*CodeBlock)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.LocalVariableDeclarationContext:
			block.Variables = append(block.Variables, v.VisitLocalVariableDeclaration(node).([]*Variable)...)
		case *javaAntlr.StatementContext:
			v.currentNode = block
			v.VisitStatement(node)
		}
	}
	return block
}

func (v *JavaVisitor) VisitLocalVariableDeclaration(ctx *javaAntlr.LocalVariableDeclarationContext) interface{} {
	variables := []*Variable{}
	fieldType := ""
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.TypeTypeContext:
			fieldType = v.VisitTypeType(node).(string)
		case *javaAntlr.VariableDeclaratorsContext:
			variables = v.VisitVariableDeclarators(node).([]*Variable)
		}
	}
	for _, vari := range variables {
		vari.Type = fieldType
	}
	return variables
}

func (v *JavaVisitor) VisitStatement(ctx *javaAntlr.StatementContext) interface{} {
	block := v.currentNode.(*CodeBlock)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.BlockContext:
			subBlock := v.VisitBlock(node).(*CodeBlock)
			subBlock.Parent = v.currentFunc
			block.CodeBlocks = append(block.CodeBlocks, subBlock)
		case *javaAntlr.ExpressionContext:
			v.VisitExpression(node)
		case *javaAntlr.StatementContext:
			v.VisitStatement(node)
		}
	}
	return nil
}

func (v *JavaVisitor) VisitExpressionList(ctx *javaAntlr.ExpressionListContext) interface{} {
	expressions := []*Expression{}
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ExpressionContext:
			expressions = append(expressions, v.VisitExpression(node).(*Expression))
		}
	}
	return expressions
}

func (v *JavaVisitor) VisitExpression(ctx *javaAntlr.ExpressionContext) interface{} {
	expressionList := []*Expression{}
	expression := new(Expression)
	isAssignExpr := false
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.MethodCallContext:
			expression = v.VisitMethodCall(node).(*Expression)
		case *javaAntlr.PrimaryContext:
			content := v.VisitPrimary(node).(*Primary)
			expression.Content = content
		case *javaAntlr.ExpressionContext:
			expression = v.VisitExpression(node).(*Expression)
			expressionList = append(expressionList, expression)

			// <assoc=right> expression
			// bop=('=' | '+=' | '-=' | '*=' | '/=' | '&=' | '|=' | '^=' | '>>=' | '>>>=' | '<<=' | '%=')
			// expression
			// 表达式赋值，是通过两个表达式和赋值符号进行赋值
			if !isAssignExpr {
				continue
			}
			varName := expressionList[0].Content.Value
			if expressionList[1].Content != nil {
				varValue := expressionList[1].Content

				newVar := &Variable{Name: varName, Value: varValue}
				if v.currentClass != nil {
					v.currentClass.UpdateVariable(newVar)
				}
				if v.currentFunc != nil {
					v.currentFunc.UpdateVariable(newVar)
				}
			}

		case *antlr.TerminalNodeImpl:
			if child.(antlr.ParseTree).GetText() == ASSIGN {
				isAssignExpr = true
			}
		}
	}
	return expression
}

func (v *JavaVisitor) VisitPrimary(ctx *javaAntlr.PrimaryContext) interface{} {
	primary := new(Primary)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.IdentifierContext:
			primary.Value = node.GetText()
			primary.RuleIndex = node.RuleIndex
		case *javaAntlr.LiteralContext:
			primary.Value = node.GetText()
			primary.RuleIndex = node.RuleIndex
		}
	}
	return primary
}

func (v *JavaVisitor) VisitArguments(ctx *javaAntlr.ArgumentsContext) interface{} {
	expressions := []*Expression{}
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ExpressionListContext:
			expressions = v.VisitExpressionList(node).([]*Expression)
		}
	}
	return expressions
	// return ctx.GetText()
}

func (v *JavaVisitor) VisitMethodCall(ctx *javaAntlr.MethodCallContext) interface{} {
	expression := new(Expression)
	expression.Node = v.currentNode
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.IdentifierContext:
			expression.MethodCall = v.VisitIdentifier(node).(string)
		case *javaAntlr.ArgumentsContext:
			expression.Arguments = v.VisitArguments(node).([]*Expression)
		}
	}
	for _, jdbcExecFunc := range jdbcExecSqlFuncs {
		if jdbcExecFunc == expression.MethodCall {
			v.ExecSqlExpressions = append(v.ExecSqlExpressions, expression)
		}
	}
	return expression
}
