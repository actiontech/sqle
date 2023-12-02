package parser

import (
	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"

	"github.com/antlr4-go/antlr/v4"
)

func (v *JavaVisitor) VisitCompilationUnit(ctx *javaAntlr.CompilationUnitContext) interface{} {
	fileBlock := new(FileBlock)
	v.Tree = fileBlock
	v.currentFile = fileBlock
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.TypeDeclarationContext:
			v.VisitTypeDeclaration(node)
		case *javaAntlr.ImportDeclarationContext:
			v.VisitImportDeclaration(node)
		}
	}
	return nil
}

func (v *JavaVisitor) VisitTypeDeclaration(ctx *javaAntlr.TypeDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ClassDeclarationContext:
			v.VisitClassDeclaration(node)
		case *javaAntlr.EnumDeclarationContext:
		}
	}
	return nil
}

func (v *JavaVisitor) VisitImportDeclaration(ctx *javaAntlr.ImportDeclarationContext) interface{} {
	// todo: 兼容import多个库
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.QualifiedNameContext:
			v.currentFile.ImportClass = append(v.currentFile.ImportClass, node.GetText())
		}
	}
	return nil
}

func (v *JavaVisitor) VisitClassDeclaration(ctx *javaAntlr.ClassDeclarationContext) interface{} {
	classBlock := new(ClassBlock)
	v.currentFile.ClassBlocks = append(v.currentFile.ClassBlocks, classBlock)
	v.currentClass = classBlock

	// 第一次循环只申明类的变量以及方法名
	classBlock.isDeclare = true
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.IdentifierContext:
			classBlock.Name = v.VisitIdentifier(node).(string)
		case *javaAntlr.ClassBodyContext:
			v.currentNode = classBlock
			v.VisitClassBody(node)
		}
	}
	// 第二次循环获取变量值并解析函数中的代码
	classBlock.isDeclare = false
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ClassBodyContext:
			v.currentNode = classBlock
			v.VisitClassBody(node)
		}
	}
	for _, funcBlock := range classBlock.FuncBlocks {
		v.currentFunc = funcBlock
		v.VisitMethodDeclaration(funcBlock.Pointer)
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
			v.currentNode = v.currentClass
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
	classBlock := v.currentClass
	for _, child := range ctx.GetChildren() {
		v.currentNode = classBlock
		switch node := child.(type) {
		case *javaAntlr.FieldDeclarationContext:
			v.VisitFieldDeclaration(node)
		case *javaAntlr.MethodDeclarationContext:
			if classBlock.isDeclare {
				classBlock.FuncBlocks = append(classBlock.FuncBlocks, v.VisitMethodDeclaration(node).(*FuncBlock))
			}
		}
	}
	// if !classBlock.isDeclare {
	// 	for _, funcBlock := range classBlock.FuncBlocks {
	// 		v.currentFunc = funcBlock
	// 		v.VisitMethodDeclaration(funcBlock.Pointer)
	// 	}
	// }

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
			if !v.currentClass.isDeclare {
				expr = v.VisitVariableInitializer(node).(*Expression)
			}
		}
	}
	// 如果是申明状态直接返回，不做任何校验
	if v.currentClass.isDeclare {
		variable.Pointer = ctx
		v.currentClass.Variables = append(v.currentClass.Variables, variable)
		return variable
	}
	variable.Value = expr
	switch v.currentNode.(type) {
	case *ClassBlock:
		// 非申请状态下，只能更新class中的值
		v.currentClass.UpdateVariable(variable)
	case *CodeBlock:
		v.currentFunc.Variables = append(v.currentFunc.Variables, variable)
	}
	return variable
}

func (v *JavaVisitor) VisitVariableDeclaratorId(ctx *javaAntlr.VariableDeclaratorIdContext) interface{} {
	return ctx.GetText()
}

func (v *JavaVisitor) VisitVariableInitializer(ctx *javaAntlr.VariableInitializerContext) interface{} {
	// 部分变量定义为表达式，如: ResultSet resultSet = statement.executeQuery(sqlQuery);
	// 需要解析statement.executeQuery(sqlQuery)以获取执行函数
	expr := new(Expression)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.ExpressionContext:
			expr = v.VisitExpression(node).(*Expression)
		}
	}
	// 直接返回变量值
	return expr
}

