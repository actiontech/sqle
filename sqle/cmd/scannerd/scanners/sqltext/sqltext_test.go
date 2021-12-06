package sqltext

import (
	"context"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSqlText(t *testing.T) {
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
	var sqlBuf []scanners.SQL
	for v := range sqlCh {
		//fmt.Printf("%+v\n", v)
		sqlBuf = append(sqlBuf, v)
	}
	assert.Len(t, sqlBuf, 12)

	params = &Params{
		SQL: "select * from user where id = 113;insert into user (id,name,age) values (1,'xiaoxi',12);insert into user (id,name,age) values (1,'xiaoxi',12),(2,'xiaoxi2',13);insert into user (id,name,age) values (1,'xiaoxi',12),(2,'xiaoxi2',13),(3,'xiaoxi5',15);",
	}
	scanner, err = New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)
	go scanner.Run(context.TODO())
	var sqlCh2 = scanner.SQLs()
	var sqlBuf2 []scanners.SQL
	for v := range sqlCh2 {
		//fmt.Printf("%+v\n", v)
		sqlBuf2 = append(sqlBuf2, v)
	}
	assert.Len(t, sqlBuf2, 4)

	// test sqltext scanner will hang until caller called ctx.Cancel().
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
