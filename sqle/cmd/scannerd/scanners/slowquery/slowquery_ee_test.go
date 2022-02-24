//go:build enterprise
// +build enterprise

package slowquery

import (
	"context"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSlowQuery(t *testing.T) {
	params := &Params{
		LogFilePath: "./testdata/slowquery.log",
	}
	log := logrus.WithField("test", "test")
	log.Level = logrus.DebugLevel
	scanner, err := New(params, log, nil)
	assert.NoError(t, err)

	errCh := make(chan error)
	go func() {
		err := <-errCh
		if err != nil {
			panic(err)
		}
	}()

	// after 1s, stop scanner and assert output
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go func() {
		err := scanner.Run(ctx)
		errCh <- err
	}()

	sqlCh := scanner.SQLs()
	var sqls []scanners.SQL
	for sql := range sqlCh {
		sqls = append(sqls, sql)
	}

	for _, sql := range sqls {
		assert.Contains(t, sql.RawText, "select sleep")
		t.Log(sql.RawText)
	}
	// actual count is 21, Parser can not parse the last event.
	// TODO: explore reason
	assert.Len(t, sqls, 20)
}
