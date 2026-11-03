package auditplan

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/sirupsen/logrus"
)

const tencentEventKeySep = "\x1f"

// tencentMySQLEventKeyHash 组合键 version=1（S2）：
// db_instance_id + Timestamp + Md5(+SqlText 降级) + UserName + Database + UserHost → SHA256 hex。
func tencentMySQLEventKeyHash(instanceID string, item tencentSlowLogItem) string {
	md5 := strings.TrimSpace(item.Md5)
	if md5 == "" {
		md5 = "sql:" + sha256HexStr(item.SQLText)
	}
	ts := ""
	if !item.ExecutionStartTime.IsZero() {
		ts = item.ExecutionStartTime.UTC().Format(time.RFC3339)
	}
	raw := strings.Join([]string{
		instanceID,
		ts,
		md5,
		item.UserName,
		item.SchemaName,
		item.UserHost,
	}, tencentEventKeySep)
	return sha256HexStr(raw)
}

// tencentPGEventKeyHash 组合键 version=1（S2）：
// 优先 instance + SessionStartTime + SessionId + ProcessId + SQL 指纹；
// SessionId/ProcessId 皆空时降级 + 可观测日志。
func tencentPGEventKeyHash(logger *logrus.Entry, instanceID string, item tencentSlowLogItem) string {
	fp := tencentPGSQLFingerprint(item.SQLText)
	ts := ""
	if !item.ExecutionStartTime.IsZero() {
		ts = item.ExecutionStartTime.UTC().Format(time.RFC3339)
	}
	sid := strings.TrimSpace(item.SessionId)
	pid := strings.TrimSpace(item.ProcessId)
	var raw string
	if sid != "" || pid != "" {
		raw = strings.Join([]string{
			instanceID,
			ts,
			sid,
			pid,
			fp,
		}, tencentEventKeySep)
	} else {
		if logger != nil {
			logger.Warnf("tencent pg event key fallback: empty SessionId/ProcessId instance=%s ts=%s", instanceID, ts)
		}
		raw = strings.Join([]string{
			instanceID,
			ts,
			fp,
			item.UserName,
			item.SchemaName,
			item.UserHost,
		}, tencentEventKeySep)
	}
	return sha256HexStr(raw)
}

func tencentPGSQLFingerprint(sqlText string) string {
	if f, err := util.Fingerprint(sqlText, true); err == nil && f != "" {
		return f
	}
	return sqlText
}

func sha256HexStr(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
