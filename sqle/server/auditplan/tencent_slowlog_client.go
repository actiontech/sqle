package auditplan

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PG Duration 诊断键（S3 AC-018：未校准前禁止写入 query_time_*）。
const (
	metricTencentPGDurationRaw  = "tencent_pg_duration_raw"
	metricTencentPGDurationUnit = "tencent_pg_duration_unit"
	tencentPGDurationUnitPending = "pending"
)

// tencentSlowLogItem 统一采集条目。
// 可选数值用指针：nil=云侧缺字段（不写 Metrics）；非 nil（含 0）=显式值。
type tencentSlowLogItem struct {
	SQLText            string
	SchemaName         string
	UserName           string
	UserHost           string
	ExecutionStartTime time.Time
	QueryTimeSec       *float64 // MySQL QueryTime（秒）
	LockTimeSec        *float64 // MySQL LockTime（秒）
	RowsExamined       *int64   // MySQL；PG 无
	RowsSent           *int64   // MySQL；PG 无（现网无 rows_sent metric，可不落）
	Md5                string
	DurationRaw        string // PG Duration 原始字符串；单位见校准门禁
	SessionId          string
	ProcessId          string
}

func tencentFloat64(v float64) *float64 { return &v }
func tencentInt64(v int64) *int64       { return &v }

type tencentMySQLSlowLogReq struct {
	InstanceId string
	StartTime  int64
	EndTime    int64
	Offset     uint64
	Limit      uint64
}

type tencentMySQLSlowLogResp struct {
	TotalCount int64
	Items      []tencentSlowLogItem
}

type tencentPGSlowLogReq struct {
	DBInstanceId string
	StartTime    string // 2006-01-02 15:04:05
	EndTime      string
	Offset       uint64
	Limit        uint64
	DatabaseName string // optional
}

type tencentPGSlowLogResp struct {
	TotalCount int64
	Items      []tencentSlowLogItem
}

type tencentMySQLSlowLogAPI interface {
	DescribeSlowLogData(req tencentMySQLSlowLogReq) (*tencentMySQLSlowLogResp, error)
}

type tencentPGSlowLogAPI interface {
	DescribeSlowQueryList(req tencentPGSlowLogReq) (*tencentPGSlowLogResp, error)
}

type tencentMySQLOpenAPIClient struct {
	tc3 *tencentTC3Client
}

func newTencentMySQLOpenAPIClient(secretId, secretKey, region string) *tencentMySQLOpenAPIClient {
	host := fmt.Sprintf("cdb.%s.tencentcloudapi.com", region)
	return &tencentMySQLOpenAPIClient{
		tc3: newTencentTC3Client(secretId, secretKey, region, "cdb", host, "2017-03-20"),
	}
}

func (c *tencentMySQLOpenAPIClient) DescribeSlowLogData(req tencentMySQLSlowLogReq) (*tencentMySQLSlowLogResp, error) {
	limit := req.Limit
	if limit == 0 || limit > 800 {
		limit = 100
	}
	payload := map[string]interface{}{
		"InstanceId": req.InstanceId,
		"StartTime":  req.StartTime,
		"EndTime":    req.EndTime,
		"Offset":     req.Offset,
		"Limit":      limit,
		"SortBy":     "Timestamp",
		"OrderBy":    "ASC",
	}
	var raw struct {
		TotalCount int64 `json:"TotalCount"`
		Items      []struct {
			Timestamp    int64    `json:"Timestamp"`
			SqlText      string   `json:"SqlText"`
			Database     string   `json:"Database"`
			UserName     string   `json:"UserName"`
			UserHost     string   `json:"UserHost"`
			QueryTime    *float64 `json:"QueryTime"`
			LockTime     *float64 `json:"LockTime"`
			RowsExamined *int64   `json:"RowsExamined"`
			RowsSent     *int64   `json:"RowsSent"`
			Md5          string   `json:"Md5"`
		} `json:"Items"`
	}
	if err := c.tc3.doJSON("DescribeSlowLogData", payload, &raw); err != nil {
		return nil, err
	}
	out := &tencentMySQLSlowLogResp{TotalCount: raw.TotalCount}
	for _, it := range raw.Items {
		out.Items = append(out.Items, tencentSlowLogItem{
			SQLText:            it.SqlText,
			SchemaName:         it.Database,
			UserName:           it.UserName,
			UserHost:           it.UserHost,
			ExecutionStartTime: time.Unix(it.Timestamp, 0).UTC(),
			QueryTimeSec:       it.QueryTime,
			LockTimeSec:        it.LockTime,
			RowsExamined:       it.RowsExamined,
			RowsSent:           it.RowsSent,
			Md5:                it.Md5,
		})
	}
	return out, nil
}

