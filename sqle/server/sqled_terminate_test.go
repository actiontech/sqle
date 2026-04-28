package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isConnectionTerminatedError(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected bool
	}{
		// ===== MySQL 终止错误 (正向匹配) =====
		"MySQL: raw ErrInvalidConn message": {
			err:      fmt.Errorf("invalid connection"),
			expected: true,
		},
		"MySQL: ErrInvalidConn wrapped in CodeError": {
			err:      fmt.Errorf("connect remote database error: invalid connection"),
			expected: true,
		},
		"MySQL: ErrInvalidConn wrapped in ExecBatch format": {
			err:      fmt.Errorf("exec sql failed: \nSELECT 1 \ninvalid connection"),
			expected: true,
		},
		"MySQL: ErrInvalidConn through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = invalid connection"),
			expected: true,
		},

		// ===== PostgreSQL 终止错误 (正向匹配) =====
		"PostgreSQL: 57P01 admin_shutdown from pgx": {
			err:      fmt.Errorf("FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			expected: true,
		},
		"PostgreSQL: 57P01 through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			expected: true,
		},
		"PostgreSQL: 57P01 wrapped in driver adaptor": {
			err:      fmt.Errorf("exec sql in driver adaptor: FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			expected: true,
		},
		"PostgreSQL: conn closed": {
			err:      fmt.Errorf("conn closed"),
			expected: true,
		},
		"PostgreSQL: conn closed through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = conn closed"),
			expected: true,
		},

		// ===== SQL Server 终止错误 (正向匹配) =====
		"SQL Server: connection is broken from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:The connection is broken and recovery is not possible.  The client driver attempted to recover the connection one or more times and all attempts failed."),
			expected: true,
		},
		"SQL Server: connection is broken from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:A severe error occurred on the current command. The connection is broken and recovery is not possible."),
			expected: true,
		},
		"SQL Server: connection is broken from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:The connection is broken and recovery is not possible."),
			expected: true,
		},
		"SQL Server: Timeout expired from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			expected: true,
		},
		"SQL Server: Timeout expired from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			expected: true,
		},
		"SQL Server: Timeout expired from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			expected: true,
		},
		"SQL Server: session is in the kill state from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Cannot continue the execution because the session is in the kill state."),
			expected: true,
		},
		"SQL Server: session is in the kill state from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:Cannot continue the execution because the session is in the kill state."),
			expected: true,
		},
		"SQL Server: session is in the kill state from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:Cannot continue the execution because the session is in the kill state."),
			expected: true,
		},

		// ===== SQL Server 非终止错误 (反向匹配，不应误判) =====
		"SQL Server: login failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Login failed for user 'sa'."),
			expected: false,
		},
		"SQL Server: syntax error": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Incorrect syntax near 'SELECTT'."),
			expected: false,
		},

		// ===== 非终止错误 (反向匹配，不应误判) =====
		"nil error": {
			err:      nil,
			expected: false,
		},
		"syntax error": {
			err:      fmt.Errorf("ERROR: syntax error at or near \"SELECTT\" (SQLSTATE 42601)"),
			expected: false,
		},
		"permission denied": {
			err:      fmt.Errorf("ERROR: permission denied for table users (SQLSTATE 42501)"),
			expected: false,
		},
		"connection refused": {
			err:      fmt.Errorf("connection refused"),
			expected: false,
		},
		"context deadline exceeded": {
			err:      fmt.Errorf("context deadline exceeded"),
			expected: false,
		},
		"generic exec error": {
			err:      fmt.Errorf("exec sql failed: ERROR: relation \"nonexistent\" does not exist"),
			expected: false,
		},
		"unique constraint violation": {
			err:      fmt.Errorf("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isConnectionTerminatedError(tc.err)
			assert.Equal(t, tc.expected, result, "case: %s", name)
		})
	}
}
