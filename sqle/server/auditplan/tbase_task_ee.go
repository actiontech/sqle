//go:build enterprise
// +build enterprise

package auditplan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

type TBasePgLog struct {
	*sqlCollector
}

func NewTBasePgLog(entry *logrus.Entry, ap *model.AuditPlan) Task {
	t := &TBasePgLog{newSQLCollector(entry, ap)}

	return t
}

func (at *TBasePgLog) FullSyncSQLs(sqls []*SQL) error {
	return at.baseTask.FullSyncSQLs(sqls)
}

func (at *TBasePgLog) PartialSyncSQLs(sqls []*SQL) error {
	return at.persist.UpdateSlowLogAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *TBasePgLog) Audit() (*AuditResultResp, error) {
	return auditWithSchema(at.logger, at.persist, at.ap)
}

func (at *TBasePgLog) GetSQLs(args map[string]interface{}) (
	[]Head, []map[string] /* head name */ string, uint64, error) {

	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: "SQL指纹",
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: "SQL",
			Type: "sql",
		},
		{
			Name: "counter",
			Desc: "数量",
		},
		{
			Name: "last_receive_timestamp",
			Desc: "最后匹配时间",
		},
		{
			Name: "average_query_time",
			Desc: "平均执行时间",
		},
		{
			Name: "max_query_time",
			Desc: "最长执行时间",
		},
		{
			Name: "db_user",
			Desc: "用户",
		},
		{
			Name: "schema",
			Desc: "Schema",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64   `json:"counter"`
			LastReceiveTimestamp string   `json:"last_receive_timestamp"`
			AverageQueryTime     *float64 `json:"query_time_avg"`
			MaxQueryTime         *float64 `json:"query_time_max"`
			DBUser               string   `json:"db_user"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, nil, 0, err
		}
		row := map[string]string{
			"sql":                    sql.SQLContent,
			"fingerprint":            sql.Fingerprint,
			"counter":                strconv.FormatUint(info.Counter, 10),
			"last_receive_timestamp": info.LastReceiveTimestamp,
			"db_user":                info.DBUser,
			"schema":                 sql.Schema,
		}

		// 兼容之前没有平均执行时间和最长执行时间的数据，没有数据的时候不会在前端显示0.00000导致误解
		if info.AverageQueryTime != nil {
			row["average_query_time"] = fmt.Sprintf("%.6f", *info.AverageQueryTime)
		}
		if info.MaxQueryTime != nil {
			row["max_query_time"] = fmt.Sprintf("%.6f", *info.MaxQueryTime)
		}
		rows = append(rows, row)
	}
	return head, rows, count, nil
}
