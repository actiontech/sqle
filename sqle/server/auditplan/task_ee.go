//go:build enterprise
// +build enterprise

package auditplan

const (
	OBMySQLIndicatorCPUTime     = "cpu_time"
	OBMySQLIndicatorIOWait      = "io_wait"
	OBMySQLIndicatorElapsedTime = "avg_elapsed_time_ms"
	SlowLogQueryNums            = 1000
)

// type SlowLogTask struct {
// 	*sqlCollector
// }

// func NewSlowLogTask(entry *logrus.Entry, ap *AuditPlan) Task {
// 	t := &SlowLogTask{newSQLCollector(entry, ap)}
// 	t.do = t.collectorDo

// 	return t
// }

const (
	slowlogCollectInputLogFile = 0 // FILE: mysql-slow.log
	slowlogCollectInputTable   = 1 // TABLE: mysql.slow_log
)

// func (at *SlowLogTask) FullSyncSQLs(sqls []*SQL) error {
// 	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() == slowlogCollectInputTable {
// 		return at.sqlCollector.FullSyncSQLs(sqls)
// 	}
// 	return at.baseTask.FullSyncSQLs(sqls)
// }

// func (at *SlowLogTask) PartialSyncSQLs(sqls []*SQL) error {
// 	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() == slowlogCollectInputTable {
// 		return at.sqlCollector.PartialSyncSQLs(sqls)
// 	}
// 	return at.persist.UpdateSlowLogAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
// }

// func (at *SlowLogTask) collectorDo() {

// 	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() != slowlogCollectInputTable {
// 		return
// 	}

// 	if at.ap.InstanceName == "" {
// 		at.logger.Warnf("instance is not configured")
// 		return
// 	}

// 	instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
// 	if err != nil {
// 		return
// 	}

// 	db, err := executor.NewExecutor(at.logger, &driverV2.DSN{
// 		Host:             instance.Host,
// 		Port:             instance.Port,
// 		User:             instance.User,
// 		Password:         instance.Password,
// 		AdditionalParams: instance.AdditionalParams,
// 		DatabaseName:     at.ap.InstanceDatabase,
// 	},
// 		at.ap.InstanceDatabase)
// 	if err != nil {
// 		at.logger.Errorf("connect to instance fail, error: %v", err)
// 		return
// 	}
// 	defer db.Db.Close()

// 	queryStartTime, err := at.persist.GetLatestStartTimeAuditPlanSQLV2(at.ap.ID)
// 	if err != nil {
// 		at.logger.Errorf("get start time failed, error: %v", err)
// 		return
// 	}
// 	querySQL := `
// 	SELECT sql_text,db,TIME_TO_SEC(query_time) AS query_time, start_time, rows_examined
// 	FROM mysql.slow_log
// 	WHERE sql_text != ''
// 		AND db NOT IN ('information_schema','performance_schema','mysql','sys')
// 	`

// 	sqlInfos := []*sqlFromSlowLog{}

// 	for {
// 		extraCondition := fmt.Sprintf(" AND start_time>'%s' ORDER BY start_time LIMIT %d", queryStartTime, SlowLogQueryNums)
// 		execQuerySQL := querySQL + extraCondition

// 		res, err := db.Db.Query(execQuerySQL)

// 		if err != nil {
// 			at.logger.Errorf("query slow log failed, error: %v", err)
// 			break
// 		}

// 		for i := range res {
// 			sqlInfo := &sqlFromSlowLog{
// 				sql:       res[i]["sql_text"].String,
// 				schema:    res[i]["db"].String,
// 				startTime: res[i]["start_time"].String,
// 			}
// 			queryTime, err := strconv.Atoi(res[i]["query_time"].String)
// 			if err != nil {
// 				at.logger.Warnf("unexpected data format: %v, ", res[i]["query_time"].String)
// 				continue
// 			}
// 			sqlInfo.queryTimeSeconds = queryTime
// 			sqlInfo.rowExamined, err = strconv.ParseFloat(res[i]["rows_examined"].String, 64)
// 			if err != nil {
// 				at.logger.Warnf("unexpected data format: %v, ", res[i]["rows_examined"].String)
// 				continue
// 			}

// 			sqlInfos = append(sqlInfos, sqlInfo)
// 		}

// 		if len(res) < SlowLogQueryNums {
// 			break
// 		}

// 		queryStartTime = res[len(res)-1]["start_time"].String

