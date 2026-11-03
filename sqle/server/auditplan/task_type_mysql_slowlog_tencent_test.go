package auditplan

import (
	"testing"

	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/stretchr/testify/assert"
)

func TestMySQLSlowLogTencentTaskV2_InstanceTypeAndParams(t *testing.T) {
	task := &MySQLSlowLogTencentTaskV2{}
	assert.Equal(t, InstanceTypeMySQL, task.InstanceType())

	ps := task.Params()
	assert.Equal(t, []string{
		paramKeyAccessKeyId,
		paramKeyAccessKeySecret,
		paramKeyRegion,
		paramKeyDBInstanceId,
		paramKeyFirstSqlsScrappedInLastPeriodHours,
		paramKeyCollectIntervalMinute,
	}, paramKeys(ps))

	assert.Equal(t, params.ParamTypePassword, ps.GetParam(paramKeyAccessKeyId).Type)
	assert.Equal(t, params.ParamTypePassword, ps.GetParam(paramKeyAccessKeySecret).Type)
	assert.Equal(t, "24", ps.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Value)
	assert.Equal(t, "60", ps.GetParam(paramKeyCollectIntervalMinute).Value)
}

func TestMetaBuilderList_Contains_TencentCdbMySQLSlowLog(t *testing.T) {
	var found *MetaBuilder
	for i := range MetaBuilderList {
		if MetaBuilderList[i].Type == TypeTencentCdbMySQLSlowLog {
			found = &MetaBuilderList[i]
			break
		}
	}
	assert.NotNil(t, found, "tencent_cdb_mysql_slow_log must be in MetaBuilderList")
	assert.Equal(t, TypeTencentCdbMySQLSlowLog, found.Type)

	meta, err := GetMeta(TypeTencentCdbMySQLSlowLog)
	assert.NoError(t, err)
	assert.Equal(t, InstanceTypeMySQL, meta.InstanceType)
}

func paramKeys(ps params.Params) []string {
	keys := make([]string, 0, len(ps))
	for _, p := range ps {
		keys = append(keys, p.Key)
	}
	return keys
}