func (v *JavaVisitor) VisitMethodDeclaration(ctx *javaAntlr.MethodDeclarationContext) interface{} {
	funcBlock := new(FuncBlock)
	if !v.currentClass.isDeclare {
		funcBlock = v.currentFunc
	}
	funcBlock.Parent = v.currentNode
	v.currentFunc = funcBlock
	if v.currentClass.isDeclare {
		funcBlock.isDeclare = true
		funcBlock.Pointer = ctx
	} else {
		funcBlock.isDeclare = false
	}

	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.MethodBodyContext:
			if !funcBlock.isDeclare {
				v.currentNode = funcBlock
				funcBlock.CodeBlock = v.VisitMethodBody(node).(*CodeBlock)
			}
		case *javaAntlr.IdentifierContext:
			funcBlock.Name = v.VisitIdentifier(node).(string)
		case *javaAntlr.FormalParametersContext:
			v.VisitFormalParameters(node)
		}
	}
	return funcBlock
}

func (v *JavaVisitor) VisitFormalParameters(ctx *javaAntlr.FormalParametersContext) interface{} {
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		// ReceiverParameter接受一个对象值，如：void start(Myclass this)
		case *javaAntlr.ReceiverParameterContext:
			v.VisitReceiverParameter(node)
		case *javaAntlr.FormalParameterListContext:
			v.VisitFormalParameterList(node)
		}
	}
	return nil
}

func (v *JavaVisitor) VisitReceiverParameter(ctx *javaAntlr.ReceiverParameterContext) interface{} {
	param := new(Param)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.TypeTypeContext:
			param.Type = node.GetText()
		// 根据解析规则，遇到TerminalNodeImpl代表 参数值属于嵌套类 如:void myMethod(InnerClass.InnerClass2 this) {}
		// 不支持
		case *antlr.TerminalNodeImpl:
			return nil
		}
	}
	v.currentFunc.Params = append(v.currentFunc.Params, param)
	return nil
}

func (v *JavaVisitor) VisitFormalParameterList(ctx *javaAntlr.FormalParameterListContext) interface{} {
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.FormalParameterContext:
			v.VisitFormalParameter(node)
		// 暂不支持
		case *javaAntlr.LastFormalParameterContext:
		}
	}
	return nil
}

