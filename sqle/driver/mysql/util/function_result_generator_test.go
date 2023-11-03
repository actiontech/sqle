package util

import (
	"testing"

	"github.com/pingcap/parser"
	"github.com/stretchr/testify/assert"
)

func TestConCatFunc(t *testing.T) {
	testCases := []struct {
		SQL          string
		ExpectResult string
	}{
		{
			SQL:          "SELECT CONCAT('a','b','c');",
			ExpectResult: "abc",
		},
		{
			SQL:          "SELECT CONCAT('a_',UPPER('b'),'_c');",
			ExpectResult: "a_B_c",
		},
		{
			SQL:          "SELECT CONCAT(CONCAT('a_',UPPER('b'),'_c'),'_','a_',UPPER('b'),'_c');",
			ExpectResult: "a_B_c_a_B_c",
		},
	}
	conCatGenerator := ConCatFunc{}
	for _, testCase := range testCases {
		funcCallVisitor := FuncCallExprVisitor{}
		stmts, _, err := parser.New().PerfectParse(testCase.SQL, "", "")
		assert.NoError(t, err)
		for _, stmt := range stmts {
			stmt.Accept(&funcCallVisitor)
		}
		if assert.NotEmpty(t, funcCallVisitor.FuncCallList) {
			conCatGenerator.ConCat = funcCallVisitor.FuncCallList[0]
			assert.Equal(t, testCase.ExpectResult, conCatGenerator.GenerateResult())
		}
	}
}

func TestUpperFunc(t *testing.T) {
	testCases := []struct {
		SQL          string
		ExpectResult string
	}{
		{
			SQL:          "SELECT UPPER('a');",
			ExpectResult: "A",
		},
		{
			SQL:          "SELECT UPPER(CONCAT('a_',UPPER('b'),'_c'));",
			ExpectResult: "A_B_C",
		},
		{
			SQL:          "SELECT UPPER(CONCAT(CONCAT('a_',UPPER('b'),'_c'),'_','a_',UPPER('b'),'_c'));",
			ExpectResult: "A_B_C_A_B_C",
		},
	}
	upperGenerator := UpperFunc{}
	for _, testCase := range testCases {
		funcCallVisitor := FuncCallExprVisitor{}
		stmts, _, err := parser.New().PerfectParse(testCase.SQL, "", "")
		assert.NoError(t, err)
		for _, stmt := range stmts {
			stmt.Accept(&funcCallVisitor)
		}
		if assert.NotEmpty(t, funcCallVisitor.FuncCallList) {
			upperGenerator.Upper = funcCallVisitor.FuncCallList[0]
			assert.Equal(t, testCase.ExpectResult, upperGenerator.GenerateResult())
		}
	}
}
