package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBinlogPrefix_ShouldReturnCorrect(t *testing.T) {
	binlog := "mysql-bin.000001"
	assert.Equal(t, "mysql-bin", GetBinlogPrefix(binlog), "binlog file(%v) prefix is mysql-bin", binlog)
}

func TestGetBinlogPrefix_ShouldReturnSame(t *testing.T) {
	binlog := "mysql-bin.0000001"
	assert.Equal(t, "mysql-bin.0000001", GetBinlogPrefix(binlog), "binlog file(%v) prefix is mysql-bin.0000001", binlog)

	binlog = "mysql-bin"
	assert.Equal(t, "mysql-bin", GetBinlogPrefix(binlog), "binlog file(%v) prefix is mysql-bin", binlog)
}
