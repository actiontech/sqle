//go:build !enterprise
// +build !enterprise

package v1

import (
	"context"
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportSqlDevRecord = errors.New(errors.EnterpriseEditionFeatures, e.New("sql dev record is enterprise version feature"))

func getSqlDEVRecordList(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlDevRecord
}

func SyncSqlDevRecord(ctx context.Context, task *model.Task, creator string) {}
