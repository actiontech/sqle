//go:build !enterprise
// +build !enterprise

package slowquery

import (
	"context"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/sirupsen/logrus"
)

var errSlowQueryNotImplemented = errors.NewNotImplemented("SlowQuery Scanner")

type SlowQuery struct{}

type Params struct {
	LogFilePath    string
	APName         string
	IncludeUsers   string
	ExcludeUsers   string
	IncludeSchemas string
	ExcludeSchemas string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SlowQuery, error) {
	return &SlowQuery{}, errSlowQueryNotImplemented
}

func (s *SlowQuery) Run(ctx context.Context) error {
	return errSlowQueryNotImplemented
}

func (sq *SlowQuery) SQLs() <-chan scanners.SQL {
	return nil
}

func (sq *SlowQuery) Upload(ctx context.Context, sqls []scanners.SQL) error {
	return errSlowQueryNotImplemented
}
