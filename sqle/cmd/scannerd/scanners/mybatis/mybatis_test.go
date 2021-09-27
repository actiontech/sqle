package mybatis

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMyBatis(t *testing.T) {
	params := &Params{
		XMLDir: "./not-exist-directory/",
	}
	scanner, err := New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)

	err = scanner.Run(context.TODO())
	assert.Error(t, err)

	params = &Params{
		XMLDir: "./testdata/",
	}
	scanner, err = New(params, logrus.New().WithField("test", "test"), nil)
	assert.NoError(t, err)

	go scanner.Run(context.TODO())

	var sqlCh = scanner.SQLs()
	var sqlBuf []scanners.SQL
	for v := range sqlCh {
		sqlBuf = append(sqlBuf, v)
	}
	assert.Len(t, sqlBuf, 10)
}
