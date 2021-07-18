package inspector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type FpCase struct {
	input  string
	expect string
}

func TestFingerprint(t *testing.T) {
	cases := []FpCase{
		{
			input:  `update  tb1 set a = "2" where a = "3" and b = 4`,
			expect: "UPDATE `tb1` SET `a`=? WHERE `a`=? AND `b`=?",
		},
		{
			input:  "select * from tb1 where a in (select a from tb2 where b = 2) and c = 100",
			expect: "SELECT * FROM `tb1` WHERE `a` IN (SELECT `a` FROM `tb2` WHERE `b`=?) AND `c`=?",
		},
		{
			input:  "REPLACE INTO `tb1` (a, b, c, d, e) VALUES (1, 1, '小明', 'F', 99)",
			expect: "REPLACE INTO `tb1` (`a`,`b`,`c`,`d`,`e`) VALUES (?,?,?,?,?)",
		},
		{
			input:  "CREATE TABLE `tb1` SELECT * FROM `tb2` WHERE a=1",
			expect: "CREATE TABLE `tb1`  AS SELECT * FROM `tb2` WHERE `a`=?",
		},
		{
			input:  "CREATE TABLE `tb1` AS SELECT * FROM `tb2` WHERE a=1",
			expect: "CREATE TABLE `tb1`  AS SELECT * FROM `tb2` WHERE `a`=?",
		},
		// newline
		{
			input:  "CREATE TABLE `tb1` (\n    a BIGINT NOT NULL AUTO_INCREMENT,\n    b BIGINT NOT NULL,\n    c DOUBLE NOT NULL,\n    PRIMARY KEY (a)\n)",
			expect: "CREATE TABLE `tb1` (`a` BIGINT NOT NULL AUTO_INCREMENT,`b` BIGINT NOT NULL,`c` DOUBLE NOT NULL,PRIMARY KEY(`a`))",
		},

		// whitespace
		{
			input:  "select * from `tb1` where a='my_db'  and  b='test1'",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},

		// comment
		{
			input:  "create database database_x -- this is a comment ",
			expect: "CREATE DATABASE `database_x`",
		},
		{
			input:  "select * from tb1 where a='my_db' and b='test1'/*this is a comment*/",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},
		{
			input:  "select * from tb1 where a='my_db' and b='test1'# this is a comment",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},
	}
	for _, c := range cases {
		testFingerprint(t, c.input, c.expect)
	}
}

func testFingerprint(t *testing.T, input, expect string) {
	acutal, err := Fingerprint(input, true)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, expect, acutal)
}
