package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		input               string
		expectedFingerPrint string
	}{
		{
			` /*some comments*/  ROLLBACK        TO/*some comments*//*some comments*/ savePoint/*some comments*/ sp `,
			`ROLLBACK TO savePoint sp`,
		}, // 未能解析 带有注释
		{
			`/*2:4536734829-293874657380*/update some_value set value=12128921313213213213,/*2:4536734829-293874657380*/ dt = current_timestamp() /*2:4536734829-293874657380*/where code = 'some_value' and/*2:4536734829-293874657380*/name = 'some_value' and/*2:4536734829-293874657380*/value = 345678945678444444444444567`, "UPDATE `some_value` SET `value`=?, `dt`=CURRENT_TIMESTAMP() WHERE `code`=? AND `name`=? AND `value`=?",
		}, // 能够解析 带有注释
		{
			`/*3456789-:56789*/  /*231456789-2786"@@@*/`,
			``,
		}, // 纯注释
		{
			"SELECT * FROM /*multiline\ncomment*/table",
			"SELECT * FROM table",
		}, // 能够解析 带有多行注释
		{
			" /*some comments*/  ROLLBACK    /*multiline\ncomment*/    TO/*some comments*//*some comments*/ savePoint/*some comments*/ sp ;",
			`ROLLBACK TO savePoint sp`,
		}, // 未能解析 带有多行注释
		{
			"SELECT /**/ * FROM table",
			"SELECT * FROM table",
		}, // 能够解析 带有空注释
		{
			" /*some comments*/  ROLLBACK    /**/    TO/*some comments*//*some comments*/ savePoint /*some comments*/sp ;",
			`ROLLBACK TO savePoint sp`,
		}, // 未能解析 带有空注释

		{
			"/*comments1 */ select CUST_NO,count(CUST_NO) as `count(CUST_NO)`,some_code, COLLATION(CUST_NO), COLLATION(some_code) from `table1`.`com_cust_risk_info` group by CUST_NO,some_code order by CUST)NO , some_code",
			"select CUST_NO,count(CUST_NO) as `count(CUST_NO)`,some_code, COLLATION(CUST_NO), COLLATION(some_code) from `table1`.`com_cust_risk_info` group by CUST_NO,some_code order by CUST)NO , some_code",
		}, // 未能解析 带有注释1
		{
			"/*comments2 */ select CUST_NO,count(CUST_NO) as `count(CUST_NO)`,some_code, COLLATION(CUST_NO), COLLATION(some_code) from `table1`.`com_cust_risk_info` group by CUST_NO,some_code order by CUST)NO , some_code",
			"select CUST_NO,count(CUST_NO) as `count(CUST_NO)`,some_code, COLLATION(CUST_NO), COLLATION(some_code) from `table1`.`com_cust_risk_info` group by CUST_NO,some_code order by CUST)NO , some_code",
		}, // 未能解析 带有注释2
	}

	for _, tc := range testCases {
		ns, err := Parse(context.TODO(), tc.input)
		assert.NoError(t, err)
		if len(ns) > 0 {
			assert.Equal(t, tc.expectedFingerPrint, ns[0].Fingerprint)
		}
		if len(ns) == 0 {
			assert.Equal(t, tc.expectedFingerPrint, ``)
		}
	}
}
func TestClearComments(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"SELECT * FROM table /*comment*/;",
			"SELECT * FROM table",
		},
		{
			"SELECT * FROM table /*comment1*/ /*comment2*/",
			"SELECT * FROM table",
		},
		{
			"SELECT /*comment*/ * FROM /*comment*/ table",
			"SELECT * FROM table",
		},
		{
			"SELECT * /*inline comment*/ FROM table",
			"SELECT * FROM table",
		},
		{
			"SELECT * FROM /*multiline\ncomment*/table",
			"SELECT * FROM table",
		},
		{
			"  SELECT *  FROM   table   WHERE   id = /*comment*//*comment1*/1",
			"SELECT * FROM table WHERE id = 1",
		},
		{
			"SELECT /**/ * FROM table",
			"SELECT * FROM table",
		},
		{
			"SELECT * FROM table /*inline*/ JOIN other_table /*on id*/ ON table.id = other_table.id",
			"SELECT * FROM table JOIN other_table ON table.id = other_table.id",
		},
		{
			"  /*comment*/SELECT *  FROM   table   WHERE   id = 1;/*comment*/show tables;/*comment*/other string;other string",
			"SELECT * FROM table WHERE id = 1",
		},
	}

	for _, tc := range testCases {
		actual := clearComments(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}
