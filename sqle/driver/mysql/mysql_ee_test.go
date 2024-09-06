//go:build enterprise
// +build enterprise

package mysql

// import (
// 	"context"
// 	"testing"

// 	"github.com/actiontech/sqle/sqle/driver"

// 	"github.com/stretchr/testify/assert"
// )

// func TestQueryPrepare(t *testing.T) {
// 	type T struct {
// 		sql      string
// 		conf     *driver.QueryPrepareConf
// 		result   *driver.QueryPrepareResult
// 		hasError bool
// 	}

// 	examples := []T{
// 		{ // 错误sql
// 			sql:      "aaaaa",
// 			hasError: true,
// 		}, { // 非查询sql
// 			sql: "insert into a values (1)",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  0,
// 				Offset: 0,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				ErrorType: driver.ErrorTypeNotQuery,
// 				Error:     driver.ErrorTypeNotQuery,
// 			},
// 			hasError: false,
// 		}, { // 没改写配置不改写(无限制sql)
// 			sql:  "select * from a",
// 			conf: nil,
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 没改写配置不改写(有限制sql)
// 			sql:  "select * from a limit 2, 3",
// 			conf: nil,
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 2, 3",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 无限制sql增加限制
// 			sql: "select * from a",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  2,
// 				Offset: 3,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 3, 2",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 限制和改写配置相同,结果和原SQL应该等效
// 			sql: "select * from a limit 3, 2",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  2,
// 				Offset: 3,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 3, 2",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 改写后的查询结果如果超过原sql改写结果则不改写
// 			sql: "select * from a limit 1, 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  2,
// 				Offset: 10,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 1, 8",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 改写配置的范围在原SQL查询范围内
// 			sql: "select * from a limit 1, 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  2,
// 				Offset: 3,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 4, 2",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 改写配置的范围超过了原SQL查询范围
// 			sql: "select * from a limit 1, 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  8,
// 				Offset: 2,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 3, 6",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 能识别offset ? limit ?
// 			sql: "select * from a limit 8 offset 1",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  8,
// 				Offset: 2,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 3, 6",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 只有limit没有offset
// 			sql: "select * from a limit 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  8,
// 				Offset: 2,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a limit 2, 6",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 不会误改写子句
// 			sql: "select * from a where id = (select id from b limit 1, 8) limit 1, 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  8,
// 				Offset: 2,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a where id = (select id from b limit 1, 8) limit 3, 6",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		}, { // 带join的sql
// 			sql: "select * from a join b on a.id = b.id limit 1, 8",
// 			conf: &driver.QueryPrepareConf{
// 				Limit:  8,
// 				Offset: 2,
// 			},
// 			result: &driver.QueryPrepareResult{
// 				NewSQL:    "select * from a join b on a.id = b.id limit 3, 6",
// 				ErrorType: driver.ErrorTypeNotError,
// 			},
// 			hasError: false,
// 		},
// 	}

// 	ctx := context.TODO()
// 	for _, e := range examples {
// 		res, err := QueryPrepare(ctx, e.sql, e.conf)
// 		if !e.hasError {
// 			assert.NoError(t, err, e.sql)
// 		}
// 		assert.EqualValues(t, e.result, res, e.sql)
// 	}

// }
