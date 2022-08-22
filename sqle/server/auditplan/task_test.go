package auditplan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTidbCompletionSchema(t *testing.T) {
	// https://github.com/actiontech/sqle-ee/issues/395
	sql := "INSERT INTO t1(a1,a2,a3,a4) VALUES('','','Y',CURRENT_DATE)"
	newSQL, err := tidbCompletionSchema(sql, "test")
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `test`.`t1` (`a1`,`a2`,`a3`,`a4`) VALUES ('','','Y',CURRENT_DATE())", newSQL)
}
