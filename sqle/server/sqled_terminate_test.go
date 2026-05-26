package server

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/stretchr/testify/assert"
)

// Test_terminatedExecResult 锁定 ExecResult 文案格式：用户在工单 SQL 详情
// 看到的不是裸的 `hive exec failed (sql=...): EOF` / `invalid connection`，
// 而是「因中止上线中断: <原文本>」，前缀稳定、原 driver/plugin 信息保留。
func Test_terminatedExecResult(t *testing.T) {
	t.Run("nil error returns bare prefix", func(t *testing.T) {
		got := terminatedExecResult(nil)
		if got != terminatedExecResultPrefix {
			t.Fatalf("expected %q, got %q", terminatedExecResultPrefix, got)
		}
	})
	t.Run("Hive exec err preserves original text", func(t *testing.T) {
		orig := errors.New("hive connection terminated: hive exec wait failed (sql=\"...\"): Connection not open")
		got := terminatedExecResult(orig)
		if !strings.HasPrefix(got, terminatedExecResultPrefix+": ") {
			t.Fatalf("expected prefix %q+\":\", got %q", terminatedExecResultPrefix, got)
		}
		if !strings.Contains(got, "Connection not open") {
			t.Fatalf("expected original error preserved, got %q", got)
		}
	})
	t.Run("MySQL invalid connection preserves original text", func(t *testing.T) {
		orig := errors.New("invalid connection")
		got := terminatedExecResult(orig)
		if !strings.HasPrefix(got, terminatedExecResultPrefix+": ") {
			t.Fatalf("expected prefix, got %q", got)
		}
		if !strings.Contains(got, "invalid connection") {
			t.Fatalf("expected original mysql err preserved, got %q", got)
		}
	})
}

