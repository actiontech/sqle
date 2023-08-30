//go:build enterprise
// +build enterprise

package tidb_audit_log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO 只用了SQL_TEXT,所以只测SQL_TEXT内容没问题,后续再补充其他值是否正确
func Test_lexerParse(t *testing.T) {
	ts := []struct {
		Line      string
		SQL       string
		Databases []string
	}{
		{
			Line:      "[2022/07/06 02:37:54.387 +00:00] [INFO] [logger.go:76] [ID=16570750740] [TIMESTAMP=2022/07/06 02:37:54.387 +00:00] [EVENT_CLASS=GENERAL] [EVENT_SUBCLASS=] [STATUS_CODE=0] [COST_TIME=0] [HOST=127.0.0.1] [CLIENT_IP=127.0.0.1] [USER=root] [DATABASES=\"[test,db1]\"] [TABLES=\"[]\"] [SQL_TEXT=\"select `mysql` . `user` . `user` , `t1` . `id` from `mysql` . `user` join `t1` on `mysql` . `user` . `user` = `t1` . `id`\"] [ROWS=0] [CONNECTION_ID=13] [CLIENT_PORT=51098] [PID=25586] [COMMAND=Query] [SQL_STATEMENTS=]",
			SQL:       "select `mysql` . `user` . `user` , `t1` . `id` from `mysql` . `user` join `t1` on `mysql` . `user` . `user` = `t1` . `id`",
			Databases: []string{"test", "db1"},
		},
	}

	parser := GetLexerParser()

	for _, s := range ts {
		sql, err := parser.Parse(s.Line)
		assert.NoError(t, err)
		assert.Equalf(t, s.SQL, sql.SQLText, "err line: %v", s.Line)
		assert.Equalf(t, s.Databases, sql.Databases, "err line: %v", s.Line)
	}
}
