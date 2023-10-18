package main

import (

	// "github.com/antlr4-go/antlr/v4"
	"os"

	parser "github.com/actiontech/java-sql-extractor/parser"
	// mysqlParser "anltrexample/parser/mysql"
	"fmt"
)


func main() {
	// Setup the input
	// javaParser, _ := parser.CreateJavaParser("/root/javaexample/test/Test1.java")

	// antlr.ParseTreeWalkerDefault.Walk(parser.NewJavaListener(), javaParser.CompilationUnit())


	p, err := parser.CreateJavaParser("/root/javaexample/test/Test7.java")
	if err != nil {
		os.Exit(-1)
	}

	v := parser.NewJavaVisitor()
	
	a:=p.CompilationUnit()
	a.Accept(v)

	fmt.Println(parser.GetSqlsFromVisitor(v))
	
}
