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

	uploadSQLCnt   int
	generateSQLCnt int

	runFailed bool
}

func getMockScanner() *mockScanner {
	return &mockScanner{
		testSQLCh: make(chan scanners.SQL, 10240),
	}
}

func (mc *mockScanner) Run(ctx context.Context) error {
	if mc.runFailed {
		return fmt.Errorf("mock scanner run failed")
	}

	mc.isRunning = true

	for i := 0; i < mc.generateSQLCnt; i++ {
		mc.testSQLCh <- scanners.SQL{RawText: fmt.Sprintf("select * from t1 where id = %v", i)}
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
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

func (mc *mockScanner) Upload(ctx context.Context, sqls []scanners.SQL, errorMessage string) error {
	mc.uploadSQLCnt += len(sqls)
	return nil
}

func Test_start(t *testing.T) {
	errCh := make(chan error, 1)
	leastPushSecond := 1
	pushBufferSize := 1024

	mc := getMockScanner()
	mc.generateSQLCnt = pushBufferSize / 2
	go func() {
		errCh <- Start(context.TODO(), mc, leastPushSecond, pushBufferSize)
	}()
	time.Sleep(time.Duration(leastPushSecond*2) * time.Second)
	assert.True(t, mc.isRunning)
	close(mc.testSQLCh)
	assert.NoError(t, <-errCh)
	assert.Equal(t, pushBufferSize/2, mc.uploadSQLCnt)

	mc = getMockScanner()
	mc.generateSQLCnt = pushBufferSize * 2
	go func() {
		errCh <- Start(context.TODO(), mc, leastPushSecond, pushBufferSize)
	}()
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

	mc = getMockScanner()
	mc.runFailed = true
	err := Start(context.TODO(), mc, leastPushSecond, pushBufferSize)
	assert.Error(t, err)
}
