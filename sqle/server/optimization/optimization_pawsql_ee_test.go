//go:build enterprise
// +build enterprise

package optimization

import (
	"reflect"
	"testing"
)

func Test_getIndexesRecommendedFromMD(t *testing.T) {
	type args struct {
		md string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test1", args: args{md: "## 推荐的索引\n   ``` sql\n  CREATE INDEX PAWSQL_IDX0073695002 ON LWQ.JOBS(MIN_SALARY,JOB_ID);\n  -- 当QB_1中引用的表JOBS作为驱动表时, 索引PAWSQL_IDX0073695002可以被用来进行索引扫描; 对(order by j.min_salary asc)避免排序; 该索引是个覆盖索引，可以避免回表.\n   ```\n\n## 性能验证详细信息"},
			want: []string{"CREATE INDEX PAWSQL_IDX0073695002 ON LWQ.JOBS(MIN_SALARY,JOB_ID);"}},
		{name: "test2", args: args{md: "\n\n## 推荐的索引\n   ``` sql\n  create index PAWSQL_IDX0512980777 ON LWQ.LOCATIONS(LOCATION_ID,CITY);\n  -- 当LOCATIONS作为被驱动表时, 索引PAWSQL_IDX0512980777可以被用来进行索引查找, 过滤条件为(d.location_id = l.location_id); 该索引是个覆盖索引，可以避免回表.\n  CREATE INDEX PAWSQL_IDX0200743195 ON LWQ.DEPARTMENTS(DEPARTMENT_NAME,DEPARTMENT_ID);\n  -- 当DEPARTMENTS作为被驱动表时, 索引PAWSQL_IDX0200743195可以被用来进行索引查找, 过滤条件为(subquery_e2.department_id = subquery_d2.department_id and subquery_d2.department_name = 'SomeDepartment'); 该索引是个覆盖索引，可以避免回表.\n  -- 当DEPARTMENTS作为驱动表时, 索引PAWSQL_IDX0200743195可以被用来进行索引查找, 过滤条件为(subquery_d2.department_name = 'SomeDepartment'); 该索引是个覆盖索引，可以避免回表.\n  CREATE INDEX PAWSQL_IDX1111883205 ON LWQ.EMPLOYEES(EMPLOYEE_ID,DEPARTMENT_ID,LAST_NAME,FIRST_NAME);\n  -- 当EMPLOYEES作为被驱动表时, 索引PAWSQL_IDX1111883205可以被用来进行索引查找, 过滤条件为(subquery_e2.department_id = subquery_d2.department_id and e.employee_id = subquery_e2.employee_id); 该索引是个覆盖索引，可以避免回表.\n  -- 当EMPLOYEES作为被驱动表时, 索引PAWSQL_IDX1111883205可以被用来进行索引查找, 过滤条件为(e.department_id = d.department_id and e.employee_id = subquery_e2.employee_id); 该索引是个覆盖索引，可以避免回表.\n  -- 当EMPLOYEES作为驱动表时, 索引PAWSQL_IDX1111883205可以被用来进行索引扫描; 对(order by e.employee_id)避免排序; 该索引是个覆盖索引，可以避免回表.\n  CREATE INDEX PAWSQL_IDX1444881373 ON LWQ.DEPARTMENTS(DEPARTMENT_ID,LOCATION_ID,DEPARTMENT_NAME);\n  -- 当DEPARTMENTS作为被驱动表时, 索引PAWSQL_IDX1444881373可以被用来进行索引查找, 过滤条件为(e.department_id = d.department_id and d.location_id = l.location_id); 该索引是个覆盖索引，可以避免回表.\n   ```\n\n## 性能验证详细信息\n ### ✅ 本次优化实施后，预计本SQL的性能将提升 85.71%"},
			want: []string{"create index PAWSQL_IDX0512980777 ON LWQ.LOCATIONS(LOCATION_ID,CITY);", "CREATE INDEX PAWSQL_IDX0200743195 ON LWQ.DEPARTMENTS(DEPARTMENT_NAME,DEPARTMENT_ID);", "CREATE INDEX PAWSQL_IDX1111883205 ON LWQ.EMPLOYEES(EMPLOYEE_ID,DEPARTMENT_ID,LAST_NAME,FIRST_NAME);", "CREATE INDEX PAWSQL_IDX1444881373 ON LWQ.DEPARTMENTS(DEPARTMENT_ID,LOCATION_ID,DEPARTMENT_NAME);"}},
		{name: "test3", args: args{md: ""}, want: []string{}},
		{name: "test4", args: args{md: "推荐的索引"}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIndexesRecommendedFromMD(tt.args.md); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIndexesRecommendedFromMD() = %v, want %v", got, tt.want)
			}
		})
	}
}
