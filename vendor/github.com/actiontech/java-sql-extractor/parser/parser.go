package parser

import (
	"github.com/antlr4-go/antlr/v4"

	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"
)

type JavaVisitor struct {
	*javaAntlr.BaseJavaParserVisitor
	Tree               *FileBlock
	ExecSqlExpressions []*Expression

	currentNode  interface{}
	currentClass *ClassBlock
	currentFunc  *FuncBlock

	isStatic bool
}

type FileBlock struct {
	ClassBlocks []*ClassBlock
}

type ClassBlock struct {
	Name       string
	FuncBlocks []*FuncBlock
	Variables  []*Variable // 静态变量
}

func (v *ClassBlock) UpdateVariable(variable *Variable) {
	for _, vari := range v.Variables {
		if vari.Name == variable.Name {
			vari.Value = variable.Value
		}
	}
}

type FuncBlock struct {
	Name      string
	CodeBlock *CodeBlock
	Variables []*Variable // 函数中所有的变量

	Parent interface{}
}

func (v *FuncBlock) UpdateVariable(variable *Variable) {
	for _, vari := range v.Variables {
		if vari.Name == variable.Name {
			vari.Value = variable.Value
		}
	}
}

type CodeBlock struct {
	Variables  []*Variable
	CodeBlocks []*CodeBlock

	Parent *FuncBlock
}

type Variable struct {
	Name  string
	Value *Primary
	Level string
	Type  string
}

type Expression struct {
	SubExpr    *Expression
	Content    *Primary
	MethodCall string
	Arguments  []*Expression

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