func Test_isConnectionTerminatedError(t *testing.T) {
	tests := map[string]struct {
		err      error
		dbType   string
		expected bool
	}{
		// ===== MySQL 终止错误 (正向匹配) =====
		"MySQL: raw ErrInvalidConn message": {
			err:      fmt.Errorf("invalid connection"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: true,
		},
		"MySQL: ErrInvalidConn wrapped in CodeError": {
			err:      fmt.Errorf("connect remote database error: invalid connection"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: true,
		},
		"MySQL: ErrInvalidConn wrapped in ExecBatch format": {
			err:      fmt.Errorf("exec sql failed: \nSELECT 1 \ninvalid connection"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: true,
		},
		"MySQL: ErrInvalidConn through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = invalid connection"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: true,
		},

		// ===== PostgreSQL 终止错误 (正向匹配) =====
		"PostgreSQL: 57P01 admin_shutdown from pgx": {
			err:      fmt.Errorf("FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: true,
		},
		"PostgreSQL: 57P01 through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: true,
		},
		"PostgreSQL: 57P01 wrapped in driver adaptor": {
			err:      fmt.Errorf("exec sql in driver adaptor: FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: true,
		},
		"PostgreSQL: conn closed": {
			err:      fmt.Errorf("conn closed"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: true,
		},
		"PostgreSQL: conn closed through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = conn closed"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: true,
		},

		// ===== GaussDB 终止错误 (正向匹配) =====
		"GaussDB: canceling statement due to user request": {
			err:      fmt.Errorf("pq: canceling statement due to user request"),
			dbType:   driverV2.DriverTypeGaussDB,
			expected: true,
		},
		"GaussDB: canceling statement through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = pq: canceling statement due to user request"),
			dbType:   driverV2.DriverTypeGaussDB,
			expected: true,
		},

		// ===== SQL Server 终止错误 (正向匹配) =====
		"SQL Server: connection is broken from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:The connection is broken and recovery is not possible.  The client driver attempted to recover the connection one or more times and all attempts failed."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: connection is broken from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:A severe error occurred on the current command. The connection is broken and recovery is not possible."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: connection is broken from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:The connection is broken and recovery is not possible."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: Timeout expired from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: Timeout expired from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: Timeout expired from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: session is in the kill state from exec failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Cannot continue the execution because the session is in the kill state."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: session is in the kill state from exec batch failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec batch failed:Cannot continue the execution because the session is in the kill state."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},
		"SQL Server: session is in the kill state from tx failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = tx failed:Cannot continue the execution because the session is in the kill state."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: true,
		},

		// ===== 跨库误判 (dbType 不匹配时不应识别为终止) =====
		"MySQL error with PostgreSQL dbType": {
			err:      fmt.Errorf("invalid connection"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},
		"PostgreSQL 57P01 with MySQL dbType": {
			err:      fmt.Errorf("FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: false,
		},
		"GaussDB cancel with PostgreSQL dbType": {
			err:      fmt.Errorf("pq: canceling statement due to user request"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},
		"SQL Server broken connection with MySQL dbType": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:The connection is broken and recovery is not possible."),
			dbType:   driverV2.DriverTypeMySQL,
			expected: false,
		},

		// ===== SQL Server 非终止错误 (反向匹配，不应误判) =====
		"SQL Server: login failed": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Login failed for user 'sa'."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: false,
		},
		"SQL Server: syntax error": {
			err:      fmt.Errorf("rpc error: code = Internal desc = exec failed:Incorrect syntax near 'SELECTT'."),
			dbType:   driverV2.DriverTypeSQLServer,
			expected: false,
		},

		// ===== 非终止错误 (反向匹配，不应误判) =====
		"nil error": {
			err:      nil,
			dbType:   driverV2.DriverTypeMySQL,
			expected: false,
		},
		"syntax error": {
			err:      fmt.Errorf("ERROR: syntax error at or near \"SELECTT\" (SQLSTATE 42601)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},
		"permission denied": {
			err:      fmt.Errorf("ERROR: permission denied for table users (SQLSTATE 42501)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},
		"connection refused": {
			err:      fmt.Errorf("connection refused"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: false,
		},
		"context deadline exceeded": {
			err:      fmt.Errorf("context deadline exceeded"),
			dbType:   driverV2.DriverTypeMySQL,
			expected: false,
		},
		"generic exec error": {
			err:      fmt.Errorf("exec sql failed: ERROR: relation \"nonexistent\" does not exist"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},
		"unique constraint violation": {
			err:      fmt.Errorf("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"),
			dbType:   driverV2.DriverTypePostgreSQL,
			expected: false,
		},

		// ===== Hive 终止错误 (正向匹配) =====
		// 与 sqle-hive-plugin/driver/hive.go::terminatedExecErrPrefix 之间的契约：
		// KillProcess 成功 issue cancel 之后，wrapHiveExecErrWithTermination 给
		// 后续所有 wrap 出口加 `hive connection terminated:` 前缀；sqled 端
		// 用一条规则把这类 SQL 行打成 terminate_succ。
		"Hive: terminated prefix bare": {
			err:      fmt.Errorf("hive connection terminated: hive exec wait failed (sql=%q): EOF", "SELECT reflect(...)"),
			expected: true,
		},
		"Hive: terminated prefix through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = hive connection terminated: hive exec failed (sql=\"...\"): Connection not open"),
			expected: true,
		},
		"Hive: terminated prefix from ExecBatch wrap": {
			err:      fmt.Errorf("exec sql failed: hive connection terminated: hive exec wait failed (sql=\"...\"): operation in state CANCELED"),
			expected: true,
		},
		// 即便 plugin 端没来得及加前缀（极端时序），sqled 也应在 hasTermination=true
		// 时直接通过下面的 raw 子串识别这几类常见 Cancel 后派生错误。
		"Hive: raw Connection not open": {
			err:      fmt.Errorf("hive exec wait failed (sql=\"reflect(...)\"): Connection not open"),
			expected: true,
		},
		"Hive: raw Connection not open through gRPC": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = Connection not open"),
			expected: true,
		},
		"Hive: raw Context was done before query executed": {
			err:      fmt.Errorf("hive exec failed (sql=\"reflect(...)\"): Context was done before the query was executed"),
			expected: true,
		},
		"Hive: raw Context is done": {
			err:      fmt.Errorf("rpc error: code = Unknown desc = Context is done"),
			expected: true,
		},
		"Hive: raw operation in state CANCELED": {
			err:      fmt.Errorf("hive exec failed (sql=\"...\"): Invalid OperationHandle: OperationHandle [opType=EXECUTE_STATEMENT, ...]: operation in state CANCELED"),
			expected: true,
		},

		// ===== Hive 非终止错误 (反向匹配，不应误判) =====
		// 普通连接握手失败（auth/transport 不匹配）—— EOF / connection reset by peer，
		// 既不在 Cancel 路径上，也不该被识别为终止成功。
		"Hive: connection reset (not after cancel)": {
			err:      fmt.Errorf("连接 Hive 失败 (auth=NONE): EOF; hint: 服务端在握手阶段立即关闭连接"),
			expected: false,
		},
		// "operation in state RUNNING / FINISHED" 等非 CANCELED 状态不应被识别为终止。
		"Hive: operation in state RUNNING": {
			err:      fmt.Errorf("hive query failed (sql=\"...\"): Invalid OperationHandle: operation in state RUNNING"),
			expected: false,
		},
		"Hive: operation in state FINISHED": {
			err:      fmt.Errorf("hive query failed (sql=\"...\"): operation in state FINISHED"),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isConnectionTerminatedError(tc.err, tc.dbType)
			assert.Equal(t, tc.expected, result, "case: %s", name)
		})
	}
}
