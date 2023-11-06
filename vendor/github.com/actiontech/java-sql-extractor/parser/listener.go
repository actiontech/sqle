package parser

import (

	"github.com/antlr4-go/antlr/v4"

	javaAntlr "github.com/actiontech/java-sql-extractor/java_antlr"
	"fmt"
)

type javaListener struct {
	*javaAntlr.BaseJavaParserListener //继承Listener基类
	*antlr.DefaultErrorListener    //继承错误基类
}

func (s *javaListener) EnterFieldDeclaration(ctx *javaAntlr.FieldDeclarationContext) {
	for _, child := range ctx.GetChildren() {
		node, ok := child.(*javaAntlr.TypeTypeContext)
		if !ok {
			continue
		}
		ruleIndex := node.GetRuleIndex()
		if ruleIndex != javaAntlr.JavaParserRULE_typeType {
			continue
		}
		
	}
}

func (s *javaListener) EnterExpression(ctx *javaAntlr.ExpressionContext) {
}

func (this *javaListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(javaAntlr.JavaParserParserStaticData.RuleNames[ctx.GetRuleIndex()], ctx.GetText())
	
}

func NewJavaListener() *javaListener {
	return new(javaListener)
}