// 		time.Sleep(500 * time.Millisecond)
// 	}

// 	if len(sqlInfos) != 0 {

// 		sqlFingerprintInfos, err := sqlFromSlowLogs(sqlInfos).mergeByFingerprint(at.logger)
// 		if err != nil {
// 			at.logger.Errorf("merge finger sqls failed, error: %v", err)
// 			return
// 		}

// 		auditPlanSQLs := make([]*model.OriginManageSQL, len(sqlFingerprintInfos))
// 		{
// 			now := time.Now()
// 			for i := range sqlFingerprintInfos {
// 				fp := sqlFingerprintInfos[i]
// 				fpInfo := fmt.Sprintf(`{"counter":%v,"last_receive_timestamp":"%v","schema":"%v","average_query_time":%d,"start_time":"%v","row_examined_avg":%v}`,
// 					fp.counter, now.Format(time.RFC3339), fp.schema, fp.queryTimeSeconds, fp.startTime, fp.rowExaminedAvg)
// 				auditPlanSQLs[i] = &model.OriginManageSQL{
// 					Source:         at.ap.Type,
// 					SourceId:       at.ap.ID,
// 					ProjectId:      at.ap.ProjectId,
// 					InstanceName:   at.ap.InstanceName,
// 					SchemaName:     at.ap.SchemaName,
// 					SqlFingerprint: fp.fingerprint,
// 					SqlText:        fp.sql,
// 					Info:           []byte(fpInfo),
// 					// EndPoint:                  "",
// 				}
// 			}
// 		}

// 		if err = at.persist.UpdateSlowLogCollectAuditPlanSQLsV2(at.ap.ID, auditPlanSQLs); err != nil {
// 			at.logger.Errorf("save mysql slow log to storage failed, error: %v", err)
// 			return
// 		}
// 	}

// 	_, err = at.Audit()
// 	if err != nil {
// 		at.logger.Errorf("audit audit plan failed,error: %v", err)
// 	}
// }

// type sqlFromSlowLog struct {
// 	sql              string
// 	schema           string
// 	queryTimeSeconds int
// 	startTime        string
// 	rowExamined      float64
// }

// type sqlFromSlowLogs []*sqlFromSlowLog

// type sqlFingerprintInfo struct {
// 	lastSql               string
// 	lastSqlSchema         string
// 	sqlCount              int
// 	totalQueryTimeSeconds int
// 	startTime             string
// 	totalExaminedRows     float64
// }

// func (s *sqlFingerprintInfo) queryTime() int {
// 	return s.totalQueryTimeSeconds / s.sqlCount
// }

// func (s *sqlFingerprintInfo) rowExaminedAvg() float64 {
// 	return s.totalExaminedRows / float64(s.sqlCount)
// }

// func (s sqlFromSlowLogs) mergeByFingerprint(entry *logrus.Entry) ([]sqlInfo, error) {

// 	sqlInfos := []sqlInfo{}
// 	sqlInfosMap := map[string] /*sql fingerprint*/ *sqlFingerprintInfo{}

// 	for i := range s {
// 		sqlItem := s[i]
// 		fp, err := util.Fingerprint(sqlItem.sql, true)
// 		if err != nil {
// 			entry.Warnf("get sql finger print failed, err: %v", err)
// 		}
// 		if fp == "" {
// 			continue
// 		}

// 		if sqlInfosMap[fp] != nil {
// 			sqlInfosMap[fp].lastSql = sqlItem.sql
// 			sqlInfosMap[fp].lastSqlSchema = sqlItem.schema
// 			sqlInfosMap[fp].sqlCount++
// 			sqlInfosMap[fp].totalQueryTimeSeconds += sqlItem.queryTimeSeconds
// 			sqlInfosMap[fp].startTime = sqlItem.startTime
// 			sqlInfosMap[fp].totalExaminedRows += sqlItem.rowExamined
// 		} else {
// 			sqlInfosMap[fp] = &sqlFingerprintInfo{
// 				sqlCount:              1,
// 				lastSql:               sqlItem.sql,
// 				lastSqlSchema:         sqlItem.schema,
// 				totalQueryTimeSeconds: sqlItem.queryTimeSeconds,
// 				startTime:             sqlItem.startTime,
// 				totalExaminedRows:     sqlItem.rowExamined,
// 			}
// 			sqlInfos = append(sqlInfos, sqlInfo{fingerprint: fp})
// 		}
// 	}

