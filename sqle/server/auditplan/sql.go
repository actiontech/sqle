package auditplan

import (
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/percona/go-mysql/query"
)

type RawSQL struct {
	rawSQL string
	schema string
}

type mergedSQL struct {
	*RawSQL

	counter     int
	fingerprint string
}

func rawSQLsToModel(sqls []*RawSQL) (res []*model.AuditPlanSQLV2) {
	return mergedSQLsToModel(mergeSQLsBasedOnFingerprints(sqls))
}

func mergeSQLsBasedOnFingerprints(sqls []*RawSQL) (res []*mergedSQL) {
	res = make([]*mergedSQL, 0)
	counter := map[string] /*sql fingerprint*/ int /*slice index*/ {}
	for i := range sqls {
		sql := sqls[i]
		fp := query.Fingerprint(sql.rawSQL)
		if sqlIndex, exist := counter[fp]; exist {
			res[sqlIndex].counter += 1
			res[sqlIndex].fingerprint = fp
			res[sqlIndex].RawSQL = sql
		} else {
			res = append(res, &mergedSQL{
				counter:     1,
				fingerprint: fp,
				RawSQL:      sql,
			})
		}
	}
	return res
}

func mergedSQLsToModel(sqls []*mergedSQL) (
	res []*model.AuditPlanSQLV2) {

	res = make([]*model.AuditPlanSQLV2, len(sqls))
	now := time.Now()
	for i := range sqls {
		sql := sqls[i]
		modelInfo := fmt.Sprintf(
			`{"counter":%v,"last_receive_timestamp":"%v"}`,
			sql.counter, now.Format(time.RFC3339))
		res[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql.fingerprint,
			SQLContent:  sql.rawSQL,
			Info:        []byte(modelInfo),
			Schema:      sql.schema,
		}
	}

	return
}
