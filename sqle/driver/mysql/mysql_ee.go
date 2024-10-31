//go:build enterprise
// +build enterprise

package mysql

import (
	"context"
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

// func (*MysqlDriverImpl) QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
// 	return QueryPrepare(ctx, sql, conf)
// }

// func QueryPrepare(ctx context.Context, sql string, conf *driverV2.QueryPrepareConf) (*driverV2.QueryPrepareResult, error) {
// 	node, err := sqlparser.Parse(sql)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// check is query sql
// 	stmt, ok := node.(*sqlparser.Select)
// 	if !ok {
// 		return &driver.QueryPrepareResult{
// 			ErrorType: driver.ErrorTypeNotQuery,
// 			Error:     driver.ErrorTypeNotQuery,
// 		}, nil
// 	}

// 	// Generate new limit
// 	limit, offset := -1, 0
// 	if stmt.Limit != nil {
// 		if stmt.Limit.Rowcount != nil {
// 			limit, _ = strconv.Atoi(stmt.Limit.Rowcount.(*sqlparser.Literal).Val)
// 		}
// 		if stmt.Limit.Offset != nil {
// 			offset, _ = strconv.Atoi(stmt.Limit.Offset.(*sqlparser.Literal).Val)
// 		} else if limit != -1 {
// 			offset = 0
// 		}
// 	}
// 	appendLimit, appendOffset := -1, -1
// 	if conf != nil {
// 		appendLimit, appendOffset = int(conf.Limit), int(conf.Offset)
// 	}
// 	if appendLimit != -1 && appendOffset == -1 {
// 		appendLimit = 0
// 	}

// 	newLimit, newOffset := CalculateOffset(limit, offset, appendLimit, appendOffset)

// 	if newLimit != -1 {
// 		l := &sqlparser.Limit{
// 			Offset: &sqlparser.Literal{
// 				Type: sqlparser.IntVal,
// 				Val:  strconv.Itoa(newOffset),
// 			},
// 			Rowcount: &sqlparser.Literal{
// 				Type: sqlparser.IntVal,
// 				Val:  strconv.Itoa(newLimit),
// 			},
// 		}
// 		stmt.SetLimit(l)
// 	}

// 	// rewrite
// 	return &driver.QueryPrepareResult{
// 		NewSQL:    sqlparser.String(stmt),
// 		ErrorType: driver.ErrorTypeNotError,
// 	}, nil
// }

// // 1 means this item has no value or no limit
// func CalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
// 	if checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset) {
// 		return oldLimit, oldOffset
// 	}
// 	return calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset)
// }

// func calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
// 	if oldLimit == -1 {
// 		return appendLimit, appendOffset
// 	}
// 	newOffset = oldOffset + appendOffset
// 	newLimit = appendLimit
// 	if newOffset+newLimit > oldLimit+oldOffset {
// 		newLimit = oldLimit - appendOffset
// 	}

// 	return newLimit, newOffset
// }

// func checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) bool {
// 	if appendLimit == -1 {
// 		return true
// 	}
// 	if oldLimit != -1 && appendOffset > oldLimit+oldOffset {
// 		return true
// 	}

// 	return false
// }

func (i *MysqlDriverImpl) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	// add timeout
	cancel := func() {}
	if conf != nil && conf.TimeOutSecond > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.TimeOutSecond)*time.Second)
		defer cancel()
	}

	columns, rows, err := conn.Db.QueryWithContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	// generate result
	res := &driverV2.QueryResult{
		Column: params.Params{},
		Rows:   []*driverV2.QueryResultRow{},
	}
	for _, column := range columns {
		res.Column = append(res.Column, &params.Param{
			Key:   column,
			Value: column,
		})
	}
	for _, row := range rows {
		r := &driverV2.QueryResultRow{
			Values: []*driverV2.QueryResultValue{},
		}
		for _, s := range row {
			r.Values = append(r.Values, &driverV2.QueryResultValue{
				Value: s.String,
			})
		}
		res.Rows = append(res.Rows, r)
	}
	return res, nil
}

func (i *MysqlDriverImpl) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {

	return nil, nil
}

func (i *MysqlDriverImpl) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabasSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {

	return nil, nil
}
