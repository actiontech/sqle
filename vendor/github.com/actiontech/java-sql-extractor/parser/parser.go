package parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"
)

type JavaVisitor struct {
	*javaAntlr.BaseJavaParserVisitor
	Tree               *FileBlock
	ExecSqlExpressions []*Expression

	currentNode  interface{}
	currentFile  *FileBlock
	currentClass *ClassBlock
	currentFunc  *FuncBlock
	currentExpr  *Expression

	isStatic bool
}

// 获取变量类型，并判断变量类型是否为import引入的类
func (v *JavaVisitor) getClassFromVar(variableName string) (string, bool) {
	var variableType string
	if v.currentFunc != nil {
		for _, vari := range v.currentFunc.Variables {
			if vari.Name == variableName {
				variableType = vari.Type
			}
		}
	}
	if variableType == "" {
		return variableType, false
	}
	for _, importStr := range v.currentFile.ImportClass {
		strSlice := strings.Split(importStr, ".")
		if strSlice[len(strSlice)-1] == variableName {
			return importStr, true
		}
	}
	return variableType, false
}

func (v *JavaVisitor) getClassFromSingleFile(className string) *ClassBlock {
	for _, classBlock := range v.currentFile.ClassBlocks {
		if classBlock.Name == className {
			return classBlock
		}
	}
	return nil
}

type FileBlock struct {
	ClassBlocks []*ClassBlock
	ImportClass []string
}

type ClassBlock struct {
	Name       string
	FuncBlocks []*FuncBlock
	Variables  []*Variable // 静态变量

	isDeclare bool // 判断该类是在申明状态还是解析状态
}

func (v *ClassBlock) UpdateVariable(variable *Variable) {
	for _, vari := range v.Variables {
		if vari.Name == variable.Name {
			vari.Value = variable.Value
		}
	}
}

func (c *ClassBlock) getFuncFromClass(funcName string) *FuncBlock {
	for _, funcBlock := range c.FuncBlocks {
		if funcBlock.Name == funcName {
			return funcBlock
		}
	}
	return nil
}

type FuncBlock struct {
	Name      string
	CodeBlock *CodeBlock
	Variables []*Variable // 函数中所有的变量
	Pointer   *javaAntlr.MethodDeclarationContext
	isDeclare bool // 判断该方法是在申明状态还是解析状态
	Params    []*Param // 函数参数列表

	CalledExprs []*Expression // 调用该函数的表达式列表

	Parent interface{}
}

func (v *FuncBlock) UpdateVariable(variable *Variable) {
	for _, vari := range v.Variables {
		if vari.Name == variable.Name {
			vari.Value = variable.Value
		}
	}
}

type Param struct {
	Content string
	Type    string
}

type CodeBlock struct {
	Variables  []*Variable
	CodeBlocks []*CodeBlock

	Parent *FuncBlock
}

type Variable struct {
	Name  string
	Value *Expression
	Level string
	Type  string

	Pointer *javaAntlr.VariableDeclaratorContext
}

type Expression struct {
	Next       *Expression // 子节点
	Before     *Expression // 父节点
	Content    string
	RuleIndex  int
	Arguments  []*Expression
	MethodCall *FuncBlock
	Depth      int    // 调用深度;例如a.b.c  a的调用深度为1 b的调用深度为2 c的调用深度为3
	Symbol     string // 表达式后面跟的符号；例如+和.

	IsImport   bool
	ImportName string
	LocalClass *ClassBlock // 指向表达式a.c()中a的class

	Node interface{}
}

type Primary struct {
	Value     string
	RuleIndex int
}

func NewJavaVisitor() *JavaVisitor {
	return &JavaVisitor{
		BaseJavaParserVisitor: new(javaAntlr.BaseJavaParserVisitor),
	}
}

func CreateJavaParser(file string) (*javaAntlr.JavaParser, error) {
	fileStream, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := javaAntlr.NewJavaLexer(fileStream)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := javaAntlr.NewJavaParser(stream)
	p.BuildParseTrees = true
	return p, nil
}

func GetSqlFromJavaFile(file string) (sqls []string, err error) {
	p, err := CreateJavaParser(file)
	if err != nil {
		return
	}
	visitor := NewJavaVisitor()
	p.CompilationUnit().Accept(visitor)
	sqls = GetSqlsFromVisitor(visitor)
	return
}