type tencentPGOpenAPIClient struct {
	tc3 *tencentTC3Client
}

func newTencentPGOpenAPIClient(secretId, secretKey, region string) *tencentPGOpenAPIClient {
	host := fmt.Sprintf("postgres.%s.tencentcloudapi.com", region)
	return &tencentPGOpenAPIClient{
		tc3: newTencentTC3Client(secretId, secretKey, region, "postgres", host, "2017-03-12"),
	}
}

func (c *tencentPGOpenAPIClient) DescribeSlowQueryList(req tencentPGSlowLogReq) (*tencentPGSlowLogResp, error) {
	limit := req.Limit
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	payload := map[string]interface{}{
		"DBInstanceId": req.DBInstanceId,
		"StartTime":    req.StartTime,
		"EndTime":      req.EndTime,
		"Offset":       req.Offset,
		"Limit":        limit,
		"OrderBy":      "SessionStartTime",
		"OrderByType":  "asc",
	}
	if strings.TrimSpace(req.DatabaseName) != "" {
		payload["DatabaseName"] = req.DatabaseName
	}
	var raw struct {
		TotalCount int64 `json:"TotalCount"`
		Items      []struct {
			RawQuery         string      `json:"RawQuery"`
			DatabaseName     string      `json:"DatabaseName"`
			UserName         string      `json:"UserName"`
			ClientAddr       string      `json:"ClientAddr"`
			SessionStartTime string      `json:"SessionStartTime"`
			Duration         interface{} `json:"Duration"`
			SessionId        string      `json:"SessionId"`
			ProcessId        string      `json:"ProcessId"`
		} `json:"Items"`
	}
	if err := c.tc3.doJSON("DescribeSlowQueryList", payload, &raw); err != nil {
		return nil, err
	}
	out := &tencentPGSlowLogResp{TotalCount: raw.TotalCount}
	for _, it := range raw.Items {
		start, _ := parseTencentPGSessionStartTime(it.SessionStartTime)
		out.Items = append(out.Items, tencentSlowLogItem{
			SQLText:            it.RawQuery,
			SchemaName:         it.DatabaseName,
			UserName:           it.UserName,
			UserHost:           it.ClientAddr,
			ExecutionStartTime: start,
			DurationRaw:        formatTencentDurationRaw(it.Duration),
			SessionId:          it.SessionId,
			ProcessId:          it.ProcessId,
		})
	}
	return out, nil
}

func formatTencentDurationRaw(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(x, 10)
	case int:
		return strconv.Itoa(x)
	default:
		return fmt.Sprint(x)
	}
}

func parseTencentPGSessionStartTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty SessionStartTime")
	}
	// 常见：2006-01-02 15:04:05 / 带 CST / RFC3339
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05 MST",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	// 去尾缀 CST 再按本地墙钟按 UTC 解释（采集侧绝对化；时区细化可后续校准）
	trimmed := strings.TrimSpace(strings.TrimSuffix(s, "CST"))
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", trimmed, time.FixedZone("CST", 8*3600)); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unsupported SessionStartTime format")
}

func formatTencentPGTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}
