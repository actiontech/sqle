package sqlFile

import (
	"context"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMyBatis(t *testing.T) {
	params := &Params{
		SQLDir: "./not-exist-directory/",
	}
	scanner, err := New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)

	err = scanner.Run(context.TODO())
	assert.Error(t, err)

	params = &Params{
		SQLDir: "./testdata/",
	}
	scanner, err = New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)

	go scanner.Run(context.TODO())

	var sqlCh = scanner.SQLs()
	sqlBuf := []scanners.SQL{}

	for v := range sqlCh {
		sqlBuf = append(sqlBuf, v)
	}
	assert.Len(t, sqlBuf, 10)

	// test SQLFile scanner will hang until caller called ctx.Cancel().
	scanner, err = New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})
	go func() {
		scanner.Run(ctx)
		close(exitCh)
	}()

	time.Sleep(1 * time.Second)
	ok := true
	select {
	case _, ok = <-exitCh:
	default:
		assert.True(t, ok)
	}

	cancel()
	_, ok = <-exitCh
	assert.False(t, ok)
}
