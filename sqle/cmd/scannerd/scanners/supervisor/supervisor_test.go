package supervisor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockScanner struct {
	isRunning bool
	testSQLCh chan scanners.SQL

	uploadSQLCnt int
}

func getMockScanner() *mockScanner {
	return &mockScanner{
		testSQLCh: make(chan scanners.SQL, 10240),
	}
}

func (mc *mockScanner) generateSQL(cnt int) {
	for i := 0; i < cnt; i++ {
		mc.testSQLCh <- scanners.SQL{RawText: fmt.Sprintf("select * from t1 where id = %v", i)}
	}
}

func (mc *mockScanner) Run(ctx context.Context) error {
	mc.isRunning = true
	return nil
}

func (mc *mockScanner) SQLs() <-chan scanners.SQL {
	sqlCh := make(chan scanners.SQL, 1024)
	go func() {
		for sql := range mc.testSQLCh {
			sqlCh <- sql
		}
		logrus.StandardLogger().Infoln("Call SQls close channel")
		close(sqlCh)
	}()
	return sqlCh
}

func (mc *mockScanner) Upload(ctx context.Context, sqls []scanners.SQL) error {
	mc.uploadSQLCnt += len(sqls)
	return nil
}

func Test_start(t *testing.T) {
	errCh := make(chan error, 1)
	leastPushSecond := 1
	pushBufferSize := 1024

	mc := getMockScanner()
	go func() {
		errCh <- Start(context.TODO(), mc, leastPushSecond, pushBufferSize)
	}()
	mc.generateSQL(pushBufferSize / 2)
	time.Sleep(time.Duration(leastPushSecond*2) * time.Second)
	assert.True(t, mc.isRunning)
	close(mc.testSQLCh)
	assert.NoError(t, <-errCh)
	assert.Equal(t, pushBufferSize/2, mc.uploadSQLCnt)

	mc = getMockScanner()
	go func() {
		errCh <- Start(context.TODO(), mc, leastPushSecond, pushBufferSize)
	}()
	mc.generateSQL(pushBufferSize * 2)
	time.Sleep(time.Duration(leastPushSecond*2) * time.Second)
	assert.True(t, mc.isRunning)
	close(mc.testSQLCh)
	assert.NoError(t, <-errCh)
	assert.Equal(t, pushBufferSize*2, mc.uploadSQLCnt)

	mc = getMockScanner()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		errCh <- Start(ctx, mc, leastPushSecond, pushBufferSize)
	}()
	cancel()
	assert.NoError(t, <-errCh)
}