// 	for i := range sqlInfos {
// 		fp := sqlInfos[i].fingerprint
// 		sqlInfo := sqlInfosMap[fp]
// 		if sqlInfo != nil {
// 			sqlInfos[i].counter = sqlInfo.sqlCount
// 			sqlInfos[i].sql = sqlInfo.lastSql
// 			sqlInfos[i].schema = sqlInfo.lastSqlSchema
// 			sqlInfos[i].queryTimeSeconds = sqlInfo.queryTime()
// 			sqlInfos[i].startTime = sqlInfo.startTime
// 			sqlInfos[i].rowExaminedAvg = utils.Round(sqlInfo.rowExaminedAvg(), 6)
// 		}
// 	}

// 	return sqlInfos, nil
// }

// func (at *SlowLogTask) GetSQLs(args map[string]interface{}) (
// 	[]Head, []map[string] /* head name */ string, uint64, error) {

// 	auditPlanSQLs, count, err := at.persist.GetInstanceAuditPlanSQLsByReq(args)
// 	if err != nil {
// 		return nil, nil, count, err
// 	}
// 	head := []Head{
// 		{
// 			Name: "fingerprint",
// 			Desc: "SQL指纹",
// 			Type: "sql",
// 		},
// 		{
// 			Name: "sql",
// 			Desc: "SQL",
// 			Type: "sql",
// 		},
// 		{
// 			Name: "counter",
// 			Desc: "数量",
// 		},
// 		{
// 			Name: "last_receive_timestamp",
// 			Desc: "最后匹配时间",
// 		},
// 		{
// 			Name: "average_query_time",
// 			Desc: "平均执行时间",
// 		},
// 		{
// 			Name: "max_query_time",
// 			Desc: "最长执行时间",
// 		},
// 		{
// 			Name: "row_examined_avg",
// 			Desc: "平均扫描行数",
// 		},
// 		{
// 			Name: "db_user",
// 			Desc: "用户",
// 		},
// 		{
// 			Name: "schema",
// 			Desc: "Schema",
// 		},
// 	}
// 	rows := make([]map[string]string, 0, len(auditPlanSQLs))
// 	for _, sql := range auditPlanSQLs {
// 		var info = struct {
// 			Counter              uint64   `json:"counter"`
// 			LastReceiveTimestamp string   `json:"last_receive_timestamp"`
// 			AverageQueryTime     *float64 `json:"query_time_avg"`
// 			MaxQueryTime         *float64 `json:"query_time_max"`
// 			RowExaminedAvg       *float64 `json:"row_examined_avg"`
// 			DBUser               string   `json:"db_user"`
// 		}{}
// 		err := json.Unmarshal(sql.Info, &info)
// 		if err != nil {
// 			return nil, nil, 0, err
// 		}
// 		row := map[string]string{
// 			"sql":                    sql.SQLContent,
// 			"fingerprint":            sql.Fingerprint,
// 			"counter":                strconv.FormatUint(info.Counter, 10),
// 			"last_receive_timestamp": info.LastReceiveTimestamp,
// 			"db_user":                info.DBUser,
// 			"schema":                 sql.Schema,
// 		}

// 		if info.RowExaminedAvg != nil {
// 			row["row_examined_avg"] = fmt.Sprintf("%.6f", *info.RowExaminedAvg)
// 		}
// 		// 兼容之前没有平均执行时间和最长执行时间的数据，没有数据的时候不会在前端显示0.00000导致误解
// 		if info.AverageQueryTime != nil {
// 			row["average_query_time"] = fmt.Sprintf("%.6f", *info.AverageQueryTime)
// 		}
// 		if info.MaxQueryTime != nil {
// 			row["max_query_time"] = fmt.Sprintf("%.6f", *info.MaxQueryTime)
// 		}
// 		rows = append(rows, row)
// 	}
// 	return head, rows, count, nil
// }

// // HACK: slow SQLs may be executed in different Schemas.
// // Before auditing sql, we need to insert a Schema switching statement.
// // And need to manually execute server.ReplenishTaskStatistics() to recalculate
// // real task object score
// func (at *SlowLogTask) Audit() (*AuditResultResp, error) {
// 	return auditWithSchema(at.logger, at.persist, at.ap)
// }

// PostgreSQLSchemaMetaTask : PostgreSQL库表元数据
