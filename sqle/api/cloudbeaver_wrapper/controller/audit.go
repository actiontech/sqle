package controller

import (
	"context"
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
		msg := "the audit level is not allowed to perform sql query"
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
						StackTrace: &result,
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
