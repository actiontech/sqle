package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
)

func (r *MutationResolverImpl) AsyncSQLExecuteQuery(ctx context.Context, connectionID string, contextID string, sql string, resultID *string, filter *model.SQLDataFilter, dataFormat *model.ResultDataFormat) (*model.AsyncTaskInfo, error) {
	success, result, err := service.AuditSQL(sql, connectionID)
	if err != nil {
		return nil, err
	}
	if !success {
		name := "SQL Audit Failed"
		msg := fmt.Sprintf("[SQLE] sql statements are not allowed to excute, caused by: \nthe highest error level in audit results is %v,  which reaches the error level limit (%v) set in SQLE.", result.AuditLevel, result.LimitLevel)
		return nil, r.Ctx.JSON(http.StatusOK, struct {
			Data struct {
				TaskInfo model.AsyncTaskInfo `json:"taskInfo"`
			} `json:"data"`
		}{
			struct {
				TaskInfo model.AsyncTaskInfo `json:"taskInfo"`
			}{
				TaskInfo: model.AsyncTaskInfo{
					Name:    &name,
					Running: false,
					Status:  &sql,
					Error: &model.ServerError{
						Message:    &msg,
						StackTrace: &result.Result,
					},
				},
			},
		})
	}

	_, err = r.Next(r.Ctx)
	if err != nil {
		return nil, err
	}

	return nil, err
}