func (v *JavaVisitor) VisitFormalParameter(ctx *javaAntlr.FormalParameterContext) interface{} {
	param := new(Param)
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.TypeTypeContext:
			param.Type = node.GetText()
		case *javaAntlr.VariableDeclaratorIdContext:
			param.Content = node.GetText()
		}
	}
	v.currentFunc.Params = append(v.currentFunc.Params, param)
	return nil
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
	/*
		只判断 a.b.c(), a = A, a = select(), new A(), a="select", a=b.c
		六类表达式
	*/
	expression := new(Expression)
	expression.Node = v.currentNode
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		// 到方法调用停止解析表达式，例如 select().update()，只解析到select()
		case *javaAntlr.MethodCallContext:
			expr := v.CustomVisitMethodCall(node, expression)
			expr.RuleIndex = node.GetRuleIndex()
			if expression.Content == "" {
				expression = expr
			} else {
				expr.Before = v.currentExpr
				v.currentExpr.Next = expr
			}
			v.currentExpr = expr
			return expression
		case *javaAntlr.PrimaryContext:
			// primary出现在表达式中，代表没有其他子节点，只有primary
			pri := v.VisitPrimary(node).(*Primary)
			expression.Content = pri.Value
			expression.RuleIndex = pri.RuleIndex
		case *javaAntlr.ExpressionContext:
			expr := v.VisitExpression(node).(*Expression)
			if expression.Content == "" {
				expression = expr
			} else {
				expr.Before = v.currentExpr
				v.currentExpr.Next = expr
			}
			v.currentExpr = expr
			// 当出现赋值表达式，a=b，更新变量a的值
			if expr.Before != nil && expr.Before.Symbol == ASSIGN {
				newVar := &Variable{Name: expr.Before.Content, Value: expr}
				if v.currentClass != nil {
					v.currentClass.UpdateVariable(newVar)
				}
				if v.currentFunc != nil {
					v.currentFunc.UpdateVariable(newVar)
				}
			}

			// if !isAssignExpr {
			// 	continue
			// }
			// varName := expressionList[0].Content
			// if expressionList[1].Content != nil {
			// 	varValue := expressionList[1].Content

			// 	newVar := &Variable{Name: varName, Value: varValue}
			// 	if v.currentClass != nil {
			// 		v.currentClass.UpdateVariable(newVar)
			// 	}
			// 	if v.currentFunc != nil {
			// 		v.currentFunc.UpdateVariable(newVar)
			// 	}
			// }
		// 变量名或者类名
		case *javaAntlr.IdentifierContext:
			expr := &Expression{Content: node.GetText(), RuleIndex: node.GetRuleIndex(), Node: v.currentNode}
			// 赋值为当前表达式的子表达式
			v.currentExpr.Next = expr
			// 赋值父节点
			expr.Before = v.currentExpr
			// 赋值为当前表达式，继续向下寻找子表达式
			v.currentExpr = expr
		case *antlr.TerminalNodeImpl:
			symbol := node.GetText()
			if symbol == ASSIGN || symbol == DOT || symbol == PLUS {
				v.currentExpr.Symbol = symbol
				if symbol == DOT {
					// 判断该表达式是否为第一个调用的表达式 如: a.b 判断表达式是否为a
					if expression.Before == nil || (expression.Before != nil && expression.Before.Depth == 0) {
						expression.Depth = 1
						variType, isImport := v.getClassFromVar(v.currentExpr.Content)
						if isImport {
							expression.ImportName = variType
							expression.IsImport = true
						} else {
							classBlock := v.getClassFromSingleFile(variType)
							if classBlock != nil {
								expression.LocalClass = classBlock
							}
						}
					} else if expression.Before != nil {
						expression.Depth = expression.Before.Depth + 1
					}
				}
			} else {
				// 非赋值符号(=)和调用符号(.)的表达式不做解析, 例如 a++, a || b
				return &Expression{}
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

func (v *JavaVisitor) CustomVisitMethodCall(ctx *javaAntlr.MethodCallContext, expr *Expression) *Expression {
	expression := new(Expression)
	expression.Node = v.currentNode
	expression.RuleIndex = ctx.GetRuleIndex()
	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *javaAntlr.IdentifierContext:
			expression.Content = v.VisitIdentifier(node).(string)
		case *javaAntlr.ArgumentsContext:
			expression.Arguments = v.VisitArguments(node).([]*Expression)
		}
	}
	// 判断before调用深度是否为0或者1, 获取类中的函数
	// 深层次调用无法解析
	if expr.Depth == 0 {
		funcBlock := v.currentClass.getFuncFromClass(expression.Content)
		if funcBlock != nil {
			expression.MethodCall = funcBlock
			funcBlock.CalledExprs = append(funcBlock.CalledExprs, expression)
		}
	} else if expr.Depth == 1 && expr.LocalClass != nil {
		funcBlock := expr.LocalClass.getFuncFromClass(expression.Content)
		if funcBlock != nil {
			expression.MethodCall = funcBlock
			funcBlock.CalledExprs = append(funcBlock.CalledExprs, expression)
		}
	}

	// 当执行sql类的引用方式为import java.sql.*;
	// expression中的IsImport为false ImportName为空 localClass为nil,说明当前的expr属于一个未知类，可能是通过import java.sql.*；这类方法引用
	// 所以没有办法确定当前表达式的类是否为jdbc中的类，只能模糊地判断expr.Content是否是jdbc中的调用方法
	if !expr.IsImport && expr.ImportName == "" && expr.LocalClass == nil {
		isJdbc := judgeIsJdbcFunc(expression.Content)
		if isJdbc {
			v.ExecSqlExpressions = append(v.ExecSqlExpressions, expression)
		}
	} else if expr.ImportName != "" {
		// todo: 对sql查询需要可拓展性封装支持更多的sql查询方式
		// 临时方案：只判断是否用了jdbc中的方法
		isJdbc := judgeIsJdbcFunc(expression.Content)
		if isJdbc {
			v.ExecSqlExpressions = append(v.ExecSqlExpressions, expression)
		}
	}
	return expression
}
