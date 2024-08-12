package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditPlanListQuery(t *testing.T) {
	// args := map[string]interface{}{}
	getQuery := func(args map[string]interface{}) string {
		q, err := getSelectQuery(instanceAuditPlanSQLBodyTpl, instanceAuditPlanSQLQueryTpl, args)
		if err != nil {
			panic(err)
		}
		return q
	}
	type testCase struct {
		desc        string
		args        map[string]interface{}
		expectQuery string
	}

	cases := []testCase{
		{
			desc: "测试没有额外入参",
			args: map[string]interface{}{},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id


ORDER BY audit_plan_sqls.id`,
		},

		{
			desc: "测试limit",
			args: map[string]interface{}{
				"limit": 1,
			},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id


ORDER BY audit_plan_sqls.id
LIMIT :limit OFFSET :offset`,
		},

		{
			desc: "测试order by",
			args: map[string]interface{}{
				"order_by": "schema_name",
			},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id

ORDER BY schema_name
DESC`,
		},

		{
			desc: "测试 asc",
			args: map[string]interface{}{
				"order_by": "schema_name",
				"is_asc":   true,
			},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id

ORDER BY schema_name
ASC`,
		},

		{
			desc: "测试 last_receive_timestamp_from",
			args: map[string]interface{}{
				"last_receive_timestamp_from": "xxxx",
				"order_by":                    "schema_name",
				"is_asc":                      true,
			},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id
AND JSON_EXTRACT(audit_plan_sqls.info, '$.last_receive_timestamp') >= :last_receive_timestamp_from

ORDER BY schema_name
ASC`,
		},

		{
			desc: "测试 last_receive_timestamp_to",
			args: map[string]interface{}{
				"last_receive_timestamp_to": "xxxx",
			},
			expectQuery: `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
audit_plan_sqls.priority

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id
AND JSON_EXTRACT(audit_plan_sqls.info, '$.last_receive_timestamp') <= :last_receive_timestamp_to


ORDER BY audit_plan_sqls.id`,
		},
	}
	for _, c := range cases {
		q := getQuery(c.args)
		assert.Equal(t, c.expectQuery, q)
	}
}
